package gosocket

import (
	"net"
)

type syncClient struct {
	handler HandleFunc
	conn    *Conn
	uuid    interface{}
}

func (client *syncClient) Info() ClientInfo {
	return ClientInfo{
		Addr:       client.conn.RemoteAddr().String(),
		ReadBytes:  client.conn.readTotal,
		WriteBytes: client.conn.writeTotal,
		ConnectAt:  client.conn.Time(),
	}
}

func (client *syncClient) SetId(v interface{}) {
	client.uuid = v
}

func (client *syncClient) GetId() interface{} {
	return client.uuid
}

func (client *syncClient) Handle() error {
	return client.handler(client.conn)
}

func (client *syncClient) Close() error {
	return client.conn.Close()
}

type HandleFunc func(conn net.Conn) error

func NewSyncClientCreator(handleFunc HandleFunc) ClientCreator {
	return func(cc *Conn) Client {
		return &syncClient{
			conn:    cc,
			handler: handleFunc,
		}
	}
}
