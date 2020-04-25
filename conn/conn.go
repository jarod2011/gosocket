package conn

import (
	"net"
	"sync/atomic"
	"time"
)

type Conn interface {
	net.Conn
	WriteBytes() uint64
	ReadBytes() uint64
	Created() time.Time
}

type conn struct {
	net.Conn
	writeBytes uint64
	readBytes  uint64
	sec        int64
	nsec       int64
}

func (c *conn) Read(b []byte) (int, error) {
	n, err := c.Conn.Read(b)
	atomic.AddUint64(&c.readBytes, uint64(n))
	return n, err
}

func (c *conn) Write(b []byte) (int, error) {
	n, err := c.Conn.Write(b)
	atomic.AddUint64(&c.writeBytes, uint64(n))
	return n, err
}

func (c *conn) WriteBytes() uint64 {
	return atomic.LoadUint64(&c.writeBytes)
}

func (c *conn) ReadBytes() uint64 {
	return atomic.LoadUint64(&c.readBytes)
}

func (c conn) Created() time.Time {
	return time.Unix(c.sec, c.nsec)
}

func New(cc net.Conn) Conn {
	conn := new(conn)
	conn.Conn = cc
	tm := time.Now()
	conn.sec = tm.Unix()
	conn.nsec = int64(tm.Nanosecond())
	return conn
}
