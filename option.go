package gosocket

import (
	"net"
)

type Options struct {
	Addr net.Addr
	Log  Logger
}

type Option func(options *Options)

var DefaultPort = 8080

func WithAddr(addr net.Addr) Option {
	return func(options *Options) {
		options.Addr = addr
	}
}

func WithPort(port uint32) Option {
	return func(options *Options) {
		options.Addr = &net.TCPAddr{Port: int(port)}
	}
}

func WithLogger(log Logger) Option {
	return func(options *Options) {
		options.Log = log
	}
}
