package async

import (
	"context"
	"errors"
	"github.com/jarod2011/gosocket/conn"
	"github.com/jarod2011/gosocket/server"
	"io"
	"sync"
	"time"
)

type handle struct {
	opt       Option
	cc        conn.Conn
	wg        sync.WaitGroup
	writeChan chan []byte
	workChan  chan []byte
	errChan   chan error
	ctx       context.Context
	cancel    context.CancelFunc
	hdl       Handler
}

func (h *handle) readProcess() {
	defer h.wg.Done()
	defer close(h.workChan)
	defer h.hdl.OnReadProcessStop()
	var last []byte
	for {
		select {
		case <-h.ctx.Done():
			return
		default:
			by, err := h.cc.ReadUntil(time.Now().Add(h.opt.ReadTimeout), false)
			if len(by) > 0 {
				last = append(last, by...)
			}
			if len(last) > 0 {
				cnt := h.hdl.SliceIndex(last)
				if cnt > 0 {
					h.workChan <- last[0:cnt]
					last = last[cnt:]
				}
			}
			if err != nil {
				if err == io.EOF {
					return
				}
				if err == conn.ErrContextDeadline {
					continue
				}
				h.errChan <- err
				return
			}
		}
	}
}

func (h *handle) writeProcess() {
	defer h.wg.Done()
	defer h.Stop()
	defer h.hdl.OnWriteProcessStop()
	for b := range h.writeChan {
		_, err := h.cc.WriteUntil(b, time.Now().Add(h.opt.WriteTimeout))
		h.hdl.OnWriteFinish(b)
		if err != nil {
			if err == io.EOF {
				return
			}
			if err == conn.ErrContextDeadline {
				continue
			}
			if h.hdl.OnWriteError(b, err) {
				h.errChan <- err
				return
			}
		}
	}
}

func (h *handle) workProcess() {
	defer h.wg.Done()
	defer close(h.writeChan)
	defer h.hdl.OnWorkProcessStop()
	for b := range h.workChan {
		if err := h.hdl.OnWork(b, h.writeChan); err != nil {
			h.errChan <- err
			return
		}
	}
}

func (h *handle) Handle(cc conn.Conn) error {
	h.cc = cc
	h.hdl.OnConnect(cc)
	h.wg.Add(3)
	go h.readProcess()
	go h.workProcess()
	go h.writeProcess()
	h.wg.Wait()
	h.hdl.OnClose()
	h.cc.Close()
	h.cc = nil
	select {
	case err := <-h.errChan:
		return err
	default:
		return nil
	}
}

func (h *handle) Stop() {
	h.cancel()
}

func New(opt Option) server.Handler {
	return func(ctx context.Context, cc conn.Conn) error {
		if opt.Context == nil {
			opt.Context = context.TODO()
		}
		h := handle{
			opt:       opt,
			errChan:   make(chan error, 3),
			writeChan: make(chan []byte, 100),
			workChan:  make(chan []byte, 100),
		}
		h.ctx, h.cancel = context.WithCancel(opt.Context)
		ch := make(chan error, 1)
		if h.opt.Creator == nil {
			return errors.New("nil handler creator")
		}
		h.hdl = h.opt.Creator()
		if h.hdl == nil {
			return errors.New("nil handler")
		}
		if h.opt.ReadTimeout <= 0 {
			h.opt.ReadTimeout = time.Second
		}
		if h.opt.WriteTimeout <= 0 {
			h.opt.WriteTimeout = time.Second
		}
		go func() {
			ch <- h.Handle(cc)
		}()
		select {
		case <-ctx.Done():
			h.Stop()
			return nil
		case err := <-ch:
			return err
		}
	}
}
