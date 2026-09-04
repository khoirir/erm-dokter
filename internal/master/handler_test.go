package master_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"erm-dokter/internal/master"
)

type mockMasterService struct {
	daftarPenjaminFn   func(ctx context.Context) ([]master.Penjamin, error)
	daftarDepoFn       func(ctx context.Context) ([]master.Depo, error)
	daftarPoliklinikFn func(ctx context.Context) ([]master.Poliklinik, error)
}

func (m *mockMasterService) DaftarPenjamin(ctx context.Context) ([]master.Penjamin, error) {
	if m.daftarPenjaminFn != nil {
		return m.daftarPenjaminFn(ctx)
	}
	return nil, nil
}

func (m *mockMasterService) DaftarDepo(ctx context.Context) ([]master.Depo, error) {
	if m.daftarDepoFn != nil {
		return m.daftarDepoFn(ctx)
	}
	return nil, nil
}

func (m *mockMasterService) DaftarPoliklinik(ctx context.Context) ([]master.Poliklinik, error) {
	if m.daftarPoliklinikFn != nil {
		return m.daftarPoliklinikFn(ctx)
	}
	return nil, nil
}

func TestMasterHandler(t *testing.T) {
	mockSvc := &mockMasterService{
		daftarPenjaminFn: func(ctx context.Context) ([]master.Penjamin, error) {
			return []master.Penjamin{{ItemMaster: master.ItemMaster{Kode: "BPJ", Nama: "BPJS Kesehatan"}}}, nil
		},
		daftarDepoFn: func(ctx context.Context) ([]master.Depo, error) {
			return []master.Depo{{ItemMaster: master.ItemMaster{Kode: "DPRJ", Nama: "Depo Rawat Jalan"}}}, nil
		},
		daftarPoliklinikFn: func(ctx context.Context) ([]master.Poliklinik, error) {
			return []master.Poliklinik{{ItemMaster: master.ItemMaster{Kode: "INT", Nama: "Poli Penyakit Dalam"}}}, nil
		},
	}

	handler := master.NewHandler(mockSvc)
	mux := http.NewServeMux()
	noOpMw := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, noOpMw, noOpMw)

	t.Run("DaftarPenjamin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/master/penjamin", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
	})

	t.Run("DaftarDepo", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/master/depo", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
	})

	t.Run("DaftarPoliklinik", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/master/poliklinik", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		var resp struct {
			Success bool               `json:"success"`
			Data    []master.Poliklinik `json:"data"`
		}
		_ = json.Unmarshal(rr.Body.Bytes(), &resp)
		if len(resp.Data) != 1 || resp.Data[0].Kode != "INT" {
			t.Errorf("Unexpected data: %+v", resp.Data)
		}
	})
}
