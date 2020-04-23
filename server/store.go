package server

import "sync"

type ConnectionStore interface {
	Store(cc Conn)
	Remove(cc Conn)
	Summary() StoreSummary
	Range(callback func(Conn))
}

type StoreSummary struct {
	ReadBytes  int64
	WriteBytes int64
	Count      int
}

type mapStore struct {
	store sync.Map
}

func (m *mapStore) Store(cc Conn) {
	m.store.Store(cc.UUID(), cc)
}

func (m *mapStore) Remove(cc Conn) {
	m.store.Delete(cc.UUID())
}

func (m *mapStore) Range(callback func(Conn)) {
	m.store.Range(func(key, value interface{}) bool {
		callback(value.(Conn))
		return true
	})
}

func (m *mapStore) Summary() StoreSummary {
	sm := StoreSummary{}
	m.store.Range(func(key, value interface{}) bool {
		sm.ReadBytes += value.(Conn).ReadBytes()
		sm.WriteBytes += value.(Conn).WriteBytes()
		sm.Count++
		return true
	})
	return sm
}
