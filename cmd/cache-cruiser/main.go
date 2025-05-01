package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	port := flag.Int("port", 0, "Port on which CacheCruiser listens")
	origin := flag.String("origin", "", "Origin server URL")
	clearCache := flag.Bool("clear-cache", false, "Clear cache and exit")
	flag.Parse()

	if *clearCache {
		fmt.Println("Cache cleared | exiting program...")
		return
	}

	if *port == 0 || *origin == "" {
		flag.Usage()
		log.Fatal("Error: both --port and --origin are required")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Received %s request for %s\n", r.Method, r.URL.RequestURI())

		targetUrl := *origin + r.URL.RequestURI()

		fmt.Println("targetUrl", targetUrl)

		resp, err := http.Get(targetUrl)
		if err != nil {
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
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

	})

	addr := fmt.Sprintf(":%d", *port)

	log.Printf("Starting CacheCruiser on %s, forwarding to %s", addr, *origin)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Error starting CacheCruiser server: %v", err)
	}
}
