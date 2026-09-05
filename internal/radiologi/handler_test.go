package radiologi_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/radiologi"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

const testEncryptionKey = "bafaa956-751d-4f59-98cc-574ee9dfe9f6"

type mockRadiologiService struct {
	getRiwayatRadiologiKunjunganFn func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter radiologi.FilterRiwayatRadiologi) ([]radiologi.HasilRadiologi, shared.PaginationMeta, error)
	getRiwayatRadiologiPasienFn    func(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter radiologi.FilterRiwayatRadiologi) ([]radiologi.HasilRadiologi, shared.PaginationMeta, error)
	getDetailHasilRadiologiFn      func(ctx context.Context, idHasil radiologi.IdHasilRadiologi, statusLanjut shared.StatusLanjut) (*radiologi.HasilRadiologi, error)
}

func (m *mockRadiologiService) GetRiwayatRadiologiKunjungan(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter radiologi.FilterRiwayatRadiologi) ([]radiologi.HasilRadiologi, shared.PaginationMeta, error) {
	if m.getRiwayatRadiologiKunjunganFn != nil {
		return m.getRiwayatRadiologiKunjunganFn(ctx, noRawat, statusLanjut, filter)
	}
	return nil, shared.PaginationMeta{}, nil
}

func (m *mockRadiologiService) GetRiwayatRadiologiPasien(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter radiologi.FilterRiwayatRadiologi) ([]radiologi.HasilRadiologi, shared.PaginationMeta, error) {
	if m.getRiwayatRadiologiPasienFn != nil {
		return m.getRiwayatRadiologiPasienFn(ctx, noRM, statusLanjut, filter)
	}
	return nil, shared.PaginationMeta{}, nil
}

func (m *mockRadiologiService) GetDetailHasilRadiologi(ctx context.Context, idHasil radiologi.IdHasilRadiologi, statusLanjut shared.StatusLanjut) (*radiologi.HasilRadiologi, error) {
	if m.getDetailHasilRadiologiFn != nil {
		return m.getDetailHasilRadiologiFn(ctx, idHasil, statusLanjut)
	}
	return nil, nil
}

func setupRadiologiRouter(svc radiologi.Service) *http.ServeMux {
	mux := http.NewServeMux()
	handler := radiologi.NewHandler(svc, testEncryptionKey)
	passthroughMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return next
	}
	handler.RegisterRoutes(mux, passthroughMiddleware, passthroughMiddleware)
	return mux
}

func TestHandler_DaftarHasilRadiologi_Success(t *testing.T) {
	noRawat := "2026/04/22/000001"
	encKunjungan, _ := crypto.Encrypt(noRawat, testEncryptionKey)

	svc := &mockRadiologiService{
		getRiwayatRadiologiKunjunganFn: func(ctx context.Context, nr string, statusLanjut shared.StatusLanjut, filter radiologi.FilterRiwayatRadiologi) ([]radiologi.HasilRadiologi, shared.PaginationMeta, error) {
			if nr != noRawat {
				t.Errorf("Expected noRawat %s, got %s", noRawat, nr)
			}
			return []radiologi.HasilRadiologi{
				{
					NoRawat:        nr,
					KodeTindakan:   "RAD001",
					NamaTindakan:   "Rontgen Thorax",
					Status:         "Ralan",
					TanggalPeriksa: "2026-04-22",
					JamPeriksa:     "10:00:00",
					Hasil:          "Cor dan Pulmo normal",
					GambarPACS:     []string{"http://pacs.example.com/viewer?token=abc"},
				},
			}, shared.NewPaginationMeta(1, 1, 10), nil
		},
	}

	router := setupRadiologiRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/radiologi/"+encKunjungan+"/Ralan?page=1&limit=10", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp struct {
		Success bool                       `json:"success"`
		Message string                     `json:"message"`
		Data    []radiologi.HasilRadiologi `json:"data"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(resp.Data) != 1 {
		t.Fatalf("Expected 1 record, got %d", len(resp.Data))
	}
	if resp.Data[0].Id == "" || resp.Data[0].IdKunjungan == "" {
		t.Errorf("Expected encrypted Id and IdKunjungan, got id=%s, idKunjungan=%s", resp.Data[0].Id, resp.Data[0].IdKunjungan)
	}
}

func TestHandler_DaftarHasilRadiologi_InvalidToken(t *testing.T) {
	svc := &mockRadiologiService{}
	router := setupRadiologiRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/radiologi/invalid-token/Ralan", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_DaftarHasilRadiologi_InvalidStatusLanjut(t *testing.T) {
	encKunjungan, _ := crypto.Encrypt("2026/04/22/000001", testEncryptionKey)
	svc := &mockRadiologiService{}
	router := setupRadiologiRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/radiologi/"+encKunjungan+"/BukanStatus", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_DaftarHasilRadiologi_InvalidFilter(t *testing.T) {
	encKunjungan, _ := crypto.Encrypt("2026/04/22/000001", testEncryptionKey)
	svc := &mockRadiologiService{}
	router := setupRadiologiRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/radiologi/"+encKunjungan+"/Ralan?tanggal=bukan-tanggal", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_DaftarHasilRadiologi_ServiceError(t *testing.T) {
	encKunjungan, _ := crypto.Encrypt("2026/04/22/000001", testEncryptionKey)
	svc := &mockRadiologiService{
		getRiwayatRadiologiKunjunganFn: func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter radiologi.FilterRiwayatRadiologi) ([]radiologi.HasilRadiologi, shared.PaginationMeta, error) {
			return nil, shared.PaginationMeta{}, errors.New("db error")
		},
	}
	router := setupRadiologiRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/radiologi/"+encKunjungan+"/Ralan", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("Expected status 500, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_DaftarHasilRadiologiByRM_Success(t *testing.T) {
	noRM := "123456"
	encPasien, _ := crypto.Encrypt(noRM, testEncryptionKey)

	svc := &mockRadiologiService{
		getRiwayatRadiologiPasienFn: func(ctx context.Context, rm string, statusLanjut shared.StatusLanjut, filter radiologi.FilterRiwayatRadiologi) ([]radiologi.HasilRadiologi, shared.PaginationMeta, error) {
			if rm != noRM {
				t.Errorf("Expected noRM %s, got %s", noRM, rm)
			}
			return []radiologi.HasilRadiologi{
				{
					NoRawat:        "2026/04/22/000001",
					KodeTindakan:   "RAD001",
					NamaTindakan:   "Rontgen Thorax",
					Status:         "Ralan",
					TanggalPeriksa: "2026-04-22",
					JamPeriksa:     "10:00:00",
					Hasil:          "Normal",
				},
			}, shared.NewPaginationMeta(1, 1, 10), nil
		},
	}

	router := setupRadiologiRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/radiologi/pasien/"+encPasien+"/Semua", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_DaftarHasilRadiologiByRM_InvalidToken(t *testing.T) {
	svc := &mockRadiologiService{}
	router := setupRadiologiRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/radiologi/pasien/invalid-token/Semua", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_DetailHasilRadiologi_Success(t *testing.T) {
	noRawat := "2026/04/22/000001"
	kodeTindakan := "RAD001"
	tanggalPeriksa := "2026-04-22"
	jamPeriksa := "10:00:00"

	encKunjungan, _ := crypto.Encrypt(noRawat, testEncryptionKey)
	composite := noRawat + "~" + kodeTindakan + "~" + tanggalPeriksa + "~" + jamPeriksa
	encRadiologi, _ := crypto.Encrypt(composite, testEncryptionKey)

	svc := &mockRadiologiService{
		getDetailHasilRadiologiFn: func(ctx context.Context, idHasil radiologi.IdHasilRadiologi, statusLanjut shared.StatusLanjut) (*radiologi.HasilRadiologi, error) {
			return &radiologi.HasilRadiologi{
				NoRawat:        idHasil.NoRawat,
				KodeTindakan:   idHasil.KodeTindakan,
				NamaTindakan:   "Rontgen Thorax",
				Status:         string(statusLanjut),
				TanggalPeriksa: idHasil.TanggalPeriksa,
				JamPeriksa:     idHasil.JamPeriksa,
				Hasil:          "Cor dan Pulmo normal",
				GambarPACS:     []string{"http://pacs.example.com/view"},
			}, nil
		},
	}

	router := setupRadiologiRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/radiologi/"+encKunjungan+"/Ralan/"+encRadiologi, nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp struct {
		Success bool                     `json:"success"`
		Data    radiologi.HasilRadiologi `json:"data"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if resp.Data.Hasil != "Cor dan Pulmo normal" {
		t.Errorf("Expected hasil Cor dan Pulmo normal, got %s", resp.Data.Hasil)
	}
}

