package client

import (
	"context"
	"github.com/jarod2011/gosocket/conn"
	"net"
	"reflect"
	"testing"
	"time"
)

func genServer(ctx context.Context, t *testing.T) error {
	lst, err := net.Listen("tcp", ":8899")
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				if err := lst.Close(); err != nil {
					t.Error(err)
				}
				return
			default:
				cc, err := lst.Accept()
				if err != nil {
					t.Errorf("accept failure: %v", err)
					continue
				}
				con := conn.New(cc)
				rr, err := con.ReadUntil(time.Now().Add(time.Second))
				if reflect.DeepEqual(rr, []byte{0x01, 0x02}) {
					t.Logf("write..")
					cc.Write(make([]byte, 100))
					cc.Write(make([]byte, 100))
				}
				t.Logf("%v", rr)
			}
		}
	}()
	return nil
}

func TestClient_Init(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := genServer(ctx, t); err != nil {
		t.Fatal(err)
	}
	cli := New(WithReadTimeout(time.Second * 5))
	if err := cli.Init(); err == nil {
		t.Fatal("should missing address")
	}
	if err := cli.Init(WithServerAddr("localhost:8899")); err != nil {
		t.Fatal(err)
	}
	if cli.Options().ServerAddr != "localhost:8899" {
		t.Errorf("address %v", cli.Options().ServerAddr)
	}
	n, err := cli.Send([]byte{0x01, 0x02})
	if err != nil {
		t.Errorf("write %v", err)
	}
	t.Logf("client send %d", n)
	if n != 2 {
		t.Errorf("write %d", n)
	}
	by, err := cli.Recv()
	if err != nil {
		t.Errorf("read err %v", err)
	}
	t.Log(by)
	if err := cli.Close(); err != nil {
		t.Fatal(err)
	}
	cancel()
}
