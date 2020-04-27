package conn

import (
	"net"
	"sync"
	"testing"
	"time"
)

func TestConn(t *testing.T) {
	c1, c2 := net.Pipe()
	con1 := New(c1)
	con2 := New(c2)

	b1 := bufferPool.Get()
	defer bufferPool.Put(b1)
	b2 := bufferPool.Get()
	defer bufferPool.Put(b2)
	var nr1, nw1, nw2, nr2 int
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		defer con1.Close()
		n, err := con1.WriteUntil([]byte{0x11}, time.Now().Add(time.Microsecond*500))
		t.Log(err)
		nw1 += n
		nr1, err = con1.ReadToUntil(b1, time.Now().Add(time.Second))
		t.Log(err)
	}()
	go func() {
		defer wg.Done()
		defer con2.Close()
		n, err := con2.ReadToUntil(b2, time.Now().Add(time.Millisecond*500))
		if n > 0 {
			nw2, err = con2.WriteUntil([]byte{0x11, 0x12}, time.Now().Add(time.Second))
			t.Log(err)
		}
		nr2 += n
		t.Log(err)
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

func genClientServer(t *testing.T) (net.Conn, net.Conn) {
	lis, err := net.Listen("tcp", ":19934")
	if err != nil {
		t.Fatal(err)
	}
	client, err := net.Dial(lis.Addr().Network(), lis.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	cc1, err := lis.Accept()
	if err != nil {
		t.Fatal(err)
	}
	return cc1, client
}

func TestConn2(t *testing.T) {
	c1, c2 := genClientServer(t)
	con1 := New(c1)
	con2 := New(c2)
	defer con1.Close()
	defer con2.Close()
	b1 := bufferPool.Get()
	defer bufferPool.Put(b1)
	b2 := bufferPool.Get()
	defer bufferPool.Put(b2)
	con1.SetSeparator(0xfe)
	t.Log(con1.GetSeparator())
	con2.SetSeparator(0xfe)
	t.Log(con2.GetSeparator())
	t1 := time.Now()
	var wg1 sync.WaitGroup
	t.Log("start check")
	wg1.Add(1)
	go func() {
		t.Log("read in routine 5 seconds")
		con2.ReadUntil(time.Now().Add(time.Second * 5))
		wg1.Done()
	}()
	t.Log("write in one second")
	con1.WriteUntil([]byte{0x11, 0x22}, time.Now().Add(time.Second))
	t.Log("write ok")
	wg1.Wait()
	t.Log("read ok")
	if time.Since(t1).Seconds() > 4 {
		t.Errorf("use separator should not use more time")
	}
	con1.ReadToUntil(b1, time.Now().Add(time.Millisecond*100))
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		n, err := con2.ReadToUntil(b2, time.Now().Add(time.Millisecond*200))
		t.Log(n, err)
		wg.Done()
	}()
	con1.WriteUntil([]byte{0x11, 0x22}, time.Now().Add(time.Millisecond*100))
	con1.Close()
	con1.WriteUntil([]byte{0x11, 0x22}, time.Now().Add(time.Millisecond*100))
	time.Sleep(time.Millisecond * 300)
	con1.WriteUntil([]byte{0x11, 0x22}, time.Now().Add(time.Millisecond*100))
	wg.Wait()
	t.Log(b2.Bytes())
	go con1.ReadUntil(time.Now().Add(time.Millisecond * 100))
	con2.Write([]byte{0x11})
	con2.Close()
	con1.WriteUntil([]byte{0x11}, time.Now().Add(time.Millisecond*100))
	time.Sleep(time.Second)
}
