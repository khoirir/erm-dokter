package routes_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"erm-dokter/internal/auth"
	"erm-dokter/internal/berkasdigital"
	"erm-dokter/internal/diagnosa"
	"erm-dokter/internal/docs"
	"erm-dokter/internal/health"
	"erm-dokter/internal/laboratorium"
	"erm-dokter/internal/master"
	"erm-dokter/internal/obat"
	"erm-dokter/internal/pasien"
	"erm-dokter/internal/pemeriksaan"
	"erm-dokter/internal/penilaianmedis"
	"erm-dokter/internal/radiologi"
	"erm-dokter/internal/rawatinap"
	"erm-dokter/internal/rawatjalan"
	"erm-dokter/internal/resep"
	"erm-dokter/internal/resumepasien"
	"erm-dokter/internal/routes"
	"erm-dokter/internal/rujukaninternal"
	"erm-dokter/internal/tindakan"
)

func TestRouteConfig_SetupAndBuildHandler(t *testing.T) {
	routeCfg := routes.NewRouteConfig(
		health.NewHandler(),
		docs.NewHandler(),
		auth.NewHandler(nil),
		master.NewHandler(nil),
		obat.NewHandler(nil, "key"),
		rawatjalan.NewHandler(nil, "key"),
		rawatinap.NewHandler(nil, "key"),
		pemeriksaan.NewHandler(nil, "key"),
		resep.NewHandler(nil, "key"),
		rujukaninternal.NewHandler(nil, "key"),
		penilaianmedis.NewHandler(nil, "key"),
		tindakan.NewHandler(nil, "key"),
		laboratorium.NewHandler(nil, "key"),
		radiologi.NewHandler(nil, "key"),
		berkasdigital.NewHandler(nil, "key"),
		resumepasien.NewHandler(nil, "key"),
		pasien.NewHandler(nil, "key"),
		diagnosa.NewHandler(nil, "key"),
		"test-jwt-secret",
		"test-service-api-key",
	)

	routeCfg.Setup()

	handler := routeCfg.BuildHandler("*")
	if handler == nil {
		t.Fatal("Expected non-nil http.Handler")
	}

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200 on /health, got %d", rr.Code)
	}

	rateLimitHit := false
	routeCfg.LoginRateLimitMiddleware = func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			rateLimitHit = true
			w.WriteHeader(http.StatusTooManyRequests)
		}
	}
	routeCfg.Mux = http.NewServeMux()
	routeCfg.Setup()
	handlerWithRateLimit := routeCfg.BuildHandler("*")

	loginReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	loginRR := httptest.NewRecorder()
	handlerWithRateLimit.ServeHTTP(loginRR, loginReq)

	if !rateLimitHit {
		t.Errorf("Expected LoginRateLimitMiddleware to be executed on /api/v1/auth/login")
	}
	if loginRR.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status 429 on /api/v1/auth/login, got %d", loginRR.Code)
	}
}
