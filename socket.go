package gosocket

import (
	"net"
	"sync"
	"sync/atomic"
)

type Service interface {
	Init(options ...Option) error
	Run() error
	Close() bool
}

type service struct {
	opt Options
	sync.RWMutex
	listen    net.Listener
	closed    uint32
	clients   sync.Map
	closeChan chan struct{}
}

func (s *service) Init(options ...Option) error {
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

func (s *service) Run() error {
	s.Lock()
	listen, err := net.Listen(s.opt.Addr.Network(), s.opt.Addr.String())
	if err != nil {
		s.Unlock()
		return err
	}
	s.listen = listen
	s.closeChan = make(chan struct{})
	s.Unlock()
	return s.run()
}

func (s *service) Close() bool {
	if atomic.CompareAndSwapUint32(&s.closed, 0, 1) {
		s.Lock()
		close(s.closeChan)
		// close all online clients
		s.clients.Range(func(key, value interface{}) bool {
			value.(Client).Close()
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

func (s *service) run() error {
	var wg sync.WaitGroup
	for {
		select {
		case <-s.closeChan:
			break
		default:
			// handle clients
		}
	}
	wg.Wait()
	return ErrServerClosed
}

func NewSocket(option ...Option) Service {
	options := Options{}
	for _, o := range option {
		o(&options)
	}
	return NewSocketWithOption(options)
}

func NewSocketWithOption(options Options) Service {
	return &service{opt: options}
}
