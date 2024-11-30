package net

import (
	"sync"
)

type Cache[K comparable, V any] struct {
	data sync.Map
}

func NewCache[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{}
}

func (c *Cache[K, V]) Get(key K) (value V, ok bool) {
	rawValue, ok := c.data.Load(key)
	if !ok {
		var zeroValue V
		return zeroValue, false
	}
	return rawValue.(V), true
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.data.Store(key, value)
}

func (c *Cache[K, V]) Flush() {
	c.data = sync.Map{}
}
