package cache

import (
	"sync"
	"time"
)

type entry struct {
	value     any
	expiresAt time.Time
}

type Cache struct {
	mu      sync.RWMutex
	storage map[string]entry
}

func New() *Cache {
	return &Cache{
		storage: make(map[string]entry),
	}
}

func (c *Cache) Set(key string, value any, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	e := entry{
		value:     value,
		expiresAt: time.Now().Add(duration),
	}

	c.storage[key] = e
}

func (c *Cache) Get(key string) (any, bool) {
	c.mu.RLock()
	e, exists := c.storage[key]
	c.mu.RUnlock()

	if !exists {
		return nil, false
	}

	if time.Now().Before(e.expiresAt) {
		return e.value, true
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	checkedEntry, ok := c.storage[key]
	if !ok {
		return nil, false
	}

	if time.Now().Before(checkedEntry.expiresAt) {
		return checkedEntry.value, true
	}

	delete(c.storage, key)
	return nil, false
}
