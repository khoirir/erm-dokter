package resumepasien_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"erm-dokter/internal/middleware"
	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/pkg/token"
	"erm-dokter/internal/resumepasien"
	"erm-dokter/internal/shared"
)

const testEncKey = "bafaa956-751d-4f59-98cc-574ee9dfe9f6"

type mockResumePasienService struct {
	detailFn  func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) (*resumepasien.ResumePasien, error)
	riwayatFn func(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut) ([]resumepasien.ResumePasien, error)
	simpanFn  func(ctx context.Context, noRawat, kodeDokter string, statusLanjut shared.StatusLanjut, req resumepasien.SimpanResumePasienRequest) (*resumepasien.ResumePasien, error)
	updateFn  func(ctx context.Context, kodeDokterLogin, noRawat string, statusLanjut shared.StatusLanjut, req resumepasien.UpdateResumePasienRequest) (*resumepasien.ResumePasien, error)
	hapusFn   func(ctx context.Context, kodeDokterLogin, noRawat string, statusLanjut shared.StatusLanjut) error
}

func (m *mockResumePasienService) DetailResumePasien(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) (*resumepasien.ResumePasien, error) {
	if m.detailFn != nil {
		return m.detailFn(ctx, noRawat, statusLanjut)
	}
	return nil, nil
}

func (m *mockResumePasienService) RiwayatResumePasienByNoRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut) ([]resumepasien.ResumePasien, error) {
	if m.riwayatFn != nil {
		return m.riwayatFn(ctx, noRM, statusLanjut)
	}
	return nil, nil
}

func (m *mockResumePasienService) SimpanResumePasien(ctx context.Context, noRawat, kodeDokter string, statusLanjut shared.StatusLanjut, req resumepasien.SimpanResumePasienRequest) (*resumepasien.ResumePasien, error) {
	if m.simpanFn != nil {
		return m.simpanFn(ctx, noRawat, kodeDokter, statusLanjut, req)
	}
	return nil, nil
}

func (m *mockResumePasienService) UpdateResumePasien(ctx context.Context, kodeDokterLogin, noRawat string, statusLanjut shared.StatusLanjut, req resumepasien.UpdateResumePasienRequest) (*resumepasien.ResumePasien, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, kodeDokterLogin, noRawat, statusLanjut, req)
	}
	return nil, nil
}

func (m *mockResumePasienService) HapusResumePasien(ctx context.Context, kodeDokterLogin, noRawat string, statusLanjut shared.StatusLanjut) error {
	if m.hapusFn != nil {
		return m.hapusFn(ctx, kodeDokterLogin, noRawat, statusLanjut)
	}
	return nil
}

func dummyAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := &token.Claims{
			KodeDokter: "DR01",
			NamaUser:   "dr. Handi",
		}
		ctx := context.WithValue(r.Context(), middleware.UserClaimKey, claims)
		next(w, r.WithContext(ctx))
	}
}

func dummyTimeoutMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
	}
}

