package tindakan_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
	"erm-dokter/internal/tindakan"
)

const testEncryptionKey = "12345678901234567890123456789012"

type mockService struct {
	getDaftarTindakanLabFn     func(ctx context.Context, kategori shared.KategoriLab, filter tindakan.FilterDaftarTindakanLab) ([]tindakan.TindakanLab, shared.PaginationMeta, error)
	getDetailTindakanLabFn     func(ctx context.Context, kategori shared.KategoriLab, kodeTindakan string) (*tindakan.DetailTindakanLab, error)
	cekKeberadaanTindakanLabFn func(ctx context.Context, kategori shared.KategoriLab, listKodeTindakan []string) (map[string]bool, error)
	cekKeberadaanTemplateLabFn func(ctx context.Context, listKodeTindakan []string, templateMap map[string][]int) (map[string]map[int]bool, error)

	getDaftarTindakanRadiologiFn     func(ctx context.Context, filter tindakan.FilterDaftarTindakanRadiologi) ([]tindakan.TindakanRadiologi, shared.PaginationMeta, error)
	getDetailTindakanRadiologiFn     func(ctx context.Context, kodeTindakan string) (*tindakan.TindakanRadiologi, error)
	cekKeberadaanTindakanRadiologiFn func(ctx context.Context, listKodeTindakan []string) (map[string]bool, error)
}

func (m *mockService) GetDaftarTindakanLab(ctx context.Context, kategori shared.KategoriLab, filter tindakan.FilterDaftarTindakanLab) ([]tindakan.TindakanLab, shared.PaginationMeta, error) {
	if m.getDaftarTindakanLabFn != nil {
		return m.getDaftarTindakanLabFn(ctx, kategori, filter)
	}
	return nil, shared.PaginationMeta{}, nil
}

func (m *mockService) GetDetailTindakanLab(ctx context.Context, kategori shared.KategoriLab, kodeTindakan string) (*tindakan.DetailTindakanLab, error) {
	if m.getDetailTindakanLabFn != nil {
		return m.getDetailTindakanLabFn(ctx, kategori, kodeTindakan)
	}
	return nil, nil
}

func (m *mockService) CekKeberadaanTindakanLab(ctx context.Context, kategori shared.KategoriLab, listKodeTindakan []string) (map[string]bool, error) {
	if m.cekKeberadaanTindakanLabFn != nil {
		return m.cekKeberadaanTindakanLabFn(ctx, kategori, listKodeTindakan)
	}
	return make(map[string]bool), nil
}

func (m *mockService) CekKeberadaanTemplateLab(ctx context.Context, listKodeTindakan []string, templateMap map[string][]int) (map[string]map[int]bool, error) {
	if m.cekKeberadaanTemplateLabFn != nil {
		return m.cekKeberadaanTemplateLabFn(ctx, listKodeTindakan, templateMap)
	}
	return make(map[string]map[int]bool), nil
}

func (m *mockService) GetDaftarTindakanRadiologi(ctx context.Context, filter tindakan.FilterDaftarTindakanRadiologi) ([]tindakan.TindakanRadiologi, shared.PaginationMeta, error) {
	if m.getDaftarTindakanRadiologiFn != nil {
		return m.getDaftarTindakanRadiologiFn(ctx, filter)
	}
	return nil, shared.PaginationMeta{}, nil
}

func (m *mockService) GetDetailTindakanRadiologi(ctx context.Context, kodeTindakan string) (*tindakan.TindakanRadiologi, error) {
	if m.getDetailTindakanRadiologiFn != nil {
		return m.getDetailTindakanRadiologiFn(ctx, kodeTindakan)
	}
	return nil, nil
}

func (m *mockService) CekKeberadaanTindakanRadiologi(ctx context.Context, listKodeTindakan []string) (map[string]bool, error) {
	if m.cekKeberadaanTindakanRadiologiFn != nil {
		return m.cekKeberadaanTindakanRadiologiFn(ctx, listKodeTindakan)
	}
	return make(map[string]bool), nil
}


