package master_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"erm-dokter/internal/master"
	"erm-dokter/internal/shared"
)

type mockMasterService struct {
	daftarPenjaminFn   func(ctx context.Context) ([]master.Penjamin, error)
	daftarDepoFn       func(ctx context.Context) ([]master.Depo, error)
	daftarPoliklinikFn func(ctx context.Context) ([]master.Poliklinik, error)
	daftarBangsalFn    func(ctx context.Context) ([]master.Bangsal, error)
	daftarKelasFn      func(ctx context.Context) ([]master.KelasKamar, error)
	daftarICD10Fn      func(ctx context.Context, filter master.FilterMasterICD) ([]master.ICD10, shared.PaginationMeta, error)
	daftarICD9Fn       func(ctx context.Context, filter master.FilterMasterICD) ([]master.ICD9, shared.PaginationMeta, error)
	cekKeberadaanICD10Fn func(ctx context.Context, listKode []string) (map[string]bool, error)
	cekKeberadaanICD9Fn  func(ctx context.Context, listKode []string) (map[string]bool, error)
	syncICDFn            func(ctx context.Context) (*master.SyncICDResult, error)
}

func (m *mockMasterService) SyncICD(ctx context.Context) (*master.SyncICDResult, error) {
	if m.syncICDFn != nil {
		return m.syncICDFn(ctx)
	}
	return nil, nil
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

func (m *mockMasterService) DaftarBangsal(ctx context.Context) ([]master.Bangsal, error) {
	if m.daftarBangsalFn != nil {
		return m.daftarBangsalFn(ctx)
	}
	return nil, nil
}

func (m *mockMasterService) DaftarKelas(ctx context.Context) ([]master.KelasKamar, error) {
	if m.daftarKelasFn != nil {
		return m.daftarKelasFn(ctx)
	}
	return nil, nil
}

func (m *mockMasterService) DaftarICD10(ctx context.Context, filter master.FilterMasterICD) ([]master.ICD10, shared.PaginationMeta, error) {
	if m.daftarICD10Fn != nil {
		return m.daftarICD10Fn(ctx, filter)
	}
	return nil, shared.PaginationMeta{}, nil
}

func (m *mockMasterService) DaftarICD9(ctx context.Context, filter master.FilterMasterICD) ([]master.ICD9, shared.PaginationMeta, error) {
	if m.daftarICD9Fn != nil {
		return m.daftarICD9Fn(ctx, filter)
	}
	return nil, shared.PaginationMeta{}, nil
}

func (m *mockMasterService) CekKeberadaanICD10(ctx context.Context, listKode []string) (map[string]bool, error) {
	if m.cekKeberadaanICD10Fn != nil {
		return m.cekKeberadaanICD10Fn(ctx, listKode)
	}
	return nil, nil
}

func (m *mockMasterService) CekKeberadaanICD9(ctx context.Context, listKode []string) (map[string]bool, error) {
	if m.cekKeberadaanICD9Fn != nil {
		return m.cekKeberadaanICD9Fn(ctx, listKode)
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
		daftarBangsalFn: func(ctx context.Context) ([]master.Bangsal, error) {
			return []master.Bangsal{{ItemMaster: master.ItemMaster{Kode: "B01", Nama: "Melati"}}}, nil
		},
		daftarKelasFn: func(ctx context.Context) ([]master.KelasKamar, error) {
			return []master.KelasKamar{{ItemMaster: master.ItemMaster{Kode: "Kelas 1", Nama: "Kelas 1"}}}, nil
		},
		daftarICD10Fn: func(ctx context.Context, filter master.FilterMasterICD) ([]master.ICD10, shared.PaginationMeta, error) {
			return []master.ICD10{{ItemMaster: master.ItemMaster{Kode: "I63.9", Nama: "Cerebral infarction"}}}, shared.NewPaginationMeta(1, 1, 20), nil
		},
		daftarICD9Fn: func(ctx context.Context, filter master.FilterMasterICD) ([]master.ICD9, shared.PaginationMeta, error) {
			return []master.ICD9{{ItemMaster: master.ItemMaster{Kode: "87.03", Nama: "CT Head"}}}, shared.NewPaginationMeta(1, 1, 20), nil
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
			Success bool                `json:"success"`
			Data    []master.Poliklinik `json:"data"`
		}
		_ = json.Unmarshal(rr.Body.Bytes(), &resp)
		if len(resp.Data) != 1 || resp.Data[0].Kode != "INT" {
			t.Errorf("Unexpected data: %+v", resp.Data)
		}
	})

	t.Run("DaftarBangsal", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/master/bangsal", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		var resp struct {
			Success bool             `json:"success"`
			Data    []master.Bangsal `json:"data"`
		}
		_ = json.Unmarshal(rr.Body.Bytes(), &resp)
		if len(resp.Data) != 1 || resp.Data[0].Kode != "B01" || resp.Data[0].Nama != "Melati" {
			t.Errorf("Unexpected data: %+v", resp.Data)
		}
	})

	t.Run("DaftarKelas", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/master/kelas", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		var resp struct {
			Success bool                `json:"success"`
			Data    []master.KelasKamar `json:"data"`
		}
		_ = json.Unmarshal(rr.Body.Bytes(), &resp)
		if len(resp.Data) != 1 || resp.Data[0].Kode != "Kelas 1" || resp.Data[0].Nama != "Kelas 1" {
			t.Errorf("Unexpected data: %+v", resp.Data)
		}
	})

	t.Run("DaftarICD10_Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/master/icd10?keyword=cerebral&page=1&limit=20", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		var resp struct {
			Success bool           `json:"success"`
			Data    []master.ICD10 `json:"data"`
			Meta    shared.PaginationMeta `json:"meta"`
		}
		_ = json.Unmarshal(rr.Body.Bytes(), &resp)
		if len(resp.Data) != 1 || resp.Data[0].Kode != "I63.9" {
			t.Errorf("Unexpected data: %+v", resp.Data)
		}
	})

	t.Run("DaftarICD10_InvalidKeywordLength", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/master/icd10?keyword=ab", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}
	})

	t.Run("DaftarICD9_Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/master/icd9?keyword=head&page=1&limit=20", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		var resp struct {
			Success bool          `json:"success"`
			Data    []master.ICD9 `json:"data"`
			Meta    shared.PaginationMeta `json:"meta"`
		}
		_ = json.Unmarshal(rr.Body.Bytes(), &resp)
		if len(resp.Data) != 1 || resp.Data[0].Kode != "87.03" {
			t.Errorf("Unexpected data: %+v", resp.Data)
		}
	})

	t.Run("DaftarICD9_InvalidKeywordLength", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/master/icd9?keyword=x", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}
	})

	t.Run("SyncICD_Success", func(t *testing.T) {
		mockSvc.syncICDFn = func(ctx context.Context) (*master.SyncICDResult, error) {
			return &master.SyncICDResult{
				TotalICD10: 15000,
				TotalICD9:  4000,
				SyncedAt:   time.Now(),
			}, nil
		}
		req := httptest.NewRequest(http.MethodPost, "/api/v1/master/sync-icd", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		var resp struct {
			Success bool                  `json:"success"`
			Data    *master.SyncICDResult `json:"data"`
		}
		_ = json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp.Data == nil || resp.Data.TotalICD10 != 15000 || resp.Data.TotalICD9 != 4000 {
			t.Errorf("Unexpected data: %+v", resp.Data)
		}
	})
}
