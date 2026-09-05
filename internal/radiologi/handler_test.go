package radiologi_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"erm-dokter/internal/middleware"
	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/pkg/token"
	"erm-dokter/internal/radiologi"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

const testEncryptionKey = "bafaa956-751d-4f59-98cc-574ee9dfe9f6"

type mockRadiologiService struct {
	getRiwayatRadiologiKunjunganFn func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter radiologi.FilterRiwayatRadiologi) ([]radiologi.HasilRadiologi, shared.PaginationMeta, error)
	getRiwayatRadiologiPasienFn    func(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter radiologi.FilterRiwayatRadiologi) ([]radiologi.HasilRadiologi, shared.PaginationMeta, error)
	getDetailHasilRadiologiFn      func(ctx context.Context, idHasil radiologi.IdHasilRadiologi, statusLanjut shared.StatusLanjut) (*radiologi.HasilRadiologi, error)

	simpanPermintaanRadiologiFn       func(ctx context.Context, kodeDokterLogin string, statusLanjut shared.StatusLanjut, req radiologi.SimpanPermintaanRadiologiRequest) (*radiologi.DetailPermintaanRadiologi, error)
	getDaftarPermintaanRadiologiFn    func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) ([]radiologi.DetailPermintaanRadiologi, error)
	getRiwayatPermintaanRadiologiByRMFn func(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter radiologi.FilterRiwayatPermintaanRadiologi) ([]radiologi.DetailPermintaanRadiologi, shared.PaginationMeta, error)
	getDetailPermintaanRadiologiFn    func(ctx context.Context, noRawat string, noPermintaan string, statusLanjut shared.StatusLanjut) (*radiologi.DetailPermintaanRadiologi, error)
	updatePermintaanRadiologiFn       func(ctx context.Context, kodeDokterLogin, noRawat, noPermintaan string, statusLanjut shared.StatusLanjut, req radiologi.SimpanPermintaanRadiologiRequest) (*radiologi.DetailPermintaanRadiologi, error)
	hapusPermintaanRadiologiFn        func(ctx context.Context, noRawat string, noPermintaan string, statusLanjut shared.StatusLanjut, kodeDokterLogin string) error
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

func (m *mockRadiologiService) SimpanPermintaanRadiologi(ctx context.Context, kodeDokterLogin string, statusLanjut shared.StatusLanjut, req radiologi.SimpanPermintaanRadiologiRequest) (*radiologi.DetailPermintaanRadiologi, error) {
	if m.simpanPermintaanRadiologiFn != nil {
		return m.simpanPermintaanRadiologiFn(ctx, kodeDokterLogin, statusLanjut, req)
	}
	return nil, nil
}

func (m *mockRadiologiService) GetDaftarPermintaanRadiologi(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) ([]radiologi.DetailPermintaanRadiologi, error) {
	if m.getDaftarPermintaanRadiologiFn != nil {
		return m.getDaftarPermintaanRadiologiFn(ctx, noRawat, statusLanjut)
	}
	return nil, nil
}

func (m *mockRadiologiService) GetRiwayatPermintaanRadiologiByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter radiologi.FilterRiwayatPermintaanRadiologi) ([]radiologi.DetailPermintaanRadiologi, shared.PaginationMeta, error) {
	if m.getRiwayatPermintaanRadiologiByRMFn != nil {
		return m.getRiwayatPermintaanRadiologiByRMFn(ctx, noRM, statusLanjut, filter)
	}
	return nil, shared.PaginationMeta{}, nil
}

func (m *mockRadiologiService) GetDetailPermintaanRadiologi(ctx context.Context, noRawat string, noPermintaan string, statusLanjut shared.StatusLanjut) (*radiologi.DetailPermintaanRadiologi, error) {
	if m.getDetailPermintaanRadiologiFn != nil {
		return m.getDetailPermintaanRadiologiFn(ctx, noRawat, noPermintaan, statusLanjut)
	}
	return nil, nil
}

