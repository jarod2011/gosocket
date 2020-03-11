package gosocket

import (
	"sync"
	"sync/atomic"
)

type ClientStore interface {
	Online() int
	Range(func(client Client) bool)
	Save(client Client)
	Delete(client Client)
}

type clientStore struct {
	index  uint64
	online int64
	store  sync.Map
}

func (store *clientStore) Online() int {
	return int(atomic.LoadInt64(&store.online))
}

func (store *clientStore) Range(fn func(client Client) bool) {
	store.store.Range(func(key, value interface{}) bool {
		return fn(value.(Client))
	})
}

func (store *clientStore) Save(client Client) {
	client.SetId(atomic.AddUint64(&store.index, 1))
	store.store.Store(client.GetId(), client)
	atomic.AddInt64(&store.online, 1)
}

func (store *clientStore) Delete(client Client) {
	store.store.Delete(client.GetId())
	atomic.AddInt64(&store.online, -1)
}
