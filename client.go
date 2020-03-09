package gosocket

import (
	"context"
	"net"
	"sync"
	"time"
)

type Client interface {
	Handle() error
	Close() error
}

type ClientCreator func(cc net.Conn) Client

var ReadTimeout = time.Second * 10
var WriteTimeout = time.Second * 10

type client struct {
	sync.RWMutex
	conn         net.Conn
	ctx          context.Context
	cancel       context.CancelFunc
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func (c *client) Handle() error {

}

func (c *client) Close() error {
	c.Lock()
	defer c.Unlock()
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func newClient(ctx context.Context, cc net.Conn) Client {
	c, cancel := context.WithCancel(ctx)
	return &client{conn: cc, ctx: c, cancel: cancel}
}

func NewDefaultClientCreator() ClientCreator {
	return func(cc net.Conn) Client {

	}
}
