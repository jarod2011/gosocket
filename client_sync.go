package gosocket

import (
	"net"
)

type syncClient struct {
	handler HandleFunc
	conn    net.Conn
}

func (client *syncClient) Handle() error {
	return client.handler(client.conn)
}

func (client *syncClient) Close() error {
	return client.conn.Close()
}

type HandleFunc func(conn net.Conn) error

func NewSyncClientCreator(handleFunc HandleFunc) ClientCreator {
	return func(cc net.Conn) Client {
		return &syncClient{
			conn:    cc,
			handler: handleFunc,
		}
	}
}
