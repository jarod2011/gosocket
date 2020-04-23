package server2

import "errors"

var (
	ErrConnectionClosed = errors.New("connection is closed")
	ErrServerClosed     = errors.New("server is closed")
)
