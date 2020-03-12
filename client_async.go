package gosocket

import (
	"context"
	"sync"
	"time"
)

type asyncClient struct {
	conn      *Conn
	scheduler *schedule
	opt       AsyncOptions
	handle    Handler
	uuid      interface{}
	errChan   chan error
}

func (client *asyncClient) Info() ClientInfo {
	return ClientInfo{
		Addr:       client.conn.RemoteAddr().String(),
		ReadBytes:  client.conn.readTotal,
		WriteBytes: client.conn.writeTotal,
		ConnectAt:  client.conn.Time(),
	}
}

func (client *asyncClient) SetId(v interface{}) {
	client.uuid = v
}

func (client *asyncClient) GetId() interface{} {
	return client.uuid
}

func (client *asyncClient) Handle() error {
	var wg sync.WaitGroup
	wg.Add(2)
	go client.readProcess(&wg)
	go client.writeProcess(&wg)
	for i := 0; i < int(client.opt.WorkerNum); i++ {
		wg.Add(1)
		go client.workProcess(&wg)
	}
	wg.Wait()
	var lastErr error
	select {
	case lastErr = <-client.errChan:
	default:
	}
	client.handle.OnClose(client.conn, lastErr)
	client.conn.Close()
	return nil
}

func (client *asyncClient) readProcess(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-client.scheduler.ctx.Done():
			return
		default:
			// TODO read bytes from conn
		}
	}
}

func (client *asyncClient) writeProcess(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-client.scheduler.ctx.Done():
			return
		case buf := <-client.scheduler.writerChan:
			rsp, err := client.handle.OnWrite(buf)
			if err != nil {
				// TODO log err
				continue
			}
			if len(rsp) > 0 {
				// TODO write to conn
			}
		}
	}
}

func (client *asyncClient) workProcess(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-client.scheduler.ctx.Done():
			return
		case buf := <-client.scheduler.workerChan:
			client.handle.OnWork(buf, client.scheduler)
		}
	}
}

func (client *asyncClient) Close() error {
	// TODO close client
	panic("implement me")
}

type schedule struct {
	ctx        context.Context
	cancel     context.CancelFunc
	workerChan chan []byte
	writerChan chan []byte
}

func (sch *schedule) SendToWorker(buf []byte) error {
	sch.workerChan <- buf
	return nil
}

func (sch *schedule) SendToWriter(buf []byte) error {
	sch.writerChan <- buf
	return nil
}

func (sch *schedule) StopClient() error {
	sch.cancel()
	return nil
}

func (sch *schedule) Context() context.Context {
	return sch.ctx
}

type Handler interface {
	OnRead(buf []byte, scheduler Scheduler) bool
	OnWork(buf []byte, scheduler Scheduler)
	OnWrite(buf []byte) ([]byte, error)
	OnClose(conn *Conn, lastErr error)
}

type Scheduler interface {
	SendToWorker(buf []byte) error
	SendToWriter(buf []byte) error
	StopClient() error
	Context() context.Context
}

var (
	DefaultReadTimeout        = time.Second * 10
	DefaultWriteTimeout       = time.Second * 10
	DefaultWorkerNum    uint8 = 1
)

type AsyncOptions struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	WorkerNum    uint8
}

func NewAsyncClient(handler Handler, options AsyncOptions) ClientCreator {
	return func(cc *Conn) Client {
		s := schedule{}
		if options.ReadTimeout <= 0 {
			options.ReadTimeout = DefaultReadTimeout
		}
		if options.WriteTimeout <= 0 {
			options.WriteTimeout = DefaultWriteTimeout
		}
		if options.WorkerNum == 0 {
			options.WorkerNum = DefaultWorkerNum
		}
		return &asyncClient{
			conn:      cc,
			scheduler: &s,
			opt:       options,
			handle:    handler,
		}
	}
}
