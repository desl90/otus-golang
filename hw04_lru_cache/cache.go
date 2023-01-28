package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	Cache
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	ptxNewItem := &cacheItem{
		key:   key,
		value: value,
	}

	if item, ok := c.items[key]; ok {
		item.Value = ptxNewItem
		c.queue.MoveToFront(item)

		return true
	}

	item := c.queue.PushFront(ptxNewItem)

	if c.queue.Len() > c.capacity {
		last := c.queue.Back()
		c.queue.Remove(last)

		delete(c.items, last.Value.(*cacheItem).key)
	}

	c.items[key] = item

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]

	if !ok {
		return nil, false
	}

	c.queue.MoveToFront(item)

	return item.Value.(*cacheItem).value, true
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
