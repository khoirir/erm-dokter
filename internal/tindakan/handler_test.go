package tindakan_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
	"erm-dokter/internal/tindakan"
)

type mockService struct {
	getDaftarTindakanLabFn func(ctx context.Context, kategori shared.KategoriLab, filter tindakan.FilterDaftarTindakanLab) ([]tindakan.TindakanLab, shared.PaginationMeta, error)
	getDetailTindakanLabFn func(ctx context.Context, kategori shared.KategoriLab, encryptedId string) (*tindakan.DetailTindakanLab, error)
}

func (m *mockService) GetDaftarTindakanLab(ctx context.Context, kategori shared.KategoriLab, filter tindakan.FilterDaftarTindakanLab) ([]tindakan.TindakanLab, shared.PaginationMeta, error) {
	if m.getDaftarTindakanLabFn != nil {
		return m.getDaftarTindakanLabFn(ctx, kategori, filter)
	}
	return nil, shared.PaginationMeta{}, nil
}

func (m *mockService) GetDetailTindakanLab(ctx context.Context, kategori shared.KategoriLab, encryptedId string) (*tindakan.DetailTindakanLab, error) {
	if m.getDetailTindakanLabFn != nil {
		return m.getDetailTindakanLabFn(ctx, kategori, encryptedId)
	}
	return nil, nil
}

func TestHandler_GetDaftarTindakanLab_Success(t *testing.T) {
	mockSvc := &mockService{
		getDaftarTindakanLabFn: func(ctx context.Context, kategori shared.KategoriLab, filter tindakan.FilterDaftarTindakanLab) ([]tindakan.TindakanLab, shared.PaginationMeta, error) {
			return []tindakan.TindakanLab{
				{
					Id:           "enc-pk001",
					KodeTindakan: "PK001",
					NamaTindakan: "Darah Rutin",
					Biaya:        50000,
				},
			}, shared.NewPaginationMeta(1, 1, 20), nil
		},
	}

	handler := tindakan.NewHandler(mockSvc)
	mux := http.NewServeMux()
	dummyMiddleware := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, dummyMiddleware, dummyMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tindakan/lab/pk?page=1&limit=20&keyword=Darah", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp struct {
		Code    int                   `json:"code"`
		Status  string                `json:"status"`
		Message string                `json:"message"`
		Data    []tindakan.TindakanLab `json:"data"`
		Meta    shared.PaginationMeta `json:"meta"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(resp.Data) != 1 || resp.Data[0].NamaTindakan != "Darah Rutin" {
		t.Errorf("Unexpected response data: %+v", resp.Data)
	}
}

func TestHandler_GetDaftarTindakanLab_InvalidKategori(t *testing.T) {
	handler := tindakan.NewHandler(&mockService{})
	mux := http.NewServeMux()
	dummyMiddleware := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, dummyMiddleware, dummyMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tindakan/lab/radiologi", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400 for invalid kategori, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_GetDaftarTindakanLab_InvalidKeyword(t *testing.T) {
	handler := tindakan.NewHandler(&mockService{})
	mux := http.NewServeMux()
	dummyMiddleware := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, dummyMiddleware, dummyMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tindakan/lab/pk?keyword=da", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400 for keyword < 3 chars, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_GetDetailTindakanLab_Success(t *testing.T) {
	mockSvc := &mockService{
		getDetailTindakanLabFn: func(ctx context.Context, kategori shared.KategoriLab, encryptedId string) (*tindakan.DetailTindakanLab, error) {
			return &tindakan.DetailTindakanLab{
				TindakanLab: tindakan.TindakanLab{
					Id:           "enc-pk001",
					KodeTindakan: "PK001",
					NamaTindakan: "Darah Lengkap",
					Biaya:        75000,
				},
				Templates: []tindakan.TemplateLab{
					{
						IdTemplate:      "enc-template-1",
						NamaPemeriksaan: "Hemoglobin",
						Satuan:          "g/dL",
						NilaiRujukanLD:  "13.5 - 17.5",
					},
				},
			}, nil
		},
	}

	handler := tindakan.NewHandler(mockSvc)
	mux := http.NewServeMux()
	dummyMiddleware := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, dummyMiddleware, dummyMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tindakan/lab/pk/enc-pk001", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp struct {
		Code    int                         `json:"code"`
		Status  string                      `json:"status"`
		Message string                      `json:"message"`
		Data    tindakan.DetailTindakanLab `json:"data"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Data.NamaTindakan != "Darah Lengkap" || len(resp.Data.Templates) != 1 {
		t.Errorf("Unexpected detail response: %+v", resp.Data)
	}
}

func TestHandler_GetDetailTindakanLab_NotFound(t *testing.T) {
	mockSvc := &mockService{
		getDetailTindakanLabFn: func(ctx context.Context, kategori shared.KategoriLab, encryptedId string) (*tindakan.DetailTindakanLab, error) {
			return nil, apperror.NewNotFoundError("Data tindakan laboratorium tidak ditemukan")
		},
	}

	handler := tindakan.NewHandler(mockSvc)
	mux := http.NewServeMux()
	dummyMiddleware := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, dummyMiddleware, dummyMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tindakan/lab/pk/enc-notfound", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("Expected status 404, got %d: %s", rr.Code, rr.Body.String())
	}
}
