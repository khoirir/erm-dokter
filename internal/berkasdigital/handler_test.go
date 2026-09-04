package berkasdigital_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"erm-dokter/internal/berkasdigital"
	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/shared/apperror"
)

const testEncryptionKey = "secret-key-32-bytes-testing-12345"

type mockService struct {
	getMasterBerkasFn    func(ctx context.Context) ([]berkasdigital.MasterBerkasDigital, error)
	getBerkasByNoRawatFn func(ctx context.Context, noRawat string, kodeList []string) ([]berkasdigital.BerkasDigitalPerawatan, error)
	getBerkasStreamFn    func(ctx context.Context, targetURL string) (*berkasdigital.BerkasStream, error)
	buildFullURLFn       func(lokasiFile string) string
}

func (m *mockService) GetMasterBerkas(ctx context.Context) ([]berkasdigital.MasterBerkasDigital, error) {
	if m.getMasterBerkasFn != nil {
		return m.getMasterBerkasFn(ctx)
	}
	return nil, nil
}

func (m *mockService) GetBerkasByNoRawat(ctx context.Context, noRawat string, kodeList []string) ([]berkasdigital.BerkasDigitalPerawatan, error) {
	if m.getBerkasByNoRawatFn != nil {
		return m.getBerkasByNoRawatFn(ctx, noRawat, kodeList)
	}
	return nil, nil
}

func (m *mockService) GetBerkasStream(ctx context.Context, targetURL string) (*berkasdigital.BerkasStream, error) {
	if m.getBerkasStreamFn != nil {
		return m.getBerkasStreamFn(ctx, targetURL)
	}
	return nil, nil
}

func (m *mockService) BuildFullURL(lokasiFile string) string {
	if m.buildFullURLFn != nil {
		return m.buildFullURLFn(lokasiFile)
	}
	return ""
}

func TestHandler_MasterBerkasDigital_Success(t *testing.T) {
	mockSvc := &mockService{
		getMasterBerkasFn: func(ctx context.Context) ([]berkasdigital.MasterBerkasDigital, error) {
			return []berkasdigital.MasterBerkasDigital{
				{Kode: "001", Nama: "Berkas SEP"},
				{Kode: "015", Nama: "PATOLOGI ANATOMI"},
			}, nil
		},
	}

	handler := berkasdigital.NewHandler(mockSvc, testEncryptionKey)
	mux := http.NewServeMux()
	dummyMiddleware := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, dummyMiddleware, dummyMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/master/berkas-digital", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp struct {
		Code    int                                 `json:"code"`
		Status  string                              `json:"status"`
		Message string                              `json:"message"`
		Data    []berkasdigital.MasterBerkasDigital `json:"data"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(resp.Data) != 2 {
		t.Errorf("Expected 2 items, got %d", len(resp.Data))
	}
}

func TestHandler_StreamBerkasDigital_InvalidToken(t *testing.T) {
	mockSvc := &mockService{}

	handler := berkasdigital.NewHandler(mockSvc, testEncryptionKey)
	mux := http.NewServeMux()
	dummyMiddleware := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, dummyMiddleware, dummyMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/berkas-digital/invalid-token", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandler_StreamBerkasDigital_Success(t *testing.T) {
	targetURL := "http://192.168.30.24/webapps/berkasrawat/pages/upload/fnab.pdf"
	encToken, err := crypto.Encrypt(targetURL, testEncryptionKey)
	if err != nil {
		t.Fatalf("Failed to encrypt token: %v", err)
	}

	dummyContent := []byte("%PDF-1.4 dummy pdf content")
	mockSvc := &mockService{
		getBerkasStreamFn: func(ctx context.Context, url string) (*berkasdigital.BerkasStream, error) {
			if url != targetURL {
				return nil, apperror.NewNotFoundError("Berkas tidak ditemukan")
			}
			return &berkasdigital.BerkasStream{
				Body:          io.NopCloser(bytes.NewReader(dummyContent)),
				ContentType:   "application/pdf",
				ContentLength: "26",
				Filename:      "fnab.pdf",
			}, nil
		},
	}

	handler := berkasdigital.NewHandler(mockSvc, testEncryptionKey)
	mux := http.NewServeMux()
	dummyMiddleware := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, dummyMiddleware, dummyMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/berkas-digital/"+encToken, nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}

	if ct := rr.Header().Get("Content-Type"); ct != "application/pdf" {
		t.Errorf("Expected Content-Type application/pdf, got %s", ct)
	}

	if rr.Body.String() != string(dummyContent) {
		t.Errorf("Unexpected body content: %s", rr.Body.String())
	}
}

