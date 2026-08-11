package middleware

import (
	"log"
	"net/http"
	"os"
	"time"
)

// statusResponseWriter adalah wrapper untuk menangkap status code response.
type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (sw *statusResponseWriter) WriteHeader(code int) {
	sw.statusCode = code
	sw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware mencatat setiap incoming request: method, path, status code, duration, dan client IP.
func LoggingMiddleware(next http.Handler) http.Handler {
	logger := log.New(os.Stdout, "[HTTP]\t", log.Ldate|log.Ltime)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		sw := &statusResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(sw, r)

		duration := time.Since(start)

		clientIP := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			clientIP = forwarded
		}

		logger.Printf("%s %s %d %v %s",
			r.Method,
			r.URL.Path,
			sw.statusCode,
			duration,
			clientIP,
		)
	})
}
