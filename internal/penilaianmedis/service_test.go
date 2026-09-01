package penilaianmedis

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"
	"time"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatjalan"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type mockRepository struct {
	// Ralan
	detailData      *PenilaianMedisRalan
	detailErr       error
	riwayatData     []PenilaianMedisRalan
	riwayatErr      error
	adaResult       bool
	adaErr          error
	simpanErr       error
	updateErr       error
	hapusErr        error
	simpanCalledReq SimpanPenilaianMedisRalanRequest
	updateCalledReq UpdatePenilaianMedisRalanRequest
	hapusCalledNo   string

	// IGD
	detailIGDData      *PenilaianMedisIGD
	detailIGDErr       error
	riwayatIGDData     []PenilaianMedisIGD
	riwayatIGDErr      error
	adaIGDResult       bool
	adaIGDErr          error
	simpanIGDErr       error
	updateIGDErr       error
	hapusIGDErr        error
	simpanIGDCalledReq SimpanPenilaianMedisIGDRequest
	updateIGDCalledReq UpdatePenilaianMedisIGDRequest
	hapusIGDCalledNo   string
}

// Ralan Mock Methods
func (m *mockRepository) DetailPenilaianMedisRalan(ctx context.Context, noRawat string) (*PenilaianMedisRalan, error) {
	return m.detailData, m.detailErr
}

func (m *mockRepository) RiwayatPenilaianMedisRalanByNoRM(ctx context.Context, noRM string) ([]PenilaianMedisRalan, error) {
	return m.riwayatData, m.riwayatErr
}

func (m *mockRepository) CekPenilaianMedisRalanAda(ctx context.Context, noRawat string) (bool, error) {
	return m.adaResult, m.adaErr
}

func (m *mockRepository) SimpanPenilaianMedisRalan(ctx context.Context, noRawat, kodeDokter string, req SimpanPenilaianMedisRalanRequest) error {
	m.simpanCalledReq = req
	return m.simpanErr
}

func (m *mockRepository) UpdatePenilaianMedisRalan(ctx context.Context, noRawat string, req UpdatePenilaianMedisRalanRequest) error {
	m.updateCalledReq = req
	return m.updateErr
}

func (m *mockRepository) HapusPenilaianMedisRalan(ctx context.Context, noRawat string) error {
	m.hapusCalledNo = noRawat
	return m.hapusErr
}

// IGD Mock Methods
func (m *mockRepository) DetailPenilaianMedisIGD(ctx context.Context, noRawat string) (*PenilaianMedisIGD, error) {
	return m.detailIGDData, m.detailIGDErr
}

func (m *mockRepository) RiwayatPenilaianMedisIGDByNoRM(ctx context.Context, noRM string) ([]PenilaianMedisIGD, error) {
	return m.riwayatIGDData, m.riwayatIGDErr
}

func (m *mockRepository) CekPenilaianMedisIGDAda(ctx context.Context, noRawat string) (bool, error) {
	return m.adaIGDResult, m.adaIGDErr
}

func (m *mockRepository) SimpanPenilaianMedisIGD(ctx context.Context, noRawat, kodeDokter string, req SimpanPenilaianMedisIGDRequest) error {
	m.simpanIGDCalledReq = req
	return m.simpanIGDErr
}

func (m *mockRepository) UpdatePenilaianMedisIGD(ctx context.Context, noRawat string, req UpdatePenilaianMedisIGDRequest) error {
	m.updateIGDCalledReq = req
	return m.updateIGDErr
}

func (m *mockRepository) HapusPenilaianMedisIGD(ctx context.Context, noRawat string) error {
	m.hapusIGDCalledNo = noRawat
	return m.hapusIGDErr
}

type mockRawatJalanService struct {
	tglReg, jamReg string
	exists         bool
	err            error
}

