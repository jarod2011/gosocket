package client

import (
	"testing"
	"time"
)

func TestOption(t *testing.T) {
	opt := Options{}
	WithReadTimeout(time.Second * 33)(&opt)
	if opt.ReadTimeout != time.Second*33 {
		t.Error("read time != 33 second")
	}
	WithServerAddr("11111")(&opt)
	if opt.ServerAddr != "11111" {
		t.Error("server addr error")
	}
	WithServerNetwork("udp")(&opt)
	if opt.ServerNetwork != "udp" {
		t.Error("server network error")
	}
	WithOptions(Options{
		ServerAddr:    "123",
		ServerNetwork: "tcp",
	})(&opt)
	if opt.ServerAddr != "123" || opt.ServerNetwork != "tcp" {
		t.Error("addr or network failure")
	}
	WithWriteTimeout(time.Second * 13)(&opt)
	if opt.WriteTimeout != time.Second*13 {
		t.Error("write timeout")
	}
}
