package client

import (
	"time"
)

// Options
// Client Options
type Options struct {
	ServerAddr    string        // Server Address
	ServerNetwork string        // Server Network
	ReadTimeout   time.Duration // read timeout duration
	WriteTimeout  time.Duration // write timeout duration
}

// Option
// callback to set Options value
type Option func(options *Options)

// WithOptions
// Set new Options overwrite current Options
func WithOptions(opt Options) Option {
	return func(options *Options) {
		options.ServerNetwork = opt.ServerNetwork
		options.WriteTimeout = opt.WriteTimeout
		options.ReadTimeout = opt.ReadTimeout
		options.ServerAddr = opt.ServerAddr
	}
}

// WithServerAddr
// Set new address overwrite current Options.ServerAddr
func WithServerAddr(addr string) Option {
	return func(options *Options) {
		options.ServerAddr = addr
	}
}

// WithServerNetwork
// Set new address overwrite current Options.ServerNetwork
func WithServerNetwork(network string) Option {
	return func(options *Options) {
		options.ServerNetwork = network
	}
}

// WithReadTimeout
// Set new read timeout duration overwrite current Options.ReadTimeout
func WithReadTimeout(duration time.Duration) Option {
	return func(options *Options) {
		options.ReadTimeout = duration
	}
}

// WithWriteTimeout
// Set new write timeout duration overwrite current Options.WriteTimeout
func WithWriteTimeout(duration time.Duration) Option {
	return func(options *Options) {
		options.WriteTimeout = duration
	}
}
