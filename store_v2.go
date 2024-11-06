package main

import "sync"

type Store[K comparable, V any] struct {
	data map[K]V
	rw   sync.RWMutex
	keys []K
}

func NewStore[K comparable, V any]() *Store[K, V] {
	return &Store[K, V]{
		data: make(map[K]V),
		rw:   sync.RWMutex{},
		keys: make([]K, 0, 42),
	}
}

func (s *Store[K, V]) Set(key K, value V) {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.data[key] = value
	s.keys = append(s.keys, key) //here is a bug
}

func (s *Store[K, V]) Get(key K) (V, error) {
	s.rw.RLock()
	defer s.rw.RUnlock()
	value, ok := s.data[key]
	if !ok {
		var zero V
		return zero, ErrNotFound
	}
	return value, nil
}

func (s *Store[K, V]) Keys() []K {
	s.rw.RLock()
	defer s.rw.RUnlock()
	keys := make([]K, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}

func (s *Store[K, V]) KeysCached() []K {
	s.rw.RLock()
	defer s.rw.RUnlock()
	r := make([]K, len(s.keys))
	copy(r, s.keys)
	return r
}
