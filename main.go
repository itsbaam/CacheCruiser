package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := flag.Int("port", 0, "Port on which CacheCruiser listens")
	origin := flag.String("origin", "", "Origin server URL")
	clearCache := flag.Bool("clear-cache", false, "Clear cache and exit")
	flag.Parse()

	fmt.Println(">> port:", *port)
	fmt.Println(">> origin:", *origin)
	fmt.Println(">> clear-cache:", *clearCache)

	if *clearCache {
		fmt.Println("Cache cleared | exiting program...")
		return
	}

	if *port == 0 || *origin == "" {
		flag.Usage()
		log.Fatal("Error: both --port and --origin are required")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Received %s response for %s\n", r.Method, r.URL.RequestURI())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	addr := fmt.Sprintf(":%d", *port)

	log.Printf("Starting CacheCruiser on %s, forwarding to %s", addr, *origin)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Error starting CacheCruiser: %v", err)
	}

	// log.Fatal(http.ListenAndServe(":8080", nil))
}
