package server

import "time"

type Summary interface {
	ReadBytes() int64
	WriteBytes() int64
	ConnectionAt() time.Time
}
