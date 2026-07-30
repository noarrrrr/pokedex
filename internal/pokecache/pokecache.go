package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	created time.Time
	val     []byte
}

type Cache struct {
	entries map[string]cacheEntry
	mux     *sync.Mutex
}

func CreateCache(interval time.Duration) *Cache {
	mutex := &sync.Mutex{}
	cache := Cache{
		map[string]cacheEntry{},
		mutex,
	}
	go cache.reapLoop(interval)
	return &cache
}

func (c *Cache) Add(key string, val []byte) {
	c.entries[key] = cacheEntry{
		created: time.Now(),
		val:     val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	entry, ok := c.entries[key]
	if !ok {
		return []byte{}, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		time.Sleep(interval)
		for key, entry := range c.entries {
			if time.Since(entry.created) > interval {
				c.mux.Lock()
				delete(c.entries, key)
				c.mux.Unlock()
			}
		}
	}
}
