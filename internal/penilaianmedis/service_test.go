package penilaianmedis

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"
	"time"

	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatjalan"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type mockRepository struct {
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
}

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

func TestService_Referensi(t *testing.T) {
	log := logger.New()
	svc := NewService(&mockRepository{}, &mockRawatJalanService{}, log, testKey)
	ref := svc.Referensi(context.Background())

	if len(ref.Anamnesis) == 0 || len(ref.Keadaan) == 0 || len(ref.Kesadaran) == 0 || len(ref.StatusFisik) == 0 {
		t.Errorf("expected all reference options to be non-empty, got: %+v", ref)
	}
}

func TestService_DetailPenilaianMedisRalan_Success(t *testing.T) {
	log := logger.New()
	repo := &mockRepository{
		detailData: &PenilaianMedisRalan{
			NoRawat: "2026/04/22/000001",
			NoRM:    "123456",
		},
	}
	svc := NewService(repo, &mockRawatJalanService{}, log, testKey)

	result, err := svc.DetailPenilaianMedisRalan(context.Background(), "2026/04/22/000001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IdKunjungan == "" || result.IdPasien == "" {
		t.Fatal("expected encrypted IdKunjungan and IdPasien")
	}

	decryptedNoRawat, _ := crypto.Decrypt(result.IdKunjungan, testKey)
	if decryptedNoRawat != "2026/04/22/000001" {
		t.Errorf("unexpected decrypted no_rawat: %s", decryptedNoRawat)
	}
}

func TestService_DetailPenilaianMedisRalan_NotFound(t *testing.T) {
	log := logger.New()
	repo := &mockRepository{
		detailErr: sql.ErrNoRows,
	}
	svc := NewService(repo, &mockRawatJalanService{}, log, testKey)

	_, err := svc.DetailPenilaianMedisRalan(context.Background(), "2026/04/22/000001")
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
			{NoRawat: "2026/04/22/000001", NoRM: "123456"},
		},
	}
	svc := NewService(repo, &mockRawatJalanService{}, log, testKey)

	result, err := svc.RiwayatPenilaianMedisRalanByNoRM(context.Background(), "123456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0].IdKunjungan == "" {
		t.Fatalf("expected 1 result with encrypted IDs, got: %+v", result)
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
			KeluhanUtama: "Demam",
			Diagnosis:    "Febris",
			TataLaksana:  "Paracetamol",
		},
	}

	_, err := svc.SimpanPenilaianMedisRalan(context.Background(), "DR001", "2026/04/22/000001", req)
	if err == nil {
		t.Fatal("expected duplicate error, got nil")
	}
	var bizErr *apperror.BusinessError
	if !errors.As(err, &bizErr) {
		t.Fatalf("expected BusinessError, got %v", err)
	}
}

func TestService_SimpanPenilaianMedisRalan_Exceeds48Hours(t *testing.T) {
	log := logger.New()
	pastDate := time.Now().Add(-72 * time.Hour).Format("2006-01-02")
	repo := &mockRepository{adaResult: false}
	rjSvc := &mockRawatJalanService{tglReg: pastDate, jamReg: "08:00:00", exists: true}
	svc := NewService(repo, rjSvc, log, testKey)

	req := SimpanPenilaianMedisRalanRequest{
		NoRawat: "2026/04/22/000001",
		DataPenilaianMedisRalan: DataPenilaianMedisRalan{
			KeluhanUtama: "Demam",
			Diagnosis:    "Febris",
			TataLaksana:  "Paracetamol",
		},
	}

	_, err := svc.SimpanPenilaianMedisRalan(context.Background(), "DR001", "2026/04/22/000001", req)
	if err == nil {
		t.Fatal("expected 48 hours forbidden error, got nil")
	}
	var forbidden *apperror.ForbiddenError
	if !errors.As(err, &forbidden) {
		t.Fatalf("expected ForbiddenError, got %v", err)
	}
}

