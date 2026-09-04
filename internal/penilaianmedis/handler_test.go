package penilaianmedis_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"erm-dokter/internal/middleware"
	"erm-dokter/internal/penilaianmedis"
	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/pkg/token"
)

const testEncKey = "bafaa956-751d-4f59-98cc-574ee9dfe9f6"

type mockPenilaianMedisService struct {
	referensiFn func(ctx context.Context) penilaianmedis.ReferensiPenilaianMedis

	detailRalanFn  func(ctx context.Context, noRawat string) (*penilaianmedis.PenilaianMedisRalan, error)
	riwayatRalanFn func(ctx context.Context, noRM string) ([]penilaianmedis.PenilaianMedisRalan, error)
	simpanRalanFn  func(ctx context.Context, kodeDokterLogin, noRawat string, req penilaianmedis.SimpanPenilaianMedisRalanRequest) (*penilaianmedis.PenilaianMedisRalan, error)
	updateRalanFn  func(ctx context.Context, kodeDokterLogin, noRawat string, req penilaianmedis.UpdatePenilaianMedisRalanRequest) (*penilaianmedis.PenilaianMedisRalan, error)
	hapusRalanFn   func(ctx context.Context, kodeDokterLogin, noRawat string) error

	detailIGDFn  func(ctx context.Context, noRawat string) (*penilaianmedis.PenilaianMedisIGD, error)
	riwayatIGDFn func(ctx context.Context, noRM string) ([]penilaianmedis.PenilaianMedisIGD, error)
	simpanIGDFn  func(ctx context.Context, kodeDokterLogin, noRawat string, req penilaianmedis.SimpanPenilaianMedisIGDRequest) (*penilaianmedis.PenilaianMedisIGD, error)
	updateIGDFn  func(ctx context.Context, kodeDokterLogin, noRawat string, req penilaianmedis.UpdatePenilaianMedisIGDRequest) (*penilaianmedis.PenilaianMedisIGD, error)
	hapusIGDFn   func(ctx context.Context, kodeDokterLogin, noRawat string) error
}

func (m *mockPenilaianMedisService) Referensi(ctx context.Context) penilaianmedis.ReferensiPenilaianMedis {
	if m.referensiFn != nil {
		return m.referensiFn(ctx)
	}
	return penilaianmedis.ReferensiPenilaianMedis{}
}

func (m *mockPenilaianMedisService) DetailPenilaianMedisRalan(ctx context.Context, noRawat string) (*penilaianmedis.PenilaianMedisRalan, error) {
	if m.detailRalanFn != nil {
		return m.detailRalanFn(ctx, noRawat)
	}
	return nil, nil
}

func (m *mockPenilaianMedisService) RiwayatPenilaianMedisRalanByNoRM(ctx context.Context, noRM string) ([]penilaianmedis.PenilaianMedisRalan, error) {
	if m.riwayatRalanFn != nil {
		return m.riwayatRalanFn(ctx, noRM)
	}
	return nil, nil
}

func (m *mockPenilaianMedisService) SimpanPenilaianMedisRalan(ctx context.Context, kodeDokterLogin, noRawat string, req penilaianmedis.SimpanPenilaianMedisRalanRequest) (*penilaianmedis.PenilaianMedisRalan, error) {
	if m.simpanRalanFn != nil {
		return m.simpanRalanFn(ctx, kodeDokterLogin, noRawat, req)
	}
	return nil, nil
}

func (m *mockPenilaianMedisService) UpdatePenilaianMedisRalan(ctx context.Context, kodeDokterLogin, noRawat string, req penilaianmedis.UpdatePenilaianMedisRalanRequest) (*penilaianmedis.PenilaianMedisRalan, error) {
	if m.updateRalanFn != nil {
		return m.updateRalanFn(ctx, kodeDokterLogin, noRawat, req)
	}
	return nil, nil
}

func (m *mockPenilaianMedisService) HapusPenilaianMedisRalan(ctx context.Context, kodeDokterLogin, noRawat string) error {
	if m.hapusRalanFn != nil {
		return m.hapusRalanFn(ctx, kodeDokterLogin, noRawat)
	}
	return nil
}

func (m *mockPenilaianMedisService) DetailPenilaianMedisIGD(ctx context.Context, noRawat string) (*penilaianmedis.PenilaianMedisIGD, error) {
	if m.detailIGDFn != nil {
		return m.detailIGDFn(ctx, noRawat)
	}
	return nil, nil
}

