package tindakan_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
	"erm-dokter/internal/tindakan"
)

type mockRepository struct {
	daftarTindakanLabFn        func(ctx context.Context, kategori shared.KategoriLab, filter tindakan.FilterDaftarTindakanLab) ([]tindakan.TindakanLab, int, error)
	getDetailTindakanLabFn     func(ctx context.Context, kategori shared.KategoriLab, kodeTindakan string) (*tindakan.TindakanLab, []tindakan.TemplateLabDB, error)
	cekKeberadaanTindakanLabFn func(ctx context.Context, kategori shared.KategoriLab, listKodeTindakan []string) (map[string]bool, error)
	cekKeberadaanTemplateLabFn func(ctx context.Context, listKodeTindakan []string, templateMap map[string][]int) (map[string]map[int]bool, error)
}

func (m *mockRepository) DaftarTindakanLab(ctx context.Context, kategori shared.KategoriLab, filter tindakan.FilterDaftarTindakanLab) ([]tindakan.TindakanLab, int, error) {
	if m.daftarTindakanLabFn != nil {
		return m.daftarTindakanLabFn(ctx, kategori, filter)
	}
	return nil, 0, nil
}

func (m *mockRepository) GetDetailTindakanLab(ctx context.Context, kategori shared.KategoriLab, kodeTindakan string) (*tindakan.TindakanLab, []tindakan.TemplateLabDB, error) {
	if m.getDetailTindakanLabFn != nil {
		return m.getDetailTindakanLabFn(ctx, kategori, kodeTindakan)
	}
	return nil, nil, nil
}

func (m *mockRepository) CekKeberadaanTindakanLab(ctx context.Context, kategori shared.KategoriLab, listKodeTindakan []string) (map[string]bool, error) {
	if m.cekKeberadaanTindakanLabFn != nil {
		return m.cekKeberadaanTindakanLabFn(ctx, kategori, listKodeTindakan)
	}
	return make(map[string]bool), nil
}

func (m *mockRepository) CekKeberadaanTemplateLab(ctx context.Context, listKodeTindakan []string, templateMap map[string][]int) (map[string]map[int]bool, error) {
	if m.cekKeberadaanTemplateLabFn != nil {
		return m.cekKeberadaanTemplateLabFn(ctx, listKodeTindakan, templateMap)
	}
	return make(map[string]map[int]bool), nil
}

const testJWTSecret = "secret-key-32-bytes-testing-12345"

func TestService_GetDaftarTindakanLab_Success(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		daftarTindakanLabFn: func(ctx context.Context, kategori shared.KategoriLab, filter tindakan.FilterDaftarTindakanLab) ([]tindakan.TindakanLab, int, error) {
			if kategori != shared.KategoriLabPK {
				t.Errorf("Expected category PK, got %s", kategori)
			}
			return []tindakan.TindakanLab{
				{
					KodeTindakan: "PK001",
					NamaTindakan: "DARAH LENGKAP",
					Biaya:        75000,
				},
			}, 1, nil
		},
	}

	svc := tindakan.NewService(mockRepo, testJWTSecret, log)
	list, meta, err := svc.GetDaftarTindakanLab(context.Background(), shared.KategoriLabPK, tindakan.FilterDaftarTindakanLab{
		Page:  1,
		Limit: 20,
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(list))
	}

	if list[0].Id == "" {
		t.Errorf("Expected encrypted Id, got empty string")
	}

	decryptedId, errDec := crypto.Decrypt(list[0].Id, testJWTSecret)
	if errDec != nil || decryptedId != "PK001" {
		t.Errorf("Failed to decrypt Id or mismatch value: %s, err: %v", decryptedId, errDec)
	}

	if meta.TotalRecords != 1 {
		t.Errorf("Expected total records 1, got %d", meta.TotalRecords)
	}
}

