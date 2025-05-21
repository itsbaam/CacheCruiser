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
	cacheType := flag.String("cache-type", "memory", "Cache type: 'memory' or 'disk'")
	cacheDir := flag.String("cache-dir", "./disk-cache-data", "Directory for disk cache (only used with --cache-type=disk)")
	flag.Parse()

	var cacheImpl cache.Cache
	var err error

	switch *cacheType {
	case "memory":
		cacheImpl = cache.NewMemoryCache()
	case "disk":
		cacheImpl, err = cache.NewDiskCache(*cacheDir)
		if err != nil {
			log.Fatalf("Failed to create disk cache: %v", err)
		}
	default:
		log.Fatalf("Unknown cache type: %s", *cacheType)
	}

	if *clearCache {
		cacheImpl.Clear()
		fmt.Println("Cache cleared | exiting program...")
		return
	}

	if *port == 0 || *origin == "" {
		flag.Usage()
		log.Fatal("Error: both --port and --origin are required")
	}

	// Create and start proxy server
	proxyServer, err := proxy.NewProxyServer(*port, *origin, cacheImpl)
	if err != nil {
		log.Fatalf("Failed to create proxy server: %v", err)
	}

	// Start the server (this blocks until the server exits)
	if err := proxyServer.Start(); err != nil {
		log.Fatalf("Error starting CacheCruiser server: %v", err)
	}
}
