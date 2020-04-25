package buffers

import (
	"bytes"
	"sync"
)

type Buffers interface {
	Get() *bytes.Buffer
	Put(buf *bytes.Buffer)
}

type buffers struct {
	pool sync.Pool
}

func (b buffers) Get() *bytes.Buffer {
	return b.pool.Get().(*bytes.Buffer)
}

func (b buffers) Put(buf *bytes.Buffer) {
	b.pool.Put(buf)
}

func New() Buffers {
	pool := new(buffers)
	pool.pool = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
	return pool
}
