package client

import (
	"errors"
	"github.com/jarod2011/gosocket/conn"
	"net"
)

// Client
// socket client interface
type Client interface {
	Options() Options // get client current options
	Connect() error
	Handling() error
	Close() error
}

type client struct {
	opt Options
	cc  conn.Conn
}

func (c *client) Handling() error {
	if c.opt.Handler == nil {
		return errors.New("handler required")
	}
	c.opt.Handler(c.cc)
	return nil
}

func (c *client) Close() error {
	return c.cc.Close()
}

func (c client) Options() Options {
	return c.opt
}

func (c *client) Connect() error {
	cc, err := net.Dial("tcp", c.opt.ServerAddr)
	if err != nil {
		return err
	}
	c.cc = conn.New(cc)
	return nil
}

// New
// create socket client by optional Option
func New(opts ...Option) Client {
	c := new(client)
	for _, o := range opts {
		o(&c.opt)
	}
	return c
}
