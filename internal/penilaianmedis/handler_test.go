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

	detailRanapFn  func(ctx context.Context, noRawat string) (*penilaianmedis.PenilaianMedisRanap, error)
	riwayatRanapFn func(ctx context.Context, noRM string) ([]penilaianmedis.PenilaianMedisRanap, error)
	simpanRanapFn  func(ctx context.Context, kodeDokterLogin, noRawat string, req penilaianmedis.SimpanPenilaianMedisRanapRequest) (*penilaianmedis.PenilaianMedisRanap, error)
	updateRanapFn  func(ctx context.Context, kodeDokterLogin, noRawat string, req penilaianmedis.UpdatePenilaianMedisRanapRequest) (*penilaianmedis.PenilaianMedisRanap, error)
	hapusRanapFn   func(ctx context.Context, kodeDokterLogin, noRawat string) error

	detailRalanKandunganFn  func(ctx context.Context, noRawat string) (*penilaianmedis.PenilaianMedisRalanKandungan, error)
	riwayatRalanKandunganFn func(ctx context.Context, noRM string) ([]penilaianmedis.PenilaianMedisRalanKandungan, error)
	simpanRalanKandunganFn  func(ctx context.Context, kodeDokterLogin, noRawat string, req penilaianmedis.SimpanPenilaianMedisRalanKandunganRequest) (*penilaianmedis.PenilaianMedisRalanKandungan, error)
	updateRalanKandunganFn  func(ctx context.Context, kodeDokterLogin, noRawat string, req penilaianmedis.UpdatePenilaianMedisRalanKandunganRequest) (*penilaianmedis.PenilaianMedisRalanKandungan, error)
	hapusRalanKandunganFn   func(ctx context.Context, kodeDokterLogin, noRawat string) error

	detailRanapKandunganFn  func(ctx context.Context, noRawat string) (*penilaianmedis.PenilaianMedisRanapKandungan, error)
	riwayatRanapKandunganFn func(ctx context.Context, noRM string) ([]penilaianmedis.PenilaianMedisRanapKandungan, error)
	simpanRanapKandunganFn  func(ctx context.Context, kodeDokterLogin, noRawat string, req penilaianmedis.SimpanPenilaianMedisRanapKandunganRequest) (*penilaianmedis.PenilaianMedisRanapKandungan, error)
	updateRanapKandunganFn  func(ctx context.Context, kodeDokterLogin, noRawat string, req penilaianmedis.UpdatePenilaianMedisRanapKandunganRequest) (*penilaianmedis.PenilaianMedisRanapKandungan, error)
	hapusRanapKandunganFn   func(ctx context.Context, kodeDokterLogin, noRawat string) error
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

func (m *mockPenilaianMedisService) DetailPenilaianMedisRanap(ctx context.Context, noRawat string) (*penilaianmedis.PenilaianMedisRanap, error) {
	if m.detailRanapFn != nil {
		return m.detailRanapFn(ctx, noRawat)
	}
	return nil, nil
}

func (m *mockPenilaianMedisService) RiwayatPenilaianMedisRanapByNoRM(ctx context.Context, noRM string) ([]penilaianmedis.PenilaianMedisRanap, error) {
	if m.riwayatRanapFn != nil {
		return m.riwayatRanapFn(ctx, noRM)
	}
	return nil, nil
}

func (m *mockPenilaianMedisService) SimpanPenilaianMedisRanap(ctx context.Context, kodeDokterLogin, noRawat string, req penilaianmedis.SimpanPenilaianMedisRanapRequest) (*penilaianmedis.PenilaianMedisRanap, error) {
	if m.simpanRanapFn != nil {
		return m.simpanRanapFn(ctx, kodeDokterLogin, noRawat, req)
	}
	return nil, nil
}

func (m *mockPenilaianMedisService) UpdatePenilaianMedisRanap(ctx context.Context, kodeDokterLogin, noRawat string, req penilaianmedis.UpdatePenilaianMedisRanapRequest) (*penilaianmedis.PenilaianMedisRanap, error) {
	if m.updateRanapFn != nil {
		return m.updateRanapFn(ctx, kodeDokterLogin, noRawat, req)
	}
	return nil, nil
}

