package rujukaninternal

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
	opsiList          []OpsiPoliDokter
	opsiErr           error
	rujukanList       []RujukanInternal
	rujukanErr        error
	adaResult         bool
	adaErr            error
	simpanErr         error
	hapusErr          error
	simpanCalledWith  struct{ noRawat, kdPoli, kdDokter string }
	hapusCalledWith   struct{ noRawat, kdDokter string }
}

func (m *mockRepository) DaftarOpsiPoliDokter(ctx context.Context, kodeDokterLogin, keyword string) ([]OpsiPoliDokter, error) {
	return m.opsiList, m.opsiErr
}

func (m *mockRepository) DaftarRujukanInternalByNoRawat(ctx context.Context, noRawat string) ([]RujukanInternal, error) {
	return m.rujukanList, m.rujukanErr
}

func (m *mockRepository) CekRujukanInternalAda(ctx context.Context, noRawat, kdDokter string) (bool, error) {
	return m.adaResult, m.adaErr
}

func (m *mockRepository) SimpanRujukanInternal(ctx context.Context, noRawat, kdPoli, kdDokter string) error {
	m.simpanCalledWith.noRawat = noRawat
	m.simpanCalledWith.kdPoli = kdPoli
	m.simpanCalledWith.kdDokter = kdDokter
	return m.simpanErr
}

func (m *mockRepository) HapusRujukanInternal(ctx context.Context, noRawat, kdDokter string) error {
	m.hapusCalledWith.noRawat = noRawat
	m.hapusCalledWith.kdDokter = kdDokter
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

func TestDaftarOpsiPoliDokter_Success(t *testing.T) {
	log := logger.New()
	repo := &mockRepository{
		opsiList: []OpsiPoliDokter{
			{KodePoli: "INT", NamaPoli: "Penyakit Dalam", KodeDokter: "DR002", NamaDokter: "dr. Sp.PD"},
		},
	}
	svc := NewService(repo, &mockRawatJalanService{}, log, testKey)

	result, err := svc.DaftarOpsiPoliDokter(context.Background(), "DR001", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result))
	}
	if result[0].Id == "" {
		t.Fatal("expected encrypted Id, got empty")
	}

	decrypted, err := crypto.Decrypt(result[0].Id, testKey)
	if err != nil {
		t.Fatalf("failed to decrypt ID: %v", err)
	}
	if decrypted != "INT~DR002" {
		t.Fatalf("expected decrypted ID 'INT~DR002', got %s", decrypted)
	}
}

func TestDaftarOpsiPoliDokter_WithKeyword(t *testing.T) {
	log := logger.New()
	repo := &mockRepository{
		opsiList: []OpsiPoliDokter{
			{KodePoli: "INT", NamaPoli: "Penyakit Dalam", KodeDokter: "DR002", NamaDokter: "dr. Sp.PD"},
		},
	}
	svc := NewService(repo, &mockRawatJalanService{}, log, testKey)

	result, err := svc.DaftarOpsiPoliDokter(context.Background(), "DR001", "Penyakit")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result))
	}
}

func TestDaftarRujukanInternal_Success(t *testing.T) {
	log := logger.New()
	repo := &mockRepository{
		rujukanList: []RujukanInternal{
			{NoRawat: "2026/04/22/000001", KodePoli: "INT", NamaPoli: "Penyakit Dalam", KodeDokter: "DR002", NamaDokter: "dr. Sp.PD"},
		},
	}
	svc := NewService(repo, &mockRawatJalanService{}, log, testKey)

	result, err := svc.DaftarRujukanInternal(context.Background(), "2026/04/22/000001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result))
	}
	if result[0].Id == "" || result[0].IdKunjungan == "" {
		t.Fatal("expected encrypted Id and IdKunjungan, got empty")
	}
}

func TestSimpanRujukanInternal_CannotReferToSelf(t *testing.T) {
	log := logger.New()
	repo := &mockRepository{}
	svc := NewService(repo, &mockRawatJalanService{}, log, testKey)

	targetEncrypted, _ := crypto.Encrypt("INT~DR001", testKey)
	req := SimpanRujukanRequest{IdTujuan: targetEncrypted}

	_, err := svc.SimpanRujukanInternal(context.Background(), "DR001", "2026/04/22/000001", req)
	if err == nil {
		t.Fatal("expected error when referring to self, got nil")
	}
	var valErr apperror.ValidationError
	if !errors.As(err, &valErr) || valErr["id_tujuan"] == "" {
		t.Fatalf("expected validation error on id_tujuan, got %v", err)
	}
}

