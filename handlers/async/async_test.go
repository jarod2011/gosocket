package async

import (
	"context"
	"github.com/jarod2011/gosocket/conn"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var n1, n2, n3, n4, n5, n6 int32

type mock struct {
	t *testing.T
}

func (m *mock) OnWorkProcessStop() {
	m.t.Logf("worker process stoped")
}

func (m *mock) OnReadProcessStop() {
	m.t.Logf("reader process stoped")
}

func (m *mock) OnWriteProcessStop() {
	m.t.Logf("writer process stoped")
}

func (m *mock) OnConnect(cc conn.Conn) {
	m.t.Logf("%s connected", cc.RemoteAddr())
	atomic.AddInt32(&n1, 1)
}

func (m *mock) SliceIndex(b []byte) int {
	m.t.Logf("slice %x", b)
	atomic.AddInt32(&n2, int32(len(b)))
	if len(b) >= 2 {
		return 2
	}
	return len(b)
}

func (m *mock) OnWork(b []byte, writeChan chan<- []byte) error {
	atomic.AddInt32(&n3, int32(len(b)))
	atomic.AddInt32(&n6, 1)
	writeChan <- b[0:1]
	return nil
}

func (m *mock) OnWriteFinish(b []byte) {
	m.t.Logf("write %x", b)
	atomic.AddInt32(&n4, int32(len(b)))
}

func (m *mock) OnWriteError(b []byte, err error) bool {
	m.t.Error(err)
	return false
}

func (m *mock) OnClose() {
	atomic.AddInt32(&n5, 1)
}

func newMockHandle(t *testing.T) Handler {
	return &mock{t: t}
}

func TestNew(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	h := New(Option{
		ReadTimeout:  0,
		WriteTimeout: 0,
		Creator:      nil,
		Context:      ctx,
	})
	c1, c2 := net.Pipe()
	cc1 := conn.New(c1)
	cc2 := conn.New(c2)
	if err := h(ctx, cc1); err == nil {
		t.Error("should nil handle error")
	}
	h = New(Option{
		Creator: func() Handler {
			return newMockHandle(t)
		},
		ReadTimeout: time.Millisecond * 100,
	})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := h(ctx, cc1)
		if err != nil {
			t.Error(err)
		}
	}()
	time.Sleep(time.Millisecond * 50)
	if n1 != 1 {
		t.Errorf("err n1 %d", n1)
	}
	cnt, err := cc2.WriteUntil([]byte{0x01, 0x00, 0x00, 0x00, 0x01, 0x01, 0x05}, time.Now().Add(time.Second*2))
	if err != nil && err != conn.ErrContextDeadline {
		t.Error(err)
	}
	if cnt != 7 {
		t.Errorf("err cnt %d", cnt)
	}
	by, err := cc2.ReadUntil(time.Now().Add(time.Second*10), true)
	if err != nil && err != conn.ErrContextDeadline {
		t.Error(err)
	}
	if len(by) != 2 {
		t.Errorf("read err count bytes %d", len(by))
	}
	time.Sleep(time.Second * 3)
	if n2 != 7+5+3+1 {
		t.Errorf("err n2 %d", n2)
	}
	if n3 != 7 {
		t.Errorf("err n3 %d", n3)
	}
	if n4 != n6 {
		t.Errorf("err n4 %d", n4)
	}
	if n5 != 0 {
		t.Errorf("err n5 %d", n5)
	}
	time.Sleep(time.Second * 1)
	cancel()
	wg.Wait()
	time.Sleep(time.Second)
	if n5 != 1 {
		t.Errorf("err n5 %d", n5)
	}
}
