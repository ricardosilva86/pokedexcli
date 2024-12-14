package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	Cache map[string]CacheEntry
	mu    sync.RWMutex
}

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	return &Cache{
		Cache: make(map[string]CacheEntry),
	}
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Cache[key] = CacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.Cache[key]
	if !exists {
		return nil, false
	}

	return entry.val, true
}

func (c *Cache) ReapLoop(internal time.Duration) {
	ticker := time.NewTicker(internal)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			for key, entry := range c.Cache {
				if time.Since(entry.createdAt) > internal {
					delete(c.Cache, key)
				}
			}
			c.mu.Unlock()
		}
	}
}
