# CacheCruiser

_A CLI caching proxy server._

## Overview
CacheCruiser forwards incoming HTTP requests to a specified origin server, caches the responses on disk, and returns cached responses on repeat requests to improve performance and reduce load on the origin.

## Prerequisites
- Go 1.18+ installed
- Internet connectivity to reach the origin server

## Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/cachecruiser.git
   cd cachecruiser
   ```
2. Build the binary:
   ```bash
   go build -o cachecruiser main.go
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

CacheCruiser makes it easy to speed up repeat requests by transparently caching responses—no configuration beyond port and origin needed.

