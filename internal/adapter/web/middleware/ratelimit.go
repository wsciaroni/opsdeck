package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type client struct {
	count     int
	expiresAt time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]*client
	limit   int
	window  time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*client),
		limit:   limit,
		window:  window,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, c := range rl.clients {
			if time.Now().After(c.expiresAt) {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		rl.mu.Lock()
		c, exists := rl.clients[ip]
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
