package conn_repo

import (
	"github.com/jarod2011/gosocket/conn"
	"net"
	"reflect"
	"testing"
)

func TestMemoryRepo(t *testing.T) {
	mem := New()
	c1, c2 := net.Pipe()
	con1 := conn.New(c1)
	con2 := conn.New(c2)
	t.Log(con1.RemoteAddr(), con2.RemoteAddr())
	if mem.Online() > 0 {
		t.Errorf("online %d", mem.Online())
	}
	if _, err := mem.AddConn(con1); err != nil {
		t.Error(err)
	}
	if mem.Online() != 1 {
		t.Errorf("online %d", mem.Online())
	}
	if _, err := mem.AddConn(con2); err != nil {
		t.Error(err)
	}
	if mem.Online() != 2 {
		t.Errorf("online %d", mem.Online())
	}
	if err := mem.RemoveConn(con1); err != nil {
		t.Error(err)
	}
	if mem.Online() != 1 {
		t.Errorf("online %d", mem.Online())
	}
	re1, ok1 := mem.GetConn(con2.ID())
	if !ok1 {
		t.Error("get con2 failure")
	} else {
		if !reflect.DeepEqual(re1.RemoteAddr(), con2.RemoteAddr()) {
			t.Errorf("%v != %v", re1.RemoteAddr(), con2.RemoteAddr())
		}
	}
	_, ok2 := mem.GetConn(con1.ID())
	if ok2 {
		t.Errorf("get con1 should failure")
	}
}
