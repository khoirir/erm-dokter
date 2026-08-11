package middleware

import (
	"net/http"
	"sync"
	"time"

	"erm-dokter/pkg/response"
)

// rateLimiter adalah rate limiter in-memory berbasis IP menggunakan token bucket pattern.
type rateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     int           // jumlah request yang diizinkan per window
	window   time.Duration // durasi window
}

type visitor struct {
	tokens    int
	lastReset time.Time
}

// NewRateLimiter membuat rate limiter baru.
// rate: jumlah request per window. window: durasi window.
func NewRateLimiter(rate int, window time.Duration) *rateLimiter {
	rl := &rateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		window:   window,
	}

	// Goroutine untuk membersihkan visitor lama setiap 5 menit
	go rl.cleanupLoop()

	return rl
}

func (rl *rateLimiter) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastReset) > rl.window*2 {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		rl.visitors[ip] = &visitor{
			tokens:    rl.rate - 1,
			lastReset: time.Now(),
		}
		return true
	}

	// Reset tokens jika window sudah lewat
	if time.Since(v.lastReset) > rl.window {
		v.tokens = rl.rate - 1
		v.lastReset = time.Now()
		return true
	}

	if v.tokens <= 0 {
		return false
	}

	v.tokens--
	return true
}

// RateLimitMiddleware membatasi jumlah request per IP per window.
func RateLimitMiddleware(rate int, window time.Duration) func(http.HandlerFunc) http.HandlerFunc {
	limiter := NewRateLimiter(rate, window)

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr

			// Coba ambil IP asli dari header (jika di belakang reverse proxy)
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				ip = forwarded
			} else if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
				ip = realIP
			}

			if !limiter.allow(ip) {
				response.Error(w, http.StatusTooManyRequests, "Terlalu banyak percobaan. Silakan coba lagi nanti.", nil)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}
