package berkasdigital_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"erm-dokter/internal/berkasdigital"
	"erm-dokter/internal/shared/apperror"
)

type mockService struct {
	getMasterBerkasFn     func(ctx context.Context) ([]berkasdigital.MasterBerkasDigital, error)
	getBerkasByNoRawatFn  func(ctx context.Context, noRawat string, kodeList []string) ([]berkasdigital.BerkasDigitalDB, error)
	streamBerkasDigitalFn func(ctx context.Context, encryptedIdBerkas string, w http.ResponseWriter) error
	buildBerkasItemFn     func(kode string, namaBerkas string, lokasiFile string) (*berkasdigital.BerkasDigitalItem, error)
	buildFullURLFn        func(lokasiFile string) string
}

func (m *mockService) GetMasterBerkas(ctx context.Context) ([]berkasdigital.MasterBerkasDigital, error) {
	if m.getMasterBerkasFn != nil {
		return m.getMasterBerkasFn(ctx)
	}
	return nil, nil
}

func (m *mockService) GetBerkasByNoRawat(ctx context.Context, noRawat string, kodeList []string) ([]berkasdigital.BerkasDigitalDB, error) {
	if m.getBerkasByNoRawatFn != nil {
		return m.getBerkasByNoRawatFn(ctx, noRawat, kodeList)
	}
	return nil, nil
}

func (m *mockService) StreamBerkasDigital(ctx context.Context, encryptedIdBerkas string, w http.ResponseWriter) error {
	if m.streamBerkasDigitalFn != nil {
		return m.streamBerkasDigitalFn(ctx, encryptedIdBerkas, w)
	}
	return nil
}

func (m *mockService) BuildBerkasItem(kode string, namaBerkas string, lokasiFile string) (*berkasdigital.BerkasDigitalItem, error) {
	if m.buildBerkasItemFn != nil {
		return m.buildBerkasItemFn(kode, namaBerkas, lokasiFile)
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

	handler := berkasdigital.NewHandler(mockSvc)
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

func TestHandler_StreamBerkasDigital_Error(t *testing.T) {
	mockSvc := &mockService{
		streamBerkasDigitalFn: func(ctx context.Context, encryptedIdBerkas string, w http.ResponseWriter) error {
			return apperror.NewBusinessError("Token berkas digital tidak valid")
		},
	}

	handler := berkasdigital.NewHandler(mockSvc)
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