func (m *mockPenilaianMedisService) HapusPenilaianMedisRanap(ctx context.Context, kodeDokterLogin, noRawat string) error {
	if m.hapusRanapFn != nil {
		return m.hapusRanapFn(ctx, kodeDokterLogin, noRawat)
	}
	return nil
}

func (m *mockPenilaianMedisService) DetailPenilaianMedisRalanKandungan(ctx context.Context, noRawat string) (*penilaianmedis.PenilaianMedisRalanKandungan, error) {
	if m.detailRalanKandunganFn != nil {
		return m.detailRalanKandunganFn(ctx, noRawat)
	}
	return nil, nil
}

func (m *mockPenilaianMedisService) RiwayatPenilaianMedisRalanKandunganByNoRM(ctx context.Context, noRM string) ([]penilaianmedis.PenilaianMedisRalanKandungan, error) {
	if m.riwayatRalanKandunganFn != nil {
		return m.riwayatRalanKandunganFn(ctx, noRM)
	}
	return nil, nil
}

func (m *mockPenilaianMedisService) SimpanPenilaianMedisRalanKandungan(ctx context.Context, kodeDokterLogin, noRawat string, req penilaianmedis.SimpanPenilaianMedisRalanKandunganRequest) (*penilaianmedis.PenilaianMedisRalanKandungan, error) {
	if m.simpanRalanKandunganFn != nil {
		return m.simpanRalanKandunganFn(ctx, kodeDokterLogin, noRawat, req)
	}
	return nil, nil
}

func (m *mockPenilaianMedisService) UpdatePenilaianMedisRalanKandungan(ctx context.Context, kodeDokterLogin, noRawat string, req penilaianmedis.UpdatePenilaianMedisRalanKandunganRequest) (*penilaianmedis.PenilaianMedisRalanKandungan, error) {
	if m.updateRalanKandunganFn != nil {
		return m.updateRalanKandunganFn(ctx, kodeDokterLogin, noRawat, req)
	}
	return nil, nil
}

func (m *mockPenilaianMedisService) HapusPenilaianMedisRalanKandungan(ctx context.Context, kodeDokterLogin, noRawat string) error {
	if m.hapusRalanKandunganFn != nil {
		return m.hapusRalanKandunganFn(ctx, kodeDokterLogin, noRawat)
	}
	return nil
}

func (m *mockPenilaianMedisService) DetailPenilaianMedisRanapKandungan(ctx context.Context, noRawat string) (*penilaianmedis.PenilaianMedisRanapKandungan, error) {
	if m.detailRanapKandunganFn != nil {
		return m.detailRanapKandunganFn(ctx, noRawat)
	}
	return nil, nil
}

func (m *mockPenilaianMedisService) RiwayatPenilaianMedisRanapKandunganByNoRM(ctx context.Context, noRM string) ([]penilaianmedis.PenilaianMedisRanapKandungan, error) {
	if m.riwayatRanapKandunganFn != nil {
		return m.riwayatRanapKandunganFn(ctx, noRM)
	}
	return nil, nil
}

func (m *mockPenilaianMedisService) SimpanPenilaianMedisRanapKandungan(ctx context.Context, kodeDokterLogin, noRawat string, req penilaianmedis.SimpanPenilaianMedisRanapKandunganRequest) (*penilaianmedis.PenilaianMedisRanapKandungan, error) {
	if m.simpanRanapKandunganFn != nil {
		return m.simpanRanapKandunganFn(ctx, kodeDokterLogin, noRawat, req)
	}
	return nil, nil
}

func (m *mockPenilaianMedisService) UpdatePenilaianMedisRanapKandungan(ctx context.Context, kodeDokterLogin, noRawat string, req penilaianmedis.UpdatePenilaianMedisRanapKandunganRequest) (*penilaianmedis.PenilaianMedisRanapKandungan, error) {
	if m.updateRanapKandunganFn != nil {
		return m.updateRanapKandunganFn(ctx, kodeDokterLogin, noRawat, req)
	}
	return nil, nil
}

