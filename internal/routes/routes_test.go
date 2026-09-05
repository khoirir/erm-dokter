package routes_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"erm-dokter/internal/auth"
	"erm-dokter/internal/berkasdigital"
	"erm-dokter/internal/docs"
	"erm-dokter/internal/health"
	"erm-dokter/internal/laboratorium"
	"erm-dokter/internal/master"
	"erm-dokter/internal/obat"
	"erm-dokter/internal/pemeriksaan"
	"erm-dokter/internal/penilaianmedis"
	"erm-dokter/internal/radiologi"
	"erm-dokter/internal/rawatjalan"
	"erm-dokter/internal/resep"
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
		pemeriksaan.NewHandler(nil, "key"),
		resep.NewHandler(nil, "key"),
		rujukaninternal.NewHandler(nil, "key"),
		penilaianmedis.NewHandler(nil, "key"),
		tindakan.NewHandler(nil, "key"),
		laboratorium.NewHandler(nil, "key"),
		radiologi.NewHandler(nil, "key"),
		berkasdigital.NewHandler(nil, "key"),
		"test-jwt-secret",
	)

	// Must not panic
	routeCfg.Setup()

	handler := routeCfg.BuildHandler("*")
	if handler == nil {
		t.Fatal("Expected non-nil http.Handler")
	}

	// Test health route through built handler
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200 on /health, got %d", rr.Code)
	}
}
