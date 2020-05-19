package server

import (
	"context"
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
	s := New(func(ctx context.Context, cc conn.Conn) error {
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
	s = New(func(ctx context.Context, cc conn.Conn) error {
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
	s := New(func(ctx context.Context, cc conn.Conn) error {
		for {
			select {
			case <-ctx.Done():
				return errors.New("i am close")
			default:
				by, err := cc.ReadUntil(time.Now().Add(time.Millisecond*200), true)
				if len(by) > 0 {
					t.Log(by)
				}
				if err != nil {
					//t.Log(err)
				}
			}
		}
	}, WithServerAddr(":10081"), WithEnableDebug(), WithMaximumOnlineClients(15), WithOnPrint(func(ctx context.Context, repo conn_repo.Repo) {
		t.Log("online %d", repo.Online())
	}), WithOnlinePrintIntervalDuration(time.Millisecond*10))
	if err := s.Start(); err == nil {
		t.Error("should uninitialized error")
	}
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
			defer cli.Close()
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
	s.AfterStart(func(ctx context.Context) {
		t.Logf("after start call")
	})
	s.BeforeStop(func(ctx context.Context) {
		t.Logf("before stop call")
	})
	if err := s.Start(); err != nil {
		t.Error(err)
	}
	wg.Wait()
}
