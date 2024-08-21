package main

import (
	"container/list"
)

// chatgpt code guh too lazy

type Cache struct {
	capacity int
	items    map[string]*list.Element
	order    *list.List
}

type Item struct {
	key   string
	value SimplifiedData
}

func NewCache(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		order:    list.New(),
	}
}

func (c *Cache) Get(key string) (SimplifiedData, bool) {
	if elem, found := c.items[key]; found {
		c.order.MoveToFront(elem)
		return elem.Value.(*Item).value, true
	}
	return SimplifiedData{}, false
}

func (c *Cache) Put(key string, value SimplifiedData) {
	if elem, found := c.items[key]; found {
		c.order.MoveToFront(elem)
		elem.Value.(*Item).value = value
		return
	}

	if c.order.Len() >= c.capacity {
		oldest := c.order.Back()
		if oldest != nil {
			c.order.Remove(oldest)
			delete(c.items, oldest.Value.(*Item).key)
		}
	}

	item := &Item{key: key, value: value}
	elem := c.order.PushFront(item)
	c.items[key] = elem
}

func (c *Cache) Flush() {
	c.items = make(map[string]*list.Element)
	c.order = list.New()
}
