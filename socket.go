package gosocket

import (
	_ "github.com/jarod2011/gosocket/buffers"
	"github.com/jarod2011/gosocket/client"
	_ "github.com/jarod2011/gosocket/conn"
	_ "github.com/jarod2011/gosocket/conn_repo"
	"github.com/jarod2011/gosocket/server"
	_ "github.com/jarod2011/gosocket/util"
)

func NewServer(handler server.Handler, opts ...server.Option) server.Server {
	return server.New(handler, opts...)
}

func NewClient(opts ...client.Option) client.Client {
	return client.New(opts...)
}
