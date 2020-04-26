package conn

import (
	"errors"
	"github.com/jarod2011/gosocket/buffers"
	"io"
	"net"
	"sync/atomic"
	"time"
)

var bufferPool = buffers.New()

type Conn interface {
	net.Conn
	ReadUntil(expire time.Time) ([]byte, error)
	ReadToUntil(writer io.Writer, expire time.Time) (int, error)
	WriteUntil(b []byte, expire time.Time) (int, error)
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

func (c *conn) ReadToUntil(writer io.Writer, expire time.Time) (int, error) {
	by, err := c.ReadUntil(expire)
	if len(by) > 0 {
		cnt, _ := writer.Write(by)
		return cnt, err
	}
	return 0, err
}

func (c *conn) ReadUntil(expire time.Time) ([]byte, error) {
	timeout := time.After(time.Until(expire))
	b := bufferPool.Get()
	defer bufferPool.Put(b)
	var cnt int64
	var err error
	for {
		select {
		case <-timeout:
			err = errors.New("context deadline")
			goto result
		default:
			if err = c.SetReadDeadline(time.Now().Add(time.Millisecond * 50)); err != nil {
				goto result
			}
			n, err1 := b.ReadFrom(c)
			cnt += n
			if err1 != nil && err1 == io.EOF {
				goto result
			}
		}
	}
result:
	return b.Next(int(cnt)), err
}

func (c *conn) WriteUntil(b []byte, expire time.Time) (cnt int, err error) {
	timeout := time.After(time.Until(expire))
	for cnt < len(b) {
		select {
		case <-timeout:
			err = errors.New("context deadline")
			return
		default:
			if err = c.SetWriteDeadline(time.Now().Add(time.Millisecond * 50)); err != nil {
				return
			}
			n, err1 := c.Write(b[cnt:])
			cnt += n
			if err1 != nil && err1 == io.EOF {
				return
			}
		}
	}
	return
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