func (m *mockRadiologiService) UpdatePermintaanRadiologi(ctx context.Context, kodeDokterLogin, noRawat, noPermintaan string, statusLanjut shared.StatusLanjut, req radiologi.SimpanPermintaanRadiologiRequest) (*radiologi.DetailPermintaanRadiologi, error) {
	if m.updatePermintaanRadiologiFn != nil {
		return m.updatePermintaanRadiologiFn(ctx, kodeDokterLogin, noRawat, noPermintaan, statusLanjut, req)
	}
	return nil, nil
}

func (m *mockRadiologiService) HapusPermintaanRadiologi(ctx context.Context, noRawat string, noPermintaan string, statusLanjut shared.StatusLanjut, kodeDokterLogin string) error {
	if m.hapusPermintaanRadiologiFn != nil {
		return m.hapusPermintaanRadiologiFn(ctx, noRawat, noPermintaan, statusLanjut, kodeDokterLogin)
	}
	return nil
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

	testCases := []struct {
		name         string
		statusLanjut string
	}{
		{name: "RejectSemua", statusLanjut: "Semua"},
		{name: "RejectInvalidStatus", statusLanjut: "InvalidStatus"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/radiologi/"+encKunjungan+"/"+tc.statusLanjut+"/"+encRadiologi, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Fatalf("Expected status 400 for status lanjut %s, got %d: %s", tc.statusLanjut, rr.Code, rr.Body.String())
			}
		})
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

