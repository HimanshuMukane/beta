package cache

import (
	"sync"
	"time"
)

type item[T any] struct {
	value     T
	expiresAt time.Time
}

type MemoryCache[T any] struct {
	mu    sync.RWMutex
	items map[string]item[T]
}

func NewMemoryCache[T any]() *MemoryCache[T] {
	return &MemoryCache[T]{
		items: make(map[string]item[T]),
	}
}

func (c *MemoryCache[T]) Get(key string) (T, bool) {
	c.mu.RLock()
	it, found := c.items[key]
	c.mu.RUnlock()

	if !found {
		var zero T
		return zero, false
	}

	if time.Now().After(it.expiresAt) {
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
		var zero T
		return zero, false
	}

	return it.value, true
}

func (c *MemoryCache[T]) Set(key string, value T, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = item[T]{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
}

func (c *MemoryCache[T]) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

