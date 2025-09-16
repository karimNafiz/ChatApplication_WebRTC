package safemap

import (
	"sync"
)

type SafeMap[K comparable, V any] struct {
	/*
		multiple reader go-routines can access the map at the same time
		but only one writer go-routine can access the map at a time

		should be private
	*/
	readWriteMutex sync.RWMutex

	/*
		the internal map to hold the keys and values
	*/
	hashMap map[K]V
}

func New[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		hashMap: make(map[K]V),
	}
}

func (m *SafeMap[K, V]) Get(key K) (V, bool) {
	m.readWriteMutex.RLock()
	defer m.readWriteMutex.RUnlock()
	value, ok := m.hashMap[key]
	return value, ok
}

func (m *SafeMap[K, V]) Set(key K, value V) {
	m.readWriteMutex.Lock()
	defer m.readWriteMutex.Unlock()
	m.hashMap[key] = value
}
