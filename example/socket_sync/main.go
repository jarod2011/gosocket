package main

import (
	"bytes"
	log "github.com/jarod2011/go-log"
	"github.com/jarod2011/gosocket"
	"net"
	"sync"
	"time"
)

var addr = &net.TCPAddr{
	Port: 8000,
}
var serve gosocket.Socket

var timeout = time.Second * 3

func init() {
	log.SetLevel(log.Debug)
	log.SetPrefix("[example.socket.sync]")
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go server(&wg)
	time.Sleep(time.Second * 2)
	client()
	log.D().Log("start to close serve")
	serve.Close()
	wg.Wait()
}

func server(wg *sync.WaitGroup) {
	defer wg.Done()
	serve = gosocket.NewSocket(
		gosocket.WithAddr(addr),
		gosocket.WithClientCreator(gosocket.NewSyncClientCreator(func(conn *gosocket.Conn) error {
			log.D().Logf("handle the connection %v", conn.RemoteAddr())
			var wwg sync.WaitGroup
			wwg.Add(2)
			ch := make(chan []byte, 10)
			go func() {
				defer wwg.Done()
				defer close(ch)
				for {
					buf := make([]byte, 1024)
					conn.SetReadDeadline(time.Now().Add(time.Second))
					n, err := conn.Read(buf)
					if n > 0 {
						log.I().Logf("[server]%v <== %d bytes(%x) === %v[client]", conn.LocalAddr(), n, buf[0:n], conn.RemoteAddr())
						ch <- buf[0:n]
					}
					if err != nil && !gosocket.IsTimeout(err) {
						if gosocket.IsClosedConnection(err) || gosocket.IsRemoteClosedError(err) {
							break
						}
						log.E().Logf("read from client failure: %v", err)
					}
				}
			}()
			go func() {
				defer wwg.Done()
				for {
					select {
					case b, ok := <-ch:
						if !ok {
							return
						}
						conn.SetWriteDeadline(time.Now().Add(time.Second))
						by := bytes.Repeat(b, 2)
						m, err := conn.Write(by)
						if m > 0 {
							log.I().Logf("[server]%v === %d bytes(%x) ==> %v[client]", conn.LocalAddr(), m, by, conn.RemoteAddr())
						}
						if err != nil && !gosocket.IsTimeout(err) {
							log.E().Logf("write to client failure: %v", err)
							if gosocket.IsClosedConnection(err) {
								break
							}
						}
					}
				}
			}()
			wwg.Wait()
			log.D().Log("close connection")
			return nil
		})),
	)
	if err := serve.Init(); err != nil {
		log.F().Logf("init server failure: %v", err)
	}
	if err := serve.Run(); err != nil && err != gosocket.ErrServerClosed {
		log.F().Logf("run server failure: %v", err)
	}
	log.I().Log("server closed")
}

func client() {
	conn, err := net.Dial(addr.Network(), addr.String())
	if err != nil {
		log.F().Logf("dial server failure: %v", err)
	}
	msg := [...][]byte{
		[]byte{0xff, 0x11, 0x22, 0x33},
		[]byte{0xee, 0x22, 0x33, 0x44},
		[]byte{0xdd, 0x33, 0x44, 0x55},
		[]byte{0xbb, 0x44, 0x55, 0x66},
	}
	var wg sync.WaitGroup
	wg.Add(2)
	ch := make(chan int)
	go func() {
		defer wg.Done()
		var index int
		for {
			buf := make([]byte, 1024)
			conn.SetReadDeadline(time.Now().Add(timeout))
			m, err := conn.Read(buf)
			if m > 0 {
				index++
				ch <- index
				log.I().Logf("[client]%v <== %d bytes(%x) === %v[server]", conn.LocalAddr(), m, buf[0:m], conn.RemoteAddr())
			}
			if index >= len(msg) {
				return
			}
			if err != nil && gosocket.IsTimeout(err) {
				log.E().Logf("read from server failure: %v", err)
				if gosocket.IsClosedConnection(err) || gosocket.IsRemoteClosedError(err) {
					return
				}
			}
		}
	}()
	go func() {
		defer wg.Done()
		for _, b := range msg {
			conn.SetWriteDeadline(time.Now().Add(timeout))
			n, err := conn.Write(b)
			if n > 0 {
				log.I().Logf("[client]%v === %d bytes(%x) ==> %v[server]", conn.LocalAddr(), n, b, conn.RemoteAddr())
			}
			if err != nil && gosocket.IsTimeout(err) {
				log.E().Logf("write to server failure: %v", err)
				if gosocket.IsClosedConnection(err) || gosocket.IsRemoteClosedError(err) {
					return
				}
			}
			<-ch
		}
	}()
	wg.Wait()
	log.D().Log("close client")
	conn.Close()
}
