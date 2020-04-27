package conn

import (
	"errors"
	"github.com/jarod2011/gosocket/buffers"
	"github.com/jarod2011/gosocket/util"
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
	SetSeparator(separator uint8)
	GetSeparator() uint8
}

type conn struct {
	net.Conn
	writeBytes uint64
	readBytes  uint64
	sec        int64
	nsec       int64
	separator  uint8
}

func (c *conn) SetSeparator(separator uint8) {
	c.separator = separator
}

func (c *conn) GetSeparator() uint8 {
	return c.separator
}

func (c *conn) ReadToUntil(writer io.Writer, expire time.Time) (int, error) {
	by, err := c.ReadUntil(expire)
	if len(by) > 0 {
		cnt, _ := writer.Write(by)
		return cnt, err
	}
	return 0, err
}

func (c *conn) ReadUntil(expire time.Time) (buf []byte, err error) {
	timeout := time.After(time.Until(expire))
	b := bufferPool.Get()
	defer bufferPool.Put(b)
	var cnt int64
	nextout := false
	for {
		select {
		case <-timeout:
			err = errors.New("context deadline")
			return
		default:
			if err = c.SetReadDeadline(time.Now().Add(time.Millisecond * 50)); err != nil {
				return
			}
			n, err1 := b.ReadFrom(c)
			cnt += n
			if n > 0 {
				buf = append(buf, b.Next(int(n))...)
				nextout = buf[len(buf)-1] == c.separator
			}
			if err1 != nil {
				if util.IsTimeout(err1) && nextout {
					buf = buf[0 : len(buf)-1]
					return
				}
				if err1 == io.EOF {
					return
				}
			}
		}
	}
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
	if c.separator > 0 {
		c.Write([]byte{c.separator})
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
