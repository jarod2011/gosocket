package server

import (
	"context"
	"github.com/jarod2011/gosocket/conn_repo"
	"reflect"
	"testing"
	"time"
)

func TestWithOptions(t *testing.T) {
	opt := Options{}
	opt1 := Options{
		ServerAddr:    "123",
		ServerNetwork: "123",
		Log:           new(defaultLogger),
		Repo:          conn_repo.New(),
		ClientMaximum: 100,
	}
	WithOptions(opt1)(&opt)
	if !reflect.DeepEqual(opt, opt1) {
		t.Errorf("%v %v", opt, opt1)
	}
}

func TestWithLogger(t *testing.T) {
	opt := Options{}
	WithLogger(new(defaultLogger))(&opt)
	if opt.Log == nil {
		t.Error("log should not nil")
	}
}

func TestWithMaximumOnlineClients(t *testing.T) {
	opt := Options{}
	WithMaximumOnlineClients(123)(&opt)
	if opt.ClientMaximum != 123 {
		t.Error("err maximum")
	}
}

func TestWithRepo(t *testing.T) {
	opt := Options{}
	WithRepo(conn_repo.New())(&opt)
	if opt.Repo == nil {
		t.Error("err repo")
	}
}

func TestWithServerNetwork(t *testing.T) {
	opt := Options{}
	WithServerNetwork("network")(&opt)
	if opt.ServerNetwork != "network" {
		t.Error("err network")
	}
}
func TestWithServerAddr(t *testing.T) {
	opt := Options{}
	WithServerAddr("1234")(&opt)
	if opt.ServerAddr != "1234" {
		t.Error("err addr")
	}
}

func TestWithMaxFreeDuration(t *testing.T) {
	opt := Options{}
	WithMaxFreeDuration(time.Minute * 10)(&opt)
	if opt.MaxFreeDuration != time.Minute*10 {
		t.Error("err duration")
	}
}

func TestWithOnlinePrintIntervalDuration(t *testing.T) {
	opt := Options{}
	WithOnlinePrintIntervalDuration(time.Minute * 3)(&opt)
	if opt.OnlinePrintInterval != time.Minute*3 {
		t.Error("err online interval")
	}
}

func TestWithContext(t *testing.T) {
	opt := Options{}
	if opt.Ctx != nil {
		t.Error("should nil context")
	}
	WithContext(context.TODO())(&opt)
	if opt.Ctx == nil {
		t.Error("nil context")
	}
}
