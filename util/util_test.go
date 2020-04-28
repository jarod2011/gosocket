package util

import (
	"errors"
	"testing"
)

func TestIsClosedConnection(t *testing.T) {
	err := errors.New("use of closed network connection")
	if !IsClosedConnection(err) {
		t.Error("error")
	}
}

func TestIsRemoteClosedError(t *testing.T) {
	err := errors.New("closed by the remote host")
	if !IsRemoteClosedError(err) {
		t.Error("error")
	}
}

func TestIsTimeout(t *testing.T) {
	err := errors.New("is timeout")
	if IsTimeout(err) {
		t.Error("error")
	}
}
