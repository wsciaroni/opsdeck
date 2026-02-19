package middleware

import (
	"net/http"
	"os"
)

func SecurityHeaders(next http.Handler) http.Handler {
	// Check environment once at startup
	isDev := os.Getenv("APP_ENV") == "development"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent clickjacking
		w.Header().Set("X-Frame-Options", "DENY")
		// Prevent MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")
		// Control referrer information
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		// Basic XSS protection (for older browsers)
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		// Content Security Policy
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https://ui-avatars.com; font-src 'self' data:; connect-src 'self'; object-src 'none'; base-uri 'self';")
		// Permissions Policy
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=()")

		// Strict Transport Security (HSTS)
		// Only enable in production-like environments to avoid breaking local dev over HTTP
		if !isDev {
			w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		}

		next.ServeHTTP(w, r)
	})
}
