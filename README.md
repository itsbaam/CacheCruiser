# CacheCruiser

_A CLI caching proxy server._

## Overview
CacheCruiser forwards incoming HTTP requests to a specified origin server, caches the responses on the cache instance, and returns cached responses on repeat requests to improve performance and reduce load on the origin.

This project currently implements only in-memory caching. Disk caching and Redis caching are planned for future releases.

## Prerequisites
- Go 1.18+ installed
- Internet connectivity

## Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/cachecruiser.git
   cd cachecruiser
   ```
2. Build the binary:
   ```bash
   go build -o cachecruiser cmd/proxy-server/main.go
   ```

## Usage
```bash
./cachecruiser --port <port> --origin <origin-url>
```

### Flags
- `--port` &lt;number&gt;: Port on which CacheCruiser listens (required)
- `--origin` &lt;url&gt;: URL of the origin server to which requests are forwarded (required)
- `--clear-cache`: Remove all cached data and exit

### Examples
```bash
# Start the proxy on port 3000, forwarding to dummyjson.com
./cachecruiser --port 3000 --origin http://dummyjson.com

# Clear the cache and exit
./cachecruiser --clear-cache
```

## Cache Behavior
- On first request to a given path, CacheCruiser forwards to the origin, returns the response with header:
  ```
  X-Cache: MISS
  ```
- On subsequent requests to the same URL, CacheCruiser serves the cached response with header:
  ```
  X-Cache: HIT
  ```

---

CacheCruiser makes it easy to speed up repeat requests by transparently caching responses — no configuration beyond port and origin needed.

## Technical Details

### Caching Behavior
- Only GET requests are cached
- Cache keys are generated using request method + URL
- Thread-safe implementation using read-write mutexes
- Support for time-based expiration of cache entries

### Architecture
CacheCruiser is built with a modular design:
- `Cache` interface that can be implemented by different storages (in-memory, disk, Redis, etc.)
- `MemoryCache` implementation for in-memory storage
- `ProxyServer` for handling HTTP requests and responses

## Development
This project is a learning exercise that comes from [Roadmap.sh Caching Proxy project](https://roadmap.sh/projects/caching-server).

### Testing
Not implemented. May come in future iterations also as a learning exercise.

## Roadmap
- [ ] Disk-based caching
- [ ] Redis-based caching
- [ ] Configurable cache TTL via command-line flags
- [ ] Cache size limits and eviction policies
- [ ] Cache invalidation endpoints
- [ ] Tests

### Disclaimer
As this project is a learning exercise, the features listed in the roadmap may or may not be implemented in the future.
