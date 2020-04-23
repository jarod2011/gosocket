package server

import (
	"net"
	"sync/atomic"
	"time"
)

type Conn interface {
	net.Conn
	Summary
	UUID() int64
	Closed() bool
	Done() <-chan struct{}
}

type conn struct {
	net.Conn
	read   int64
	write  int64
	sec    int64
	nsec   int64
	closed uint32
	done   chan struct{}
	id     int64
}

func (c conn) Done() <-chan struct{} {
	return c.done
}

func (c conn) UUID() int64 {
	return c.id
}

func (c conn) Closed() bool {
	return StateClosed == atomic.LoadUint32(&c.closed)
}

func (c conn) ConnectionAt() time.Time {
	return time.Unix(c.sec, c.nsec)
}

func (c *conn) Read(b []byte) (n int, err error) {
	if c.Closed() {
		return 0, ErrConnectionClosed
	}
	n, err = c.Conn.Read(b)
	atomic.AddInt64(&c.read, int64(n))
	return
}

func (c *conn) Write(b []byte) (n int, err error) {
	if c.Closed() {
		return 0, ErrConnectionClosed
	}
	n, err = c.Conn.Write(b)
	atomic.AddInt64(&c.write, int64(n))
	return
}

func (c *conn) Close() error {
	if atomic.CompareAndSwapUint32(&c.closed, StateOpened, StateClosed) {
		err := c.Conn.Close()
		close(c.done)
		return err
	}
	return nil
}

func (c conn) WriteBytes() int64 {
	return atomic.LoadInt64(&c.write)
}

func (c conn) ReadBytes() int64 {
	return atomic.LoadInt64(&c.read)
}

func newConn(c net.Conn) Conn {
	t := time.Now()
	return &conn{
		Conn:   c,
		sec:    t.Unix(),
		nsec:   int64(t.Nanosecond()),
		closed: StateOpened,
		id:     atomic.AddInt64(&uniqueId, 1),
	}
}
