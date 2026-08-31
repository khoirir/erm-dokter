package routes

import (
	"net/http"
	"time"

	"erm-dokter/internal/auth"
	"erm-dokter/internal/docs"
	"erm-dokter/internal/health"
	"erm-dokter/internal/master"
	"erm-dokter/internal/middleware"
	"erm-dokter/internal/obat"
	"erm-dokter/internal/pemeriksaan"
	"erm-dokter/internal/rawatjalan"
	"erm-dokter/internal/penilaianmedis"
	"erm-dokter/internal/resep"
	"erm-dokter/internal/rujukaninternal"
)

type RouteConfig struct {
	Mux                    *http.ServeMux
	HealthHandler          *health.Handler
	DocsHandler            *docs.Handler
	AuthHandler            *auth.Handler
	MasterHandler          *master.Handler
	ObatHandler            *obat.Handler
	RawatJalanHandler      *rawatjalan.Handler
	PemeriksaanHandler     *pemeriksaan.Handler
	ResepHandler           *resep.Handler
	RujukanInternalHandler *rujukaninternal.Handler
	PenilaianMedisHandler  *penilaianmedis.Handler
	AuthMiddleware         func(http.HandlerFunc) http.HandlerFunc
	TimeoutMiddleware      func(http.HandlerFunc) http.HandlerFunc
}

func NewRouteConfig(
	healthHandler *health.Handler,
	docsHandler *docs.Handler,
	authHandler *auth.Handler,
	masterHandler *master.Handler,
	obatHandler *obat.Handler,
	rawatJalanHandler *rawatjalan.Handler,
	pemeriksaanHandler *pemeriksaan.Handler,
	resepHandler *resep.Handler,
	rujukanInternalHandler *rujukaninternal.Handler,
	penilaianMedisHandler *penilaianmedis.Handler,
	jwtSecret string,
) *RouteConfig {
	return &RouteConfig{
		Mux:                    http.NewServeMux(),
		HealthHandler:          healthHandler,
		DocsHandler:            docsHandler,
		AuthHandler:            authHandler,
		MasterHandler:          masterHandler,
		ObatHandler:            obatHandler,
		RawatJalanHandler:      rawatJalanHandler,
		PemeriksaanHandler:     pemeriksaanHandler,
		ResepHandler:           resepHandler,
		RujukanInternalHandler: rujukanInternalHandler,
		PenilaianMedisHandler:  penilaianMedisHandler,
		AuthMiddleware:         middleware.JWTMiddleware(jwtSecret),
		TimeoutMiddleware:      middleware.TimeoutMiddleware(30 * time.Second),
	}
}

func (c *RouteConfig) Setup() {
	loginRateLimit := middleware.RateLimitMiddleware(10, 1*time.Minute)

	c.HealthHandler.RegisterRoutes(c.Mux)
	c.DocsHandler.RegisterRoutes(c.Mux)
	c.AuthHandler.RegisterRoutes(c.Mux, loginRateLimit, c.AuthMiddleware, c.TimeoutMiddleware)
	c.MasterHandler.RegisterRoutes(c.Mux, c.AuthMiddleware, c.TimeoutMiddleware)
	c.ObatHandler.RegisterRoutes(c.Mux, c.AuthMiddleware, c.TimeoutMiddleware)
	c.RawatJalanHandler.RegisterRoutes(c.Mux, c.AuthMiddleware, c.TimeoutMiddleware)
	c.PemeriksaanHandler.RegisterRoutes(c.Mux, c.AuthMiddleware, c.TimeoutMiddleware)
	c.ResepHandler.RegisterRoutes(c.Mux, c.AuthMiddleware, c.TimeoutMiddleware)
	c.RujukanInternalHandler.RegisterRoutes(c.Mux, c.AuthMiddleware, c.TimeoutMiddleware)
	c.PenilaianMedisHandler.RegisterRoutes(c.Mux, c.AuthMiddleware, c.TimeoutMiddleware)
}

func (c *RouteConfig) BuildHandler(corsOrigin string) http.Handler {
	var h http.Handler = c.Mux
	h = middleware.CORSMiddleware(corsOrigin)(h)
	h = middleware.LoggingMiddleware(h)
	return h
}
