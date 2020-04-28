package server

import (
	"errors"
	"github.com/jarod2011/gosocket/client"
	"github.com/jarod2011/gosocket/conn"
	"github.com/jarod2011/gosocket/conn_repo"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	s := New(func(cc conn.Conn, ch <-chan struct{}) error {
		return nil
	})
	if !reflect.DeepEqual(s.Options(), defaultOptions) {
		t.Error("should be default options")
	}
}

func TestServer_Init(t *testing.T) {
	var s Server = new(server)
	err := s.Init()
	if err == nil || !strings.Contains(err.Error(), "address") {
		t.Errorf("should address required err %v", err)
	}
	err = s.Init(WithServerAddr(":9900"))
	if err == nil || !strings.Contains(err.Error(), "repo") {
		t.Errorf("should repo required err %v", err)
	}
	err = s.Init(WithRepo(conn_repo.New()))
	if err == nil || !strings.Contains(err.Error(), "handler") {
		t.Errorf("should handler required err %v", err)
	}
	s = New(func(cc conn.Conn, ch <-chan struct{}) error {
		return nil
	}, WithMaximumOnlineClients(1), WithServerAddr("tttt"))
	err = s.Init()
	if err == nil || !strings.Contains(err.Error(), "maximum") {
		t.Errorf("should client maximum too small %v", err)
	}
	err = s.Init(WithMaximumOnlineClients(100))
	if err == nil {
		t.Error("should listen net address err")
	}
	err = s.Init(WithServerAddr(":29909"))
	if err != nil {
		t.Error(err)
	}
}

func TestServer_Start(t *testing.T) {
	defaultSummaryPrintInterval = time.Millisecond * 100
	defaultInterval = time.Millisecond * 20
	s := New(func(cc conn.Conn, ch <-chan struct{}) error {
		for {
			select {
			case <-ch:
				return errors.New("i am close")
			default:
				byte, err := cc.ReadUntil(time.Now().Add(time.Millisecond * 200))
				if len(byte) > 0 {
					t.Log(byte)
				}
				if err != nil {
					//t.Log(err)
				}
			}
		}
	}, WithServerAddr(":10081"), WithEnableDebug(), WithMaximumOnlineClients(15))
	if err := s.Init(); err != nil {
		t.Error(err)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		<-time.After(time.Second * 3)
		t.Log("start to stop")
		s.Stop()
		wg.Done()
	}()
	for range [20][]int{} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cli := client.New(client.WithServerAddr("127.0.0.1:10081"))
			if err := cli.Init(); err != nil {
				t.Error(err)
			}
			cnt, err := cli.Send([]byte{0x00, 0x01, 0x03})
			if err != nil {
				t.Error(err)
			}
			if cnt != 3 {
				t.Errorf("cnt: %v", cnt)
			}
			time.Sleep(time.Second * 2)
		}()
	}
	defer s.Stop()
	if err := s.Start(); err != nil {
		t.Error(err)
	}
	wg.Wait()
}
