package main

import "container/list"

// chatgpt code guh too lazy

type Cache[T any] struct {
	capacity int
	items    map[string]*list.Element
	order    *list.List
}

type Item[T any] struct {
	key   string
	value T
}

func NewCache[T any](capacity int) *Cache[T] {
	return &Cache[T]{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		order:    list.New(),
	}
}

func (c *Cache[T]) Get(key string) (T, bool) {
	if elem, found := c.items[key]; found {
		c.order.MoveToFront(elem)
		return elem.Value.(*Item[T]).value, true
	}
	var zero T
	return zero, false
}

func (c *Cache[T]) Put(key string, value T) {
	if elem, found := c.items[key]; found {
		c.order.MoveToFront(elem)
		elem.Value.(*Item[T]).value = value
		return
	}

	if c.order.Len() >= c.capacity {
		oldest := c.order.Back()
		if oldest != nil {
			c.order.Remove(oldest)
			delete(c.items, oldest.Value.(*Item[T]).key)
		}
	}

	item := &Item[T]{key: key, value: value}
	elem := c.order.PushFront(item)
	c.items[key] = elem
}

func (c *Cache[T]) Flush() {
	c.items = make(map[string]*list.Element)
	c.order = list.New()
}
