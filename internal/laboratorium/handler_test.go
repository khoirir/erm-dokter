package laboratorium_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"erm-dokter/internal/laboratorium"
	"erm-dokter/internal/middleware"
	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/pkg/token"
	"erm-dokter/internal/shared"
)

const testEncryptionKey = "bafaa956-751d-4f59-98cc-574ee9dfe9f6"

type mockService struct {
	getRiwayatLabKunjunganFn        func(ctx context.Context, kategori shared.KategoriLab, noRawat string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) (*laboratorium.HasilLaboratoriumKunjungan, shared.PaginationMeta, error)
	getRiwayatLabPasienFn           func(ctx context.Context, kategori shared.KategoriLab, noRM string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratorium, shared.PaginationMeta, error)
	getDetailHasilLabFn             func(ctx context.Context, kategori shared.KategoriLab, noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa string) (*laboratorium.HasilLaboratorium, error)
	simpanPermintaanLabPKFn          func(ctx context.Context, kodeDokterLogin string, statusLanjut shared.StatusLanjut, req laboratorium.SimpanPermintaanLabPKRequest) (*laboratorium.DetailPermintaanLabPK, error)
	getDaftarPermintaanLabPKFn       func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) ([]laboratorium.PermintaanLabPK, error)
	getRiwayatPermintaanLabPKByRMFn func(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.PermintaanLabPK, shared.PaginationMeta, error)
	getDetailPermintaanLabPKFn       func(ctx context.Context, noRawat string, noPermintaan string, statusLanjut shared.StatusLanjut) (*laboratorium.DetailPermintaanLabPK, error)
	hapusPermintaanLabPKFn           func(ctx context.Context, noRawat string, noPermintaan string, statusLanjut shared.StatusLanjut, kodeDokterLogin string) error
}

func (m *mockService) GetRiwayatLabKunjungan(ctx context.Context, kategori shared.KategoriLab, noRawat string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) (*laboratorium.HasilLaboratoriumKunjungan, shared.PaginationMeta, error) {
	if m.getRiwayatLabKunjunganFn != nil {
		return m.getRiwayatLabKunjunganFn(ctx, kategori, noRawat, statusLanjut, filter)
	}
	return nil, shared.PaginationMeta{}, nil
}

func (m *mockService) GetRiwayatLabPasien(ctx context.Context, kategori shared.KategoriLab, noRM string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratorium, shared.PaginationMeta, error) {
	if m.getRiwayatLabPasienFn != nil {
		return m.getRiwayatLabPasienFn(ctx, kategori, noRM, statusLanjut, filter)
	}
	return nil, shared.PaginationMeta{}, nil
}

func (m *mockService) GetDetailHasilLab(ctx context.Context, kategori shared.KategoriLab, noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa string) (*laboratorium.HasilLaboratorium, error) {
	if m.getDetailHasilLabFn != nil {
		return m.getDetailHasilLabFn(ctx, kategori, noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa)
	}
	return nil, nil
}

func (m *mockService) SimpanPermintaanLabPK(ctx context.Context, kodeDokterLogin string, statusLanjut shared.StatusLanjut, req laboratorium.SimpanPermintaanLabPKRequest) (*laboratorium.DetailPermintaanLabPK, error) {
	if m.simpanPermintaanLabPKFn != nil {
		return m.simpanPermintaanLabPKFn(ctx, kodeDokterLogin, statusLanjut, req)
	}
	return nil, nil
}

func (m *mockService) GetDaftarPermintaanLabPK(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) ([]laboratorium.PermintaanLabPK, error) {
	if m.getDaftarPermintaanLabPKFn != nil {
		return m.getDaftarPermintaanLabPKFn(ctx, noRawat, statusLanjut)
	}
	return nil, nil
}

func (m *mockService) GetRiwayatPermintaanLabPKByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.PermintaanLabPK, shared.PaginationMeta, error) {
	if m.getRiwayatPermintaanLabPKByRMFn != nil {
		return m.getRiwayatPermintaanLabPKByRMFn(ctx, noRM, statusLanjut, filter)
	}
	return nil, shared.PaginationMeta{}, nil
}

