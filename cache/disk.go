package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type DiskCache struct {
	baseDir string
	mutex   sync.RWMutex
}

type diskCacheEntry struct {
	Response *CachedResponse `json:"response"`
	Expiry   time.Time       `json:"expiry"`
}

func NewDiskCache(baseDir string) (*DiskCache, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to cache directory: %w", err)
	}

	return &DiskCache{
		baseDir: baseDir,
	}, nil
}

func (c *DiskCache) getCacheFilePath(key CacheKey) string {
	// Convert the key to a valid filename using base64 encoding to handle special characters
	filename := filepath.Join(c.baseDir, string(key))
	return filename
}

func (c *DiskCache) Get(key CacheKey) (*CachedResponse, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	filePath := c.getCacheFilePath(key)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, false
	}

	var entry diskCacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, false
	}

	// Check if entry has expired
	if !entry.Expiry.IsZero() && time.Now().After(entry.Expiry) {
		// Remove expired entry
		c.mutex.RUnlock()
		c.mutex.Lock()
		os.Remove(filePath)
		c.mutex.Unlock()
		c.mutex.RLock()
		return nil, false
	}

	return entry.Response, true
}

func (c *DiskCache) Set(key CacheKey, response *CachedResponse) {
	c.SetWithExpiry(key, response, 0)
}

func (c *DiskCache) SetWithExpiry(key CacheKey, response *CachedResponse, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry := diskCacheEntry{
		Response: response,
	}

	if ttl > 0 {
		entry.Expiry = time.Now().Add(ttl)
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return
	}

	filePath := c.getCacheFilePath(key)
	os.WriteFile(filePath, data, 0644)
}

func (c *DiskCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Remove all file in cache dir
	os.RemoveAll(c.baseDir)
	os.MkdirAll(c.baseDir, 0755)
}
