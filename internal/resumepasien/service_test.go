package resumepasien_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatinap"
	"erm-dokter/internal/rawatjalan"
	"erm-dokter/internal/resumepasien"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type mockRepository struct {
	detailData      *resumepasien.ResumePasien
	detailErr       error
	riwayatData     []resumepasien.ResumePasien
	riwayatErr      error
	adaResult       bool
	adaErr          error
	simpanErr       error
	updateErr       error
	hapusErr        error
	simpanCalledReq resumepasien.SimpanResumePasienRequest
	updateCalledReq resumepasien.UpdateResumePasienRequest
}

func (m *mockRepository) DetailResumePasien(ctx context.Context, noRawat string) (*resumepasien.ResumePasien, error) {
	if m.detailErr != nil {
		return nil, m.detailErr
	}
	if m.detailData != nil {
		return m.detailData, nil
	}
	return nil, sql.ErrNoRows
}

func (m *mockRepository) RiwayatResumePasienByNoRM(ctx context.Context, noRM string) ([]resumepasien.ResumePasien, error) {
	if m.riwayatErr != nil {
		return nil, m.riwayatErr
	}
	return m.riwayatData, nil
}

func (m *mockRepository) CekResumePasienAda(ctx context.Context, noRawat string) (bool, error) {
	if m.adaErr != nil {
		return false, m.adaErr
	}
	return m.adaResult, nil
}

func (m *mockRepository) SimpanResumePasien(ctx context.Context, noRawat, kodeDokter string, req resumepasien.SimpanResumePasienRequest) (*resumepasien.ResumePasien, error) {
	if m.simpanErr != nil {
		return nil, m.simpanErr
	}
	m.simpanCalledReq = req
	return &resumepasien.ResumePasien{
		NoRawat:          noRawat,
		KodeDokter:       kodeDokter,
		DataResumePasien: req.DataResumePasien,
	}, nil
}

func (m *mockRepository) UpdateResumePasien(ctx context.Context, noRawat string, req resumepasien.UpdateResumePasienRequest) (*resumepasien.ResumePasien, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	m.updateCalledReq = req
	return &resumepasien.ResumePasien{
		NoRawat:          noRawat,
		DataResumePasien: req.DataResumePasien,
	}, nil
}

func (m *mockRepository) HapusResumePasien(ctx context.Context, noRawat string) error {
	return m.hapusErr
}

type mockRawatJalanService struct {
	rawatjalan.Service
	tglReg string
	jamReg string
	exists bool
	err    error
}

func (m *mockRawatJalanService) GetWaktuRegistrasi(ctx context.Context, noRawat string) (string, string, bool, error) {
	if m.err != nil {
		return "", "", false, m.err
	}
	return m.tglReg, m.jamReg, m.exists, nil
}

type mockRawatInapService struct {
	rawatinap.Service
	isAktifRanap   bool
	hasRecordKamar bool
	err            error
}

func (m *mockRawatInapService) CekStatusKamarInap(ctx context.Context, noRawat string) (bool, bool, error) {
	if m.err != nil {
		return false, false, m.err
	}
	return m.isAktifRanap, m.hasRecordKamar, nil
}

func setupTestService(repo *mockRepository, ralanSvc *mockRawatJalanService, ranapSvc *mockRawatInapService) resumepasien.Service {
	log := logger.New()
	return resumepasien.NewService(repo, ralanSvc, ranapSvc, 48, log)
}

func TestService_DetailResumePasien(t *testing.T) {
	t.Run("Sukses", func(t *testing.T) {
		repo := &mockRepository{
			detailData: &resumepasien.ResumePasien{
				NoRawat:    "2026/09/07/000001",
				KodeDokter: "DR01",
				NamaDokter: "dr. Handi",
			},
		}
		svc := setupTestService(repo, &mockRawatJalanService{}, &mockRawatInapService{})

		res, err := svc.DetailResumePasien(context.Background(), "2026/09/07/000001", shared.StatusLanjutRawatJalan)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.NoRawat != "2026/09/07/000001" {
			t.Errorf("expected NoRawat '2026/09/07/000001', got '%s'", res.NoRawat)
		}
	})

	t.Run("Status Lanjut Bukan Ralan (Ditolak)", func(t *testing.T) {
		svc := setupTestService(&mockRepository{}, &mockRawatJalanService{}, &mockRawatInapService{})

		_, err := svc.DetailResumePasien(context.Background(), "2026/09/07/000001", shared.StatusLanjutRawatInap)
		if err == nil {
			t.Fatal("expected error for non-ralan status, got nil")
		}
		var bizErr *apperror.BusinessError
		if !errors.As(err, &bizErr) {
			t.Errorf("expected BusinessError, got %T: %v", err, err)
		}
	})

	t.Run("Tidak Ditemukan (404)", func(t *testing.T) {
		repo := &mockRepository{detailData: nil}
		svc := setupTestService(repo, &mockRawatJalanService{}, &mockRawatInapService{})

		_, err := svc.DetailResumePasien(context.Background(), "2026/09/07/999999", shared.StatusLanjutRawatJalan)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var notFound *apperror.NotFoundError
		if !errors.As(err, &notFound) {
			t.Errorf("expected NotFoundError, got %T: %v", err, err)
		}
	})
}

