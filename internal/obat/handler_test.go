package obat_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"erm-dokter/internal/obat"
	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/shared"
)

const testEncKey = "bafaa956-751d-4f59-98cc-574ee9dfe9f6"

type mockObatService struct {
	daftarObatFn     func(ctx context.Context, filter obat.FilterDaftarObat) ([]obat.Obat, shared.PaginationMeta, error)
	detailObatFn     func(ctx context.Context, kodeObat string) (*obat.Obat, error)
	daftarJenisFn    func(ctx context.Context) ([]obat.JenisObat, error)
	daftarGolonganFn func(ctx context.Context) ([]obat.GolonganObat, error)
	daftarKategoriFn func(ctx context.Context) ([]obat.KategoriObat, error)
}

func (m *mockObatService) DaftarObat(ctx context.Context, filter obat.FilterDaftarObat) ([]obat.Obat, shared.PaginationMeta, error) {
	if m.daftarObatFn != nil {
		return m.daftarObatFn(ctx, filter)
	}
	return nil, shared.PaginationMeta{}, nil
}

func (m *mockObatService) DetailObat(ctx context.Context, kodeObat string) (*obat.Obat, error) {
	if m.detailObatFn != nil {
		return m.detailObatFn(ctx, kodeObat)
	}
	return nil, nil
}

func (m *mockObatService) DaftarJenis(ctx context.Context) ([]obat.JenisObat, error) {
	if m.daftarJenisFn != nil {
		return m.daftarJenisFn(ctx)
	}
	return []obat.JenisObat{{Kode: "TAB", Nama: "Tablet"}}, nil
}

func (m *mockObatService) DaftarGolongan(ctx context.Context) ([]obat.GolonganObat, error) {
	if m.daftarGolonganFn != nil {
		return m.daftarGolonganFn(ctx)
	}
	return []obat.GolonganObat{{Kode: "BBS", Nama: "Bebas"}}, nil
}

func (m *mockObatService) DaftarKategori(ctx context.Context) ([]obat.KategoriObat, error) {
	if m.daftarKategoriFn != nil {
		return m.daftarKategoriFn(ctx)
	}
	return []obat.KategoriObat{{Kode: "ATB", Nama: "Antibiotik"}}, nil
}

func (m *mockObatService) CekKeberadaanObat(ctx context.Context, listKodeObat []string) (map[string]bool, error) {
	return map[string]bool{}, nil
}

func TestObatHandler_DaftarObat(t *testing.T) {
	mockSvc := &mockObatService{
		daftarObatFn: func(ctx context.Context, filter obat.FilterDaftarObat) ([]obat.Obat, shared.PaginationMeta, error) {
			return []obat.Obat{
				{
					KodeObat: "OB001",
					NamaObat: "Paracetamol 500mg",
					Harga:    "5000",
					Stok:     100,
				},
			}, shared.NewPaginationMeta(1, 1, 20), nil
		},
	}

	handler := obat.NewHandler(mockSvc, testEncKey)
	mux := http.NewServeMux()
	noOpMw := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, noOpMw, noOpMw)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/obat?page=1&limit=20", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rr.Code)
	}

	var resp struct {
		Success bool        `json:"success"`
		Data    []obat.Obat `json:"data"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(resp.Data) != 1 || resp.Data[0].NamaObat != "Paracetamol 500mg" {
		t.Errorf("Unexpected data: %+v", resp.Data)
	}
	if resp.Data[0].Id == "" {
		t.Error("Expected encrypted ID in obat response")
	}
}

func TestObatHandler_DetailObat(t *testing.T) {
	encObat, _ := crypto.Encrypt("OB001", testEncKey)

	mockSvc := &mockObatService{
		detailObatFn: func(ctx context.Context, kodeObat string) (*obat.Obat, error) {
			if kodeObat != "OB001" {
				t.Errorf("Expected OB001, got %s", kodeObat)
			}
			return &obat.Obat{
				KodeObat:  "OB001",
				NamaObat:  "Paracetamol 500mg",
				Harga:     "5000",
				Komposisi: "Paracetamol",
			}, nil
		},
	}

	handler := obat.NewHandler(mockSvc, testEncKey)
	mux := http.NewServeMux()
	noOpMw := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, noOpMw, noOpMw)

	t.Run("Valid Encrypted ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/obat/"+encObat, nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", rr.Code)
		}
	})

	t.Run("Invalid Encrypted ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/obat/invalid-token-here", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("Expected status 400 for invalid ID, got %d", rr.Code)
		}
	})
}
