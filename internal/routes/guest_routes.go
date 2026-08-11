package routes

import (
	"time"

	"erm-dokter/internal/middleware"
)

func (c *RouteConfig) SetupGuestRoute() {
	c.Mux.HandleFunc("GET /health", c.HealthHandler.HealthCheck)

	loginRateLimit := middleware.RateLimitMiddleware(10, 1*time.Minute)
	c.Mux.HandleFunc("POST /api/v1/auth/login", loginRateLimit(c.AuthHandler.Login))
}