func TestSimpanRujukanInternal_DuplicateReferral(t *testing.T) {
	log := logger.New()
	repo := &mockRepository{
		adaResult: true, // sudah ada rujukan
	}
	today := time.Now().Format("2006-01-02")
	rjSvc := &mockRawatJalanService{tglReg: today, jamReg: "08:00:00", exists: true}
	svc := NewService(repo, rjSvc, log, testKey)

	targetEncrypted, _ := crypto.Encrypt("INT~DR002", testKey)
	req := SimpanRujukanRequest{IdTujuan: targetEncrypted}

	_, err := svc.SimpanRujukanInternal(context.Background(), "DR001", "2026/04/22/000001", req)
	if err == nil {
		t.Fatal("expected error for duplicate referral, got nil")
	}
	var valErr apperror.ValidationError
	if !errors.As(err, &valErr) || valErr["id_tujuan"] == "" {
		t.Fatalf("expected validation error on duplicate referral, got %v", err)
	}
}

func TestSimpanRujukanInternal_Exceeds48Hours(t *testing.T) {
	log := logger.New()
	repo := &mockRepository{adaResult: false}
	pastDate := time.Now().Add(-72 * time.Hour).Format("2006-01-02")
	rjSvc := &mockRawatJalanService{tglReg: pastDate, jamReg: "08:00:00", exists: true}
	svc := NewService(repo, rjSvc, log, testKey)

	targetEncrypted, _ := crypto.Encrypt("INT~DR002", testKey)
	req := SimpanRujukanRequest{IdTujuan: targetEncrypted}

	_, err := svc.SimpanRujukanInternal(context.Background(), "DR001", "2026/04/22/000001", req)
	if err == nil {
		t.Fatal("expected 48 hours forbidden error, got nil")
	}
	var forbidden *apperror.ForbiddenError
	if !errors.As(err, &forbidden) {
		t.Fatalf("expected ForbiddenError, got %v", err)
	}
}

func TestSimpanRujukanInternal_Success(t *testing.T) {
	log := logger.New()
	today := time.Now().Format("2006-01-02")
	repo := &mockRepository{
		adaResult: false,
		rujukanList: []RujukanInternal{
			{NoRawat: "2026/04/22/000001", KodePoli: "INT", NamaPoli: "Penyakit Dalam", KodeDokter: "DR002", NamaDokter: "dr. Sp.PD"},
		},
	}
	rjSvc := &mockRawatJalanService{tglReg: today, jamReg: "08:00:00", exists: true}
	svc := NewService(repo, rjSvc, log, testKey)

	targetEncrypted, _ := crypto.Encrypt("INT~DR002", testKey)
	req := SimpanRujukanRequest{IdTujuan: targetEncrypted}

	result, err := svc.SimpanRujukanInternal(context.Background(), "DR001", "2026/04/22/000001", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if repo.simpanCalledWith.kdDokter != "DR002" || repo.simpanCalledWith.kdPoli != "INT" {
		t.Fatalf("unexpected save arguments: %+v", repo.simpanCalledWith)
	}
}

func TestHapusRujukanInternal_Success(t *testing.T) {
	log := logger.New()
	today := time.Now().Format("2006-01-02")
	repo := &mockRepository{}
	rjSvc := &mockRawatJalanService{tglReg: today, jamReg: "08:00:00", exists: true}
	svc := NewService(repo, rjSvc, log, testKey)

	idRujukanEncrypted, _ := crypto.Encrypt("2026/04/22/000001~INT~DR002", testKey)

	err := svc.HapusRujukanInternal(context.Background(), "DR001", "2026/04/22/000001", idRujukanEncrypted)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.hapusCalledWith.noRawat != "2026/04/22/000001" || repo.hapusCalledWith.kdDokter != "DR002" {
		t.Fatalf("unexpected delete arguments: %+v", repo.hapusCalledWith)
	}
}

func TestHapusRujukanInternal_NoRawatMismatch(t *testing.T) {
	log := logger.New()
	repo := &mockRepository{}
	svc := NewService(repo, &mockRawatJalanService{}, log, testKey)

	idRujukanEncrypted, _ := crypto.Encrypt("2026/04/22/000002~INT~DR002", testKey)

	err := svc.HapusRujukanInternal(context.Background(), "DR001", "2026/04/22/000001", idRujukanEncrypted)
	if err == nil {
		t.Fatal("expected error on no_rawat mismatch, got nil")
	}
}

func TestHapusRujukanInternal_NotFound(t *testing.T) {
	log := logger.New()
	today := time.Now().Format("2006-01-02")
	repo := &mockRepository{hapusErr: sql.ErrNoRows}
	rjSvc := &mockRawatJalanService{tglReg: today, jamReg: "08:00:00", exists: true}
	svc := NewService(repo, rjSvc, log, testKey)

	idRujukanEncrypted, _ := crypto.Encrypt("2026/04/22/000001~INT~DR002", testKey)

	err := svc.HapusRujukanInternal(context.Background(), "DR001", "2026/04/22/000001", idRujukanEncrypted)
	if err == nil {
		t.Fatal("expected not found error, got nil")
	}
	var notFound *apperror.NotFoundError
	if !errors.As(err, &notFound) {
		t.Fatalf("expected NotFoundError, got %v", err)
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
