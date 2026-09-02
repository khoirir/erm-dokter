package berkasdigital_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"erm-dokter/internal/berkasdigital"
	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/pkg/logger"
)

type mockRepository struct {
	fetchMasterBerkasFn    func(ctx context.Context) ([]berkasdigital.MasterBerkasDigital, error)
	fetchBerkasByNoRawatFn func(ctx context.Context, noRawat string, kodeList []string) ([]berkasdigital.BerkasDigitalDB, error)
}

func (m *mockRepository) FetchMasterBerkas(ctx context.Context) ([]berkasdigital.MasterBerkasDigital, error) {
	if m.fetchMasterBerkasFn != nil {
		return m.fetchMasterBerkasFn(ctx)
	}
	return nil, nil
}

func (m *mockRepository) FetchBerkasByNoRawat(ctx context.Context, noRawat string, kodeList []string) ([]berkasdigital.BerkasDigitalDB, error) {
	if m.fetchBerkasByNoRawatFn != nil {
		return m.fetchBerkasByNoRawatFn(ctx, noRawat, kodeList)
	}
	return nil, nil
}

const testKey = "secret-key-32-bytes-testing-12345"
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

	svc := berkasdigital.NewService(mockRepo, testKey, testBaseURL, log)
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
		fetchBerkasByNoRawatFn: func(ctx context.Context, noRawat string, kodeList []string) ([]berkasdigital.BerkasDigitalDB, error) {
			return []berkasdigital.BerkasDigitalDB{
				{
					NoRawat:    "2026/04/22/000001",
					Kode:       "005",
					NamaBerkas: "HASIL LAB PK",
					LokasiFile: "pages/upload/pk.pdf",
				},
			}, nil
		},
	}

	svc := berkasdigital.NewService(mockRepo, testKey, testBaseURL, log)
	list, err := svc.GetBerkasByNoRawat(context.Background(), "2026/04/22/000001", []string{"005"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(list) != 1 || list[0].NamaBerkas != "HASIL LAB PK" {
		t.Errorf("Unexpected result: %+v", list)
	}
}

func TestService_BuildBerkasItem(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{}
	svc := berkasdigital.NewService(mockRepo, testKey, testBaseURL, log)

	item, err := svc.BuildBerkasItem("015", "PATOLOGI ANATOMI", "pages/upload/fnab.pdf")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if item == nil {
		t.Fatalf("Expected non-nil item")
	}
	if item.Kode != "015" || item.NamaBerkas != "PATOLOGI ANATOMI" {
		t.Errorf("Unexpected item data: %+v", item)
	}

	decryptedURL, errDec := crypto.Decrypt(item.IdBerkas, testKey)
	if errDec != nil || decryptedURL != "http://192.168.30.24/webapps/berkasrawat/pages/upload/fnab.pdf" {
		t.Errorf("Unexpected decrypted URL: %s, err: %v", decryptedURL, errDec)
	}
}

func TestService_StreamBerkasDigital_InvalidToken(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{}
	svc := berkasdigital.NewService(mockRepo, testKey, testBaseURL, log)

	rec := httptest.NewRecorder()
	err := svc.StreamBerkasDigital(context.Background(), "invalid-token", rec)
	if err == nil {
		t.Fatalf("Expected error for invalid token, got nil")
	}
}
