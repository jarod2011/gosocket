package gosocket

import (
	"context"
	"net"
	"sync"
	"sync/atomic"
)

type Socket interface {
	Init(options ...Option) error
	Run() error
	Close() bool
	Store() ClientStore
}

type server struct {
	opt Options
	sync.RWMutex
	listen  net.Listener
	closed  uint32
	clients *clientStore
	ctx     context.Context
	cancel  context.CancelFunc
}

func (s *server) Store() ClientStore {
	return s.clients
}

func (s *server) Init(options ...Option) error {
	s.Lock()
	defer s.Unlock()
	for _, o := range options {
		o(&s.opt)
	}
	if s.opt.Log == nil {
		s.opt.Log = newLogger()
	}
	if s.opt.Addr == nil {
		s.opt.Addr = &net.TCPAddr{Port: DefaultPort}
	}
	return nil
}

func (s *server) Run() error {
	s.Lock()
	listen, err := net.Listen(s.opt.Addr.Network(), s.opt.Addr.String())
	if err != nil {
		s.Unlock()
		return err
	}
	s.listen = listen
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.clients = &clientStore{}
	s.Unlock()
	return s.run()
}

func (s *server) Close() bool {
	if atomic.CompareAndSwapUint32(&s.closed, 0, 1) {
		s.Lock()
		s.cancel()
		// close all online clients
		s.clients.Range(func(client Client) bool {
			client.Close()
			return true
		})
		if s.listen != nil {
			s.listen.Close()
			s.listen = nil
		}
		s.Unlock()
		return true
	}
	return false
}

func (s *server) run() error {
	var wg sync.WaitGroup
loop:
	for {
		select {
		case <-s.ctx.Done():
			break loop
		default:
			// handle conn
			conn, err := s.listen.Accept()
			if err != nil {
				if IsClosedConnection(err) {
					break loop
				}
			}
			cli := s.opt.Creator(newConn(conn))
			// save to clients
			s.clients.Save(cli)
			wg.Add(1)
			go func(c Client) {
				defer wg.Done()
				defer s.clients.Delete(c)
				c.Handle()
			}(cli)
		}
	}
	wg.Wait()
	return ErrServerClosed
}

func NewSocket(option ...Option) Socket {
	options := Options{}
	for _, o := range option {
		o(&options)
	}
	return &server{opt: options}
}
