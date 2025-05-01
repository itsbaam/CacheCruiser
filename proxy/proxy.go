package proxy

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/itsbaam/cachecruiser/cache"
)

// ProxyServer represents a caching HTTP proxy server
type ProxyServer struct {
	Port       int
	Origin     *url.URL
	Cache      cache.Cache
	httpClient *http.Client
}

func NewProxyServer(port int, originURL string, cacheImpl cache.Cache) (*ProxyServer, error) {
	origin, err := url.Parse(originURL)
	if err != nil {
		return nil, fmt.Errorf("invalid origin URL: %w", err)
	}

	return &ProxyServer{
		Port:       port,
		Origin:     origin,
		Cache:      cacheImpl,
		httpClient: &http.Client{},
	}, nil
}

func (p *ProxyServer) Start() error {
	http.HandleFunc("/", p.handleRequest)
	addr := fmt.Sprintf(":%d", p.Port)

	log.Printf("Starting CacheCruiser on %s, forwarding to %s", addr, p.Origin.String())
	return http.ListenAndServe(addr, nil)
}

func (p *ProxyServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s request for %s", r.Method, r.URL.RequestURI())

	// Only cache GET requests
	if r.Method == http.MethodGet {
		cacheKey := cache.GenerateKey(r)

		// Try to get from cache
		if cachedResp, found := p.Cache.Get(cacheKey); found {
			log.Printf("Cache HIT for %s", r.URL.RequestURI())
			p.sendCachedResponse(w, cachedResp, true)
			return
		}

		// Cache miss, forward to origin
		resp, err := p.forwardToOrigin(r)
		if err != nil {
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// Read the entire response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Create a cached response
		cachedResp := &cache.CachedResponse{
			StatusCode: resp.StatusCode,
			Headers:    resp.Header.Clone(),
			Body:       body,
		}

		p.Cache.Set(cacheKey, cachedResp)

		p.sendCachedResponse(w, cachedResp, false)
		return
	}

	// Non-GET requests are not cached
	resp, err := p.forwardToOrigin(r)
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	w.Header().Set("X-Cache", "MISS")

	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("Error copying response body: %v", err)
	}
}

func (p *ProxyServer) forwardToOrigin(r *http.Request) (*http.Response, error) {
	targetURL := p.Origin.String() + r.URL.RequestURI()

	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		return nil, err
	}

	// Copy headers from the original request
	for name, values := range r.Header {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	// Send the request
	return p.httpClient.Do(req)
}

func (p *ProxyServer) sendCachedResponse(w http.ResponseWriter, resp *cache.CachedResponse, cacheHit bool) {
	for name, values := range resp.Headers {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// Set cache status header
	if cacheHit {
		w.Header().Set("X-Cache", "HIT")
	} else {
		w.Header().Set("X-Cache", "MISS")
	}

	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, bytes.NewReader(resp.Body)); err != nil {
		log.Printf("Error writing response body: %v", err)
	}
}