func authMiddlewareForTest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), middleware.UserClaimKey, &token.Claims{
			KodeDokter: "DR01",
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func setupRadiologiRouterWithAuth(svc radiologi.Service) *http.ServeMux {
	mux := http.NewServeMux()
	handler := radiologi.NewHandler(svc, testEncryptionKey)
	dummyTimeout := func(next http.HandlerFunc) http.HandlerFunc {
		return next
	}
	handler.RegisterRoutes(mux, authMiddlewareForTest, dummyTimeout)
	return mux
}

func TestHandler_SimpanPermintaanRadiologi_Success(t *testing.T) {
	noRawat := "2026/09/05/000001"
	encKunjungan, _ := crypto.Encrypt(noRawat, testEncryptionKey)
	encTindakan, _ := crypto.Encrypt("RAD001", testEncryptionKey)

	svc := &mockRadiologiService{
		simpanPermintaanRadiologiFn: func(ctx context.Context, kodeDokterLogin string, statusLanjut shared.StatusLanjut, req radiologi.SimpanPermintaanRadiologiRequest) (*radiologi.DetailPermintaanRadiologi, error) {
			if kodeDokterLogin != "DR01" {
				t.Errorf("Expected kodeDokter DR01, got %s", kodeDokterLogin)
			}
			if req.NoRawat != noRawat {
				t.Errorf("Expected noRawat %s, got %s", noRawat, req.NoRawat)
			}
			if len(req.Pemeriksaan) != 1 || req.Pemeriksaan[0].KodeTindakan != "RAD001" {
				t.Errorf("Expected decrypted KodeTindakan RAD001, got %+v", req.Pemeriksaan)
			}
			return &radiologi.DetailPermintaanRadiologi{
				PermintaanRadiologiHeader: radiologi.PermintaanRadiologiHeader{
					NoPermintaan:      "RAD202609050001",
					NoRawat:           noRawat,
					TanggalPermintaan: req.TanggalPermintaan,
					JamPermintaan:     req.JamPermintaan,
					DokterPerujuk:     radiologi.DokterInfo{KodeDokter: "DR01", NamaDokter: "dr. SpRad"},
					Status:            "Ralan",
					InformasiTambahan: req.InformasiTambahan,
					DiagnosaKlinis:    req.DiagnosaKlinis,
				},
				Pemeriksaan: []radiologi.PemeriksaanRadiologiItem{
					{KodeTindakan: "RAD001", NamaTindakan: "Thorax AP", StatusBayar: "Belum"},
				},
			}, nil
		},
	}

	router := setupRadiologiRouterWithAuth(svc)

	body := `{
		"no_rawat": "` + noRawat + `",
		"tanggal_permintaan": "2026-09-05",
		"jam_permintaan": "10:00:00",
		"informasi_tambahan": "Thorax AP",
		"diagnosa_klinis": "Batuk kronis",
		"pemeriksaan": [
			{"id_tindakan": "` + encTindakan + `"}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/radiologi/permintaan/"+encKunjungan+"/Ralan", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Expected status 201 Created, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp struct {
		Success bool                                `json:"success"`
		Message string                              `json:"message"`
		Data    radiologi.DetailPermintaanRadiologi `json:"data"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if resp.Data.NoPermintaan == "RAD202609050001" {
		t.Errorf("Expected encrypted NoPermintaan, got plain: %s", resp.Data.NoPermintaan)
	}
}

func TestHandler_SimpanPermintaanRadiologi_MismatchedNoRawat(t *testing.T) {
	encKunjungan, _ := crypto.Encrypt("2026/09/05/000001", testEncryptionKey)
	encTindakan, _ := crypto.Encrypt("RAD001", testEncryptionKey)

	svc := &mockRadiologiService{}
	router := setupRadiologiRouterWithAuth(svc)

	body := `{
		"no_rawat": "2026/09/05/999999",
		"tanggal_permintaan": "2026-09-05",
		"jam_permintaan": "10:00:00",
		"pemeriksaan": [{"id_tindakan": "` + encTindakan + `"}]
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/radiologi/permintaan/"+encKunjungan+"/Ralan", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400 Bad Request, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_SimpanPermintaanRadiologi_RejectSemua(t *testing.T) {
	encKunjungan, _ := crypto.Encrypt("2026/09/05/000001", testEncryptionKey)
	svc := &mockRadiologiService{}
	router := setupRadiologiRouterWithAuth(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/radiologi/permintaan/"+encKunjungan+"/Semua", strings.NewReader("{}"))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400 Bad Request for status lanjut 'Semua' on create, got %d", rr.Code)
	}
}

func TestHandler_DaftarPermintaanRadiologi_Success(t *testing.T) {
	noRawat := "2026/09/05/000001"
	encKunjungan, _ := crypto.Encrypt(noRawat, testEncryptionKey)

	svc := &mockRadiologiService{
		getDaftarPermintaanRadiologiFn: func(ctx context.Context, nr string, statusLanjut shared.StatusLanjut) ([]radiologi.DetailPermintaanRadiologi, error) {
			return []radiologi.DetailPermintaanRadiologi{
				{
					PermintaanRadiologiHeader: radiologi.PermintaanRadiologiHeader{
						NoPermintaan:      "RAD202609050001",
						NoRawat:           nr,
						TanggalPermintaan: "2026-09-05",
						Status:            "Ralan",
					},
					Pemeriksaan: []radiologi.PemeriksaanRadiologiItem{
						{KodeTindakan: "RAD001", NamaTindakan: "Thorax AP"},
					},
				},
			}, nil
		},
	}

	router := setupRadiologiRouterWithAuth(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/radiologi/permintaan/"+encKunjungan+"/Semua", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_DaftarPermintaanRadiologiByRM_Success(t *testing.T) {
	noRM := "123456"
	encPasien, _ := crypto.Encrypt(noRM, testEncryptionKey)

	svc := &mockRadiologiService{
		getRiwayatPermintaanRadiologiByRMFn: func(ctx context.Context, rm string, statusLanjut shared.StatusLanjut, filter radiologi.FilterRiwayatPermintaanRadiologi) ([]radiologi.DetailPermintaanRadiologi, shared.PaginationMeta, error) {
			return []radiologi.DetailPermintaanRadiologi{
				{
					PermintaanRadiologiHeader: radiologi.PermintaanRadiologiHeader{
						NoPermintaan:      "RAD202609050001",
						TanggalPermintaan: "2026-09-05",
						Status:            "Ralan",
					},
				},
			}, shared.NewPaginationMeta(1, 1, 5), nil
		},
	}

	router := setupRadiologiRouterWithAuth(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/radiologi/pasien/"+encPasien+"/permintaan/Semua?page=1&limit=5", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_DetailPermintaanRadiologi_Success(t *testing.T) {
	noRawat := "2026/09/05/000001"
	noOrder := "RAD202609050001"
	encKunjungan, _ := crypto.Encrypt(noRawat, testEncryptionKey)
	encPermintaan, _ := crypto.Encrypt(noOrder, testEncryptionKey)

	svc := &mockRadiologiService{
		getDetailPermintaanRadiologiFn: func(ctx context.Context, nr, np string, statusLanjut shared.StatusLanjut) (*radiologi.DetailPermintaanRadiologi, error) {
			return &radiologi.DetailPermintaanRadiologi{
				PermintaanRadiologiHeader: radiologi.PermintaanRadiologiHeader{
					NoPermintaan:      np,
					NoRawat:           nr,
					TanggalPermintaan: "2026-09-05",
					Status:            "Ralan",
				},
			}, nil
		},
	}

	router := setupRadiologiRouterWithAuth(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/radiologi/permintaan/"+encKunjungan+"/Ralan/"+encPermintaan, nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_DetailPermintaanRadiologi_RejectSemua(t *testing.T) {
	encKunjungan, _ := crypto.Encrypt("2026/09/05/000001", testEncryptionKey)
	encPermintaan, _ := crypto.Encrypt("RAD202609050001", testEncryptionKey)

	svc := &mockRadiologiService{}
	router := setupRadiologiRouterWithAuth(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/radiologi/permintaan/"+encKunjungan+"/Semua/"+encPermintaan, nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400 Bad Request for status lanjut 'Semua' on detail, got %d", rr.Code)
	}
}

func TestHandler_UpdatePermintaanRadiologi_RejectSemua(t *testing.T) {
	encKunjungan, _ := crypto.Encrypt("2026/09/05/000001", testEncryptionKey)
	encPermintaan, _ := crypto.Encrypt("RAD202609050001", testEncryptionKey)

	svc := &mockRadiologiService{}
	router := setupRadiologiRouterWithAuth(svc)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/radiologi/permintaan/"+encKunjungan+"/Semua/"+encPermintaan, strings.NewReader("{}"))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400 Bad Request for status lanjut 'Semua' on update, got %d", rr.Code)
	}
}

func TestHandler_HapusPermintaanRadiologi_RejectSemua(t *testing.T) {
	encKunjungan, _ := crypto.Encrypt("2026/09/05/000001", testEncryptionKey)
	encPermintaan, _ := crypto.Encrypt("RAD202609050001", testEncryptionKey)

	svc := &mockRadiologiService{}
	router := setupRadiologiRouterWithAuth(svc)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/radiologi/permintaan/"+encKunjungan+"/Semua/"+encPermintaan, nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400 Bad Request for status lanjut 'Semua' on delete, got %d", rr.Code)
	}
}

func TestHandler_HapusPermintaanRadiologi_Success(t *testing.T) {
	noRawat := "2026/09/05/000001"
	noOrder := "RAD202609050001"
	encKunjungan, _ := crypto.Encrypt(noRawat, testEncryptionKey)
	encPermintaan, _ := crypto.Encrypt(noOrder, testEncryptionKey)

	svc := &mockRadiologiService{
		hapusPermintaanRadiologiFn: func(ctx context.Context, nr, np string, statusLanjut shared.StatusLanjut, kodeDokterLogin string) error {
			if kodeDokterLogin != "DR01" {
				t.Errorf("Expected DR01, got %s", kodeDokterLogin)
			}
			return nil
		},
	}

	router := setupRadiologiRouterWithAuth(svc)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/radiologi/permintaan/"+encKunjungan+"/Ralan/"+encPermintaan, nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d: %s", rr.Code, rr.Body.String())
	}
}

