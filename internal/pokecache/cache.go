package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cache    map[string]cacheEntry
	interval time.Duration
	m        *sync.RWMutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	cache := Cache{
		cache:    make(map[string]cacheEntry),
		interval: interval,
		m:        &sync.RWMutex{},
	}

	go cache.reapLoop()

	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.m.Lock()
	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.m.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.m.RLock()
	entry, ok := c.cache[key]
	c.m.RUnlock()
	return entry.val, ok
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	for range ticker.C {
		c.m.Lock()
		for k, v := range c.cache {
			reapDeadline := v.createdAt.Add(c.interval)
			if time.Now().After(reapDeadline) {
				delete(c.cache, k)
			}
		}
		c.m.Unlock()
	}
}
