package conn

import (
	"github.com/jarod2011/gosocket/buffers"
	"net"
	"sync"
	"testing"
	"time"
)

func TestConn(t *testing.T) {
	c1, c2 := net.Pipe()
	con1 := New(c1)
	con2 := New(c2)
	buf := buffers.New()
	b1 := buf.Get()
	defer buf.Put(b1)
	b2 := buf.Get()
	defer buf.Put(b2)
	var nr1, nw1, nw2, nr2 int
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		defer con1.Close()
		con1.SetReadDeadline(time.Now().Add(time.Second * 1))
		n, err := con1.Write(append([]byte{0x11}))
		t.Log(err)
		nw1 += n
		for {
			con1.SetReadDeadline(time.Now().Add(time.Millisecond * 50))
			m, err := b1.ReadFrom(con1)
			t.Log(err)
			nr1 += int(m)
			if m > 0 {
				return
			}
		}
	}()
	go func() {
		defer wg.Done()
		defer con2.Close()
		for {
			con2.SetReadDeadline(time.Now().Add(time.Millisecond * 50))
			n, err := b2.ReadFrom(con2)
			t.Log(err)
			nr2 += int(n)
			if n > 0 {
				m, err := con2.Write([]byte{0x11, 0x12})
				t.Log(err)
				nw2 += m
				return
			}
		}
	}()
	wg.Wait()
	t.Logf("con1 read %d write %d", nr1, nw1)
	t.Logf("con2 read %d write %d", nr2, nw2)
	if con1.ReadBytes() != uint64(nr1) {
		t.Errorf("%d != %d", nr1, con1.ReadBytes())
	}
	if con1.WriteBytes() != uint64(nw1) {
		t.Errorf("%d != %d", nw1, con1.WriteBytes())
	}
	if con2.ReadBytes() != uint64(nr2) {
		t.Errorf("%d != %d", nr2, con2.ReadBytes())
	}
	if con2.WriteBytes() != uint64(nw2) {
		t.Errorf("%d != %d", nw2, con2.WriteBytes())
	}
	if con1.WriteBytes() != con2.ReadBytes() {
		t.Errorf("%d != %d", con1.WriteBytes(), con2.ReadBytes())
	}
	if con1.ReadBytes() != con2.WriteBytes() {
		t.Errorf("%d != %d", con1.ReadBytes(), con2.WriteBytes())
	}
	t.Log(con1.Created(), con2.Created())
}
