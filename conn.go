package gosocket

import (
	"net"
	"sync/atomic"
	"time"
)

type Conn struct {
	net.Conn
	readTotal  uint64
	writeTotal uint64
	sec        int64
	nsec       int64
}

func (cc *Conn) Read(b []byte) (n int, err error) {
	n, err = cc.Conn.Read(b)
	atomic.AddUint64(&cc.readTotal, uint64(n))
	return
}

func (cc *Conn) Write(b []byte) (n int, err error) {
	n, err = cc.Conn.Write(b)
	atomic.AddUint64(&cc.writeTotal, uint64(n))
	return
}

func (cc Conn) Time() time.Time {
	return time.Unix(cc.sec, cc.nsec)
}

func newConn(conn net.Conn) *Conn {
	now := time.Now()
	return &Conn{
		Conn:       conn,
		readTotal:  0,
		writeTotal: 0,
		sec:        now.Unix(),
		nsec:       int64(now.Nanosecond()),
	}
}
