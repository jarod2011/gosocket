package server

import "github.com/jarod2011/gosocket/conn_repo"

type Options struct {
	ServerAddr    string
	ServerNetwork string
	Log           Logger
	Repo          conn_repo.Repo
	ClientMaximum int
}

type Option func(options *Options)

func WithOptions(opt Options) Option {
	return func(options *Options) {
		options.ServerAddr = opt.ServerAddr
		options.ServerNetwork = opt.ServerNetwork
		options.Log = opt.Log
		options.Repo = opt.Repo
		options.ClientMaximum = opt.ClientMaximum
	}
}

func WithServerAddr(addr string) Option {
	return func(options *Options) {
		options.ServerAddr = addr
	}
}

func WithServerNetwork(network string) Option {
	return func(options *Options) {
		options.ServerNetwork = network
	}
}

func WithLogger(log Logger) Option {
	return func(options *Options) {
		options.Log = log
	}
}

func WithRepo(repo conn_repo.Repo) Option {
	return func(options *Options) {
		options.Repo = repo
	}
}

func WithMaximumOnlineClients(maximum int) Option {
	return func(options *Options) {
		options.ClientMaximum = maximum
	}
}
