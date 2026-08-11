package routes

import (
	"net/http"
	"time"

	"erm-dokter/internal/handler"
	"erm-dokter/internal/middleware"
)

type RouteConfig struct {
	Mux               *http.ServeMux
	HealthHandler     *handler.HealthHandler
	AuthHandler       *handler.AuthHandler
	PenjaminHandler   *handler.PenjaminHandler
	RawatJalanHandler *handler.RawatJalanHandler
	AuthMiddleware    func(http.HandlerFunc) http.HandlerFunc
	TimeoutMiddleware func(http.HandlerFunc) http.HandlerFunc
}

func NewRouteConfig(
	healthHandler *handler.HealthHandler,
	authHandler *handler.AuthHandler,
	penjaminHandler *handler.PenjaminHandler,
	rawatJalanHandler *handler.RawatJalanHandler,
	jwtSecret string,
) *RouteConfig {
	return &RouteConfig{
		Mux:               http.NewServeMux(),
		HealthHandler:     healthHandler,
		AuthHandler:       authHandler,
		PenjaminHandler:   penjaminHandler,
		RawatJalanHandler: rawatJalanHandler,
		AuthMiddleware:    middleware.JWTMiddleware(jwtSecret),
		TimeoutMiddleware: middleware.TimeoutMiddleware(30 * time.Second),
	}
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupAuthRoute() {
	c.setupPenjaminRoutes()
	c.setupRawatJalanRoutes()
}

func (c *RouteConfig) BuildHandler(corsOrigin string) http.Handler {
	var h http.Handler = c.Mux
	h = middleware.CORSMiddleware(corsOrigin)(h)
	h = middleware.LoggingMiddleware(h)
	return h
}
