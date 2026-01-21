package usecase

import (
	"sync"
	"time"
)

type CacheService struct {
	mu    sync.RWMutex
	cache map[string]*cacheEntry
	ttl   time.Duration
}

type cacheEntry struct {
	value     interface{}
	createdAt time.Time
}

func NewCacheService(ttl time.Duration) *CacheService {
	return &CacheService{
		mu:    sync.RWMutex{},
		cache: make(map[string]*cacheEntry),
		ttl:   ttl,
	}
}

func (c *CacheService) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.cache[key]
	if !exists {
		return nil, false
	}

	if entry.createdAt.Add(c.ttl).Before(time.Now()) {
		return nil, false
	}

	return entry.value, true
}

func (c *CacheService) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry := &cacheEntry{
		value:     value,
		createdAt: time.Now(),
	}

	c.cache[key] = entry
}

func (c *CacheService) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.cache, key)
}

func (c *CacheService) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]*cacheEntry)
}

type CacheStats struct {
	TotalEntries int `json:"total_entries"`
	TTL          int `json:"ttl_seconds"`
}

func (c *CacheService) GetStats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return CacheStats{
		TotalEntries: len(c.cache),
		TTL:          int(c.ttl.Seconds()),
	}
}
