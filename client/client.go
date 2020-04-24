package client

// Client
// socket client interface
type Client interface {
	Options() Options // get client current options
}

type client struct {
	opt Options
}

func (c client) Options() Options {
	return c.opt
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
