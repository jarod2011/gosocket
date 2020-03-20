package gosocket

import (
	"context"
	"time"
)

type Context interface {
	context.Context
	Cancel()
	RemoteAddr() string
	ReadBytes() int64
	WriteBytes() int64
	ConnectedTime() time.Time
	WriteToClient(b []byte) error
}
