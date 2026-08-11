package middleware

import (
	"net/http"
)

// CORSMiddleware menambahkan header CORS ke setiap response.
// Mendukung preflight OPTIONS request.
func CORSMiddleware(allowOrigin string) func(http.Handler) http.Handler {
	if allowOrigin == "" {
		allowOrigin = "*"
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")
			w.Header().Set("Access-Control-Max-Age", "3600")

			// Handle preflight request
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
