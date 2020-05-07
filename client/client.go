package client

import (
	"github.com/jarod2011/gosocket/buffers"
	"github.com/jarod2011/gosocket/conn"
	"net"
	"time"
)

var buf = buffers.New()

// Client
// socket client interface
type Client interface {
	Options() Options // get client current options
	Init(opts ...Option) error
	Send([]byte) (int, error)
	Recv() ([]byte, error)
	Close() error
}

type client struct {
	opt Options
	cc  conn.Conn
}

func (c *client) Send(b []byte) (cnt int, err error) {
	return c.cc.WriteUntil(b, time.Now().Add(c.Options().WriteTimeout))
}

func (c *client) Recv() ([]byte, error) {
	return c.cc.ReadUntil(time.Now().Add(c.Options().ReadTimeout), true)
}

func (c *client) Close() error {
	return c.cc.Close()
}

func (c client) Options() Options {
	return c.opt
}

func (c *client) Init(opts ...Option) error {
	for _, o := range opts {
		o(&c.opt)
	}
	cc, err := net.Dial(c.opt.ServerNetwork, c.opt.ServerAddr)
	if err != nil {
		return err
	}
	c.cc = conn.New(cc)
	return nil
}

var DefaultTimeout = time.Second * 10

// New
// create socket client by optional Option
func New(opts ...Option) Client {
	c := client{opt: Options{
		ServerNetwork: "tcp",
		ReadTimeout:   DefaultTimeout,
		WriteTimeout:  DefaultTimeout,
	}}
	for _, o := range opts {
		o(&c.opt)
	}
	return &c
}
