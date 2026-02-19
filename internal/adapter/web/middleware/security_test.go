package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSecurityHeaders(t *testing.T) {
	// Ensure we are not in development mode
	originalEnv := os.Getenv("APP_ENV")
	defer os.Setenv("APP_ENV", originalEnv)
	os.Setenv("APP_ENV", "production")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := SecurityHeaders(nextHandler)
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	headers := map[string]string{
		"X-Frame-Options":           "DENY",
		"X-Content-Type-Options":    "nosniff",
		"Referrer-Policy":           "strict-origin-when-cross-origin",
		"X-XSS-Protection":          "1; mode=block",
		"Content-Security-Policy":   "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https://ui-avatars.com; font-src 'self' data:; connect-src 'self'; object-src 'none'; base-uri 'self';",
		"Permissions-Policy":        "camera=(), microphone=(), geolocation=(), payment=()",
		"Strict-Transport-Security": "max-age=63072000; includeSubDomains; preload",
	}

	for key, expected := range headers {
		if val := rr.Header().Get(key); val != expected {
			t.Errorf("expected header %s to be %s, got %s", key, expected, val)
		}
	}
}

func TestSecurityHeadersDevelopment(t *testing.T) {
	// Set development mode
	originalEnv := os.Getenv("APP_ENV")
	defer os.Setenv("APP_ENV", originalEnv)
	os.Setenv("APP_ENV", "development")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := SecurityHeaders(nextHandler)
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Check for standard headers (just one to be sure middleware ran)
	if val := rr.Header().Get("X-Frame-Options"); val != "DENY" {
		t.Errorf("expected X-Frame-Options to be DENY, got %s", val)
	}

	// Check that HSTS is MISSING
	if val := rr.Header().Get("Strict-Transport-Security"); val != "" {
		t.Errorf("expected Strict-Transport-Security to be missing in dev, got %s", val)
	}
}
