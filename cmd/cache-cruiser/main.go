package main

import (
	"flag"     // CLI flag parsing
	"fmt"      // formatted I/O
	"io"       // reading entire body
	"log"      // structured logging
	"net/http" // HTTP server and client
	"sync"     // concurrency primitives
)

// CachedResponse holds status, headers, and body for caching
// Used in-memory in Phase 1
type CachedResponse struct {
	Status int
	Header http.Header
	Body   []byte
}

// Global in-memory cache and mutex for thread-safe access
var (
	cache      = make(map[string]*CachedResponse)
	cacheMutex = sync.RWMutex{}
)

func main() {
	// 1. Parse CLI flags
	port := flag.Int("port", 0, "Port on which CacheCruiser listens")
	origin := flag.String("origin", "", "Origin server URL")
	clearCache := flag.Bool("clear-cache", false, "Clear cache at startup")
	flag.Parse()

	// 2. Handle cache clearing at startup (but keep the server running)
	if *clearCache {
		cacheMutex.Lock()
		cache = make(map[string]*CachedResponse) // reset the map
		cacheMutex.Unlock()
		fmt.Println("Cache cleared at startup, starting server...")
		// Note: no return, server will start normally
	}

	// 3. Validate required flags
	if *port == 0 || *origin == "" {
		flag.Usage()
		log.Fatal("Error: both --port and --origin are required")
	}

	// 4. Register HTTP handler with in-memory cache
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		key := r.Method + "|" + r.URL.RequestURI()

		// Check cache hit
		cacheMutex.RLock()
		if cr, ok := cache[key]; ok {
			for name, values := range cr.Header {
				for _, v := range values {
					w.Header().Add(name, v)
				}
			}
			w.Header().Set("X-Cache", "HIT")
			w.WriteHeader(cr.Status)
			w.Write(cr.Body)
			cacheMutex.RUnlock()
			return
		}
		cacheMutex.RUnlock()

		// Cache miss: forward to origin
		target := *origin + r.URL.RequestURI()
		log.Printf("Forwarding to %s", target)
		resp, err := http.Get(target)
		if err != nil {
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Failed to read response", http.StatusInternalServerError)
			return
		}

		// Copy response headers
		for name, values := range resp.Header {
			for _, v := range values {
				w.Header().Add(name, v)
			}
		}
		w.Header().Set("X-Cache", "MISS")

		// Send response to client
		w.WriteHeader(resp.StatusCode)
		w.Write(body)

		// Store in cache
		cacheMutex.Lock()
		cache[key] = &CachedResponse{
			Status: resp.StatusCode,
			Header: resp.Header.Clone(),
			Body:   body,
		}
		cacheMutex.Unlock()
	})

	// 5. Start HTTP server
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Starting CacheCruiser on %s, forwarding to %s", addr, *origin)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
