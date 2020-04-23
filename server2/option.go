package server2

import (
	"context"
	"net"
)

type Options struct {
	Address        net.Addr
	MaxConnection  int64
	ctx            context.Context
	cancel         context.CancelFunc
	Store          ConnectionStore
	HandlerCreator func() ConnectionHandler
}

func (opt Options) Validate() error {
	return nil
}

type Option func(options *Options)

func WithOptions(opt Options) Option {
	return func(options *Options) {
		options = &opt
	}
}

func WithAddress(addr net.Addr) Option {
	return func(options *Options) {
		options.Address = addr
	}
}

func WithContext(ctx context.Context) Option {
	return func(options *Options) {
		options.ctx, options.cancel = context.WithCancel(ctx)
	}
}

func WithMaxConnection(maximum int64) Option {
	return func(options *Options) {
		options.MaxConnection = maximum
	}
}

func WithConnectionStore(store ConnectionStore) Option {
	return func(options *Options) {
		options.Store = store
	}
}
