package cache

import (
	"net/http"
)

type CacheKey string

// CachedResponse holds what we need to store and replay an HTTP response.
type CachedResponse struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

// Cache defines the interface for caching HTTP responses
type Cache interface {
	Get(key CacheKey) (*CachedResponse, bool)
	Set(key CacheKey, response *CachedResponse)
	Clear()
}

// Creates a cache key from an HTTP request
func GenerateKey(r *http.Request) CacheKey {
	// For a simple implementation, use method + URL
	// More sophisticated implementations might include headers,
	// query parameters, or request body for POST requests
	return CacheKey(r.Method + "|" + r.URL.String())
}