func TestResumePasienHandler(t *testing.T) {
	rawNoRawat := "2026/09/07/000001"
	encNoRawat, _ := crypto.Encrypt(rawNoRawat, testEncKey)
	rawNoRM := "123456"
	encNoRM, _ := crypto.Encrypt(rawNoRM, testEncKey)

	t.Run("GET DetailResumePasien Sukses", func(t *testing.T) {
		mockSvc := &mockResumePasienService{
			detailFn: func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) (*resumepasien.ResumePasien, error) {
				return &resumepasien.ResumePasien{
					NoRawat:    noRawat,
					KodeDokter: "DR01",
					NamaDokter: "dr. Handi",
					DataResumePasien: resumepasien.DataResumePasien{
						KeluhanUtama:  "Batuk",
						DiagnosaUtama: "ISPA",
						KondisiPulang: resumepasien.KondisiPulangHidup,
					},
				}, nil
			},
		}

		handler := resumepasien.NewHandler(mockSvc, testEncKey)
		mux := http.NewServeMux()
		handler.RegisterRoutes(mux, dummyAuthMiddleware, dummyTimeoutMiddleware)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/resume/"+encNoRawat+"/Ralan", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("GET DetailResumePasien Status Ranap Ditolak (400)", func(t *testing.T) {
		handler := resumepasien.NewHandler(&mockResumePasienService{}, testEncKey)
		mux := http.NewServeMux()
		handler.RegisterRoutes(mux, dummyAuthMiddleware, dummyTimeoutMiddleware)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/resume/"+encNoRawat+"/Ranap", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400 for non-ralan status, got %d", w.Code)
		}
	})

	t.Run("GET RiwayatResumePasienByNoRM Sukses", func(t *testing.T) {
		mockSvc := &mockResumePasienService{
			riwayatFn: func(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut) ([]resumepasien.ResumePasien, error) {
				return []resumepasien.ResumePasien{
					{
						NoRawat: rawNoRawat,
						DataResumePasien: resumepasien.DataResumePasien{
							KeluhanUtama:  "Batuk",
							DiagnosaUtama: "ISPA",
						},
					},
				}, nil
			},
		}

		handler := resumepasien.NewHandler(mockSvc, testEncKey)
		mux := http.NewServeMux()
		handler.RegisterRoutes(mux, dummyAuthMiddleware, dummyTimeoutMiddleware)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/resume/pasien/"+encNoRM+"/Ralan", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("POST SimpanResumePasien 201 Created", func(t *testing.T) {
		mockSvc := &mockResumePasienService{
			simpanFn: func(ctx context.Context, noRawat, kodeDokter string, statusLanjut shared.StatusLanjut, req resumepasien.SimpanResumePasienRequest) (*resumepasien.ResumePasien, error) {
				return &resumepasien.ResumePasien{
					NoRawat:          noRawat,
					KodeDokter:       kodeDokter,
					DataResumePasien: req.DataResumePasien,
				}, nil
			},
		}

		handler := resumepasien.NewHandler(mockSvc, testEncKey)
		mux := http.NewServeMux()
		handler.RegisterRoutes(mux, dummyAuthMiddleware, dummyTimeoutMiddleware)

		payload := resumepasien.SimpanResumePasienRequest{
			NoRawat: rawNoRawat,
			DataResumePasien: resumepasien.DataResumePasien{
				KeluhanUtama:  "Demam 3 hari",
				DiagnosaUtama: "Febris suspect DHF",
				KondisiPulang: resumepasien.KondisiPulangHidup,
			},
		}
		bodyBytes, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/resume/"+encNoRawat+"/Ralan", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("POST SimpanResumePasien Mismatch NoRawat 400 Bad Request", func(t *testing.T) {
		handler := resumepasien.NewHandler(&mockResumePasienService{}, testEncKey)
		mux := http.NewServeMux()
		handler.RegisterRoutes(mux, dummyAuthMiddleware, dummyTimeoutMiddleware)

		payload := resumepasien.SimpanResumePasienRequest{
			NoRawat: "DIFFERENT_NO_RAWAT",
			DataResumePasien: resumepasien.DataResumePasien{
				KeluhanUtama:  "Demam",
				DiagnosaUtama: "Febris",
			},
		}
		bodyBytes, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/resume/"+encNoRawat+"/Ralan", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400 for mismatch NoRawat, got %d", w.Code)
		}
	})

	t.Run("PUT UpdateResumePasien 200 OK", func(t *testing.T) {
		mockSvc := &mockResumePasienService{
			updateFn: func(ctx context.Context, kodeDokterLogin, noRawat string, statusLanjut shared.StatusLanjut, req resumepasien.UpdateResumePasienRequest) (*resumepasien.ResumePasien, error) {
				return &resumepasien.ResumePasien{
					NoRawat:          noRawat,
					KodeDokter:       kodeDokterLogin,
					DataResumePasien: req.DataResumePasien,
				}, nil
			},
		}

		handler := resumepasien.NewHandler(mockSvc, testEncKey)
		mux := http.NewServeMux()
		handler.RegisterRoutes(mux, dummyAuthMiddleware, dummyTimeoutMiddleware)

		payload := resumepasien.UpdateResumePasienRequest{
			DataResumePasien: resumepasien.DataResumePasien{
				KeluhanUtama:  "Demam membaik",
				DiagnosaUtama: "DHF Grade 1",
				KondisiPulang: resumepasien.KondisiPulangHidup,
			},
		}
		bodyBytes, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/resume/"+encNoRawat+"/Ralan", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("DELETE HapusResumePasien 200 OK", func(t *testing.T) {
		mockSvc := &mockResumePasienService{
			hapusFn: func(ctx context.Context, kodeDokterLogin, noRawat string, statusLanjut shared.StatusLanjut) error {
				return nil
			},
		}

		handler := resumepasien.NewHandler(mockSvc, testEncKey)
		mux := http.NewServeMux()
		handler.RegisterRoutes(mux, dummyAuthMiddleware, dummyTimeoutMiddleware)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/resume/"+encNoRawat+"/Ralan", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})
}
