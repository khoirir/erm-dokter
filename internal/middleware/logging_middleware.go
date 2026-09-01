package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"erm-dokter/internal/pkg/token"
)

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (sw *statusResponseWriter) WriteHeader(code int) {
	sw.statusCode = code
	sw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	logger := log.New(os.Stdout, "[HTTP]\t", log.Ldate|log.Ltime)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()

		sw := &statusResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		user := extractUserFromRequest(r)

		next.ServeHTTP(sw, r)

		duration := time.Since(start)

		clientIP := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			clientIP = forwarded
		}

		logger.Printf("[%s] %s %s %d %v %s",
			user,
			r.Method,
			r.URL.Path,
			sw.statusCode,
			duration,
			clientIP,
		)
	})
}

func extractUserFromRequest(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "-"
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return "-"
	}

	claims, err := token.ExtractClaimsUnverified(parts[1])
	if err != nil || claims.KodeDokter == "" {
		return "-"
	}

	return claims.KodeDokter
}
