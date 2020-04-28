package conn_repo

import (
	"github.com/jarod2011/gosocket/conn"
	"sync"
	"sync/atomic"
)

type Repo interface {
	AddConn(cc conn.Conn) (uint64, error)
	RemoveConn(cc conn.Conn) error
	Online() int
	GetConn(id uint64) (conn.Conn, bool)
	Iterate(func(cc conn.Conn) bool)
}

type memoryRepo struct {
	store  sync.Map
	online int32
}

func (m *memoryRepo) Iterate(fn func(cc conn.Conn) bool) {
	m.store.Range(func(key, value interface{}) bool {
		return fn(value.(conn.Conn))
	})
}

func (m *memoryRepo) AddConn(cc conn.Conn) (uint64, error) {
	m.store.Store(cc.ID(), cc)
	atomic.AddInt32(&m.online, 1)
	return cc.ID(), nil
}

func (m *memoryRepo) RemoveConn(cc conn.Conn) error {
	m.store.Delete(cc.ID())
	atomic.AddInt32(&m.online, -1)
	return nil
}

func (m *memoryRepo) Online() int {
	return int(atomic.LoadInt32(&m.online))
}

func (m *memoryRepo) GetConn(id uint64) (conn.Conn, bool) {
	if res, ok := m.store.Load(id); ok {
		return res.(conn.Conn), true
	}
	return nil, false
}

func New() Repo {
	return new(memoryRepo)
}