func TestService_GetDaftarTindakanLab_RepoError(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		daftarTindakanLabFn: func(ctx context.Context, kategori shared.KategoriLab, filter tindakan.FilterDaftarTindakanLab) ([]tindakan.TindakanLab, int, error) {
			return nil, 0, errors.New("db error")
		},
	}

	svc := tindakan.NewService(mockRepo, testJWTSecret, log)
	_, _, err := svc.GetDaftarTindakanLab(context.Background(), shared.KategoriLabPK, tindakan.FilterDaftarTindakanLab{})

	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestService_GetDetailTindakanLab_Success(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		getDetailTindakanLabFn: func(ctx context.Context, kategori shared.KategoriLab, kodeTindakan string) (*tindakan.TindakanLab, []tindakan.TemplateLabDB, error) {
			if kodeTindakan != "PK001" {
				t.Errorf("Expected kodeTindakan PK001, got %s", kodeTindakan)
			}
			return &tindakan.TindakanLab{
					KodeTindakan: "PK001",
					NamaTindakan: "DARAH LENGKAP",
					Biaya:        75000,
				}, []tindakan.TemplateLabDB{
					{
						IdTemplate:      101,
						KodeTindakan:    "PK001",
						NamaPemeriksaan: "Hemoglobin",
						Satuan:          "g/dL",
						NilaiRujukanLD:  "13.5 - 17.5",
						NilaiRujukanLA:  "11.5 - 15.5",
						NilaiRujukanPD:  "12.0 - 16.0",
						NilaiRujukanPA:  "11.0 - 15.0",
					},
					{
						IdTemplate:      102,
						KodeTindakan:    "PK001",
						NamaPemeriksaan: "Leukosit",
						Satuan:          "/uL",
						NilaiRujukanLD:  "4.000 - 10.000",
						NilaiRujukanLA:  "4.000 - 10.000",
						NilaiRujukanPD:  "4.000 - 10.000",
						NilaiRujukanPA:  "4.000 - 10.000",
					},
				}, nil
		},
	}

	svc := tindakan.NewService(mockRepo, testJWTSecret, log)
	encryptedId, _ := crypto.Encrypt("PK001", testJWTSecret)

	detail, err := svc.GetDetailTindakanLab(context.Background(), shared.KategoriLabPK, encryptedId)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if detail.KodeTindakan != "PK001" {
		t.Errorf("Expected KodeTindakan PK001, got %s", detail.KodeTindakan)
	}

	if len(detail.Templates) != 2 {
		t.Fatalf("Expected 2 templates, got %d", len(detail.Templates))
	}

	decTemplateId, errDec := crypto.Decrypt(detail.Templates[0].IdTemplate, testJWTSecret)
	if errDec != nil || decTemplateId != "101" {
		t.Errorf("Expected decrypted template ID 101, got %s", decTemplateId)
	}
}

func TestService_GetDetailTindakanLab_InvalidEncryptedId(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{}

	svc := tindakan.NewService(mockRepo, testJWTSecret, log)
	_, err := svc.GetDetailTindakanLab(context.Background(), shared.KategoriLabPK, "invalid-token")

	if err == nil {
		t.Fatalf("Expected error for invalid token, got nil")
	}

	var businessErr *apperror.BusinessError
	if !errors.As(err, &businessErr) {
		t.Errorf("Expected BusinessError, got %T: %v", err, err)
	}
}

func TestService_GetDetailTindakanLab_NotFound(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		getDetailTindakanLabFn: func(ctx context.Context, kategori shared.KategoriLab, kodeTindakan string) (*tindakan.TindakanLab, []tindakan.TemplateLabDB, error) {
			return nil, nil, sql.ErrNoRows
		},
	}

	svc := tindakan.NewService(mockRepo, testJWTSecret, log)
	encryptedId, _ := crypto.Encrypt("PK999", testJWTSecret)

	_, err := svc.GetDetailTindakanLab(context.Background(), shared.KategoriLabPK, encryptedId)
	if err == nil {
		t.Fatalf("Expected error for not found, got nil")
	}

	var notFoundErr *apperror.NotFoundError
	if !errors.As(err, &notFoundErr) {
		t.Errorf("Expected NotFoundError, got %T: %v", err, err)
	}
}

func TestService_GetDetailTindakanLab_RepoError(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		getDetailTindakanLabFn: func(ctx context.Context, kategori shared.KategoriLab, kodeTindakan string) (*tindakan.TindakanLab, []tindakan.TemplateLabDB, error) {
			return nil, nil, errors.New("db error")
		},
	}

	svc := tindakan.NewService(mockRepo, testJWTSecret, log)
	encryptedId, _ := crypto.Encrypt("PK001", testJWTSecret)

	_, err := svc.GetDetailTindakanLab(context.Background(), shared.KategoriLabPK, encryptedId)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}