func TestHandler_DetailHasilRadiologi_MalformedCompositeID(t *testing.T) {
	encKunjungan, _ := crypto.Encrypt("2026/04/22/000001", testEncryptionKey)
	encRadiologi, _ := crypto.Encrypt("invalid~composite", testEncryptionKey)

	svc := &mockRadiologiService{}
	router := setupRadiologiRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/radiologi/"+encKunjungan+"/Ralan/"+encRadiologi, nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_DetailHasilRadiologi_MismatchedNoRawat(t *testing.T) {
	encKunjungan, _ := crypto.Encrypt("2026/04/22/000001", testEncryptionKey)
	composite := "2026/04/22/999999~RAD001~2026-04-22~10:00:00"
	encRadiologi, _ := crypto.Encrypt(composite, testEncryptionKey)

	svc := &mockRadiologiService{}
	router := setupRadiologiRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/radiologi/"+encKunjungan+"/Ralan/"+encRadiologi, nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_DetailHasilRadiologi_InvalidStatusLanjut(t *testing.T) {
	noRawat := "2026/04/22/000001"
	encKunjungan, _ := crypto.Encrypt(noRawat, testEncryptionKey)
	composite := noRawat + "~RAD001~2026-04-22~10:00:00"
	encRadiologi, _ := crypto.Encrypt(composite, testEncryptionKey)

	svc := &mockRadiologiService{}
	router := setupRadiologiRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/radiologi/"+encKunjungan+"/InvalidStatus/"+encRadiologi, nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400 for invalid status lanjut, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_DetailHasilRadiologi_NotFound(t *testing.T) {
	noRawat := "2026/04/22/000001"
	encKunjungan, _ := crypto.Encrypt(noRawat, testEncryptionKey)
	composite := noRawat + "~RAD001~2026-04-22~10:00:00"
	encRadiologi, _ := crypto.Encrypt(composite, testEncryptionKey)

	svc := &mockRadiologiService{
		getDetailHasilRadiologiFn: func(ctx context.Context, idHasil radiologi.IdHasilRadiologi, statusLanjut shared.StatusLanjut) (*radiologi.HasilRadiologi, error) {
			return nil, apperror.NewNotFoundError("Data hasil pemeriksaan radiologi tidak ditemukan")
		},
	}

	router := setupRadiologiRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/radiologi/"+encKunjungan+"/Ralan/"+encRadiologi, nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("Expected status 404, got %d: %s", rr.Code, rr.Body.String())
	}
}
