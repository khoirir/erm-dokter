package pemeriksaan_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"erm-dokter/internal/middleware"
	"erm-dokter/internal/pemeriksaan"
	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/pkg/token"
	"erm-dokter/internal/shared"
)

const testEncKey = "bafaa956-751d-4f59-98cc-574ee9dfe9f6"

type mockPemeriksaanService struct {
	daftarPemeriksaanFn     func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter pemeriksaan.FilterDaftarPemeriksaan) ([]pemeriksaan.Pemeriksaan, shared.PaginationMeta, error)
	daftarPemeriksaanByRMFn func(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter pemeriksaan.FilterDaftarPemeriksaan) ([]pemeriksaan.Pemeriksaan, shared.PaginationMeta, error)
	detailPemeriksaanFn     func(ctx context.Context, id pemeriksaan.IdPemeriksaan, statusLanjut shared.StatusLanjut) (*pemeriksaan.Pemeriksaan, error)
	simpanPemeriksaanFn     func(ctx context.Context, kodeDokter string, statusLanjut shared.StatusLanjut, req pemeriksaan.SimpanPemeriksaanRequest) (*pemeriksaan.Pemeriksaan, error)
	updatePemeriksaanFn     func(ctx context.Context, kodeDokter string, id pemeriksaan.IdPemeriksaan, statusLanjut shared.StatusLanjut, req pemeriksaan.UpdatePemeriksaanRequest) (*pemeriksaan.Pemeriksaan, error)
	hapusPemeriksaanFn      func(ctx context.Context, kodeDokter string, id pemeriksaan.IdPemeriksaan, statusLanjut shared.StatusLanjut) error
	daftarKesadaranFn       func(ctx context.Context) []pemeriksaan.OpsiReferensi
}

func (m *mockPemeriksaanService) DaftarPemeriksaan(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter pemeriksaan.FilterDaftarPemeriksaan) ([]pemeriksaan.Pemeriksaan, shared.PaginationMeta, error) {
	if m.daftarPemeriksaanFn != nil {
		return m.daftarPemeriksaanFn(ctx, noRawat, statusLanjut, filter)
	}
	return nil, shared.PaginationMeta{}, nil
}

func (m *mockPemeriksaanService) DaftarPemeriksaanByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter pemeriksaan.FilterDaftarPemeriksaan) ([]pemeriksaan.Pemeriksaan, shared.PaginationMeta, error) {
	if m.daftarPemeriksaanByRMFn != nil {
		return m.daftarPemeriksaanByRMFn(ctx, noRM, statusLanjut, filter)
	}
	return nil, shared.PaginationMeta{}, nil
}

func (m *mockPemeriksaanService) DetailPemeriksaan(ctx context.Context, id pemeriksaan.IdPemeriksaan, statusLanjut shared.StatusLanjut) (*pemeriksaan.Pemeriksaan, error) {
	if m.detailPemeriksaanFn != nil {
		return m.detailPemeriksaanFn(ctx, id, statusLanjut)
	}
	return nil, nil
}

func (m *mockPemeriksaanService) SimpanPemeriksaan(ctx context.Context, kodeDokter string, statusLanjut shared.StatusLanjut, req pemeriksaan.SimpanPemeriksaanRequest) (*pemeriksaan.Pemeriksaan, error) {
	if m.simpanPemeriksaanFn != nil {
		return m.simpanPemeriksaanFn(ctx, kodeDokter, statusLanjut, req)
	}
	return nil, nil
}

func (m *mockPemeriksaanService) UpdatePemeriksaan(ctx context.Context, kodeDokter string, id pemeriksaan.IdPemeriksaan, statusLanjut shared.StatusLanjut, req pemeriksaan.UpdatePemeriksaanRequest) (*pemeriksaan.Pemeriksaan, error) {
	if m.updatePemeriksaanFn != nil {
		return m.updatePemeriksaanFn(ctx, kodeDokter, id, statusLanjut, req)
	}
	return nil, nil
}