func TestService_RiwayatResumePasienByNoRM(t *testing.T) {
	repo := &mockRepository{
		riwayatData: []resumepasien.ResumePasien{
			{NoRawat: "RAWAT1"},
			{NoRawat: "RAWAT2"},
		},
	}
	svc := setupTestService(repo, &mockRawatJalanService{}, &mockRawatInapService{})

	list, err := svc.RiwayatResumePasienByNoRM(context.Background(), "123456", shared.StatusLanjutRawatJalan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 records, got %d", len(list))
	}
}

func TestService_SimpanResumePasien(t *testing.T) {
	now := time.Now()
	tglHariIni := now.Format("2006-01-02")
	jamSekarang := now.Format("15:04:05")

	validReq := resumepasien.SimpanResumePasienRequest{
		NoRawat: "2026/09/07/000001",
		DataResumePasien: resumepasien.DataResumePasien{
			KeluhanUtama:  "Batuk",
			DiagnosaUtama: "ISPA",
			KondisiPulang: resumepasien.KondisiPulangHidup,
		},
	}

	t.Run("Sukses", func(t *testing.T) {
		repo := &mockRepository{adaResult: false}
		ralanSvc := &mockRawatJalanService{
			tglReg: tglHariIni,
			jamReg: jamSekarang,
			exists: true,
		}
		ranapSvc := &mockRawatInapService{hasRecordKamar: false}
		svc := setupTestService(repo, ralanSvc, ranapSvc)

		res, err := svc.SimpanResumePasien(context.Background(), "2026/09/07/000001", "DR01", shared.StatusLanjutRawatJalan, validReq)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.NoRawat != "2026/09/07/000001" {
			t.Errorf("expected NoRawat '2026/09/07/000001', got '%s'", res.NoRawat)
		}
	})

	t.Run("Status Lanjut Bukan Ralan (Ditolak)", func(t *testing.T) {
		svc := setupTestService(&mockRepository{}, &mockRawatJalanService{}, &mockRawatInapService{})

		_, err := svc.SimpanResumePasien(context.Background(), "2026/09/07/000001", "DR01", shared.StatusLanjutRawatInap, validReq)
		if err == nil {
			t.Fatal("expected error for non-ralan status, got nil")
		}
		var bizErr *apperror.BusinessError
		if !errors.As(err, &bizErr) {
			t.Errorf("expected BusinessError, got: %v", err)
		}
	})

	t.Run("Kunjungan Tidak Ditemukan", func(t *testing.T) {
		repo := &mockRepository{}
		ralanSvc := &mockRawatJalanService{exists: false}
		svc := setupTestService(repo, ralanSvc, &mockRawatInapService{})

		_, err := svc.SimpanResumePasien(context.Background(), "2026/09/07/000001", "DR01", shared.StatusLanjutRawatJalan, validReq)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var notFound *apperror.NotFoundError
		if !errors.As(err, &notFound) {
			t.Errorf("expected NotFoundError, got: %v", err)
		}
	})

	t.Run("Duplikasi Data Kunjungan", func(t *testing.T) {
		repo := &mockRepository{adaResult: true}
		ralanSvc := &mockRawatJalanService{
			tglReg: tglHariIni,
			jamReg: jamSekarang,
			exists: true,
		}
		svc := setupTestService(repo, ralanSvc, &mockRawatInapService{})

		_, err := svc.SimpanResumePasien(context.Background(), "2026/09/07/000001", "DR01", shared.StatusLanjutRawatJalan, validReq)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var bizErr *apperror.BusinessError
		if !errors.As(err, &bizErr) {
			t.Errorf("expected BusinessError, got: %v", err)
		}
	})

	t.Run("Lewat Batas 48 Jam Rawat Jalan", func(t *testing.T) {
		repo := &mockRepository{adaResult: false}
		ralanSvc := &mockRawatJalanService{
			tglReg: "2020-01-01",
			jamReg: "08:00:00",
			exists: true,
		}
		svc := setupTestService(repo, ralanSvc, &mockRawatInapService{hasRecordKamar: false})

		_, err := svc.SimpanResumePasien(context.Background(), "2026/09/07/000001", "DR01", shared.StatusLanjutRawatJalan, validReq)
		if err == nil {
			t.Fatal("expected error for expired rawat jalan, got nil")
		}
		var forbbiden *apperror.ForbiddenError
		if !errors.As(err, &forbbiden) {
			t.Errorf("expected ForbiddenError, got: %v", err)
		}
	})

	t.Run("Pasien Sudah Terdaftar di Rawat Inap (Ditolak)", func(t *testing.T) {
		repo := &mockRepository{adaResult: false}
		ralanSvc := &mockRawatJalanService{
			tglReg: tglHariIni,
			jamReg: jamSekarang,
			exists: true,
		}
		ranapSvc := &mockRawatInapService{
			hasRecordKamar: true,
			isAktifRanap:   true,
		}
		svc := setupTestService(repo, ralanSvc, ranapSvc)

		_, err := svc.SimpanResumePasien(context.Background(), "2026/09/07/000001", "DR01", shared.StatusLanjutRawatJalan, validReq)
		if err == nil {
			t.Fatal("expected error for inpatient, got nil")
		}
		var bizErr *apperror.BusinessError
		if !errors.As(err, &bizErr) {
			t.Errorf("expected BusinessError, got: %v", err)
		}
	})
}

