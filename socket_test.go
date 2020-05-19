package gosocket

import (
	"context"
	"github.com/jarod2011/gosocket/conn"
	"net"
	"testing"
)

func TestNewServer(t *testing.T) {
	NewServer(func(ctx context.Context, cc conn.Conn) error {
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
