package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	CacheEntries map[string]cacheEntry
	mu           sync.Mutex
	interval     time.Duration
}
type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cptr := &Cache{CacheEntries: map[string]cacheEntry{}, interval: interval}
	go cptr.reapLoop()
	return cptr
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		for key, value := range c.CacheEntries {
			if value.createdAt.Before(now.Add(-c.interval)) {
				delete(c.CacheEntries, key)
			}
		}
	}

}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry := cacheEntry{time.Now(), val}
	c.CacheEntries[key] = entry
}
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, ok := c.CacheEntries[key]
	if !ok {
		return nil, false
	}
	return val.val, true
}
