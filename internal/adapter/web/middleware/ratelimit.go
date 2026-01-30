package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// maxClients prevents memory exhaustion (DoS) by limiting the map size.
const maxClients = 10000

type client struct {
	count     int
	expiresAt time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]*client
	limit   int
	window  time.Duration
	stop    chan struct{}
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*client),
		limit:   limit,
		window:  window,
		stop:    make(chan struct{}),
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) Stop() {
	close(rl.stop)
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			for ip, c := range rl.clients {
				if time.Now().After(c.expiresAt) {
					delete(rl.clients, ip)
				}
			}
			rl.mu.Unlock()
		case <-rl.stop:
			return
		}
	}
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// We use RemoteAddr instead of X-Forwarded-For to prevent IP spoofing.
		// If behind a trusted proxy, this logic should be updated to verify the proxy chain.
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		if ip == "" {
			// Should not happen for HTTP/1.1+
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		rl.mu.Lock()

		c, exists := rl.clients[ip]

		// DoS Protection: If map is full, and we need to add a new entry, reset it.
		// This is a simple strategy to prevent OOM.
		if !exists && len(rl.clients) >= maxClients {
			rl.clients = make(map[string]*client)
			exists = false // ensure we treat as new
		}
		if !exists || time.Now().After(c.expiresAt) {
			rl.clients[ip] = &client{count: 1, expiresAt: time.Now().Add(rl.window)}
		} else {
			c.count++
			if c.count > rl.limit {
				rl.mu.Unlock()
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}
		}
		rl.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
