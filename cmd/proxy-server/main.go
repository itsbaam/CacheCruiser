package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/itsbaam/cachecruiser/cache"
	"github.com/itsbaam/cachecruiser/proxy"
)

func main() {
	port := flag.Int("port", 0, "Port on which CacheCruiser listens")
	origin := flag.String("origin", "", "Origin server URL")
	clearCache := flag.Bool("clear-cache", false, "Clear cache and exit")
	flag.Parse()

	memoryCache := cache.NewMemoryCache()

	if *clearCache {
		memoryCache.Clear()
		fmt.Println("Cache cleared | exiting program...")
		return
	}

	if *port == 0 || *origin == "" {
		flag.Usage()
		log.Fatal("Error: both --port and --origin are required")
	}

	// Create and start proxy server
	proxyServer, err := proxy.NewProxyServer(*port, *origin, memoryCache)
	if err != nil {
		log.Fatalf("Failed to create proxy server: %v", err)
	}

	// Start the server (this blocks until the server exits)
	if err := proxyServer.Start(); err != nil {
		log.Fatalf("Error starting CacheCruiser server: %v", err)
	}
}
