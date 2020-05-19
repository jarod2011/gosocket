package server

import (
	"context"
	"errors"
	"github.com/jarod2011/gosocket/conn"
	"github.com/jarod2011/gosocket/conn_repo"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type Server interface {
	Start() error
	Stop() bool
	Init(opts ...Option) error
	Options() Options
	AfterStart(fn func(ctx context.Context))
	BeforeStop(fn func(ctx context.Context))
}

type ticket struct {
	p chan struct{}
}

func newTicket(cap int) *ticket {
	return &ticket{p: make(chan struct{}, cap)}
}

func (tk *ticket) Take() chan<- struct{} {
	return tk.p
}

func (tk *ticket) Repay() <-chan struct{} {
	return tk.p
}

type server struct {
	closed  int32
	opt     Options
	lst     net.Listener
	hdl     Handler
	tickets *ticket
	ctx     context.Context
	cancel  context.CancelFunc

	afterHdl  func(ctx context.Context)
	beforeHdl func(ctx context.Context)
}

func (s *server) Start() error {
	if s.lst == nil {
		return errors.New("uninitialized")
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		defer s.opt.Log.Log(debugPrefix, "monitor routine exit")
		defer func() {
			// clean online connections
			s.opt.Log.Logf(infoPrefix+"now %d online clients, clean...", s.opt.Repo.Online())
			s.opt.Repo.Iterate(func(cc conn.Conn) bool {
				s.opt.Log.Logf(infoPrefix+"close conn %v, summary:\nwrite: %d bytes\nread: %d bytes\nconnected: %v\nactiveAt: %v", cc.RemoteAddr(), cc.WriteBytes(), cc.ReadBytes(), time.Since(cc.Created()), cc.LastActive())
				return true
			})
			time.Sleep(time.Second)
			s.lst.Close()
		}()
		printOnline := time.NewTicker(s.opt.OnlinePrintInterval)
		defer printOnline.Stop()
		checkConn := time.NewTicker(time.Second * 20)
		defer checkConn.Stop()
		for {
			select {
			case <-s.ctx.Done():
				return
			case <-printOnline.C:
				s.opt.Log.Logf(infoPrefix+"now online %d connections", s.opt.Repo.Online())
			case <-checkConn.C:
				if s.opt.Repo.Online() > 0 {
					s.opt.Repo.Iterate(func(cc conn.Conn) bool {
						if time.Since(cc.LastActive()) > s.opt.MaxFreeDuration {
							s.opt.Log.Logf(infoPrefix+"close free conn %v", cc.RemoteAddr())
							cc.Close()
						}
						return true
					})
				}
			}
		}
	}()
	go func() {
		defer wg.Done()
		defer s.opt.Log.Log(debugPrefix, "listen routine exit")
		var wgg sync.WaitGroup
		defer wgg.Wait()
		go func() {
			<-time.After(time.Second)
			if s.afterHdl != nil {
				s.afterHdl(s.ctx)
			}
		}()
		for {
			select {
			case <-s.ctx.Done():
				s.opt.Log.Logf(infoPrefix + "stop listen routine")
				time.Sleep(time.Second)
				return
			case <-time.After(defaultInterval):
				s.opt.Log.Logf(warnPrefix+"online %d is maximum(%d)", s.opt.Repo.Online(), s.opt.ClientMaximum)
			case s.tickets.Take() <- struct{}{}:
				cc, err := s.lst.Accept()
				if err != nil {
					s.opt.Log.Logf(errPrefix+"accept connect failure: %v", err)
					<-s.tickets.Repay()
					continue
				}
				con := conn.New(cc)
				if _, err := s.opt.Repo.AddConn(con); err != nil {
					s.opt.Log.Logf(errPrefix+"save conn %v to repo failure: %v", con.RemoteAddr(), err)
					<-s.tickets.Repay()
					con.Close()
					continue
				}
				wgg.Add(1)
				go func(c conn.Conn) {
					defer wgg.Done()
					defer func() {
						s.opt.Log.Logf(debugPrefix+"remove conn %v", c.RemoteAddr())
						if err := s.opt.Repo.RemoveConn(c); err != nil {
							s.opt.Log.Logf(errPrefix+"remove conn %v from repo failure: %v", c.RemoteAddr(), err)
						}
						s.opt.Log.Logf(infoPrefix+"close conn %v, summary:\nwrite: %d bytes\nread: %d bytes\nconnected: %v\nactiveAt: %v", cc.RemoteAddr(), c.WriteBytes(), c.ReadBytes(), time.Since(c.Created()), c.LastActive())
						c.Close()
						<-s.tickets.Repay()
					}()
					if err := s.hdl(s.ctx, c); err != nil {
						s.opt.Log.Logf(errPrefix+"handle conn %v failure: %v", c.RemoteAddr(), err)
					}
				}(con)
			}
		}
	}()
	wg.Wait()
	return nil
}

func (s *server) Stop() bool {
	if atomic.CompareAndSwapInt32(&s.closed, 0, 1) {
		if s.beforeHdl != nil {
			s.beforeHdl(s.ctx)
		}
		s.cancel()
		return true
	}
	return false
}

func (s *server) AfterStart(fn func(ctx context.Context)) {
	s.afterHdl = fn
}

func (s *server) BeforeStop(fn func(ctx context.Context)) {
	s.beforeHdl = fn
}

func (s *server) Init(opts ...Option) (err error) {
	for _, o := range opts {
		o(&s.opt)
	}
	if len(s.opt.ServerAddr) == 0 {
		return errors.New("address required")
	}
	if s.opt.Repo == nil {
		return errors.New("repo required")
	}
	if s.hdl == nil {
		return errors.New("handler required")
	}
	if s.opt.ClientMaximum < 10 {
		return errors.New("client maximum should more than 10")
	}
	s.tickets = newTicket(s.opt.ClientMaximum)
	s.lst, err = net.Listen(s.opt.ServerNetwork, s.opt.ServerAddr)
	s.ctx, s.cancel = context.WithCancel(s.opt.Ctx)
	return
}

func (s server) Options() Options {
	return s.opt
}

var defaultOptions = Options{
	ServerNetwork:       "tcp",
	Repo:                conn_repo.New(),
	Log:                 new(defaultLogger),
	ClientMaximum:       defaultMaximumOnlineClients,
	PrintDebug:          false,
	MaxFreeDuration:     time.Minute * 10,
	OnlinePrintInterval: defaultSummaryPrintInterval,
	Ctx:                 context.Background(),
}

var defaultInterval = time.Second * 10
var defaultMaximumOnlineClients = 10000
var defaultSummaryPrintInterval = time.Minute

type Handler func(ctx context.Context, conn conn.Conn) error // handle conn

func New(handler Handler, opts ...Option) Server {
	s := new(server)
	WithOptions(defaultOptions)(&s.opt)
	s.hdl = handler
	for _, o := range opts {
		o(&s.opt)
	}
	return s
}
