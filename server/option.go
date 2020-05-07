package server

import (
	"github.com/jarod2011/gosocket/conn_repo"
	"time"
)

type Options struct {
	ServerAddr      string
	ServerNetwork   string
	Log             Logger
	Repo            conn_repo.Repo
	ClientMaximum   int
	PrintDebug      bool
	MaxFreeDuration time.Duration
}

type Option func(options *Options)

func WithOptions(opt Options) Option {
	return func(options *Options) {
		options.ServerAddr = opt.ServerAddr
		options.ServerNetwork = opt.ServerNetwork
		options.Log = opt.Log
		options.Repo = opt.Repo
		options.ClientMaximum = opt.ClientMaximum
		options.MaxFreeDuration = opt.MaxFreeDuration
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

func WithEnableDebug() Option {
	return func(options *Options) {
		options.Log.EnableDebug()
	}
}

func WithMaxFreeDuration(duration time.Duration) Option {
	return func(options *Options) {
		options.MaxFreeDuration = duration
	}
}