func TestHandler_GetDaftarTindakanLab_Success(t *testing.T) {
	mockSvc := &mockService{
		getDaftarTindakanLabFn: func(ctx context.Context, kategori shared.KategoriLab, filter tindakan.FilterDaftarTindakanLab) ([]tindakan.TindakanLab, shared.PaginationMeta, error) {
			return []tindakan.TindakanLab{
				{
					KodeTindakan: "PK001",
					NamaTindakan: "Darah Rutin",
					Biaya:        50000,
				},
			}, shared.NewPaginationMeta(1, 1, 20), nil
		},
	}

	handler := tindakan.NewHandler(mockSvc, testEncryptionKey)
	mux := http.NewServeMux()
	dummyMiddleware := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, dummyMiddleware, dummyMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tindakan/laboratorium/pk?page=1&limit=20&keyword=Darah", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp struct {
		Code    int                    `json:"code"`
		Status  string                 `json:"status"`
		Message string                 `json:"message"`
		Data    []tindakan.TindakanLab `json:"data"`
		Meta    shared.PaginationMeta  `json:"meta"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(resp.Data) != 1 || resp.Data[0].NamaTindakan != "Darah Rutin" {
		t.Errorf("Unexpected response data: %+v", resp.Data)
	}

	decryptedId, errDec := crypto.Decrypt(resp.Data[0].Id, testEncryptionKey)
	if errDec != nil || decryptedId != "PK001" {
		t.Errorf("Failed to decrypt Id in handler response: %s, err: %v", decryptedId, errDec)
	}
}

func TestHandler_GetDaftarTindakanLab_InvalidKategori(t *testing.T) {
	handler := tindakan.NewHandler(&mockService{}, testEncryptionKey)
	mux := http.NewServeMux()
	dummyMiddleware := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, dummyMiddleware, dummyMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tindakan/laboratorium/radiologi", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400 for invalid kategori, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_GetDaftarTindakanLab_InvalidKeyword(t *testing.T) {
	handler := tindakan.NewHandler(&mockService{}, testEncryptionKey)
	mux := http.NewServeMux()
	dummyMiddleware := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, dummyMiddleware, dummyMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tindakan/laboratorium/pk?keyword=da", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400 for keyword < 3 chars, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_GetDetailTindakanLab_Success(t *testing.T) {
	mockSvc := &mockService{
		getDetailTindakanLabFn: func(ctx context.Context, kategori shared.KategoriLab, kodeTindakan string) (*tindakan.DetailTindakanLab, error) {
			if kodeTindakan != "PK001" {
				t.Errorf("Expected decrypted kodeTindakan PK001, got %s", kodeTindakan)
			}
			return &tindakan.DetailTindakanLab{
				TindakanLab: tindakan.TindakanLab{
					KodeTindakan: "PK001",
					NamaTindakan: "Darah Lengkap",
					Biaya:        75000,
				},
				Templates: []tindakan.TemplateLab{
					{
						IdTemplate:      "101",
						NamaPemeriksaan: "Hemoglobin",
						Satuan:          "g/dL",
						NilaiRujukanLD:  "13.5 - 17.5",
					},
				},
			}, nil
		},
	}

	handler := tindakan.NewHandler(mockSvc, testEncryptionKey)
	mux := http.NewServeMux()
	dummyMiddleware := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, dummyMiddleware, dummyMiddleware)

	encId, _ := crypto.Encrypt("PK001", testEncryptionKey)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tindakan/laboratorium/pk/"+encId, nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp struct {
		Code    int                        `json:"code"`
		Status  string                     `json:"status"`
		Message string                     `json:"message"`
		Data    tindakan.DetailTindakanLab `json:"data"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Data.NamaTindakan != "Darah Lengkap" || len(resp.Data.Templates) != 1 {
		t.Errorf("Unexpected detail response: %+v", resp.Data)
	}

	if resp.Data.Id != encId {
		t.Errorf("Expected Id %s, got %s", encId, resp.Data.Id)
	}

	decTemplateId, errDec := crypto.Decrypt(resp.Data.Templates[0].IdTemplate, testEncryptionKey)
	if errDec != nil || decTemplateId != "101" {
		t.Errorf("Expected decrypted template ID 101, got %s", decTemplateId)
	}
}

