package middleware

import (
	"net/http"
)

// CSRFProtection is a middleware that protects against CSRF attacks by requiring
// a custom header (X-Requested-With) on state-changing requests.
// This works because standard HTML forms cannot set custom headers, and SameSite
// cookies are sent on top-level navigations (like form POSTs).
// By requiring a header that only JS can set (via fetch/XHR), we ensure the request
// was initiated by our frontend application.
func CSRFProtection(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Safe methods are exempt
		if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions || r.Method == http.MethodTrace {
			next.ServeHTTP(w, r)
			return
		}

		// Check for X-Requested-With header
		if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
			http.Error(w, "CSRF Protection: Missing X-Requested-With header", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
