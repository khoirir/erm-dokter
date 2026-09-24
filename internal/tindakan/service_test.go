package tindakan_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
	"erm-dokter/internal/tindakan"
)

type mockRepository struct {
	daftarTindakanLabFn        func(ctx context.Context, kategori shared.KategoriLab, filter tindakan.FilterDaftarTindakanLab) ([]tindakan.TindakanLab, int, error)
	getDetailTindakanLabFn     func(ctx context.Context, kategori shared.KategoriLab, kodeTindakan string) (*tindakan.TindakanLab, []tindakan.TemplateLab, error)
	cekKeberadaanTindakanLabFn func(ctx context.Context, kategori shared.KategoriLab, listKodeTindakan []string) (map[string]bool, error)
	cekKeberadaanTemplateLabFn func(ctx context.Context, listKodeTindakan []string, templateMap map[string][]int) (map[string]map[int]bool, error)

	daftarTindakanRadiologiFn        func(ctx context.Context, filter tindakan.FilterDaftarTindakanRadiologi) ([]tindakan.TindakanRadiologi, int, error)
	getDetailTindakanRadiologiFn     func(ctx context.Context, kodeTindakan string) (*tindakan.TindakanRadiologi, error)
	cekKeberadaanTindakanRadiologiFn func(ctx context.Context, listKodeTindakan []string) (map[string]bool, error)
}

func (m *mockRepository) DaftarTindakanLab(ctx context.Context, kategori shared.KategoriLab, filter tindakan.FilterDaftarTindakanLab) ([]tindakan.TindakanLab, int, error) {
	if m.daftarTindakanLabFn != nil {
		return m.daftarTindakanLabFn(ctx, kategori, filter)
	}
	return nil, 0, nil
}

