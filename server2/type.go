package server2

import "time"

const (
	StateOpened uint32 = iota
	StateClosed
)

var (
	uniqueId int64
)

func init() {
	uniqueId = time.Now().Unix()
}
