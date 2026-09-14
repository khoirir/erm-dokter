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
)

const testEncKey = "bafaa956-751d-4f59-98cc-574ee9dfe9f6"

type mockResumePasienService struct {
	detailFn  func(ctx context.Context, noRawat string) (*resumepasien.ResumePasienRalan, error)
	riwayatFn func(ctx context.Context, noRM string) ([]resumepasien.ResumePasienRalan, error)
	simpanFn  func(ctx context.Context, noRawat, kodeDokter string, req resumepasien.SimpanResumePasienRalanRequest) (*resumepasien.ResumePasienRalan, error)
	updateFn  func(ctx context.Context, kodeDokterLogin, noRawat string, req resumepasien.UpdateResumePasienRalanRequest) (*resumepasien.ResumePasienRalan, error)
	hapusFn   func(ctx context.Context, kodeDokterLogin, noRawat string) error

	// Ranap
	detailRanapFn  func(ctx context.Context, noRawat string) (*resumepasien.ResumePasienRanap, error)
	riwayatRanapFn func(ctx context.Context, noRM string) ([]resumepasien.ResumePasienRanap, error)
	simpanRanapFn  func(ctx context.Context, noRawat, kodeDokter string, req resumepasien.SimpanResumePasienRanapRequest) (*resumepasien.ResumePasienRanap, error)
	updateRanapFn  func(ctx context.Context, kodeDokterLogin, noRawat string, req resumepasien.UpdateResumePasienRanapRequest) (*resumepasien.ResumePasienRanap, error)
	hapusRanapFn   func(ctx context.Context, kodeDokterLogin, noRawat string) error

	// Referensi
	referensiRanapFn func(ctx context.Context) resumepasien.ReferensiResumeRanap
	referensiRalanFn func(ctx context.Context) resumepasien.ReferensiResumeRalan
}

func (m *mockResumePasienService) DetailResumePasienRalan(ctx context.Context, noRawat string) (*resumepasien.ResumePasienRalan, error) {
	if m.detailFn != nil {
		return m.detailFn(ctx, noRawat)
	}
	return nil, nil
}

func (m *mockResumePasienService) RiwayatResumePasienRalanByNoRM(ctx context.Context, noRM string) ([]resumepasien.ResumePasienRalan, error) {
	if m.riwayatFn != nil {
		return m.riwayatFn(ctx, noRM)
	}
	return nil, nil
}

func (m *mockResumePasienService) SimpanResumePasienRalan(ctx context.Context, noRawat, kodeDokter string, req resumepasien.SimpanResumePasienRalanRequest) (*resumepasien.ResumePasienRalan, error) {
	if m.simpanFn != nil {
		return m.simpanFn(ctx, noRawat, kodeDokter, req)
	}
	return nil, nil
}

func (m *mockResumePasienService) UpdateResumePasienRalan(ctx context.Context, kodeDokterLogin, noRawat string, req resumepasien.UpdateResumePasienRalanRequest) (*resumepasien.ResumePasienRalan, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, kodeDokterLogin, noRawat, req)
	}
	return nil, nil
}

func (m *mockResumePasienService) HapusResumePasienRalan(ctx context.Context, kodeDokterLogin, noRawat string) error {
	if m.hapusFn != nil {
		return m.hapusFn(ctx, kodeDokterLogin, noRawat)
	}
	return nil
}

func (m *mockResumePasienService) DetailResumePasienRanap(ctx context.Context, noRawat string) (*resumepasien.ResumePasienRanap, error) {
	if m.detailRanapFn != nil {
		return m.detailRanapFn(ctx, noRawat)
	}
	return nil, nil
}

func (m *mockResumePasienService) RiwayatResumePasienRanapByNoRM(ctx context.Context, noRM string) ([]resumepasien.ResumePasienRanap, error) {
	if m.riwayatRanapFn != nil {
		return m.riwayatRanapFn(ctx, noRM)
	}
	return nil, nil
}

func (m *mockResumePasienService) SimpanResumePasienRanap(ctx context.Context, noRawat, kodeDokter string, req resumepasien.SimpanResumePasienRanapRequest) (*resumepasien.ResumePasienRanap, error) {
	if m.simpanRanapFn != nil {
		return m.simpanRanapFn(ctx, noRawat, kodeDokter, req)
	}
	return nil, nil
}

func (m *mockResumePasienService) UpdateResumePasienRanap(ctx context.Context, kodeDokterLogin, noRawat string, req resumepasien.UpdateResumePasienRanapRequest) (*resumepasien.ResumePasienRanap, error) {
	if m.updateRanapFn != nil {
		return m.updateRanapFn(ctx, kodeDokterLogin, noRawat, req)
	}
	return nil, nil
}