func (m *mockPenilaianMedisService) RiwayatPenilaianMedisIGDByNoRM(ctx context.Context, noRM string) ([]penilaianmedis.PenilaianMedisIGD, error) {
	if m.riwayatIGDFn != nil {
		return m.riwayatIGDFn(ctx, noRM)
	}
	return nil, nil
}

func (m *mockPenilaianMedisService) SimpanPenilaianMedisIGD(ctx context.Context, kodeDokterLogin, noRawat string, req penilaianmedis.SimpanPenilaianMedisIGDRequest) (*penilaianmedis.PenilaianMedisIGD, error) {
	if m.simpanIGDFn != nil {
		return m.simpanIGDFn(ctx, kodeDokterLogin, noRawat, req)
	}
	return nil, nil
}

func (m *mockPenilaianMedisService) UpdatePenilaianMedisIGD(ctx context.Context, kodeDokterLogin, noRawat string, req penilaianmedis.UpdatePenilaianMedisIGDRequest) (*penilaianmedis.PenilaianMedisIGD, error) {
	if m.updateIGDFn != nil {
		return m.updateIGDFn(ctx, kodeDokterLogin, noRawat, req)
	}
	return nil, nil
}

func (m *mockPenilaianMedisService) HapusPenilaianMedisIGD(ctx context.Context, kodeDokterLogin, noRawat string) error {
	if m.hapusIGDFn != nil {
		return m.hapusIGDFn(ctx, kodeDokterLogin, noRawat)
	}
	return nil
}

func authMw(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), middleware.UserClaimKey, &token.Claims{KodeDokter: "DR01"})
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func TestPenilaianMedisHandler(t *testing.T) {
	encKunjungan, _ := crypto.Encrypt("2026/09/03/000001", testEncKey)
	now := time.Now()

	mockSvc := &mockPenilaianMedisService{
		detailRalanFn: func(ctx context.Context, noRawat string) (*penilaianmedis.PenilaianMedisRalan, error) {
			return &penilaianmedis.PenilaianMedisRalan{
				NoRawat: noRawat,
				DataPenilaianMedisRalan: penilaianmedis.DataPenilaianMedisRalan{
					TanggalPenilaian: "2026-09-03 10:00:00",
					Anamnesis:        penilaianmedis.Autoanamnesis,
					KeluhanUtama:     "Demam",
				},
			}, nil
		},
		simpanRalanFn: func(ctx context.Context, kodeDokterLogin, noRawat string, req penilaianmedis.SimpanPenilaianMedisRalanRequest) (*penilaianmedis.PenilaianMedisRalan, error) {
			return &penilaianmedis.PenilaianMedisRalan{
				NoRawat:                 noRawat,
				DataPenilaianMedisRalan: req.DataPenilaianMedisRalan,
			}, nil
		},
		detailIGDFn: func(ctx context.Context, noRawat string) (*penilaianmedis.PenilaianMedisIGD, error) {
			return &penilaianmedis.PenilaianMedisIGD{
				NoRawat: noRawat,
				DataPenilaianMedisIGD: penilaianmedis.DataPenilaianMedisIGD{
					TanggalPenilaian: "2026-09-03 10:00:00",
					Anamnesis:        penilaianmedis.Autoanamnesis,
					KeluhanUtama:     "Nyeri dada akut",
				},
			}, nil
		},
	}

	handler := penilaianmedis.NewHandler(mockSvc, testEncKey)
	mux := http.NewServeMux()
	noOpMw := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, authMw, noOpMw)

	t.Run("Referensi", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/penilaian-medis/referensi", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
	})

	t.Run("Detail Ralan", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/penilaian-medis/ralan/"+encKunjungan, nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
	})

	t.Run("Simpan Ralan", func(t *testing.T) {
		body, _ := json.Marshal(penilaianmedis.SimpanPenilaianMedisRalanRequest{
			NoRawat: "2026/09/03/000001",
			DataPenilaianMedisRalan: penilaianmedis.DataPenilaianMedisRalan{
				TanggalPenilaian: now.Format("2006-01-02 15:04:05"),
				Anamnesis:        penilaianmedis.Autoanamnesis,
				KeluhanUtama:     "Demam",
				Diagnosis:        "Febris",
				TataLaksana:      "Paracetamol",
			},
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/penilaian-medis/ralan/"+encKunjungan, bytes.NewReader(body))
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status 201 Created, got %d: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("Detail IGD", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/penilaian-medis/igd/"+encKunjungan, nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
	})
}
