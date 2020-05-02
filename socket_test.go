package gosocket

import (
	"github.com/jarod2011/gosocket/conn"
	"net"
	"testing"
)

func TestNewServer(t *testing.T) {
	NewServer(func(cc conn.Conn, notifyClose <-chan struct{}) error {
		return nil
	})
}

func TestNewClient(t *testing.T) {
	NewClient()
}

func TestNewConn(t *testing.T) {
	c1, _ := net.Pipe()
	NewConn(c1)
}

func TestNewBufferPool(t *testing.T) {
	NewBufferPool()
}
