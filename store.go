package concurrent_store

import (
	"errors"
	"sync"
)

var ErrStoreEmpty = errors.New("store is empty")

// Store is a concurrency-safe deduplicated store
type Store struct {
	lock  *sync.RWMutex
	items map[interface{}]struct{}
}

func NewStore() *Store {
	return &Store{
		lock:  &sync.RWMutex{},
		items: map[interface{}]struct{}{},
	}
}

func (l *Store) Add(v interface{}) {
	l.lock.Lock()
	l.items[v] = struct{}{}
	l.lock.Unlock()
}

func (l *Store) Contains(v interface{}) bool {
	l.lock.RLock()
	_, ok := l.items[v]
	l.lock.RUnlock()

	return ok
}

func (l *Store) Pop() (interface{}, error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if len(l.items) == 0 {
		return nil, ErrStoreEmpty
	}

	for k := range l.items {
		delete(l.items, k)
		return k, nil
	}

	panic("reached unreachable state in concurrent store")
	return nil, nil
}

func (l *Store) All() map[interface{}]struct{} {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.items
}
