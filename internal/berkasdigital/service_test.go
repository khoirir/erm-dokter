package berkasdigital_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"erm-dokter/internal/berkasdigital"
	"erm-dokter/internal/pkg/logger"
)

type mockRepository struct {
	fetchMasterBerkasFn    func(ctx context.Context) ([]berkasdigital.MasterBerkasDigital, error)
	fetchBerkasByNoRawatFn func(ctx context.Context, noRawat string, kodeList []string) ([]berkasdigital.BerkasDigitalPerawatan, error)
}

func (m *mockRepository) FetchMasterBerkas(ctx context.Context) ([]berkasdigital.MasterBerkasDigital, error) {
	if m.fetchMasterBerkasFn != nil {
		return m.fetchMasterBerkasFn(ctx)
	}
	return nil, nil
}

func (m *mockRepository) FetchBerkasByNoRawat(ctx context.Context, noRawat string, kodeList []string) ([]berkasdigital.BerkasDigitalPerawatan, error) {
	if m.fetchBerkasByNoRawatFn != nil {
		return m.fetchBerkasByNoRawatFn(ctx, noRawat, kodeList)
	}
	return nil, nil
}

const testBaseURL = "http://192.168.30.24/webapps/berkasrawat/"

func TestService_GetMasterBerkas(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		fetchMasterBerkasFn: func(ctx context.Context) ([]berkasdigital.MasterBerkasDigital, error) {
			return []berkasdigital.MasterBerkasDigital{
				{Kode: "001", Nama: "Berkas SEP"},
				{Kode: "015", Nama: "PATOLOGI ANATOMI"},
			}, nil
		},
	}

	svc := berkasdigital.NewService(mockRepo, testBaseURL, log)
	list, err := svc.GetMasterBerkas(context.Background())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("Expected 2 items, got %d", len(list))
	}
}

func TestService_GetBerkasByNoRawat(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		fetchBerkasByNoRawatFn: func(ctx context.Context, noRawat string, kodeList []string) ([]berkasdigital.BerkasDigitalPerawatan, error) {
			return []berkasdigital.BerkasDigitalPerawatan{
				{
					NoRawat:    "2026/04/22/000001",
					Kode:       "005",
					NamaBerkas: "HASIL LAB PK",
					LokasiFile: "pages/upload/pk.pdf",
				},
			}, nil
		},
	}

	svc := berkasdigital.NewService(mockRepo, testBaseURL, log)
	list, err := svc.GetBerkasByNoRawat(context.Background(), "2026/04/22/000001", []string{"005"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(list) != 1 || list[0].NamaBerkas != "HASIL LAB PK" {
		t.Errorf("Unexpected result: %+v", list)
	}
}

func TestService_BuildFullURL(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{}
	svc := berkasdigital.NewService(mockRepo, testBaseURL, log)

	url1 := svc.BuildFullURL("pages/upload/fnab.pdf")
	expected1 := "http://192.168.30.24/webapps/berkasrawat/pages/upload/fnab.pdf"
	if url1 != expected1 {
		t.Errorf("Expected %s, got %s", expected1, url1)
	}

	url2 := svc.BuildFullURL("https://storage.rs.com/file.pdf")
	if url2 != "https://storage.rs.com/file.pdf" {
		t.Errorf("Expected external URL preserved, got %s", url2)
	}
}

func TestService_GetBerkasStream_Success(t *testing.T) {
	dummyContent := []byte("%PDF-1.4 dummy pdf content")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/pdf")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(dummyContent)
	}))
	defer ts.Close()

	log := logger.New()
	mockRepo := &mockRepository{}
	svc := berkasdigital.NewService(mockRepo, ts.URL, log)

	stream, err := svc.GetBerkasStream(context.Background(), ts.URL+"/test.pdf")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer stream.Body.Close()

	if stream.ContentType != "application/pdf" {
		t.Errorf("Expected Content-Type application/pdf, got %s", stream.ContentType)
	}

	bodyBytes, errRead := io.ReadAll(stream.Body)
	if errRead != nil {
		t.Fatalf("Failed to read stream body: %v", errRead)
	}
	if string(bodyBytes) != string(dummyContent) {
		t.Errorf("Unexpected body content: %s", string(bodyBytes))
	}
}

func TestService_GetBerkasStream_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer ts.Close()

	log := logger.New()
	mockRepo := &mockRepository{}
	svc := berkasdigital.NewService(mockRepo, ts.URL, log)

	_, err := svc.GetBerkasStream(context.Background(), ts.URL+"/missing.pdf")
	if err == nil {
		t.Fatalf("Expected error for 404 response, got nil")
	}
}

