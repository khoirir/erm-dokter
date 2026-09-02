package laboratorium_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"erm-dokter/internal/laboratorium"
	"erm-dokter/internal/shared"
)

type mockService struct {
	getRiwayatLabKunjunganFn func(ctx context.Context, kategori shared.KategoriLab, encryptedIdKunjungan string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) (*laboratorium.HasilLaboratoriumKunjungan, shared.PaginationMeta, error)
	getRiwayatLabPasienFn    func(ctx context.Context, kategori shared.KategoriLab, encryptedIdPasien string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratorium, shared.PaginationMeta, error)
	getDetailHasilLabFn      func(ctx context.Context, kategori shared.KategoriLab, encryptedIdKunjungan string, encryptedIdHasil string) (*laboratorium.HasilLaboratorium, error)
}

func (m *mockService) GetRiwayatLabKunjungan(ctx context.Context, kategori shared.KategoriLab, encryptedIdKunjungan string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) (*laboratorium.HasilLaboratoriumKunjungan, shared.PaginationMeta, error) {
	if m.getRiwayatLabKunjunganFn != nil {
		return m.getRiwayatLabKunjunganFn(ctx, kategori, encryptedIdKunjungan, statusLanjut, filter)
	}
	return nil, shared.PaginationMeta{}, nil
}

func (m *mockService) GetRiwayatLabPasien(ctx context.Context, kategori shared.KategoriLab, encryptedIdPasien string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratorium, shared.PaginationMeta, error) {
	if m.getRiwayatLabPasienFn != nil {
		return m.getRiwayatLabPasienFn(ctx, kategori, encryptedIdPasien, statusLanjut, filter)
	}
	return nil, shared.PaginationMeta{}, nil
}

func (m *mockService) GetDetailHasilLab(ctx context.Context, kategori shared.KategoriLab, encryptedIdKunjungan string, encryptedIdHasil string) (*laboratorium.HasilLaboratorium, error) {
	if m.getDetailHasilLabFn != nil {
		return m.getDetailHasilLabFn(ctx, kategori, encryptedIdKunjungan, encryptedIdHasil)
	}
	return nil, nil
}

func TestHandler_DaftarHasilLab_Success(t *testing.T) {
	mockSvc := &mockService{
		getRiwayatLabKunjunganFn: func(ctx context.Context, kategori shared.KategoriLab, encryptedIdKunjungan string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) (*laboratorium.HasilLaboratoriumKunjungan, shared.PaginationMeta, error) {
			return &laboratorium.HasilLaboratoriumKunjungan{
				HasilPemeriksaan: []laboratorium.HasilLaboratorium{
					{
						Id:           "enc-hasil-1",
						NamaTindakan: "DARAH LENGKAP",
						Kategori:     "PK",
					},
				},
				BerkasDigital: []laboratorium.BerkasDigital{},
			}, shared.NewPaginationMeta(1, 1, 5), nil
		},
	}

	handler := laboratorium.NewHandler(mockSvc)
	mux := http.NewServeMux()
	dummyMiddleware := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, dummyMiddleware, dummyMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/laboratorium/pk/enc-kunjungan-1/Semua?page=1&limit=5", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp struct {
		Code    int                                       `json:"code"`
		Status  string                                    `json:"status"`
		Message string                                    `json:"message"`
		Data    laboratorium.HasilLaboratoriumKunjungan `json:"data"`
		Meta    shared.PaginationMeta                     `json:"meta"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(resp.Data.HasilPemeriksaan) != 1 || resp.Data.HasilPemeriksaan[0].NamaTindakan != "DARAH LENGKAP" {
		t.Errorf("Unexpected response data: %+v", resp.Data)
	}
}

func TestHandler_DaftarHasilLab_InvalidKategori(t *testing.T) {
	handler := laboratorium.NewHandler(&mockService{})
	mux := http.NewServeMux()
	dummyMiddleware := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, dummyMiddleware, dummyMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/laboratorium/radiologi/enc-kunjungan/Semua", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400 for invalid kategori, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_DaftarHasilLab_InvalidStatusLanjut(t *testing.T) {
	handler := laboratorium.NewHandler(&mockService{})
	mux := http.NewServeMux()
	dummyMiddleware := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, dummyMiddleware, dummyMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/laboratorium/pk/enc-kunjungan/InvalidStatus", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400 for invalid status lanjut, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_DaftarHasilLabByRM_Success(t *testing.T) {
	mockSvc := &mockService{
		getRiwayatLabPasienFn: func(ctx context.Context, kategori shared.KategoriLab, encryptedIdPasien string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratorium, shared.PaginationMeta, error) {
			return []laboratorium.HasilLaboratorium{
				{
					Id:           "enc-hasil-2",
					NamaTindakan: "BIOPSI",
					Kategori:     "PA",
				},
			}, shared.NewPaginationMeta(1, 1, 5), nil
		},
	}

	handler := laboratorium.NewHandler(mockSvc)
	mux := http.NewServeMux()
	dummyMiddleware := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, dummyMiddleware, dummyMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/laboratorium/pa/pasien/enc-pasien-1/Semua", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_DetailHasilLab_Success(t *testing.T) {
	mockSvc := &mockService{
		getDetailHasilLabFn: func(ctx context.Context, kategori shared.KategoriLab, encryptedIdKunjungan string, encryptedIdHasil string) (*laboratorium.HasilLaboratorium, error) {
			return &laboratorium.HasilLaboratorium{
				Id:           "enc-hasil-1",
				NamaTindakan: "DARAH LENGKAP",
				Kategori:     "PK",
			}, nil
		},
	}

	handler := laboratorium.NewHandler(mockSvc)
	mux := http.NewServeMux()
	dummyMiddleware := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, dummyMiddleware, dummyMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/laboratorium/pk/enc-kunjungan-1/Semua/enc-hasil-1", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}
}
