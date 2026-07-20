package main

import (
	"sync"
	"time"
)

type cacheEntry struct {
	created time.Time
	val     []byte
}

type Cache map[string]cacheEntry

func createCache(interval time.Duration) Cache {
	cache := Cache{}
	go cache.reapLoop(interval)
	return cache
}

func (c Cache) Add(key string, val []byte) {
	c[key] = cacheEntry{
		created: time.Now(),
		val:     val,
	}
}

func (c Cache) Get(key string) ([]byte, bool) {
	entry, ok := c[key]
	if !ok {
		return []byte{}, false
	}
	return entry.val, true
}

func (c Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	mux := &sync.Mutex{}
	for {
		time.Sleep(time.Second * 5)
		for entry := range c {
			if time.Since(c[entry].created) > interval {
				mux.Lock()
				delete(c, entry)
				mux.Unlock()
			}
		}
	}
}