func (m *mockPenilaianMedisService) HapusPenilaianMedisRanapKandungan(ctx context.Context, kodeDokterLogin, noRawat string) error {
	if m.hapusRanapKandunganFn != nil {
		return m.hapusRanapKandunganFn(ctx, kodeDokterLogin, noRawat)
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

	t.Run("Detail Ranap", func(t *testing.T) {
		mockSvc.detailRanapFn = func(ctx context.Context, noRawat string) (*penilaianmedis.PenilaianMedisRanap, error) {
			return &penilaianmedis.PenilaianMedisRanap{
				NoRawat: noRawat,
				DataPenilaianMedisRanap: penilaianmedis.DataPenilaianMedisRanap{
					TanggalPenilaian: "2026-09-03 10:00:00",
					Anamnesis:        penilaianmedis.Autoanamnesis,
					KeluhanUtama:     "Sesak nafas",
				},
			}, nil
		}
		req := httptest.NewRequest(http.MethodGet, "/api/v1/penilaian-medis/ranap/"+encKunjungan, nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
	})

	t.Run("Simpan Ranap", func(t *testing.T) {
		mockSvc.simpanRanapFn = func(ctx context.Context, kodeDokterLogin, noRawat string, req penilaianmedis.SimpanPenilaianMedisRanapRequest) (*penilaianmedis.PenilaianMedisRanap, error) {
			return &penilaianmedis.PenilaianMedisRanap{
				NoRawat:                 noRawat,
				DataPenilaianMedisRanap: req.DataPenilaianMedisRanap,
			}, nil
		}
		body, _ := json.Marshal(penilaianmedis.SimpanPenilaianMedisRanapRequest{
			NoRawat: "2026/09/03/000001",
			DataPenilaianMedisRanap: penilaianmedis.DataPenilaianMedisRanap{
				TanggalPenilaian: now.Format("2006-01-02 15:04:05"),
				Anamnesis:        penilaianmedis.Autoanamnesis,
				KeluhanUtama:     "Sesak nafas",
				Diagnosis:        "Pneumonia",
				TataLaksana:      "Oksigenasi",
			},
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/penilaian-medis/ranap/"+encKunjungan, bytes.NewReader(body))
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status 201 Created, got %d: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("Update Ranap", func(t *testing.T) {
		mockSvc.updateRanapFn = func(ctx context.Context, kodeDokterLogin, noRawat string, req penilaianmedis.UpdatePenilaianMedisRanapRequest) (*penilaianmedis.PenilaianMedisRanap, error) {
			return &penilaianmedis.PenilaianMedisRanap{
				NoRawat:                 noRawat,
				DataPenilaianMedisRanap: req.DataPenilaianMedisRanap,
			}, nil
		}
		body, _ := json.Marshal(penilaianmedis.UpdatePenilaianMedisRanapRequest{
			DataPenilaianMedisRanap: penilaianmedis.DataPenilaianMedisRanap{
				TanggalPenilaian: now.Format("2006-01-02 15:04:05"),
				Anamnesis:        penilaianmedis.Autoanamnesis,
				KeluhanUtama:     "Sesak nafas membaik",
				Diagnosis:        "Pneumonia perbaikan",
				TataLaksana:      "Oksigenasi",
			},
		})
		req := httptest.NewRequest(http.MethodPut, "/api/v1/penilaian-medis/ranap/"+encKunjungan, bytes.NewReader(body))
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200 OK, got %d: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("Hapus Ranap", func(t *testing.T) {
		mockSvc.hapusRanapFn = func(ctx context.Context, kodeDokterLogin, noRawat string) error {
			return nil
		}
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/penilaian-medis/ranap/"+encKunjungan, nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200 OK, got %d: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("Detail Ralan Kandungan", func(t *testing.T) {
		mockSvc.detailRalanKandunganFn = func(ctx context.Context, noRawat string) (*penilaianmedis.PenilaianMedisRalanKandungan, error) {
			return &penilaianmedis.PenilaianMedisRalanKandungan{
				NoRawat: noRawat,
				DataPenilaianMedisRalanKandungan: penilaianmedis.DataPenilaianMedisRalanKandungan{
					TanggalPenilaian: "2026-09-03 10:00:00",
					Anamnesis:        penilaianmedis.Autoanamnesis,
					KeluhanUtama:     "Perut kencang",
					Kontraksi:        penilaianmedis.KontraksiAda,
				},
			}, nil
		}
		req := httptest.NewRequest(http.MethodGet, "/api/v1/penilaian-medis/ralan-kandungan/"+encKunjungan, nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200 OK, got %d: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("Riwayat Ralan Kandungan", func(t *testing.T) {
		mockSvc.riwayatRalanKandunganFn = func(ctx context.Context, noRM string) ([]penilaianmedis.PenilaianMedisRalanKandungan, error) {
			return []penilaianmedis.PenilaianMedisRalanKandungan{
				{
					NoRawat: "2026/09/03/000001",
					DataPenilaianMedisRalanKandungan: penilaianmedis.DataPenilaianMedisRalanKandungan{
						KeluhanUtama: "Mules-mules",
						Kontraksi:    penilaianmedis.KontraksiAda,
					},
				},
			}, nil
		}
		encRM, _ := crypto.Encrypt("123456", testEncKey)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/penilaian-medis/ralan-kandungan/pasien/"+encRM, nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200 OK, got %d: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("Simpan Ralan Kandungan", func(t *testing.T) {
		mockSvc.simpanRalanKandunganFn = func(ctx context.Context, kodeDokterLogin, noRawat string, req penilaianmedis.SimpanPenilaianMedisRalanKandunganRequest) (*penilaianmedis.PenilaianMedisRalanKandungan, error) {
			return &penilaianmedis.PenilaianMedisRalanKandungan{
				NoRawat:                          noRawat,
				DataPenilaianMedisRalanKandungan: req.DataPenilaianMedisRalanKandungan,
			}, nil
		}
		body, _ := json.Marshal(penilaianmedis.SimpanPenilaianMedisRalanKandunganRequest{
			NoRawat: "2026/09/03/000001",
			DataPenilaianMedisRalanKandungan: penilaianmedis.DataPenilaianMedisRalanKandungan{
				TanggalPenilaian: now.Format("2006-01-02 15:04:05"),
				Anamnesis:        penilaianmedis.Autoanamnesis,
				KeluhanUtama:     "Perut mulas",
				Diagnosis:        "G1P0A0",
				TataLaksana:      "Observasi persalinan",
				Kontraksi:        penilaianmedis.KontraksiAda,
			},
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/penilaian-medis/ralan-kandungan/"+encKunjungan, bytes.NewReader(body))
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status 201 Created, got %d: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("Update Ralan Kandungan", func(t *testing.T) {
		mockSvc.updateRalanKandunganFn = func(ctx context.Context, kodeDokterLogin, noRawat string, req penilaianmedis.UpdatePenilaianMedisRalanKandunganRequest) (*penilaianmedis.PenilaianMedisRalanKandungan, error) {
			return &penilaianmedis.PenilaianMedisRalanKandungan{
				NoRawat:                          noRawat,
				DataPenilaianMedisRalanKandungan: req.DataPenilaianMedisRalanKandungan,
			}, nil
		}
		body, _ := json.Marshal(penilaianmedis.UpdatePenilaianMedisRalanKandunganRequest{
			DataPenilaianMedisRalanKandungan: penilaianmedis.DataPenilaianMedisRalanKandungan{
				TanggalPenilaian: now.Format("2006-01-02 15:04:05"),
				Anamnesis:        penilaianmedis.Autoanamnesis,
				KeluhanUtama:     "Mules berkurang",
				Diagnosis:        "G1P0A0 belum inpartu",
				TataLaksana:      "Rawat jalan",
				Kontraksi:        penilaianmedis.KontraksiTidak,
			},
		})
		req := httptest.NewRequest(http.MethodPut, "/api/v1/penilaian-medis/ralan-kandungan/"+encKunjungan, bytes.NewReader(body))
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200 OK, got %d: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("Hapus Ralan Kandungan", func(t *testing.T) {
		mockSvc.hapusRalanKandunganFn = func(ctx context.Context, kodeDokterLogin, noRawat string) error {
			return nil
		}
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/penilaian-medis/ralan-kandungan/"+encKunjungan, nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200 OK, got %d: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("Detail Ranap Kandungan", func(t *testing.T) {
		mockSvc.detailRanapKandunganFn = func(ctx context.Context, noRawat string) (*penilaianmedis.PenilaianMedisRanapKandungan, error) {
			return &penilaianmedis.PenilaianMedisRanapKandungan{
				NoRawat: noRawat,
				DataPenilaianMedisRanapKandungan: penilaianmedis.DataPenilaianMedisRanapKandungan{
					TanggalPenilaian: "2026-09-03 10:00:00",
					Anamnesis:        penilaianmedis.Autoanamnesis,
					KeluhanUtama:     "Perut kencang-kencang",
					Kontraksi:        penilaianmedis.KontraksiAda,
					Edukasi:          "Edukasi pendamping persalinan",
				},
			}, nil
		}
		req := httptest.NewRequest(http.MethodGet, "/api/v1/penilaian-medis/ranap-kandungan/"+encKunjungan, nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200 OK, got %d: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("Riwayat Ranap Kandungan", func(t *testing.T) {
		mockSvc.riwayatRanapKandunganFn = func(ctx context.Context, noRM string) ([]penilaianmedis.PenilaianMedisRanapKandungan, error) {
			return []penilaianmedis.PenilaianMedisRanapKandungan{
				{
					NoRawat: "2026/09/03/000001",
					DataPenilaianMedisRanapKandungan: penilaianmedis.DataPenilaianMedisRanapKandungan{
						KeluhanUtama: "Mules teratur",
						Kontraksi:    penilaianmedis.KontraksiAda,
						Edukasi:      "Tirah baring",
					},
				},
			}, nil
		}
		encRM, _ := crypto.Encrypt("123456", testEncKey)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/penilaian-medis/ranap-kandungan/pasien/"+encRM, nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200 OK, got %d: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("Simpan Ranap Kandungan", func(t *testing.T) {
		mockSvc.simpanRanapKandunganFn = func(ctx context.Context, kodeDokterLogin, noRawat string, req penilaianmedis.SimpanPenilaianMedisRanapKandunganRequest) (*penilaianmedis.PenilaianMedisRanapKandungan, error) {
			return &penilaianmedis.PenilaianMedisRanapKandungan{
				NoRawat:                          noRawat,
				DataPenilaianMedisRanapKandungan: req.DataPenilaianMedisRanapKandungan,
			}, nil
		}
		body, _ := json.Marshal(penilaianmedis.SimpanPenilaianMedisRanapKandunganRequest{
			NoRawat: "2026/09/03/000001",
			DataPenilaianMedisRanapKandungan: penilaianmedis.DataPenilaianMedisRanapKandungan{
				TanggalPenilaian: now.Format("2006-01-02 15:04:05"),
				Anamnesis:        penilaianmedis.Autoanamnesis,
				KeluhanUtama:     "Perut mulas",
				Diagnosis:        "G2P1A0 inpartu",
				TataLaksana:      "Observasi ketat",
				Kontraksi:        penilaianmedis.KontraksiAda,
				Edukasi:          "Edukasi melahirkan",
			},
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/penilaian-medis/ranap-kandungan/"+encKunjungan, bytes.NewReader(body))
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status 201 Created, got %d: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("Update Ranap Kandungan", func(t *testing.T) {
		mockSvc.updateRanapKandunganFn = func(ctx context.Context, kodeDokterLogin, noRawat string, req penilaianmedis.UpdatePenilaianMedisRanapKandunganRequest) (*penilaianmedis.PenilaianMedisRanapKandungan, error) {
			return &penilaianmedis.PenilaianMedisRanapKandungan{
				NoRawat:                          noRawat,
				DataPenilaianMedisRanapKandungan: req.DataPenilaianMedisRanapKandungan,
			}, nil
		}
		body, _ := json.Marshal(penilaianmedis.UpdatePenilaianMedisRanapKandunganRequest{
			DataPenilaianMedisRanapKandungan: penilaianmedis.DataPenilaianMedisRanapKandungan{
				TanggalPenilaian: now.Format("2006-01-02 15:04:05"),
				Anamnesis:        penilaianmedis.Autoanamnesis,
				KeluhanUtama:     "Mules bertambah sering",
				Diagnosis:        "G2P1A0 inpartu kala I fase aktif",
				TataLaksana:      "Observasi ketat",
				Kontraksi:        penilaianmedis.KontraksiAda,
				Edukasi:          "Edukasi persalinan",
			},
		})
		req := httptest.NewRequest(http.MethodPut, "/api/v1/penilaian-medis/ranap-kandungan/"+encKunjungan, bytes.NewReader(body))
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200 OK, got %d: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("Hapus Ranap Kandungan", func(t *testing.T) {
		mockSvc.hapusRanapKandunganFn = func(ctx context.Context, kodeDokterLogin, noRawat string) error {
			return nil
		}
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/penilaian-medis/ranap-kandungan/"+encKunjungan, nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200 OK, got %d: %s", rr.Code, rr.Body.String())
		}
	})
}
