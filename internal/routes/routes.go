package routes

import (
	"net/http"
	"time"

	"erm-dokter/internal/auth"
	"erm-dokter/internal/docs"
	"erm-dokter/internal/health"
	"erm-dokter/internal/middleware"
	"erm-dokter/internal/pemeriksaan"
	"erm-dokter/internal/penjamin"
	"erm-dokter/internal/rawatjalan"
)

type RouteConfig struct {
	Mux                *http.ServeMux
	HealthHandler      *health.Handler
	DocsHandler        *docs.Handler
	AuthHandler        *auth.Handler
	PenjaminHandler    *penjamin.Handler
	RawatJalanHandler  *rawatjalan.Handler
	PemeriksaanHandler *pemeriksaan.Handler
	AuthMiddleware     func(http.HandlerFunc) http.HandlerFunc
	TimeoutMiddleware  func(http.HandlerFunc) http.HandlerFunc
}

func NewRouteConfig(
	healthHandler *health.Handler,
	docsHandler *docs.Handler,
	authHandler *auth.Handler,
	penjaminHandler *penjamin.Handler,
	rawatJalanHandler *rawatjalan.Handler,
	pemeriksaanHandler *pemeriksaan.Handler,
	jwtSecret string,
) *RouteConfig {
	return &RouteConfig{
		Mux:                http.NewServeMux(),
		HealthHandler:      healthHandler,
		DocsHandler:        docsHandler,
		AuthHandler:        authHandler,
		PenjaminHandler:    penjaminHandler,
		RawatJalanHandler:  rawatJalanHandler,
		PemeriksaanHandler: pemeriksaanHandler,
		AuthMiddleware:     middleware.JWTMiddleware(jwtSecret),
		TimeoutMiddleware:  middleware.TimeoutMiddleware(30 * time.Second),
	}
}

func (c *RouteConfig) Setup() {
	loginRateLimit := middleware.RateLimitMiddleware(10, 1*time.Minute)

	c.HealthHandler.RegisterRoutes(c.Mux)
	c.DocsHandler.RegisterRoutes(c.Mux)
	c.AuthHandler.RegisterRoutes(c.Mux, loginRateLimit)
	c.PenjaminHandler.RegisterRoutes(c.Mux, c.AuthMiddleware, c.TimeoutMiddleware)
	c.RawatJalanHandler.RegisterRoutes(c.Mux, c.AuthMiddleware, c.TimeoutMiddleware)
	c.PemeriksaanHandler.RegisterRoutes(c.Mux, c.AuthMiddleware, c.TimeoutMiddleware)
}

func (c *RouteConfig) BuildHandler(corsOrigin string) http.Handler {
	var h http.Handler = c.Mux
	h = middleware.CORSMiddleware(corsOrigin)(h)
	h = middleware.LoggingMiddleware(h)
	return h
}