func (m *mockService) GetDetailPermintaanLabPK(ctx context.Context, noRawat string, noPermintaan string, statusLanjut shared.StatusLanjut) (*laboratorium.DetailPermintaanLabPK, error) {
	if m.getDetailPermintaanLabPKFn != nil {
		return m.getDetailPermintaanLabPKFn(ctx, noRawat, noPermintaan, statusLanjut)
	}
	return nil, nil
}

func (m *mockService) HapusPermintaanLabPK(ctx context.Context, noRawat string, noPermintaan string, statusLanjut shared.StatusLanjut, kodeDokterLogin string) error {
	if m.hapusPermintaanLabPKFn != nil {
		return m.hapusPermintaanLabPKFn(ctx, noRawat, noPermintaan, statusLanjut, kodeDokterLogin)
	}
	return nil
}

func authMiddlewareForTest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), middleware.UserClaimKey, &token.Claims{
			KodeDokter: "DR01",
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func TestHandler_DaftarHasilLab_Success(t *testing.T) {
	encKunjungan, _ := crypto.Encrypt("2026/04/22/000001", testEncryptionKey)
	mockSvc := &mockService{
		getRiwayatLabKunjunganFn: func(ctx context.Context, kategori shared.KategoriLab, noRawat string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) (*laboratorium.HasilLaboratoriumKunjungan, shared.PaginationMeta, error) {
			if noRawat != "2026/04/22/000001" {
				t.Errorf("Expected noRawat 2026/04/22/000001, got %s", noRawat)
			}
			return &laboratorium.HasilLaboratoriumKunjungan{
				HasilPemeriksaan: []laboratorium.HasilLaboratorium{
					{
						NoRawat:      "2026/04/22/000001",
						KodeTindakan: "LAB001",
						NamaTindakan: "DARAH LENGKAP",
						Kategori:     "PK",
					},
				},
				BerkasDigital: []laboratorium.BerkasDigital{},
			}, shared.NewPaginationMeta(1, 1, 5), nil
		},
	}

	handler := laboratorium.NewHandler(mockSvc, testEncryptionKey)
	mux := http.NewServeMux()
	dummyMiddleware := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, dummyMiddleware, dummyMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/laboratorium/pk/"+encKunjungan+"/Semua?page=1&limit=5", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp struct {
		Code    int                                     `json:"code"`
		Status  string                                  `json:"status"`
		Message string                                  `json:"message"`
		Data    laboratorium.HasilLaboratoriumKunjungan `json:"data"`
		Meta    shared.PaginationMeta                   `json:"meta"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(resp.Data.HasilPemeriksaan) != 1 || resp.Data.HasilPemeriksaan[0].NamaTindakan != "DARAH LENGKAP" {
		t.Errorf("Unexpected response data: %+v", resp.Data)
	}
	if resp.Data.HasilPemeriksaan[0].Id == "" {
		t.Errorf("Expected encrypted Id in response")
	}
}

func TestHandler_DaftarHasilLab_InvalidKategori(t *testing.T) {
	encKunjungan, _ := crypto.Encrypt("2026/04/22/000001", testEncryptionKey)
	handler := laboratorium.NewHandler(&mockService{}, testEncryptionKey)
	mux := http.NewServeMux()
	dummyMiddleware := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, dummyMiddleware, dummyMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/laboratorium/radiologi/"+encKunjungan+"/Semua", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400 for invalid kategori, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_DaftarHasilLab_InvalidStatusLanjut(t *testing.T) {
	encKunjungan, _ := crypto.Encrypt("2026/04/22/000001", testEncryptionKey)
	handler := laboratorium.NewHandler(&mockService{}, testEncryptionKey)
	mux := http.NewServeMux()
	dummyMiddleware := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, dummyMiddleware, dummyMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/laboratorium/pk/"+encKunjungan+"/InvalidStatus", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400 for invalid status lanjut, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_DaftarHasilLabByRM_Success(t *testing.T) {
	encPasien, _ := crypto.Encrypt("123456", testEncryptionKey)
	mockSvc := &mockService{
		getRiwayatLabPasienFn: func(ctx context.Context, kategori shared.KategoriLab, noRM string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratorium, shared.PaginationMeta, error) {
			if noRM != "123456" {
				t.Errorf("Expected noRM 123456, got %s", noRM)
			}
			return []laboratorium.HasilLaboratorium{
				{
					NoRawat:      "2026/04/22/000001",
					KodeTindakan: "PA001",
					NamaTindakan: "BIOPSI",
					Kategori:     "PA",
				},
			}, shared.NewPaginationMeta(1, 1, 5), nil
		},
	}

	handler := laboratorium.NewHandler(mockSvc, testEncryptionKey)
	mux := http.NewServeMux()
	dummyMiddleware := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, dummyMiddleware, dummyMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/laboratorium/pa/pasien/"+encPasien+"/Semua", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_DetailHasilLab_Success(t *testing.T) {
	encKunjungan, _ := crypto.Encrypt("2026/04/22/000001", testEncryptionKey)
	encHasil, _ := crypto.Encrypt("2026/04/22/000001~LAB001~2026-04-22~10:00:00", testEncryptionKey)

	mockSvc := &mockService{
		getDetailHasilLabFn: func(ctx context.Context, kategori shared.KategoriLab, noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa string) (*laboratorium.HasilLaboratorium, error) {
			return &laboratorium.HasilLaboratorium{
				NoRawat:        noRawat,
				KodeTindakan:   kodeTindakan,
				TanggalPeriksa: tanggalPeriksa,
				JamPeriksa:     jamPeriksa,
				NamaTindakan:   "DARAH LENGKAP",
				Kategori:       "PK",
			}, nil
		},
	}

	handler := laboratorium.NewHandler(mockSvc, testEncryptionKey)
	mux := http.NewServeMux()
	dummyMiddleware := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, dummyMiddleware, dummyMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/laboratorium/pk/"+encKunjungan+"/Semua/"+encHasil, nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_SimpanPermintaanLabPK_Success(t *testing.T) {
	encKunjungan, _ := crypto.Encrypt("2026/09/03/000001", testEncryptionKey)
	encTindakan, _ := crypto.Encrypt("TND001", testEncryptionKey)
	encTemplate, _ := crypto.Encrypt("101", testEncryptionKey)

	mockSvc := &mockService{
		simpanPermintaanLabPKFn: func(ctx context.Context, kodeDokterLogin string, statusLanjut shared.StatusLanjut, req laboratorium.SimpanPermintaanLabPKRequest) (*laboratorium.DetailPermintaanLabPK, error) {
			if req.NoRawat != "2026/09/03/000001" {
				t.Errorf("Expected NoRawat 2026/09/03/000001, got %s", req.NoRawat)
			}
			if statusLanjut != shared.StatusLanjutRawatJalan {
				t.Errorf("Expected StatusLanjut Ralan, got %s", statusLanjut)
			}
			if len(req.Pemeriksaan) != 1 || req.Pemeriksaan[0].KodeTindakan != "TND001" {
				t.Errorf("Expected decrypted KodeTindakan TND001, got %+v", req.Pemeriksaan)
			}
			return &laboratorium.DetailPermintaanLabPK{
				PermintaanLabPK: laboratorium.PermintaanLabPK{
					NoPermintaan: "PK202609030001",
					NoRawat:      req.NoRawat,
					Status:       "ralan",
				},
				Pemeriksaan: []laboratorium.PemeriksaanLabPKItem{},
			}, nil
		},
	}

	handler := laboratorium.NewHandler(mockSvc, testEncryptionKey)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, authMiddlewareForTest, func(next http.HandlerFunc) http.HandlerFunc { return next })

	body, _ := json.Marshal(laboratorium.SimpanPermintaanLabPKRequest{
		NoRawat:           "2026/09/03/000001",
		TanggalPermintaan: "2026-09-03",
		JamPermintaan:     "10:00:00",
		DiagnosaKlinis:    "Febris H-3",
		Pemeriksaan: []laboratorium.ItemPemeriksaanLabPKRequest{
			{
				IdTindakan: encTindakan,
				IdTemplate: []string{encTemplate},
			},
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/laboratorium/pk/permintaan/"+encKunjungan+"/Ralan", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Expected status 201 Created, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_SimpanPermintaanLabPK_MismatchNoRawat(t *testing.T) {
	encKunjungan, _ := crypto.Encrypt("2026/09/03/000001", testEncryptionKey)

	handler := laboratorium.NewHandler(&mockService{}, testEncryptionKey)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, authMiddlewareForTest, func(next http.HandlerFunc) http.HandlerFunc { return next })

	body, _ := json.Marshal(laboratorium.SimpanPermintaanLabPKRequest{
		NoRawat:           "2026/09/03/999999",
		TanggalPermintaan: "2026-09-03",
		JamPermintaan:     "10:00:00",
		DiagnosaKlinis:    "Febris H-3",
		Pemeriksaan: []laboratorium.ItemPemeriksaanLabPKRequest{
			{IdTindakan: "enc-tindakan-1"},
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/laboratorium/pk/permintaan/"+encKunjungan+"/Ralan", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400 Bad Request for mismatched NoRawat, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_DaftarPermintaanLabPK_Success(t *testing.T) {
	encKunjungan, _ := crypto.Encrypt("2026/09/03/000001", testEncryptionKey)
	mockSvc := &mockService{
		getDaftarPermintaanLabPKFn: func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) ([]laboratorium.PermintaanLabPK, error) {
			return []laboratorium.PermintaanLabPK{
				{
					NoPermintaan: "PK202609030001",
					NoRawat:      noRawat,
					Status:       "ralan",
					StatusProses: "Menunggu Sampel",
				},
			}, nil
		},
	}

	handler := laboratorium.NewHandler(mockSvc, testEncryptionKey)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, authMiddlewareForTest, func(next http.HandlerFunc) http.HandlerFunc { return next })

	req := httptest.NewRequest(http.MethodGet, "/api/v1/laboratorium/pk/permintaan/"+encKunjungan+"/Semua", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_DaftarPermintaanLabPKByRM_Success(t *testing.T) {
	encPasien, _ := crypto.Encrypt("123456", testEncryptionKey)
	mockSvc := &mockService{
		getRiwayatPermintaanLabPKByRMFn: func(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.PermintaanLabPK, shared.PaginationMeta, error) {
			return []laboratorium.PermintaanLabPK{
				{
					NoPermintaan: "PK202609030001",
					NoRawat:      "2026/09/03/000001",
					Status:       "ralan",
				},
			}, shared.NewPaginationMeta(1, 1, 5), nil
		},
	}

	handler := laboratorium.NewHandler(mockSvc, testEncryptionKey)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, authMiddlewareForTest, func(next http.HandlerFunc) http.HandlerFunc { return next })

	req := httptest.NewRequest(http.MethodGet, "/api/v1/laboratorium/pk/permintaan/pasien/"+encPasien+"/Semua", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_DetailPermintaanLabPK_Success(t *testing.T) {
	encKunjungan, _ := crypto.Encrypt("2026/09/03/000001", testEncryptionKey)
	encOrder, _ := crypto.Encrypt("PK202609030001", testEncryptionKey)

	mockSvc := &mockService{
		getDetailPermintaanLabPKFn: func(ctx context.Context, noRawat string, noPermintaan string, statusLanjut shared.StatusLanjut) (*laboratorium.DetailPermintaanLabPK, error) {
			return &laboratorium.DetailPermintaanLabPK{
				PermintaanLabPK: laboratorium.PermintaanLabPK{
					NoPermintaan: noPermintaan,
					NoRawat:      noRawat,
					Status:       "ralan",
				},
				Pemeriksaan: []laboratorium.PemeriksaanLabPKItem{},
			}, nil
		},
	}

	handler := laboratorium.NewHandler(mockSvc, testEncryptionKey)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, authMiddlewareForTest, func(next http.HandlerFunc) http.HandlerFunc { return next })

	req := httptest.NewRequest(http.MethodGet, "/api/v1/laboratorium/pk/permintaan/"+encKunjungan+"/Semua/"+encOrder, nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_HapusPermintaanLabPK_Success(t *testing.T) {
	encKunjungan, _ := crypto.Encrypt("2026/09/03/000001", testEncryptionKey)
	encOrder, _ := crypto.Encrypt("PK202609030001", testEncryptionKey)

	mockSvc := &mockService{
		hapusPermintaanLabPKFn: func(ctx context.Context, noRawat string, noPermintaan string, statusLanjut shared.StatusLanjut, kodeDokterLogin string) error {
			return nil
		},
	}

	handler := laboratorium.NewHandler(mockSvc, testEncryptionKey)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, authMiddlewareForTest, func(next http.HandlerFunc) http.HandlerFunc { return next })

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/laboratorium/pk/permintaan/"+encKunjungan+"/Semua/"+encOrder, nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}
}
