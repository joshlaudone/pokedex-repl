package pokecache

import (
	"sync"
	"time"
)

type PokeCache struct {
	cache map[string]cacheEntry
	mu    sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func New(interval time.Duration) *PokeCache {
	cache := PokeCache{cache: make(map[string]cacheEntry)}
	cache.reapLoop(interval)
	return &cache
}

func (pc *PokeCache) Add(key string, val []byte) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	pc.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (pc *PokeCache) Get(key string) ([]byte, bool) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	entry, exists := pc.cache[key]
	if !exists {
		return []byte{}, false
	}

	return entry.val, true
}

func (pc *PokeCache) removeIfExpired(key string, interval time.Duration) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	entry, ok := pc.cache[key]
	if !ok {
		return
	}

	if time.Since(entry.createdAt) > interval {
		delete(pc.cache, key)
	}
}

func (pc *PokeCache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			for key := range pc.cache {
				pc.removeIfExpired(key, interval)
			}
		}
	}()
}
