package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	CacheEntries map[string]cacheEntry
	mu           sync.Mutex
}
type cacheEntry struct {
	cratedAt time.Time
	val      []byte
}

func NewCache(interval time.Duration) *Cache {
	return &Cache{CacheEntries: []cacheEntry{}}
}

func (c Cache) Add(key string, val []byte) {
	entry := cacheEntry{time.Now(), val}
	c.CacheEntries[key] = entry
}
func (c Cache) Get(key string) ([]byte, bool) {
	val, ok := c.CacheEntries[key]
	if !ok {
		return nil, false
	}
	return val.val, true
}
