package async

import (
	"context"
	"github.com/jarod2011/gosocket/conn"
	"time"
)

type Option struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Handle       Handler
	Context      context.Context
}

type Handler interface {
	OnConnect(cc conn.Conn)
	SliceIndex(b []byte) int
	OnWork(b []byte, writeChan chan<- []byte) error
	OnWorkProcessStop()
	OnReadProcessStop()
	OnWriteProcessStop()
	OnWriteFinish(b []byte)
	OnWriteError(b []byte, err error) (ifExit bool)
	OnClose()
}