func (m *mockRawatJalanService) DaftarAntreanDokter(ctx context.Context, kodeDokter string, filter rawatjalan.FilterAntreanDokter) ([]rawatjalan.KunjunganRawatJalan, shared.PaginationMeta, error) {
	return nil, shared.PaginationMeta{}, nil
}
func (m *mockRawatJalanService) DetailKunjungan(ctx context.Context, noRawat string, kodeDokter string) (*rawatjalan.KunjunganRawatJalan, error) {
	return nil, nil
}
func (m *mockRawatJalanService) RiwayatKunjunganPasien(ctx context.Context, noRM string) ([]rawatjalan.KunjunganRawatJalan, error) {
	return nil, nil
}
func (m *mockRawatJalanService) GetWaktuRegistrasi(ctx context.Context, noRawat string) (string, string, bool, error) {
	return m.tglReg, m.jamReg, m.exists, m.err
}
func (m *mockRawatJalanService) DaftarStatusPemeriksaan(ctx context.Context) []rawatjalan.OpsiReferensi {
	return nil
}
func (m *mockRawatJalanService) DaftarStatusLanjut(ctx context.Context) []rawatjalan.OpsiReferensi {
	return nil
}
func (m *mockRawatJalanService) DaftarStatusBayar(ctx context.Context) []rawatjalan.OpsiReferensi {
	return nil
}
func (m *mockRawatJalanService) DaftarJenisAntrean(ctx context.Context) []rawatjalan.OpsiReferensi {
	return nil
}

const testKey = "12345678901234567890123456789012"

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestService_Referensi(t *testing.T) {
	log := logger.New()
	svc := NewService(&mockRepository{}, &mockRawatJalanService{}, log, testKey)
	ref := svc.Referensi(context.Background())

	if len(ref.Anamnesis) == 0 || len(ref.Keadaan) == 0 || len(ref.Kesadaran) == 0 || len(ref.StatusFisik) == 0 {
		t.Errorf("expected all reference options to be non-empty, got: %+v", ref)
	}
}

// ==========================================
// RAWAT JALAN TESTS
// ==========================================

func TestService_DetailPenilaianMedisRalan_Success(t *testing.T) {
	log := logger.New()
	repo := &mockRepository{
		detailData: &PenilaianMedisRalan{
			NoRawat:    "2026/04/22/000001",
			NoRM:       "123456",
			KodeDokter: "DR001",
			NamaDokter: "dr. Handi",
			DataPenilaianMedisRalan: DataPenilaianMedisRalan{
				KeluhanUtama: "Batuk",
			},
		},
	}
	svc := NewService(repo, &mockRawatJalanService{}, log, testKey)

	result, err := svc.DetailPenilaianMedisRalan(context.Background(), "2026/04/22/000001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.NoRawat != "2026/04/22/000001" {
		t.Errorf("expected no_rawat 2026/04/22/000001, got %s", result.NoRawat)
	}
	if result.IdKunjungan == "" || result.IdPasien == "" {
		t.Error("expected encrypted id_kunjungan and id_pasien to be populated")
	}
}

func TestService_DetailPenilaianMedisRalan_NotFound(t *testing.T) {
	log := logger.New()
	repo := &mockRepository{
		detailErr: sql.ErrNoRows,
	}
	svc := NewService(repo, &mockRawatJalanService{}, log, testKey)

	_, err := svc.DetailPenilaianMedisRalan(context.Background(), "2026/04/22/999999")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var notFound *apperror.NotFoundError
	if !errors.As(err, &notFound) {
		t.Fatalf("expected NotFoundError, got %v", err)
	}
}

func TestService_RiwayatPenilaianMedisRalanByNoRM_Success(t *testing.T) {
	log := logger.New()
	repo := &mockRepository{
		riwayatData: []PenilaianMedisRalan{
			{
				NoRawat: "2026/04/22/000001",
				NoRM:    "123456",
			},
			{
				NoRawat: "2026/04/21/000002",
				NoRM:    "123456",
			},
		},
	}
	svc := NewService(repo, &mockRawatJalanService{}, log, testKey)

	list, err := svc.RiwayatPenilaianMedisRalanByNoRM(context.Background(), "123456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 items, got %d", len(list))
	}
	for _, item := range list {
		if item.IdKunjungan == "" || item.IdPasien == "" {
			t.Error("expected encrypted ids for each history item")
		}
	}
}

func TestService_SimpanPenilaianMedisRalan_Expired48Hours(t *testing.T) {
	log := logger.New()
	repo := &mockRepository{}
	rjSvc := &mockRawatJalanService{tglReg: "2020-01-01", jamReg: "10:00:00", exists: true}
	svc := NewService(repo, rjSvc, log, testKey)

	req := SimpanPenilaianMedisRalanRequest{
		NoRawat: "2026/04/22/000001",
		DataPenilaianMedisRalan: DataPenilaianMedisRalan{
			TanggalPenilaian: "2026-04-22 09:30:00",
			KeluhanUtama:     "Demam",
			Diagnosis:        "Febris",
			TataLaksana:      "Paracetamol",
		},
	}

	_, err := svc.SimpanPenilaianMedisRalan(context.Background(), "DR001", "2026/04/22/000001", req)
	if err == nil {
		t.Fatal("expected 48-hour expiration error, got nil")
	}
	var forbErr *apperror.ForbiddenError
	if !errors.As(err, &forbErr) {
		t.Fatalf("expected ForbiddenError, got %v", err)
	}
}

