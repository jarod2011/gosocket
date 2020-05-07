package conn

import (
	"github.com/jarod2011/gosocket/buffers"
	"github.com/jarod2011/gosocket/util"
	"io"
	"net"
	"sync/atomic"
	"time"
)

var bufferPool = buffers.New()
var uuid uint64

type Conn interface {
	net.Conn
	ReadUntil(expire time.Time, whenEmpty bool) ([]byte, error)
	ReadToUntil(writer io.Writer, expire time.Time) (int, error)
	WriteUntil(b []byte, expire time.Time) (int, error)
	WriteBytes() uint64
	ReadBytes() uint64
	Created() time.Time
	LastActive() time.Time
	SetSeparator(separator uint8)
	GetSeparator() uint8
	ID() uint64
}

type conn struct {
	id uint64
	net.Conn
	writeBytes     uint64
	readBytes      uint64
	sec            int64
	nsec           int64
	separator      uint8
	lastActiveSec  int64
	lastActiveNSec int64
}

func (c *conn) LastActive() time.Time {
	return time.Unix(c.lastActiveSec, c.lastActiveNSec)
}

func (c *conn) SetSeparator(separator uint8) {
	c.separator = separator
}

func (c *conn) GetSeparator() uint8 {
	return c.separator
}

func (c *conn) ReadToUntil(writer io.Writer, expire time.Time) (int, error) {
	by, err := c.ReadUntil(expire, true)
	if len(by) > 0 {
		cnt, _ := writer.Write(by)
		return cnt, err
	}
	return 0, err
}

func (c *conn) ReadUntil(expire time.Time, whenEmpty bool) (buf []byte, err error) {
	timeout := time.After(time.Until(expire))
	b := bufferPool.Get()
	defer bufferPool.Put(b)
	defer func() {
		if util.IsRemoteClosedError(err) || util.IsClosedConnection(err) || err == io.EOF {
			err = io.EOF
		}
	}()
	var cnt int64
	nextout := false
	for {
		select {
		case <-timeout:
			if len(buf) == 0 {
				err = ErrContextDeadline
			}
			return
		default:
			if err = c.Conn.SetReadDeadline(time.Now().Add(time.Millisecond * 50)); err != nil {
				return
			}
			n, err1 := b.ReadFrom(c)
			cnt += n
			if n > 0 {
				buf = append(buf, b.Next(int(n))...)
				if buf[len(buf)-1] == c.separator {
					nextout = true
				} else if !whenEmpty && cnt > 0 {
					nextout = true
				} else {
					nextout = false
				}
			}
			if err1 != nil {
				if util.IsTimeout(err1) {
					if nextout {
						buf = buf[0 : len(buf)-1]
						return
					}
					continue
				}
			}
			err = err1
		}
	}
}

func (c *conn) WriteUntil(b []byte, expire time.Time) (cnt int, err error) {
	defer func() {
		if util.IsRemoteClosedError(err) || util.IsClosedConnection(err) || err == io.EOF {
			err = io.EOF
		}
	}()
	timeout := time.After(time.Until(expire))
	for cnt < len(b) {
		select {
		case <-timeout:
			err = ErrContextDeadline
			return
		default:
			if err = c.Conn.SetWriteDeadline(time.Now().Add(time.Millisecond * 50)); err != nil {
				return
			}
			n, err1 := c.Write(b[cnt:])
			cnt += n
			if err1 != nil {
				if !util.IsTimeout(err1) {
					err = err1
					return
				}
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
	if n > 0 {
		tm := time.Now()
		c.lastActiveSec = tm.Unix()
		c.lastActiveNSec = int64(tm.Nanosecond())
	}
	return n, err
}

func (c *conn) Write(b []byte) (int, error) {
	n, err := c.Conn.Write(b)
	atomic.AddUint64(&c.writeBytes, uint64(n))
	if n > 0 {
		tm := time.Now()
		c.lastActiveSec = tm.Unix()
		c.lastActiveNSec = int64(tm.Nanosecond())
	}
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

func (c conn) ID() uint64 {
	return c.id
}

func New(cc net.Conn) Conn {
	conn := new(conn)
	conn.Conn = cc
	conn.id = atomic.AddUint64(&uuid, 1)
	tm := time.Now()
	conn.sec = tm.Unix()
	conn.nsec = int64(tm.Nanosecond())
	conn.lastActiveSec = conn.sec
	conn.lastActiveNSec = conn.nsec
	return conn
}