func (m *mockResumePasienService) HapusResumePasienRanap(ctx context.Context, kodeDokterLogin, noRawat string) error {
	if m.hapusRanapFn != nil {
		return m.hapusRanapFn(ctx, kodeDokterLogin, noRawat)
	}
	return nil
}

func (m *mockResumePasienService) ReferensiRanap(ctx context.Context) resumepasien.ReferensiResumeRanap {
	if m.referensiRanapFn != nil {
		return m.referensiRanapFn(ctx)
	}
	return resumepasien.ReferensiResumeRanap{
		CaraKeluar:    resumepasien.DaftarOpsiCaraKeluar(),
		KeadaanPulang: resumepasien.DaftarOpsiKeadaanPulang(),
		Dilanjutkan:   resumepasien.DaftarOpsiDilanjutkan(),
	}
}

func (m *mockResumePasienService) ReferensiRalan(ctx context.Context) resumepasien.ReferensiResumeRalan {
	if m.referensiRalanFn != nil {
		return m.referensiRalanFn(ctx)
	}
	return resumepasien.ReferensiResumeRalan{
		KeadaanPulang: resumepasien.DaftarOpsiKeadaanPulangRalan(),
	}
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
			detailFn: func(ctx context.Context, noRawat string) (*resumepasien.ResumePasienRalan, error) {
				return &resumepasien.ResumePasienRalan{
					NoRawat:    noRawat,
					KodeDokter: "DR01",
					NamaDokter: "dr. Handi",
					DataResumePasienRalan: resumepasien.DataResumePasienRalan{
						KeluhanUtama:  "Batuk",
						DiagnosaUtama: "ISPA",
						KeadaanPulang: resumepasien.KeadaanPulangHidup,
					},
				}, nil
			},
		}

		handler := resumepasien.NewHandler(mockSvc, testEncKey)
		mux := http.NewServeMux()
		handler.RegisterRoutes(mux, dummyAuthMiddleware, dummyTimeoutMiddleware)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/resume/ralan/"+encNoRawat, nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("GET RiwayatResumePasienByNoRM Sukses", func(t *testing.T) {
		mockSvc := &mockResumePasienService{
			riwayatFn: func(ctx context.Context, noRM string) ([]resumepasien.ResumePasienRalan, error) {
				return []resumepasien.ResumePasienRalan{
					{
						NoRawat: rawNoRawat,
						DataResumePasienRalan: resumepasien.DataResumePasienRalan{
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

		req := httptest.NewRequest(http.MethodGet, "/api/v1/resume/ralan/pasien/"+encNoRM, nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("POST SimpanResumePasien 201 Created", func(t *testing.T) {
		mockSvc := &mockResumePasienService{
			simpanFn: func(ctx context.Context, noRawat, kodeDokter string, req resumepasien.SimpanResumePasienRalanRequest) (*resumepasien.ResumePasienRalan, error) {
				return &resumepasien.ResumePasienRalan{
					NoRawat:               noRawat,
					KodeDokter:            kodeDokter,
					DataResumePasienRalan: req.DataResumePasienRalan,
				}, nil
			},
		}

		handler := resumepasien.NewHandler(mockSvc, testEncKey)
		mux := http.NewServeMux()
		handler.RegisterRoutes(mux, dummyAuthMiddleware, dummyTimeoutMiddleware)

		payload := resumepasien.SimpanResumePasienRalanRequest{
			NoRawat: rawNoRawat,
			DataResumePasienRalan: resumepasien.DataResumePasienRalan{
				KeluhanUtama:  "Demam 3 hari",
				DiagnosaUtama: "Febris suspect DHF",
				KeadaanPulang: resumepasien.KeadaanPulangHidup,
			},
		}
		bodyBytes, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/resume/ralan/"+encNoRawat, bytes.NewReader(bodyBytes))
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

		payload := resumepasien.SimpanResumePasienRalanRequest{
			NoRawat: "DIFFERENT_NO_RAWAT",
			DataResumePasienRalan: resumepasien.DataResumePasienRalan{
				KeluhanUtama:  "Demam",
				DiagnosaUtama: "Febris",
			},
		}
		bodyBytes, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/resume/ralan/"+encNoRawat, bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400 for mismatch NoRawat, got %d", w.Code)
		}
	})

	t.Run("PUT UpdateResumePasien 200 OK", func(t *testing.T) {
		mockSvc := &mockResumePasienService{
			updateFn: func(ctx context.Context, kodeDokterLogin, noRawat string, req resumepasien.UpdateResumePasienRalanRequest) (*resumepasien.ResumePasienRalan, error) {
				return &resumepasien.ResumePasienRalan{
					NoRawat:               noRawat,
					KodeDokter:            kodeDokterLogin,
					DataResumePasienRalan: req.DataResumePasienRalan,
				}, nil
			},
		}

		handler := resumepasien.NewHandler(mockSvc, testEncKey)
		mux := http.NewServeMux()
		handler.RegisterRoutes(mux, dummyAuthMiddleware, dummyTimeoutMiddleware)

		payload := resumepasien.UpdateResumePasienRalanRequest{
			DataResumePasienRalan: resumepasien.DataResumePasienRalan{
				KeluhanUtama:  "Demam membaik",
				DiagnosaUtama: "DHF Grade 1",
				KeadaanPulang: resumepasien.KeadaanPulangHidup,
			},
		}
		bodyBytes, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/resume/ralan/"+encNoRawat, bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("DELETE HapusResumePasien 200 OK", func(t *testing.T) {
		mockSvc := &mockResumePasienService{
			hapusFn: func(ctx context.Context, kodeDokterLogin, noRawat string) error {
				return nil
			},
		}

		handler := resumepasien.NewHandler(mockSvc, testEncKey)
		mux := http.NewServeMux()
		handler.RegisterRoutes(mux, dummyAuthMiddleware, dummyTimeoutMiddleware)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/resume/ralan/"+encNoRawat, nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	// ==================== RANAP ====================

	t.Run("GET DetailResumePasienRanap Sukses", func(t *testing.T) {
		mockSvc := &mockResumePasienService{
			detailRanapFn: func(ctx context.Context, noRawat string) (*resumepasien.ResumePasienRanap, error) {
				return &resumepasien.ResumePasienRanap{
					NoRawat:    noRawat,
					KodeDokter: "DR01",
					NamaDokter: "dr. Handi",
					DataResumePasienRanap: resumepasien.DataResumePasienRanap{
						KeluhanUtama:  "Nyeri dada",
						DiagnosaUtama: "STEMI",
						CaraKeluar:    resumepasien.CaraKeluarAtasIzinDokter,
						KeadaanPulang: resumepasien.KeadaanPulangMembaik,
						Dilanjutkan:   resumepasien.DilanjutkanKembaliKeRS,
					},
				}, nil
			},
		}

		handler := resumepasien.NewHandler(mockSvc, testEncKey)
		mux := http.NewServeMux()
		handler.RegisterRoutes(mux, dummyAuthMiddleware, dummyTimeoutMiddleware)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/resume/ranap/"+encNoRawat, nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("GET RiwayatResumePasienRanapByNoRM Sukses", func(t *testing.T) {
		mockSvc := &mockResumePasienService{
			riwayatRanapFn: func(ctx context.Context, noRM string) ([]resumepasien.ResumePasienRanap, error) {
				return []resumepasien.ResumePasienRanap{
					{
						NoRawat: rawNoRawat,
						DataResumePasienRanap: resumepasien.DataResumePasienRanap{
							DiagnosaUtama: "STEMI",
						},
					},
				}, nil
			},
		}

		handler := resumepasien.NewHandler(mockSvc, testEncKey)
		mux := http.NewServeMux()
		handler.RegisterRoutes(mux, dummyAuthMiddleware, dummyTimeoutMiddleware)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/resume/ranap/pasien/"+encNoRM, nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("POST SimpanResumePasienRanap 201 Created", func(t *testing.T) {
		mockSvc := &mockResumePasienService{
			simpanRanapFn: func(ctx context.Context, noRawat, kodeDokter string, req resumepasien.SimpanResumePasienRanapRequest) (*resumepasien.ResumePasienRanap, error) {
				return &resumepasien.ResumePasienRanap{
					NoRawat:               noRawat,
					KodeDokter:            kodeDokter,
					DataResumePasienRanap: req.DataResumePasienRanap,
				}, nil
			},
		}

		handler := resumepasien.NewHandler(mockSvc, testEncKey)
		mux := http.NewServeMux()
		handler.RegisterRoutes(mux, dummyAuthMiddleware, dummyTimeoutMiddleware)

		payload := resumepasien.SimpanResumePasienRanapRequest{
			NoRawat: rawNoRawat,
			DataResumePasienRanap: resumepasien.DataResumePasienRanap{
				DiagnosaAwal:  "Chest Pain",
				AlasanRawat:   "Evaluasi nyeri dada",
				KeluhanUtama:  "Nyeri dada kiri menjalar",
				DiagnosaUtama: "STEMI Anterior",
				CaraKeluar:    resumepasien.CaraKeluarAtasIzinDokter,
				KeadaanPulang: resumepasien.KeadaanPulangMembaik,
				Dilanjutkan:   resumepasien.DilanjutkanKembaliKeRS,
				WaktuKontrol:  "2026-09-12 09:00:00",
			},
		}
		bodyBytes, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/resume/ranap/"+encNoRawat, bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("POST SimpanResumePasienRanap Mismatch NoRawat 400 Bad Request", func(t *testing.T) {
		handler := resumepasien.NewHandler(&mockResumePasienService{}, testEncKey)
		mux := http.NewServeMux()
		handler.RegisterRoutes(mux, dummyAuthMiddleware, dummyTimeoutMiddleware)

		payload := resumepasien.SimpanResumePasienRanapRequest{
			NoRawat: "DIFFERENT_NO_RAWAT",
			DataResumePasienRanap: resumepasien.DataResumePasienRanap{
				DiagnosaAwal:  "Chest Pain",
				AlasanRawat:   "Evaluasi nyeri dada",
				KeluhanUtama:  "Nyeri dada",
				DiagnosaUtama: "STEMI",
				CaraKeluar:    resumepasien.CaraKeluarAtasIzinDokter,
				KeadaanPulang: resumepasien.KeadaanPulangMembaik,
				Dilanjutkan:   resumepasien.DilanjutkanKembaliKeRS,
				WaktuKontrol:  "2026-09-12 09:00:00",
			},
		}
		bodyBytes, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/resume/ranap/"+encNoRawat, bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400 for mismatch NoRawat, got %d", w.Code)
		}
	})

	t.Run("PUT UpdateResumePasienRanap 200 OK", func(t *testing.T) {
		mockSvc := &mockResumePasienService{
			updateRanapFn: func(ctx context.Context, kodeDokterLogin, noRawat string, req resumepasien.UpdateResumePasienRanapRequest) (*resumepasien.ResumePasienRanap, error) {
				return &resumepasien.ResumePasienRanap{
					NoRawat:               noRawat,
					KodeDokter:            kodeDokterLogin,
					DataResumePasienRanap: req.DataResumePasienRanap,
				}, nil
			},
		}

		handler := resumepasien.NewHandler(mockSvc, testEncKey)
		mux := http.NewServeMux()
		handler.RegisterRoutes(mux, dummyAuthMiddleware, dummyTimeoutMiddleware)

		payload := resumepasien.UpdateResumePasienRanapRequest{
			DataResumePasienRanap: resumepasien.DataResumePasienRanap{
				DiagnosaAwal:  "Chest Pain",
				AlasanRawat:   "Evaluasi nyeri dada",
				KeluhanUtama:  "Nyeri dada berkurang",
				DiagnosaUtama: "STEMI Anterior Resolving",
				CaraKeluar:    resumepasien.CaraKeluarAtasIzinDokter,
				KeadaanPulang: resumepasien.KeadaanPulangSembuh,
				Dilanjutkan:   resumepasien.DilanjutkanKembaliKeRS,
				WaktuKontrol:  "2026-09-12 09:00:00",
			},
		}
		bodyBytes, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/resume/ranap/"+encNoRawat, bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("DELETE HapusResumePasienRanap 200 OK", func(t *testing.T) {
		mockSvc := &mockResumePasienService{
			hapusRanapFn: func(ctx context.Context, kodeDokterLogin, noRawat string) error {
				return nil
			},
		}

		handler := resumepasien.NewHandler(mockSvc, testEncKey)
		mux := http.NewServeMux()
		handler.RegisterRoutes(mux, dummyAuthMiddleware, dummyTimeoutMiddleware)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/resume/ranap/"+encNoRawat, nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	// ==================== REFERENSI ====================

	t.Run("GET Referensi Ralan Sukses", func(t *testing.T) {
		handler := resumepasien.NewHandler(&mockResumePasienService{}, testEncKey)
		mux := http.NewServeMux()
		handler.RegisterRoutes(mux, dummyAuthMiddleware, dummyTimeoutMiddleware)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/resume/ralan/referensi", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("GET Referensi Ranap Sukses", func(t *testing.T) {
		handler := resumepasien.NewHandler(&mockResumePasienService{}, testEncKey)
		mux := http.NewServeMux()
		handler.RegisterRoutes(mux, dummyAuthMiddleware, dummyTimeoutMiddleware)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/resume/ranap/referensi", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})
}