func (m *mockRepository) GetDetailTindakanLab(ctx context.Context, kategori shared.KategoriLab, kodeTindakan string) (*tindakan.TindakanLab, []tindakan.TemplateLab, error) {
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

func (m *mockRepository) DaftarTindakanRadiologi(ctx context.Context, filter tindakan.FilterDaftarTindakanRadiologi) ([]tindakan.TindakanRadiologi, int, error) {
	if m.daftarTindakanRadiologiFn != nil {
		return m.daftarTindakanRadiologiFn(ctx, filter)
	}
	return nil, 0, nil
}

func (m *mockRepository) GetDetailTindakanRadiologi(ctx context.Context, kodeTindakan string) (*tindakan.TindakanRadiologi, error) {
	if m.getDetailTindakanRadiologiFn != nil {
		return m.getDetailTindakanRadiologiFn(ctx, kodeTindakan)
	}
	return nil, nil
}

func (m *mockRepository) CekKeberadaanTindakanRadiologi(ctx context.Context, listKodeTindakan []string) (map[string]bool, error) {
	if m.cekKeberadaanTindakanRadiologiFn != nil {
		return m.cekKeberadaanTindakanRadiologiFn(ctx, listKodeTindakan)
	}
	return make(map[string]bool), nil
}

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

	svc := tindakan.NewService(mockRepo, log)
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

	if list[0].KodeTindakan != "PK001" {
		t.Errorf("Expected KodeTindakan PK001, got %s", list[0].KodeTindakan)
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

	svc := tindakan.NewService(mockRepo, log)
	_, _, err := svc.GetDaftarTindakanLab(context.Background(), shared.KategoriLabPK, tindakan.FilterDaftarTindakanLab{})

	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestService_GetDetailTindakanLab_Success(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		getDetailTindakanLabFn: func(ctx context.Context, kategori shared.KategoriLab, kodeTindakan string) (*tindakan.TindakanLab, []tindakan.TemplateLab, error) {
			if kodeTindakan != "PK001" {
				t.Errorf("Expected kodeTindakan PK001, got %s", kodeTindakan)
			}
			return &tindakan.TindakanLab{
					KodeTindakan: "PK001",
					NamaTindakan: "DARAH LENGKAP",
					Biaya:        75000,
				}, []tindakan.TemplateLab{
					{
						IdTemplate:      "101",
						NamaPemeriksaan: "Hemoglobin",
						Satuan:          "g/dL",
						NilaiRujukanLD:  "13.5 - 17.5",
						NilaiRujukanLA:  "11.5 - 15.5",
						NilaiRujukanPD:  "12.0 - 16.0",
						NilaiRujukanPA:  "11.0 - 15.0",
					},
					{
						IdTemplate:      "102",
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

	svc := tindakan.NewService(mockRepo, log)
	detail, err := svc.GetDetailTindakanLab(context.Background(), shared.KategoriLabPK, "PK001")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if detail.KodeTindakan != "PK001" {
		t.Errorf("Expected KodeTindakan PK001, got %s", detail.KodeTindakan)
	}

	if len(detail.Templates) != 2 {
		t.Fatalf("Expected 2 templates, got %d", len(detail.Templates))
	}

	if detail.Templates[0].IdTemplate != "101" {
		t.Errorf("Expected template ID 101, got %s", detail.Templates[0].IdTemplate)
	}
}

func TestService_GetDetailTindakanLab_NotFound(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		getDetailTindakanLabFn: func(ctx context.Context, kategori shared.KategoriLab, kodeTindakan string) (*tindakan.TindakanLab, []tindakan.TemplateLab, error) {
			return nil, nil, sql.ErrNoRows
		},
	}

	svc := tindakan.NewService(mockRepo, log)
	_, err := svc.GetDetailTindakanLab(context.Background(), shared.KategoriLabPK, "PK999")
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
		getDetailTindakanLabFn: func(ctx context.Context, kategori shared.KategoriLab, kodeTindakan string) (*tindakan.TindakanLab, []tindakan.TemplateLab, error) {
			return nil, nil, errors.New("db error")
		},
	}

	svc := tindakan.NewService(mockRepo, log)
	_, err := svc.GetDetailTindakanLab(context.Background(), shared.KategoriLabPK, "PK001")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestService_CekKeberadaanTindakanLab(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		cekKeberadaanTindakanLabFn: func(ctx context.Context, kategori shared.KategoriLab, listKodeTindakan []string) (map[string]bool, error) {
			return map[string]bool{"PK001": true}, nil
		},
	}

	svc := tindakan.NewService(mockRepo, log)
	res, err := svc.CekKeberadaanTindakanLab(context.Background(), shared.KategoriLabPK, []string{"PK001"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !res["PK001"] {
		t.Errorf("Expected PK001 to be true")
	}
}

func TestService_CekKeberadaanTemplateLab(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		cekKeberadaanTemplateLabFn: func(ctx context.Context, listKodeTindakan []string, templateMap map[string][]int) (map[string]map[int]bool, error) {
			return map[string]map[int]bool{"PK001": {101: true}}, nil
		},
	}

	svc := tindakan.NewService(mockRepo, log)
	res, err := svc.CekKeberadaanTemplateLab(context.Background(), []string{"PK001"}, map[string][]int{"PK001": {101}})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !res["PK001"][101] {
		t.Errorf("Expected PK001 template 101 to be true")
	}
}

func TestService_GetDaftarTindakanRadiologi_Success(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		daftarTindakanRadiologiFn: func(ctx context.Context, filter tindakan.FilterDaftarTindakanRadiologi) ([]tindakan.TindakanRadiologi, int, error) {
			return []tindakan.TindakanRadiologi{
				{
					KodeTindakan: "RAD001",
					NamaTindakan: "RONTGEN THORAX AP/PA",
					Biaya:        125000,
				},
			}, 1, nil
		},
	}

	svc := tindakan.NewService(mockRepo, log)
	list, meta, err := svc.GetDaftarTindakanRadiologi(context.Background(), tindakan.FilterDaftarTindakanRadiologi{Page: 1, Limit: 20})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(list))
	}
	if list[0].KodeTindakan != "RAD001" {
		t.Errorf("Expected kode tindakan RAD001, got %s", list[0].KodeTindakan)
	}
	if meta.TotalRecords != 1 {
		t.Errorf("Expected total records 1, got %d", meta.TotalRecords)
	}
}

func TestService_GetDaftarTindakanRadiologi_Error(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		daftarTindakanRadiologiFn: func(ctx context.Context, filter tindakan.FilterDaftarTindakanRadiologi) ([]tindakan.TindakanRadiologi, int, error) {
			return nil, 0, errors.New("db error")
		},
	}

	svc := tindakan.NewService(mockRepo, log)
	_, _, err := svc.GetDaftarTindakanRadiologi(context.Background(), tindakan.FilterDaftarTindakanRadiologi{Page: 1, Limit: 20})
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestService_GetDetailTindakanRadiologi_Success(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		getDetailTindakanRadiologiFn: func(ctx context.Context, kodeTindakan string) (*tindakan.TindakanRadiologi, error) {
			return &tindakan.TindakanRadiologi{
				KodeTindakan: kodeTindakan,
				NamaTindakan: "RONTGEN THORAX",
				Biaya:        125000,
			}, nil
		},
	}

	svc := tindakan.NewService(mockRepo, log)
	res, err := svc.GetDetailTindakanRadiologi(context.Background(), "RAD001")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if res.KodeTindakan != "RAD001" {
		t.Errorf("Expected RAD001, got %s", res.KodeTindakan)
	}
}

func TestService_GetDetailTindakanRadiologi_NotFound(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		getDetailTindakanRadiologiFn: func(ctx context.Context, kodeTindakan string) (*tindakan.TindakanRadiologi, error) {
			return nil, sql.ErrNoRows
		},
	}

	svc := tindakan.NewService(mockRepo, log)
	_, err := svc.GetDetailTindakanRadiologi(context.Background(), "RAD999")
	if err == nil {
		t.Fatalf("Expected not found error, got nil")
	}
	var notFoundErr *apperror.NotFoundError
	if !errors.As(err, &notFoundErr) {
		t.Errorf("Expected NotFoundError, got %T: %v", err, err)
	}
}

func TestService_GetDetailTindakanRadiologi_Error(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		getDetailTindakanRadiologiFn: func(ctx context.Context, kodeTindakan string) (*tindakan.TindakanRadiologi, error) {
			return nil, errors.New("db error")
		},
	}

	svc := tindakan.NewService(mockRepo, log)
	_, err := svc.GetDetailTindakanRadiologi(context.Background(), "RAD001")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestService_CekKeberadaanTindakanRadiologi(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		cekKeberadaanTindakanRadiologiFn: func(ctx context.Context, listKodeTindakan []string) (map[string]bool, error) {
			return map[string]bool{"RAD001": true}, nil
		},
	}

	svc := tindakan.NewService(mockRepo, log)
	res, err := svc.CekKeberadaanTindakanRadiologi(context.Background(), []string{"RAD001"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !res["RAD001"] {
		t.Errorf("Expected RAD001 to be true")
	}
}