func TestHandler_GetDetailTindakanLab_InvalidEncryptedId(t *testing.T) {
	handler := tindakan.NewHandler(&mockService{}, testEncryptionKey)
	mux := http.NewServeMux()
	dummyMiddleware := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, dummyMiddleware, dummyMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tindakan/laboratorium/pk/invalid-token-here", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400 for invalid token, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_GetDetailTindakanLab_NotFound(t *testing.T) {
	mockSvc := &mockService{
		getDetailTindakanLabFn: func(ctx context.Context, kategori shared.KategoriLab, kodeTindakan string) (*tindakan.DetailTindakanLab, error) {
			return nil, apperror.NewNotFoundError("Tindakan laboratorium tidak ditemukan")
		},
	}

	handler := tindakan.NewHandler(mockSvc, testEncryptionKey)
	mux := http.NewServeMux()
	dummyMiddleware := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, dummyMiddleware, dummyMiddleware)

	encId, _ := crypto.Encrypt("PK999", testEncryptionKey)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tindakan/laboratorium/pk/"+encId, nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("Expected status 404, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_GetDaftarTindakanRadiologi_Success(t *testing.T) {
	mockSvc := &mockService{
		getDaftarTindakanRadiologiFn: func(ctx context.Context, filter tindakan.FilterDaftarTindakanRadiologi) ([]tindakan.TindakanRadiologi, shared.PaginationMeta, error) {
			return []tindakan.TindakanRadiologi{
				{
					KodeTindakan: "RAD001",
					NamaTindakan: "Rontgen Thorax AP/PA",
					Biaya:        125000,
				},
			}, shared.NewPaginationMeta(1, 1, 20), nil
		},
	}

	handler := tindakan.NewHandler(mockSvc, testEncryptionKey)
	mux := http.NewServeMux()
	dummyMiddleware := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, dummyMiddleware, dummyMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tindakan/radiologi?keyword=Thorax&page=1&limit=20", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp struct {
		Data []tindakan.TindakanRadiologi `json:"data"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(resp.Data) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(resp.Data))
	}
	decId, err := crypto.Decrypt(resp.Data[0].Id, testEncryptionKey)
	if err != nil || decId != "RAD001" {
		t.Errorf("Expected decrypted ID RAD001, got %s (err: %v)", decId, err)
	}
}

func TestHandler_GetDaftarTindakanRadiologi_InvalidFilter(t *testing.T) {
	handler := tindakan.NewHandler(&mockService{}, testEncryptionKey)
	mux := http.NewServeMux()
	dummyMiddleware := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, dummyMiddleware, dummyMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tindakan/radiologi?keyword=Th", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400 for keyword < 3 chars, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_GetDetailTindakanRadiologi_Success(t *testing.T) {
	mockSvc := &mockService{
		getDetailTindakanRadiologiFn: func(ctx context.Context, kodeTindakan string) (*tindakan.TindakanRadiologi, error) {
			return &tindakan.TindakanRadiologi{
				KodeTindakan: kodeTindakan,
				NamaTindakan: "Rontgen Thorax",
				Biaya:        125000,
			}, nil
		},
	}

	handler := tindakan.NewHandler(mockSvc, testEncryptionKey)
	mux := http.NewServeMux()
	dummyMiddleware := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, dummyMiddleware, dummyMiddleware)

	encId, _ := crypto.Encrypt("RAD001", testEncryptionKey)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tindakan/radiologi/"+encId, nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp struct {
		Data tindakan.TindakanRadiologi `json:"data"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if resp.Data.Id != encId {
		t.Errorf("Expected encrypted ID %s, got %s", encId, resp.Data.Id)
	}
	if resp.Data.KodeTindakan != "RAD001" {
		t.Errorf("Expected KodeTindakan RAD001, got %s", resp.Data.KodeTindakan)
	}
}

func TestHandler_GetDetailTindakanRadiologi_InvalidEncryptedId(t *testing.T) {
	handler := tindakan.NewHandler(&mockService{}, testEncryptionKey)
	mux := http.NewServeMux()
	dummyMiddleware := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, dummyMiddleware, dummyMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tindakan/radiologi/invalid-token-here", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400 for invalid token, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_GetDetailTindakanRadiologi_NotFound(t *testing.T) {
	mockSvc := &mockService{
		getDetailTindakanRadiologiFn: func(ctx context.Context, kodeTindakan string) (*tindakan.TindakanRadiologi, error) {
			return nil, apperror.NewNotFoundError("Tindakan radiologi tidak ditemukan")
		},
	}

	handler := tindakan.NewHandler(mockSvc, testEncryptionKey)
	mux := http.NewServeMux()
	dummyMiddleware := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, dummyMiddleware, dummyMiddleware)

	encId, _ := crypto.Encrypt("RAD999", testEncryptionKey)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tindakan/radiologi/"+encId, nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("Expected status 404, got %d: %s", rr.Code, rr.Body.String())
	}
}