func TestService_SimpanPenilaianMedisRalan_EarlierThanRegistration(t *testing.T) {
	log := logger.New()
	today := time.Now().Format("2006-01-02")
	repo := &mockRepository{}
	rjSvc := &mockRawatJalanService{tglReg: today, jamReg: "10:00:00", exists: true}
	svc := NewService(repo, rjSvc, log, testKey)

	// Penilaian medis jam 08:00 (lebih awal dari registrasi jam 10:00)
	req := SimpanPenilaianMedisRalanRequest{
		NoRawat: "2026/04/22/000001",
		DataPenilaianMedisRalan: DataPenilaianMedisRalan{
			TanggalPenilaian: today + " 08:00:00",
			KeluhanUtama:     "Demam",
			Diagnosis:        "Febris",
			TataLaksana:      "Paracetamol",
		},
	}

	_, err := svc.SimpanPenilaianMedisRalan(context.Background(), "DR001", "2026/04/22/000001", req)
	if err == nil {
		t.Fatal("expected validation error when waktu penilaian is before waktu_registrasi")
	}
	var valErr apperror.ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
	if valErr["tanggal_penilaian"] == "" {
		t.Errorf("expected error on tanggal_penilaian, got: %+v", valErr)
	}
}

func TestService_SimpanPenilaianMedisRalan_Success(t *testing.T) {
	log := logger.New()
	today := time.Now().Format("2006-01-02")
	repo := &mockRepository{
		adaResult: false,
		detailData: &PenilaianMedisRalan{
			NoRawat:    "2026/04/22/000001",
			KodeDokter: "DR001",
		},
	}
	rjSvc := &mockRawatJalanService{tglReg: today, jamReg: "08:00:00", exists: true}
	svc := NewService(repo, rjSvc, log, testKey)

	req := SimpanPenilaianMedisRalanRequest{
		NoRawat: "2026/04/22/000001",
		DataPenilaianMedisRalan: DataPenilaianMedisRalan{
			TanggalPenilaian: today + " 09:00:00",
			KeluhanUtama:     "Demam",
			Diagnosis:        "Febris",
			TataLaksana:      "Paracetamol",
		},
	}

	result, err := svc.SimpanPenilaianMedisRalan(context.Background(), "DR001", "2026/04/22/000001", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.NoRawat != "2026/04/22/000001" {
		t.Errorf("expected no_rawat 2026/04/22/000001, got %s", result.NoRawat)
	}
}

func TestService_SimpanPenilaianMedisRalan_Duplicate(t *testing.T) {
	log := logger.New()
	today := time.Now().Format("2006-01-02")
	repo := &mockRepository{
		adaResult: true, // sudah ada
	}
	rjSvc := &mockRawatJalanService{tglReg: today, jamReg: "08:00:00", exists: true}
	svc := NewService(repo, rjSvc, log, testKey)

	req := SimpanPenilaianMedisRalanRequest{
		NoRawat: "2026/04/22/000001",
		DataPenilaianMedisRalan: DataPenilaianMedisRalan{
			TanggalPenilaian: today + " 09:00:00",
			KeluhanUtama:     "Demam",
			Diagnosis:        "Febris",
			TataLaksana:      "Paracetamol",
		},
	}

	_, err := svc.SimpanPenilaianMedisRalan(context.Background(), "DR001", "2026/04/22/000001", req)
	if err == nil {
		t.Fatal("expected duplicate error, got nil")
	}
	var busErr *apperror.BusinessError
	if !errors.As(err, &busErr) {
		t.Fatalf("expected BusinessError, got %v", err)
	}
}

func TestService_UpdatePenilaianMedisRalan_Forbidden(t *testing.T) {
	log := logger.New()
	today := time.Now().Format("2006-01-02")
	repo := &mockRepository{
		detailData: &PenilaianMedisRalan{
			NoRawat:    "2026/04/22/000001",
			KodeDokter: "DR001", // dibuat oleh DR001
			NamaDokter: "dr. Handi",
		},
	}
	rjSvc := &mockRawatJalanService{tglReg: today, jamReg: "08:00:00", exists: true}
	svc := NewService(repo, rjSvc, log, testKey)

	req := UpdatePenilaianMedisRalanRequest{
		DataPenilaianMedisRalan: DataPenilaianMedisRalan{
			TanggalPenilaian: today + " 09:00:00",
			KeluhanUtama:     "Demam berkurang",
			Diagnosis:        "Febris H2",
			TataLaksana:      "Lanjut Paracetamol",
		},
	}

	// Dicoba update oleh DR002
	_, err := svc.UpdatePenilaianMedisRalan(context.Background(), "DR002", "2026/04/22/000001", req)
	if err == nil {
		t.Fatal("expected forbidden error, got nil")
	}
	var forbErr *apperror.ForbiddenError
	if !errors.As(err, &forbErr) {
		t.Fatalf("expected ForbiddenError, got %v", err)
	}
}

func TestService_UpdatePenilaianMedisRalan_Success(t *testing.T) {
	log := logger.New()
	today := time.Now().Format("2006-01-02")
	repo := &mockRepository{
		detailData: &PenilaianMedisRalan{
			NoRawat:    "2026/04/22/000001",
			KodeDokter: "DR001",
			NamaDokter: "dr. Handi",
		},
	}
	rjSvc := &mockRawatJalanService{tglReg: today, jamReg: "08:00:00", exists: true}
	svc := NewService(repo, rjSvc, log, testKey)

	req := UpdatePenilaianMedisRalanRequest{
		DataPenilaianMedisRalan: DataPenilaianMedisRalan{
			TanggalPenilaian: today + " 09:00:00",
			KeluhanUtama:     "Demam berkurang",
			Diagnosis:        "Febris H2",
			TataLaksana:      "Lanjut Paracetamol",
		},
	}

	result, err := svc.UpdatePenilaianMedisRalan(context.Background(), "DR001", "2026/04/22/000001", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.NoRawat != "2026/04/22/000001" {
		t.Errorf("expected no_rawat 2026/04/22/000001, got %s", result.NoRawat)
	}
}

func TestService_HapusPenilaianMedisRalan_Success(t *testing.T) {
	log := logger.New()
	today := time.Now().Format("2006-01-02")
	repo := &mockRepository{
		detailData: &PenilaianMedisRalan{
			NoRawat:    "2026/04/22/000001",
			KodeDokter: "DR001",
			NamaDokter: "dr. Handi",
		},
	}
	rjSvc := &mockRawatJalanService{tglReg: today, jamReg: "08:00:00", exists: true}
	svc := NewService(repo, rjSvc, log, testKey)

	err := svc.HapusPenilaianMedisRalan(context.Background(), "DR001", "2026/04/22/000001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.hapusCalledNo != "2026/04/22/000001" {
		t.Errorf("expected hapusCalledNo to be 2026/04/22/000001, got %s", repo.hapusCalledNo)
	}
}

func TestService_HapusPenilaianMedisRalan_Forbidden(t *testing.T) {
	log := logger.New()
	today := time.Now().Format("2006-01-02")
	repo := &mockRepository{
		detailData: &PenilaianMedisRalan{
			NoRawat:    "2026/04/22/000001",
			KodeDokter: "DR001", // dibuat oleh DR001
			NamaDokter: "dr. Handi",
		},
	}
	rjSvc := &mockRawatJalanService{tglReg: today, jamReg: "08:00:00", exists: true}
	svc := NewService(repo, rjSvc, log, testKey)

	// Dicoba hapus oleh DR002
	err := svc.HapusPenilaianMedisRalan(context.Background(), "DR002", "2026/04/22/000001")
	if err == nil {
		t.Fatal("expected forbidden error, got nil")
	}
	var forbErr *apperror.ForbiddenError
	if !errors.As(err, &forbErr) {
		t.Fatalf("expected ForbiddenError, got %v", err)
	}
}

// ==========================================
// GAWAT DARURAT (IGD) TESTS
// ==========================================

func TestService_DetailPenilaianMedisIGD_Success(t *testing.T) {
	log := logger.New()
	repo := &mockRepository{
		detailIGDData: &PenilaianMedisIGD{
			NoRawat:    "2026/04/22/000002",
			NoRM:       "654321",
			KodeDokter: "DR001",
			NamaDokter: "dr. Handi",
			DataPenilaianMedisIGD: DataPenilaianMedisIGD{
				KeluhanUtama: "Sesak Napas Akut",
				EKG:          "Sinus Takikardia",
			},
		},
	}
	svc := NewService(repo, &mockRawatJalanService{}, log, testKey)

	result, err := svc.DetailPenilaianMedisIGD(context.Background(), "2026/04/22/000002")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.NoRawat != "2026/04/22/000002" {
		t.Errorf("expected no_rawat 2026/04/22/000002, got %s", result.NoRawat)
	}
	if result.IdKunjungan == "" || result.IdPasien == "" {
		t.Error("expected encrypted id_kunjungan and id_pasien to be populated")
	}
}

func TestService_DetailPenilaianMedisIGD_NotFound(t *testing.T) {
	log := logger.New()
	repo := &mockRepository{
		detailIGDErr: sql.ErrNoRows,
	}
	svc := NewService(repo, &mockRawatJalanService{}, log, testKey)

	_, err := svc.DetailPenilaianMedisIGD(context.Background(), "2026/04/22/999999")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var notFound *apperror.NotFoundError
	if !errors.As(err, &notFound) {
		t.Fatalf("expected NotFoundError, got %v", err)
	}
}

func TestService_RiwayatPenilaianMedisIGDByNoRM_Success(t *testing.T) {
	log := logger.New()
	repo := &mockRepository{
		riwayatIGDData: []PenilaianMedisIGD{
			{NoRawat: "2026/04/22/000002", NoRM: "654321"},
		},
	}
	svc := NewService(repo, &mockRawatJalanService{}, log, testKey)

	list, err := svc.RiwayatPenilaianMedisIGDByNoRM(context.Background(), "654321")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 item, got %d", len(list))
	}
	if list[0].IdKunjungan == "" || list[0].IdPasien == "" {
		t.Error("expected encrypted ids for history item")
	}
}

func TestService_SimpanPenilaianMedisIGD_Success(t *testing.T) {
	log := logger.New()
	today := time.Now().Format("2006-01-02")
	repo := &mockRepository{
		adaIGDResult: false,
		detailIGDData: &PenilaianMedisIGD{
			NoRawat:    "2026/04/22/000002",
			KodeDokter: "DR001",
		},
	}
	rjSvc := &mockRawatJalanService{tglReg: today, jamReg: "08:00:00", exists: true}
	svc := NewService(repo, rjSvc, log, testKey)

	req := SimpanPenilaianMedisIGDRequest{
		NoRawat: "2026/04/22/000002",
		DataPenilaianMedisIGD: DataPenilaianMedisIGD{
			TanggalPenilaian: today + " 09:00:00",
			KeluhanUtama:     "Nyeri dada hebat",
			Diagnosis:        "STEMI",
			TataLaksana:      "Oksigen + Aspilet",
			EKG:              "ST Elevasi V1-V4",
		},
	}

	result, err := svc.SimpanPenilaianMedisIGD(context.Background(), "DR001", "2026/04/22/000002", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.NoRawat != "2026/04/22/000002" {
		t.Errorf("expected no_rawat 2026/04/22/000002, got %s", result.NoRawat)
	}
}

func TestService_SimpanPenilaianMedisIGD_Duplicate(t *testing.T) {
	log := logger.New()
	today := time.Now().Format("2006-01-02")
	repo := &mockRepository{
		adaIGDResult: true, // sudah ada di IGD
	}
	rjSvc := &mockRawatJalanService{tglReg: today, jamReg: "08:00:00", exists: true}
	svc := NewService(repo, rjSvc, log, testKey)

	req := SimpanPenilaianMedisIGDRequest{
		NoRawat: "2026/04/22/000002",
		DataPenilaianMedisIGD: DataPenilaianMedisIGD{
			TanggalPenilaian: today + " 09:00:00",
			KeluhanUtama:     "Nyeri dada hebat",
			Diagnosis:        "STEMI",
			TataLaksana:      "Oksigen",
		},
	}

	_, err := svc.SimpanPenilaianMedisIGD(context.Background(), "DR001", "2026/04/22/000002", req)
	if err == nil {
		t.Fatal("expected duplicate error for IGD, got nil")
	}
	var busErr *apperror.BusinessError
	if !errors.As(err, &busErr) {
		t.Fatalf("expected BusinessError, got %v", err)
	}
}

func TestService_UpdatePenilaianMedisIGD_Forbidden(t *testing.T) {
	log := logger.New()
	today := time.Now().Format("2006-01-02")
	repo := &mockRepository{
		detailIGDData: &PenilaianMedisIGD{
			NoRawat:    "2026/04/22/000002",
			KodeDokter: "DR001", // dibuat oleh DR001
			NamaDokter: "dr. Handi",
		},
	}
	rjSvc := &mockRawatJalanService{tglReg: today, jamReg: "08:00:00", exists: true}
	svc := NewService(repo, rjSvc, log, testKey)

	req := UpdatePenilaianMedisIGDRequest{
		DataPenilaianMedisIGD: DataPenilaianMedisIGD{
			TanggalPenilaian: today + " 09:00:00",
			KeluhanUtama:     "Nyeri dada membaik",
			Diagnosis:        "STEMI Post PCI",
			TataLaksana:      "Lanjut terapi",
		},
	}

	// Update dicoba oleh DR002
	_, err := svc.UpdatePenilaianMedisIGD(context.Background(), "DR002", "2026/04/22/000002", req)
	if err == nil {
		t.Fatal("expected forbidden error for IGD, got nil")
	}
	var forbErr *apperror.ForbiddenError
	if !errors.As(err, &forbErr) {
		t.Fatalf("expected ForbiddenError, got %v", err)
	}
}

func TestService_UpdatePenilaianMedisIGD_Success(t *testing.T) {
	log := logger.New()
	today := time.Now().Format("2006-01-02")
	repo := &mockRepository{
		detailIGDData: &PenilaianMedisIGD{
			NoRawat:    "2026/04/22/000002",
			KodeDokter: "DR001",
			NamaDokter: "dr. Handi",
		},
	}
	rjSvc := &mockRawatJalanService{tglReg: today, jamReg: "08:00:00", exists: true}
	svc := NewService(repo, rjSvc, log, testKey)

	req := UpdatePenilaianMedisIGDRequest{
		DataPenilaianMedisIGD: DataPenilaianMedisIGD{
			TanggalPenilaian: today + " 09:00:00",
			KeluhanUtama:     "Nyeri dada membaik",
			Diagnosis:        "STEMI Post PCI",
			TataLaksana:      "Lanjut terapi",
		},
	}

	result, err := svc.UpdatePenilaianMedisIGD(context.Background(), "DR001", "2026/04/22/000002", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.NoRawat != "2026/04/22/000002" {
		t.Errorf("expected no_rawat 2026/04/22/000002, got %s", result.NoRawat)
	}
}

func TestService_HapusPenilaianMedisIGD_Success(t *testing.T) {
	log := logger.New()
	today := time.Now().Format("2006-01-02")
	repo := &mockRepository{
		detailIGDData: &PenilaianMedisIGD{
			NoRawat:    "2026/04/22/000002",
			KodeDokter: "DR001",
			NamaDokter: "dr. Handi",
		},
	}
	rjSvc := &mockRawatJalanService{tglReg: today, jamReg: "08:00:00", exists: true}
	svc := NewService(repo, rjSvc, log, testKey)

	err := svc.HapusPenilaianMedisIGD(context.Background(), "DR001", "2026/04/22/000002")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.hapusIGDCalledNo != "2026/04/22/000002" {
		t.Errorf("expected hapusIGDCalledNo to be 2026/04/22/000002, got %s", repo.hapusIGDCalledNo)
	}
}

func TestService_HapusPenilaianMedisIGD_Forbidden(t *testing.T) {
	log := logger.New()
	today := time.Now().Format("2006-01-02")
	repo := &mockRepository{
		detailIGDData: &PenilaianMedisIGD{
			NoRawat:    "2026/04/22/000002",
			KodeDokter: "DR001", // dibuat oleh DR001
			NamaDokter: "dr. Handi",
		},
	}
	rjSvc := &mockRawatJalanService{tglReg: today, jamReg: "08:00:00", exists: true}
	svc := NewService(repo, rjSvc, log, testKey)

	// Dicoba hapus oleh DR002
	err := svc.HapusPenilaianMedisIGD(context.Background(), "DR002", "2026/04/22/000002")
	if err == nil {
		t.Fatal("expected forbidden error for IGD, got nil")
	}
	var forbErr *apperror.ForbiddenError
	if !errors.As(err, &forbErr) {
		t.Fatalf("expected ForbiddenError, got %v", err)
	}
}
