package gosocket

import (
	"time"
)

type Client interface {
	SetId(v interface{})
	GetId() interface{}
	Handle() error
	Close() error
	Info() ClientInfo
}

type ClientInfo struct {
	Addr       string
	ReadBytes  uint64
	WriteBytes uint64
	ConnectAt  time.Time
}

type ClientCreator func(cc *Conn) Client
