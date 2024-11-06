package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStoreV2(t *testing.T) {
	t.Run("success on create", func(t *testing.T) {
		got := NewStore[int, int]()
		assert.NotNil(t, got)
	})

	t.Run("success on set", func(t *testing.T) {
		got := NewStore[int, int]()
		got.Set(1, 2)
		assert.Equal(t, 2, got.data[1])
	})

	t.Run("success on get", func(t *testing.T) {
		got := NewStore[int, int]()
		got.Set(1, 2)
		v, err := got.Get(1)
		assert.NoError(t, err)
		assert.Equal(t, 2, v)
	})

	t.Run("struct", func(t *testing.T) {
		type Mystruct struct {
			Field1 string
			Field2 int
		}
		got := NewStore[string, Mystruct]()
		got.Set("a", Mystruct{Field1: "hello", Field2: 42})
		v, err := got.Get("a")
		assert.NoError(t, err)
		assert.Equal(t, Mystruct{Field1: "hello", Field2: 42}, v)
		got1 := NewStore[string, map[string]int]()
		got1.Set("a", map[string]int{"a": 1})
		v1, err := got1.Get("a")
		assert.NoError(t, err)
		assert.Equal(t, map[string]int{"a": 1}, v1)
	})
}

func TestKeysV2(t *testing.T) {
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
			got := NewStore[int, int]()
			for _, k := range tt.keys {
				got.Set(k, k)
			}
			gotKeys := got.Keys()
			cachedKeys := got.KeysCached()
			if len(tt.want) == 0 {
				assert.Empty(t, gotKeys)
				assert.Empty(t, cachedKeys)
				return
			}
			cachedKeys[0] = 0
			assert.ElementsMatch(t, tt.want, gotKeys)
			// assert.Equal(t, tt.want, cachedKeys)
			assert.Equal(t, tt.want, got.keys)
		})
	}
}

func BenchmarkStore(b *testing.B) {
	b.Run("set", storeSet)
	b.Run("get", storeGet)
	b.Run("set get", storeSetGet)
	b.Run("keys", storeKeys)
	b.Run("keys cached", storeKeysCached)
}

func storeKeysCached(b *testing.B) {
	store := NewStore[int, int]()
	for i := 0; i < b.N; i++ {
		store.Set(i, i)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			store.KeysCached()
		}
	})
}

func storeKeys(b *testing.B) {
	store := NewStore[int, int]()
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

func storeSet(b *testing.B) {
	store := NewStore[int, int]()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			store.Set(1, 1)
		}
	})
}

func storeGet(b *testing.B) {
	store := NewStore[int, int]()
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

func storeSetGet(b *testing.B) {
	store := NewStore[int, int]()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			store.Set(1, 1)
			_, _ = store.Get(1)
		}
	})
}