func TestService_UpdateResumePasien(t *testing.T) {
	now := time.Now()
	tglHariIni := now.Format("2006-01-02")
	jamSekarang := now.Format("15:04:05")

	validUpdateReq := resumepasien.UpdateResumePasienRequest{
		DataResumePasien: resumepasien.DataResumePasien{
			KeluhanUtama:  "Update keluhan",
			DiagnosaUtama: "Update diagnosa",
			KondisiPulang: resumepasien.KondisiPulangHidup,
		},
	}

	t.Run("Sukses", func(t *testing.T) {
		repo := &mockRepository{
			detailData: &resumepasien.ResumePasien{
				NoRawat:    "2026/09/07/000001",
				KodeDokter: "DR01",
			},
		}
		ralanSvc := &mockRawatJalanService{
			tglReg: tglHariIni,
			jamReg: jamSekarang,
			exists: true,
		}
		svc := setupTestService(repo, ralanSvc, &mockRawatInapService{})

		res, err := svc.UpdateResumePasien(context.Background(), "DR01", "2026/09/07/000001", shared.StatusLanjutRawatJalan, validUpdateReq)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.NoRawat != "2026/09/07/000001" {
			t.Errorf("expected NoRawat '2026/09/07/000001', got '%s'", res.NoRawat)
		}
	})

	t.Run("Proteksi Dokter Lain (403 Forbidden)", func(t *testing.T) {
		repo := &mockRepository{
			detailData: &resumepasien.ResumePasien{
				NoRawat:    "2026/09/07/000001",
				KodeDokter: "DR01",
			},
		}
		svc := setupTestService(repo, &mockRawatJalanService{}, &mockRawatInapService{})

		_, err := svc.UpdateResumePasien(context.Background(), "DR99", "2026/09/07/000001", shared.StatusLanjutRawatJalan, validUpdateReq)
		if err == nil {
			t.Fatal("expected ForbiddenError for different doctor, got nil")
		}
		var forbidden *apperror.ForbiddenError
		if !errors.As(err, &forbidden) {
			t.Errorf("expected ForbiddenError, got: %v", err)
		}
	})

	t.Run("Data Tidak Ditemukan", func(t *testing.T) {
		repo := &mockRepository{detailData: nil}
		svc := setupTestService(repo, &mockRawatJalanService{}, &mockRawatInapService{})

		_, err := svc.UpdateResumePasien(context.Background(), "DR01", "2026/09/07/000001", shared.StatusLanjutRawatJalan, validUpdateReq)
		if err == nil {
			t.Fatal("expected NotFoundError, got nil")
		}
		var notFound *apperror.NotFoundError
		if !errors.As(err, &notFound) {
			t.Errorf("expected NotFoundError, got: %v", err)
		}
	})
}

func TestService_HapusResumePasien(t *testing.T) {
	now := time.Now()
	tglHariIni := now.Format("2006-01-02")
	jamSekarang := now.Format("15:04:05")

	t.Run("Sukses", func(t *testing.T) {
		repo := &mockRepository{
			detailData: &resumepasien.ResumePasien{
				NoRawat:    "2026/09/07/000001",
				KodeDokter: "DR01",
			},
		}
		ralanSvc := &mockRawatJalanService{
			tglReg: tglHariIni,
			jamReg: jamSekarang,
			exists: true,
		}
		svc := setupTestService(repo, ralanSvc, &mockRawatInapService{})

		err := svc.HapusResumePasien(context.Background(), "DR01", "2026/09/07/000001", shared.StatusLanjutRawatJalan)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("Proteksi Dokter Lain (403 Forbidden)", func(t *testing.T) {
		repo := &mockRepository{
			detailData: &resumepasien.ResumePasien{
				NoRawat:    "2026/09/07/000001",
				KodeDokter: "DR01",
			},
		}
		svc := setupTestService(repo, &mockRawatJalanService{}, &mockRawatInapService{})

		err := svc.HapusResumePasien(context.Background(), "DR99", "2026/09/07/000001", shared.StatusLanjutRawatJalan)
		if err == nil {
			t.Fatal("expected ForbiddenError for different doctor, got nil")
		}
		var forbidden *apperror.ForbiddenError
		if !errors.As(err, &forbidden) {
			t.Errorf("expected ForbiddenError, got: %v", err)
		}
	})
}