func (m *mockPemeriksaanService) HapusPemeriksaan(ctx context.Context, kodeDokter string, id pemeriksaan.IdPemeriksaan, statusLanjut shared.StatusLanjut) error {
	if m.hapusPemeriksaanFn != nil {
		return m.hapusPemeriksaanFn(ctx, kodeDokter, id, statusLanjut)
	}
	return nil
}

func (m *mockPemeriksaanService) DaftarKesadaran(ctx context.Context) []pemeriksaan.OpsiReferensi {
	if m.daftarKesadaranFn != nil {
		return m.daftarKesadaranFn(ctx)
	}
	return []pemeriksaan.OpsiReferensi{{Label: "Compos Mentis", Value: "Compos Mentis"}}
}

func authMw(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), middleware.UserClaimKey, &token.Claims{KodeDokter: "DR01"})
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func TestPemeriksaanHandler_DaftarPemeriksaan(t *testing.T) {
	encKunjungan, _ := crypto.Encrypt("2026/09/03/000001", testEncKey)
	mockSvc := &mockPemeriksaanService{
		daftarPemeriksaanFn: func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter pemeriksaan.FilterDaftarPemeriksaan) ([]pemeriksaan.Pemeriksaan, shared.PaginationMeta, error) {
			return []pemeriksaan.Pemeriksaan{
				{
					NoRawat:            noRawat,
					TanggalPemeriksaan: "2026-09-03",
					JamPemeriksaan:     "09:00:00",
					SuhuTubuh:          "36.5",
					Tensi:              "120/80",
					Kesadaran:          "Compos Mentis",
				},
			}, shared.NewPaginationMeta(1, 1, 20), nil
		},
	}

	handler := pemeriksaan.NewHandler(mockSvc, testEncKey)
	mux := http.NewServeMux()
	noOpMw := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, authMw, noOpMw)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pemeriksaan/"+encKunjungan+"/Ralan?page=1&limit=20", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rr.Code)
	}

	var resp struct {
		Success bool                      `json:"success"`
		Data    []pemeriksaan.Pemeriksaan `json:"data"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(resp.Data) != 1 || resp.Data[0].Id == "" {
		t.Errorf("Expected encrypted Id in pemeriksaan response, got %+v", resp.Data)
	}
}

func TestPemeriksaanHandler_SimpanPemeriksaan(t *testing.T) {
	encKunjungan, _ := crypto.Encrypt("2026/09/03/000001", testEncKey)
	now := time.Now()

	mockSvc := &mockPemeriksaanService{
		simpanPemeriksaanFn: func(ctx context.Context, kodeDokter string, statusLanjut shared.StatusLanjut, req pemeriksaan.SimpanPemeriksaanRequest) (*pemeriksaan.Pemeriksaan, error) {
			return &pemeriksaan.Pemeriksaan{
				NoRawat:            req.NoRawat,
				TanggalPemeriksaan: req.TanggalPemeriksaan,
				JamPemeriksaan:     req.JamPemeriksaan,
				Keluhan:            req.Keluhan,
			}, nil
		},
	}

	handler := pemeriksaan.NewHandler(mockSvc, testEncKey)
	mux := http.NewServeMux()
	noOpMw := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, authMw, noOpMw)

	body, _ := json.Marshal(pemeriksaan.SimpanPemeriksaanRequest{
		NoRawat: "2026/09/03/000001",
		DataPemeriksaan: pemeriksaan.DataPemeriksaan{
			TanggalPemeriksaan: now.Format("2006-01-02"),
			JamPemeriksaan:     now.Format("15:04:05"),
			Kesadaran:          pemeriksaan.KesadaranComposMentis,
			Keluhan:            "Nyeri dada",
			Pemeriksaan:        "Cor dbn",
			Penilaian:           "Atypical chest pain",
			Instruksi:           "Bed rest",
			RencanaTindakLanjut: "Observasi IGD",
			Evaluasi:            "Keluhan berkurang",
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/pemeriksaan/"+encKunjungan+"/Ralan", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Expected status 201 Created, got %d: %s", rr.Code, rr.Body.String())
	}
}
