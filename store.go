package main

import (
	"errors"
	"sync"
)

var ErrNotFound = errors.New("not found")

type SyncStore[K comparable, V any] struct {
	data sync.Map
}

func NewSyncStore[K comparable, V any]() *SyncStore[K, V] {
	return &SyncStore[K, V]{
		data: sync.Map{},
	}
}

func (s *SyncStore[K, V]) Set(key K, value V) {
	s.data.Store(key, value)
}

func (s *SyncStore[K, V]) Get(key K) (V, error) {
	v, ok := s.data.Load(key)
	if !ok {
		var zero V
		return zero, ErrNotFound
	}
	vr, ok := v.(V)
	if !ok {
		var zero V
		return zero, ErrNotFound
	}
	return vr, nil
}

func (s *SyncStore[K, V]) Keys() []K {
	keys := make([]K, 0, 42)
	s.data.Range(func(key, value any) bool {
		keys = append(keys, key.(K))
		return true
	})
	return keys
}

func (s *SyncStore[K, V]) GetAny(key K) (any, error) {
	v, ok := s.data.Load(key)
	if !ok {
		var zero V
		return zero, ErrNotFound
	}
	return v, nil
}