func TestService_SimpanPenilaianMedisRalan_WaktuPemeriksaanBeforeRegistrasi(t *testing.T) {
	log := logger.New()
	today := time.Now().Format("2006-01-02")
	repo := &mockRepository{adaResult: false}
	rjSvc := &mockRawatJalanService{tglReg: today, jamReg: "10:00:00", exists: true}
	svc := NewService(repo, rjSvc, log, testKey)

	// Pemeriksaan jam 08:00 (lebih awal dari registrasi jam 10:00)
	req := SimpanPenilaianMedisRalanRequest{
		NoRawat: "2026/04/22/000001",
		DataPenilaianMedisRalan: DataPenilaianMedisRalan{
			TanggalPemeriksaan: today + " 08:00:00",
			KeluhanUtama:       "Demam",
			Diagnosis:          "Febris",
			TataLaksana:        "Paracetamol",
		},
	}

	_, err := svc.SimpanPenilaianMedisRalan(context.Background(), "DR001", "2026/04/22/000001", req)
	if err == nil {
		t.Fatal("expected validation error when waktu_pemeriksaan is before waktu_registrasi")
	}
	var valErr apperror.ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
	if valErr["tanggal_pemeriksaan"] == "" {
		t.Errorf("expected error on tanggal_pemeriksaan, got: %+v", valErr)
	}
}

func TestService_SimpanPenilaianMedisRalan_Success(t *testing.T) {
	log := logger.New()
	today := time.Now().Format("2006-01-02")
	repo := &mockRepository{
		adaResult: false,
		detailData: &PenilaianMedisRalan{
			NoRawat:    "2026/04/22/000001",
			NoRM:       "123456",
			KodeDokter: "DR001",
		},
	}
	rjSvc := &mockRawatJalanService{tglReg: today, jamReg: "08:00:00", exists: true}
	svc := NewService(repo, rjSvc, log, testKey)

	req := SimpanPenilaianMedisRalanRequest{
		NoRawat: "2026/04/22/000001",
		DataPenilaianMedisRalan: DataPenilaianMedisRalan{
			KeluhanUtama: "Demam",
			Diagnosis:    "Febris",
			TataLaksana:  "Paracetamol",
		},
	}

	result, err := svc.SimpanPenilaianMedisRalan(context.Background(), "DR001", "2026/04/22/000001", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || result.IdKunjungan == "" {
		t.Fatal("expected result with encrypted ID, got nil")
	}
	if repo.simpanCalledReq.KeluhanUtama != "Demam" {
		t.Errorf("expected repo called with sanitized request, got: %+v", repo.simpanCalledReq)
	}
}

func TestService_UpdatePenilaianMedisRalan_ForbiddenOtherDoctor(t *testing.T) {
	log := logger.New()
	repo := &mockRepository{
		detailData: &PenilaianMedisRalan{
			NoRawat:    "2026/04/22/000001",
			KodeDokter: "DR002", // Dibuat dokter lain
			NamaDokter: "dr. Lain",
		},
	}
	svc := NewService(repo, &mockRawatJalanService{}, log, testKey)

	req := UpdatePenilaianMedisRalanRequest{
		DataPenilaianMedisRalan: DataPenilaianMedisRalan{
			KeluhanUtama: "Demam",
			Diagnosis:    "Febris",
			TataLaksana:  "Paracetamol",
		},
	}

	_, err := svc.UpdatePenilaianMedisRalan(context.Background(), "DR001", "2026/04/22/000001", req)
	if err == nil {
		t.Fatal("expected forbidden error, got nil")
	}
	var forbidden *apperror.ForbiddenError
	if !errors.As(err, &forbidden) {
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
			KeluhanUtama: "Demam berdarah",
			Diagnosis:    "DHF",
			TataLaksana:  "RL 20 tpm",
		},
	}

	result, err := svc.UpdatePenilaianMedisRalan(context.Background(), "DR001", "2026/04/22/000001", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || result.IdKunjungan == "" {
		t.Fatal("expected result with encrypted ID, got nil")
	}
	if repo.updateCalledReq.KeluhanUtama != "Demam berdarah" {
		t.Errorf("expected repo called with updated request, got: %+v", repo.updateCalledReq)
	}
}

func TestService_HapusPenilaianMedisRalan_Success(t *testing.T) {
	log := logger.New()
	today := time.Now().Format("2006-01-02")
	repo := &mockRepository{
		detailData: &PenilaianMedisRalan{
			NoRawat:    "2026/04/22/000001",
			KodeDokter: "DR001",
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

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
