package server

import (
	"context"
	log "github.com/jarod2011/go-log"
	"github.com/jarod2011/gosocket/util"
	"net"
	"sync"
	"sync/atomic"
)

type Socket interface {
	Start() error
	Stop() bool
	Closed() bool
	Options() Options
}

type socket struct {
	option Options
	closed uint32
}

func (s socket) Options() Options {
	return s.option
}

func (s *socket) clean() {
	s.option.Store.Range(func(c Conn) {
		c.Close()
	})
}

func (s *socket) Start() error {
	if s.Closed() {
		return ErrServerClosed
	}
	if err := s.option.Validate(); err != nil {
		return err
	}
	listen, err := net.Listen(s.option.Address.Network(), s.option.Address.String())
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
loop:
	for {
		select {
		case <-s.option.ctx.Done():
			break loop
		default:
			cc, err := listen.Accept()
			if err != nil {
				if util.IsClosedConnection(err) {
					break loop
				}
				log.D().Logf("accept failure: %v", err)
				continue
			}
			log.D().Logf("accept connection: %v", cc.RemoteAddr())
			c := newConn(cc)
			s.option.Store.Store(c)
			wg.Add(1)
			go func(conn Conn) {
				handler := s.option.HandlerCreator()
				if err := handler(conn); err != nil {
					log.E().Logf("handling connection %v failure: %v", conn.RemoteAddr(), err)
				}
				log.D().Logf("connection %v read %d bytes write %d bytes", conn.RemoteAddr(), conn.ReadBytes(), conn.WriteBytes())
				s.option.Store.Remove(conn)
				wg.Done()
			}(c)
		}
	}
	s.clean()
	defer wg.Wait()
	return ErrServerClosed
}

func (s *socket) Stop() bool {
	if atomic.CompareAndSwapUint32(&s.closed, StateOpened, StateClosed) {
		s.option.cancel()
		return true
	}
	return false
}

func (s *socket) Closed() bool {
	return StateClosed == atomic.LoadUint32(&s.closed)
}

func NewSocket(opts ...Option) Socket {
	opt := Options{
		Address:       &net.TCPAddr{Port: 8080},
		MaxConnection: 10000,
		Store:         new(mapStore),
	}
	opt.ctx, opt.cancel = context.WithCancel(context.Background())
	for _, o := range opts {
		o(&opt)
	}
	return &socket{closed: StateOpened, option: opt}
}
