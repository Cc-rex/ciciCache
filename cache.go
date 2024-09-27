package ciciCache

import (
	"ciciCache/lru"
	"sync"
)

// TODO: Add concurrency features to lru.

type cache struct {
	mu       sync.Mutex
	lru      *lru.Cache
	capacity int64 // the maximum capacity of cache
}

func newCache(capacity int64) *cache {
	return &cache{capacity: capacity}
}

func (c *cache) Add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.capacity, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) Get(key string) (ByteView, bool) {
	if c.lru == nil {
		return ByteView{}, false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.lru.Get(key)
	if ok {
		return value.(ByteView), true
	}
	return ByteView{}, false

}
