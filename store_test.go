package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStore(t *testing.T) {
	t.Run("success on create", func(t *testing.T) {
		got := NewSyncStore[int, int]()
		assert.Equal(t, &SyncStore[int, int]{}, got)
	})

	t.Run("success on set", func(t *testing.T) {
		got := NewSyncStore[int, int]()
		got.Set(1, 2)
		v, ok := got.data.Load(1)
		assert.True(t, ok)
		assert.Equal(t, 2, v)
	})

	t.Run("success on get", func(t *testing.T) {
		got := NewSyncStore[int, int]()
		got.Set(1, 2)
		v, err := got.Get(1)
		assert.NoError(t, err)
		assert.Equal(t, 2, v)
	})
}

func TestKeys(t *testing.T) {
	tests := []struct {
		name string
		keys []int
		want []int
	}{
		{
			name: "empty",
			keys: []int{},
			want: []int{},
		},
		{
			name: "one",
			keys: []int{1},
			want: []int{1},
		},
		{
			name: "two",
			keys: []int{1, 2},
			want: []int{1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSyncStore[int, int]()
			for _, k := range tt.keys {
				got.Set(k, k)
			}
			gotKeys := got.Keys()
			assert.ElementsMatch(t, tt.want, gotKeys)
		})
	}
}

func BenchmarkStoreAsyncSet(b *testing.B) {
	b.Run("set", asyncSet)
	b.Run("get", asyncGet)
	b.Run("get any", asyncGetAny)
	b.Run("set get", asyncSetGet)
	b.Run("keys", asyncKeys)
}

func asyncKeys(b *testing.B) {
	store := NewSyncStore[int, int]()
	for i := 0; i < b.N; i++ {
		store.Set(i, i)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			store.Keys()
		}
	})
}

func asyncSet(b *testing.B) {
	store := NewSyncStore[int, int]()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			store.Set(1, 1)
		}
	})
}

func asyncGet(b *testing.B) {
	store := NewSyncStore[int, int]()
	for i := 0; i < b.N; i++ {
		store.Set(i, i)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = store.Get(1)
		}
	})
}

func asyncGetAny(b *testing.B) {
	store := NewSyncStore[int, int]()
	for i := 0; i < b.N; i++ {
		store.Set(i, i)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = store.GetAny(1)
		}
	})
}

func asyncSetGet(b *testing.B) {
	store := NewSyncStore[int, int]()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			store.Set(1, 1)
			_, _ = store.Get(1)
		}
	})
}
