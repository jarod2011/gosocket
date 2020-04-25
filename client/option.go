package client

import (
	"github.com/jarod2011/gosocket/conn"
	"time"
)

// Options
// Client Options
type Options struct {
	ServerAddr     string            // Server Address
	ReadTimeout    time.Duration     // read timeout duration
	WriteTimeout   time.Duration     // write timeout duration
	OnClosedHandle OnClosed          // handle on closed connection
	Handler        ConnectionHandler // handle on connection
}

// Option
// callback to set Options value
type Option func(options *Options)

// OnHandler
// handle the connection
type ConnectionHandler func(cc conn.Conn, opt Options)

// OnClosed
// handle the connection closed
type OnClosed func(cc conn.Conn, closed bool, err error)

// WithOptions
// Set new Options overwrite current Options
func WithOptions(opt Options) Option {
	return func(options *Options) {
		options = &opt
	}
}

// WithServerAddr
// Set new address overwrite current Options.ServerAddr
func WithServerAddr(addr string) Option {
	return func(options *Options) {
		options.ServerAddr = addr
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

func WithOnClosed(handle OnClosed) Option {
	return func(options *Options) {
		options.OnClosedHandle = handle
	}
}

func WithOnHandleConnection(handle ConnectionHandler) Option {
	return func(options *Options) {
		options.Handler = handle
	}
}
