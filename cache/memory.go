package cache

import (
	"sync"
	"time"
)

type MemoryCache struct {
	entries map[CacheKey]*cacheEntry
	mutex   sync.RWMutex
}

// cacheEntry is a struct that holds a CachedResponse and an expiry time
// This is used to implement the on-demand expiration strategy. If the expiry functionality is not needed, MemoryCache entries can map directly to CachedResponse without the cacheEntry struct.
type cacheEntry struct {
	response *CachedResponse
	expiry   time.Time
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		entries: make(map[CacheKey]*cacheEntry),
	}
}

func (c *MemoryCache) Get(key CacheKey) (*CachedResponse, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	entry, exists := c.entries[key]
	if !exists {
		return nil, false
	}

	// This a lazy deletion approach, also called on-demand expiration (only occurs if the expired cached content is requested)
	// Check if expired (if expiry is set)
	if !entry.expiry.IsZero() && time.Now().After(entry.expiry) {
		// Expired entry, remove it
		c.mutex.RUnlock()
		c.mutex.Lock()
		delete(c.entries, key)
		c.mutex.Unlock()
		c.mutex.RLock()
		return nil, false
	}

	return entry.response, true
}

func (c *MemoryCache) Set(key CacheKey, response *CachedResponse) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Store with no expiry by default
	c.entries[key] = &cacheEntry{
		response: response,
	}
}

func (c *MemoryCache) SetWithExpiry(key CacheKey, response *CachedResponse, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entries[key] = &cacheEntry{
		response: response,
		expiry:   time.Now().Add(ttl),
	}
}

func (c *MemoryCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entries = make(map[CacheKey]*cacheEntry)
}
