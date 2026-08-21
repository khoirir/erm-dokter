package pemeriksaan_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"erm-dokter/internal/pemeriksaan"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatjalan"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type mockRepository struct {
	daftarPemeriksaanFunc func(ctx context.Context, listNoRawat []string, statusLanjut shared.StatusLanjut, filter pemeriksaan.FilterDaftarPemeriksaan) ([]pemeriksaan.Pemeriksaan, int, error)
	detailPemeriksaanFunc func(ctx context.Context, id pemeriksaan.IdPemeriksaan, statusLanjut shared.StatusLanjut) (*pemeriksaan.Pemeriksaan, error)
	simpanPemeriksaanFunc func(ctx context.Context, kodeDokter string, statusLanjut shared.StatusLanjut, req pemeriksaan.SimpanPemeriksaanRequest) error
	hapusPemeriksaanFunc  func(ctx context.Context, id pemeriksaan.IdPemeriksaan, statusLanjut shared.StatusLanjut) error
}

func (m *mockRepository) DaftarPemeriksaan(ctx context.Context, listNoRawat []string, statusLanjut shared.StatusLanjut, filter pemeriksaan.FilterDaftarPemeriksaan) ([]pemeriksaan.Pemeriksaan, int, error) {
	if m.daftarPemeriksaanFunc != nil {
		return m.daftarPemeriksaanFunc(ctx, listNoRawat, statusLanjut, filter)
	}
	return []pemeriksaan.Pemeriksaan{}, 0, nil
}

func (m *mockRepository) DetailPemeriksaan(ctx context.Context, id pemeriksaan.IdPemeriksaan, statusLanjut shared.StatusLanjut) (*pemeriksaan.Pemeriksaan, error) {
	if m.detailPemeriksaanFunc != nil {
		return m.detailPemeriksaanFunc(ctx, id, statusLanjut)
	}
	return nil, nil
}

func (m *mockRepository) SimpanPemeriksaan(ctx context.Context, kodeDokter string, statusLanjut shared.StatusLanjut, req pemeriksaan.SimpanPemeriksaanRequest) error {
	if m.simpanPemeriksaanFunc != nil {
		return m.simpanPemeriksaanFunc(ctx, kodeDokter, statusLanjut, req)
	}
	return nil
}

func (m *mockRepository) HapusPemeriksaan(ctx context.Context, id pemeriksaan.IdPemeriksaan, statusLanjut shared.StatusLanjut) error {
	if m.hapusPemeriksaanFunc != nil {
		return m.hapusPemeriksaanFunc(ctx, id, statusLanjut)
	}
	return nil
}

type mockRawatJalanService struct {
	rawatjalan.Service
	riwayatKunjunganPasienFunc func(ctx context.Context, noRM string) ([]rawatjalan.KunjunganRawatJalan, error)
	getWaktuRegistrasiFunc     func(ctx context.Context, noRawat string) (string, string, bool, error)
}

func (m *mockRawatJalanService) RiwayatKunjunganPasien(ctx context.Context, noRM string) ([]rawatjalan.KunjunganRawatJalan, error) {
	if m.riwayatKunjunganPasienFunc != nil {
		return m.riwayatKunjunganPasienFunc(ctx, noRM)
	}
	return []rawatjalan.KunjunganRawatJalan{}, nil
}

func (m *mockRawatJalanService) GetWaktuRegistrasi(ctx context.Context, noRawat string) (string, string, bool, error) {
	if m.getWaktuRegistrasiFunc != nil {
		return m.getWaktuRegistrasiFunc(ctx, noRawat)
	}
	return "2020-01-01", "00:00:00", true, nil
}

func TestGetDaftarKesadaran(t *testing.T) {
	repo := &mockRepository{}
	rjRepo := &mockRawatJalanService{}
	log := logger.New()
	svc := pemeriksaan.NewService(repo, rjRepo, 48, log)

	daftar := svc.GetDaftarKesadaran(context.Background())
	if len(daftar) != 11 {
		t.Fatalf("expected 11 items in daftar kesadaran, got %d", len(daftar))
	}

	foundComposMentis := false
	for _, item := range daftar {
		if item.Value == "Compos Mentis" && item.Label == "Compos Mentis" {
			foundComposMentis = true
			break
		}
	}
	if !foundComposMentis {
		t.Error("expected 'Compos Mentis' to be in daftar kesadaran")
	}
}

func TestSimpanPemeriksaan_SuccessRalan(t *testing.T) {
	called := false
	repo := &mockRepository{
		simpanPemeriksaanFunc: func(ctx context.Context, kodeDokter string, statusLanjut shared.StatusLanjut, req pemeriksaan.SimpanPemeriksaanRequest) error {
			called = true
			if kodeDokter != "DK001" {
				t.Errorf("expected kodeDokter 'DK001', got '%s'", kodeDokter)
			}
			if statusLanjut != shared.StatusLanjutRawatJalan {
				t.Errorf("expected statusLanjut 'Ralan', got '%s'", statusLanjut)
			}
			if req.NoRawat != "2026/04/22/036934" {
				t.Errorf("expected no_rawat '2026/04/22/036934', got '%s'", req.NoRawat)
			}
			return nil
		},
	}
	rjRepo := &mockRawatJalanService{}
	log := logger.New()
	svc := pemeriksaan.NewService(repo, rjRepo, 48, log)

	req := pemeriksaan.SimpanPemeriksaanRequest{
		NoRawat:             "2026/04/22/036934",
		TanggalPemeriksaan:  "2026-04-23",
		JamPemeriksaan:      "12:10:00",
		Kesadaran:           pemeriksaan.KesadaranComposMentis,
		Keluhan:             "Demam",
		Pemeriksaan:         "Normal",
		Penilaian:           "Febris",
		RencanaTindakLanjut: "Istirahat",
		Instruksi:           "Minum obat",
		Evaluasi:            "Stabil",
	}

	res, err := svc.SimpanPemeriksaan(context.Background(), "DK001", shared.StatusLanjutRawatJalan, req)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !called {
		t.Error("expected repository SimpanPemeriksaan to be called")
	}
	if res == nil {
		t.Fatal("expected returned Pemeriksaan object, got nil")
	}
	if res.NoRawat != "2026/04/22/036934" {
		t.Errorf("expected NoRawat '2026/04/22/036934', got '%s'", res.NoRawat)
	}
}

func TestSimpanPemeriksaan_ValidationError(t *testing.T) {
	repo := &mockRepository{}
	rjRepo := &mockRawatJalanService{}
	log := logger.New()
	svc := pemeriksaan.NewService(repo, rjRepo, 48, log)

	// Missing required fields (no_rawat, keluhan, dll) and invalid kesadaran
	req := pemeriksaan.SimpanPemeriksaanRequest{
		NoRawat:   "",
		Kesadaran: "KesadaranPalsu",
	}

	_, err := svc.SimpanPemeriksaan(context.Background(), "DK001", shared.StatusLanjutRawatJalan, req)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	valErr, ok := err.(apperror.ValidationError)
	if !ok {
		t.Fatalf("expected apperror.ValidationError, got %T", err)
	}

	if _, exists := valErr["no_rawat"]; !exists {
		t.Error("expected validation error for 'no_rawat'")
	}
	if _, exists := valErr["kesadaran"]; !exists {
		t.Error("expected validation error for 'kesadaran'")
	}
	if _, exists := valErr["keluhan"]; !exists {
		t.Error("expected validation error for 'keluhan'")
	}
}

func TestSimpanPemeriksaan_InvalidStatusLanjut(t *testing.T) {
	repo := &mockRepository{}
	rjRepo := &mockRawatJalanService{}
	log := logger.New()
	svc := pemeriksaan.NewService(repo, rjRepo, 48, log)

	req := pemeriksaan.SimpanPemeriksaanRequest{
		NoRawat:             "2026/04/22/036934",
		TanggalPemeriksaan:  "2026-04-23",
		JamPemeriksaan:      "12:10:00",
		Kesadaran:           pemeriksaan.KesadaranComposMentis,
		Keluhan:             "Demam",
		Pemeriksaan:         "Normal",
		Penilaian:           "Febris",
		RencanaTindakLanjut: "Istirahat",
		Instruksi:           "Minum obat",
		Evaluasi:            "Stabil",
	}

	_, err := svc.SimpanPemeriksaan(context.Background(), "DK001", "StatusGhoib", req)
	if err == nil {
		t.Fatal("expected error for invalid status lanjut, got nil")
	}
}

func TestSimpanPemeriksaan_RepoError(t *testing.T) {
	repo := &mockRepository{
		simpanPemeriksaanFunc: func(ctx context.Context, kodeDokter string, statusLanjut shared.StatusLanjut, req pemeriksaan.SimpanPemeriksaanRequest) error {
			return errors.New("db error")
		},
	}
	rjRepo := &mockRawatJalanService{}
	log := logger.New()
	svc := pemeriksaan.NewService(repo, rjRepo, 48, log)

	req := pemeriksaan.SimpanPemeriksaanRequest{
		NoRawat:             "2026/04/22/036934",
		TanggalPemeriksaan:  "2026-04-23",
		JamPemeriksaan:      "12:10:00",
		Kesadaran:           pemeriksaan.KesadaranComposMentis,
		Keluhan:             "Demam",
		Pemeriksaan:         "Normal",
		Penilaian:           "Febris",
		RencanaTindakLanjut: "Istirahat",
		Instruksi:           "Minum obat",
		Evaluasi:            "Stabil",
	}

	_, err := svc.SimpanPemeriksaan(context.Background(), "DK001", shared.StatusLanjutRawatJalan, req)
	if err == nil {
		t.Fatal("expected repo error, got nil")
	}
}

func TestSimpanPemeriksaan_SuhuTubuhValidation(t *testing.T) {
	repo := &mockRepository{}
	rjRepo := &mockRawatJalanService{}
	log := logger.New()
	svc := pemeriksaan.NewService(repo, rjRepo, 48, log)

	baseReq := pemeriksaan.SimpanPemeriksaanRequest{
		NoRawat:             "2026/04/22/036934",
		TanggalPemeriksaan:  "2026-04-23",
		JamPemeriksaan:      "12:10:00",
		Kesadaran:           pemeriksaan.KesadaranComposMentis,
		Keluhan:             "Demam",
		Pemeriksaan:         "Normal",
		Penilaian:           "Febris",
		RencanaTindakLanjut: "Istirahat",
		Instruksi:           "Minum obat",
		Evaluasi:            "Stabil",
	}

	// 1. Non-numeric
	reqNonNumeric := baseReq
	reqNonNumeric.SuhuTubuh = "36.A"
	_, err := svc.SimpanPemeriksaan(context.Background(), "DK001", shared.StatusLanjutRawatJalan, reqNonNumeric)
	if err == nil {
		t.Fatal("expected error for non-numeric suhu tubuh, got nil")
	}

	// 2. Terlalu rendah (< 25)
	reqTooLow := baseReq
	reqTooLow.SuhuTubuh = "20.0"
	_, err = svc.SimpanPemeriksaan(context.Background(), "DK001", shared.StatusLanjutRawatJalan, reqTooLow)
	if err == nil {
		t.Fatal("expected error for too low suhu tubuh, got nil")
	}

	// 3. Terlalu tinggi (> 45)
	reqTooHigh := baseReq
	reqTooHigh.SuhuTubuh = "48.5"
	_, err = svc.SimpanPemeriksaan(context.Background(), "DK001", shared.StatusLanjutRawatJalan, reqTooHigh)
	if err == nil {
		t.Fatal("expected error for too high suhu tubuh, got nil")
	}

	// 4. Valid dengan koma (otomatis diubah ke titik oleh Sanitize)
	reqValidComma := baseReq
	reqValidComma.SuhuTubuh = "36,8"
	_, err = svc.SimpanPemeriksaan(context.Background(), "DK001", shared.StatusLanjutRawatJalan, reqValidComma)
	if err != nil {
		t.Fatalf("expected nil error for valid suhu tubuh with comma, got %v", err)
	}
}

func TestSimpanPemeriksaan_TTVValidation(t *testing.T) {
	repo := &mockRepository{}
	rjRepo := &mockRawatJalanService{}
	log := logger.New()
	svc := pemeriksaan.NewService(repo, rjRepo, 48, log)

	baseReq := pemeriksaan.SimpanPemeriksaanRequest{
		NoRawat:             "2026/04/22/036934",
		TanggalPemeriksaan:  "2026-04-23",
		JamPemeriksaan:      "12:10:00",
		Kesadaran:           pemeriksaan.KesadaranComposMentis,
		Keluhan:             "Demam",
		Pemeriksaan:         "Normal",
		Penilaian:           "Febris",
		RencanaTindakLanjut: "Istirahat",
		Instruksi:           "Minum obat",
		Evaluasi:            "Stabil",
	}

	// 1. Tensi format salah (tanpa slash)
	reqTensiNoSlash := baseReq
	reqTensiNoSlash.Tensi = "12080"
	_, err := svc.SimpanPemeriksaan(context.Background(), "DK001", shared.StatusLanjutRawatJalan, reqTensiNoSlash)
	if err == nil {
		t.Fatal("expected error for tensi without slash, got nil")
	}

	// 2. Tensi sistolik <= diastolik
	reqTensiInvalid := baseReq
	reqTensiInvalid.Tensi = "80/120"
	_, err = svc.SimpanPemeriksaan(context.Background(), "DK001", shared.StatusLanjutRawatJalan, reqTensiInvalid)
	if err == nil {
		t.Fatal("expected error for sistolik <= diastolik, got nil")
	}

	// 3. Tensi valid
	reqTensiValid := baseReq
	reqTensiValid.Tensi = "120/80"
	_, err = svc.SimpanPemeriksaan(context.Background(), "DK001", shared.StatusLanjutRawatJalan, reqTensiValid)
	if err != nil {
		t.Fatalf("expected nil error for valid tensi, got %v", err)
	}

	// 4. Nadi di luar rentang
	reqNadiInvalid := baseReq
	reqNadiInvalid.Nadi = "400"
	_, err = svc.SimpanPemeriksaan(context.Background(), "DK001", shared.StatusLanjutRawatJalan, reqNadiInvalid)
	if err == nil {
		t.Fatal("expected error for nadi > 300, got nil")
	}

	// 5. Respirasi di luar rentang
	reqRespInvalid := baseReq
	reqRespInvalid.Respirasi = "150"
	_, err = svc.SimpanPemeriksaan(context.Background(), "DK001", shared.StatusLanjutRawatJalan, reqRespInvalid)
	if err == nil {
		t.Fatal("expected error for respirasi > 100, got nil")
	}

	// 6. Tinggi badan di luar rentang
	reqTBInvalid := baseReq
	reqTBInvalid.TinggiBadan = "350"
	_, err = svc.SimpanPemeriksaan(context.Background(), "DK001", shared.StatusLanjutRawatJalan, reqTBInvalid)
	if err == nil {
		t.Fatal("expected error for tinggi badan > 250, got nil")
	}

	// 7. Berat badan valid dengan koma (Sanitize)
	reqBBValid := baseReq
	reqBBValid.BeratBadan = "65,5"
	_, err = svc.SimpanPemeriksaan(context.Background(), "DK001", shared.StatusLanjutRawatJalan, reqBBValid)
	if err != nil {
		t.Fatalf("expected nil error for valid berat badan, got %v", err)
	}
}

func TestSimpanPemeriksaan_WaktuMasaDepanValidation(t *testing.T) {
	repo := &mockRepository{}
	rjRepo := &mockRawatJalanService{}
	log := logger.New()
	svc := pemeriksaan.NewService(repo, rjRepo, 48, log)

	baseReq := pemeriksaan.SimpanPemeriksaanRequest{
		NoRawat:             "2026/04/22/036934",
		Kesadaran:           pemeriksaan.KesadaranComposMentis,
		Keluhan:             "Demam",
		Pemeriksaan:         "Normal",
		Penilaian:           "Febris",
		RencanaTindakLanjut: "Istirahat",
		Instruksi:           "Minum obat",
		Evaluasi:            "Stabil",
	}

	// 1. Tanggal besok (masa depan)
	besok := time.Now().AddDate(0, 0, 1)
	reqFutureDate := baseReq
	reqFutureDate.TanggalPemeriksaan = besok.Format("2006-01-02")
	reqFutureDate.JamPemeriksaan = "10:00:00"

	_, err := svc.SimpanPemeriksaan(context.Background(), "DK001", shared.StatusLanjutRawatJalan, reqFutureDate)
	if err == nil {
		t.Fatal("expected error for future examination date, got nil")
	}

	valErr, ok := err.(apperror.ValidationError)
	if !ok {
		t.Fatalf("expected apperror.ValidationError, got %T", err)
	}
	if _, exists := valErr["tanggal_pemeriksaan"]; !exists {
		t.Error("expected validation error on 'tanggal_pemeriksaan'")
	}

	// 2. Hari ini tapi jam 2 jam ke depan (masa depan)
	nanti := time.Now().Add(2 * time.Hour)
	reqFutureTime := baseReq
	reqFutureTime.TanggalPemeriksaan = nanti.Format("2006-01-02")
	reqFutureTime.JamPemeriksaan = nanti.Format("15:04:05")

	_, err = svc.SimpanPemeriksaan(context.Background(), "DK001", shared.StatusLanjutRawatJalan, reqFutureTime)
	if err == nil {
		t.Fatal("expected error for future examination time today, got nil")
	}
}

func TestSimpanPemeriksaan_WaktuSebelumRegistrasi(t *testing.T) {
	repo := &mockRepository{}
	rjRepo := &mockRawatJalanService{
		getWaktuRegistrasiFunc: func(ctx context.Context, noRawat string) (string, string, bool, error) {
			// Registrasi tanggal 2026-04-23 jam 10:00:00
			return "2026-04-23", "10:00:00", true, nil
		},
	}
	log := logger.New()
	svc := pemeriksaan.NewService(repo, rjRepo, 48, log)

	req := pemeriksaan.SimpanPemeriksaanRequest{
		NoRawat:             "2026/04/22/036934",
		TanggalPemeriksaan:  "2026-04-23",
		JamPemeriksaan:      "09:30:00", // Lebih awal dari jam registrasi 10:00:00
		Kesadaran:           pemeriksaan.KesadaranComposMentis,
		Keluhan:             "Demam",
		Pemeriksaan:         "Normal",
		Penilaian:           "Febris",
		RencanaTindakLanjut: "Istirahat",
		Instruksi:           "Minum obat",
		Evaluasi:            "Stabil",
	}

	_, err := svc.SimpanPemeriksaan(context.Background(), "DK001", shared.StatusLanjutRawatJalan, req)
	if err == nil {
		t.Fatal("expected error when examination time is before registration time, got nil")
	}

	valErr, ok := err.(apperror.ValidationError)
	if !ok {
		t.Fatalf("expected apperror.ValidationError, got %T", err)
	}

	if _, exists := valErr["tanggal_pemeriksaan"]; !exists {
		t.Error("expected validation error on 'tanggal_pemeriksaan'")
	}
}

func TestSimpanPemeriksaan_RegistrasiNotFound(t *testing.T) {
	repo := &mockRepository{}
	rjRepo := &mockRawatJalanService{
		getWaktuRegistrasiFunc: func(ctx context.Context, noRawat string) (string, string, bool, error) {
			return "", "", false, nil
		},
	}
	log := logger.New()
	svc := pemeriksaan.NewService(repo, rjRepo, 48, log)

	req := pemeriksaan.SimpanPemeriksaanRequest{
		NoRawat:             "2026/04/22/036934",
		TanggalPemeriksaan:  "2026-04-23",
		JamPemeriksaan:      "12:00:00",
		Kesadaran:           pemeriksaan.KesadaranComposMentis,
		Keluhan:             "Demam",
		Pemeriksaan:         "Normal",
		Penilaian:           "Febris",
		RencanaTindakLanjut: "Istirahat",
		Instruksi:           "Minum obat",
		Evaluasi:            "Stabil",
	}

	_, err := svc.SimpanPemeriksaan(context.Background(), "DK001", shared.StatusLanjutRawatJalan, req)
	if err == nil {
		t.Fatal("expected NotFoundError when registration not found, got nil")
	}
}

func TestSimpanPemeriksaan_DuplicateEntry(t *testing.T) {
	repo := &mockRepository{
		simpanPemeriksaanFunc: func(ctx context.Context, kodeDokter string, statusLanjut shared.StatusLanjut, req pemeriksaan.SimpanPemeriksaanRequest) error {
			return errors.New("Error 1062 (23000): Duplicate entry '2026/04/22/036934-2026-04-23-12:10:10' for key 'PRIMARY'")
		},
	}
	rjRepo := &mockRawatJalanService{}
	log := logger.New()
	svc := pemeriksaan.NewService(repo, rjRepo, 48, log)

	req := pemeriksaan.SimpanPemeriksaanRequest{
		NoRawat:             "2026/04/22/036934",
		TanggalPemeriksaan:  "2026-04-23",
		JamPemeriksaan:      "12:10:10",
		Kesadaran:           pemeriksaan.KesadaranComposMentis,
		Keluhan:             "Demam",
		Pemeriksaan:         "Normal",
		Penilaian:           "Febris",
		RencanaTindakLanjut: "Istirahat",
		Instruksi:           "Minum obat",
		Evaluasi:            "Stabil",
	}

	_, err := svc.SimpanPemeriksaan(context.Background(), "DK001", shared.StatusLanjutRawatJalan, req)
	if err == nil {
		t.Fatal("expected validation error for duplicate entry, got nil")
	}

	valErr, ok := err.(apperror.ValidationError)
	if !ok {
		t.Fatalf("expected apperror.ValidationError, got %T", err)
	}

	if _, exists := valErr["jam_pemeriksaan"]; !exists {
		t.Error("expected validation error for 'jam_pemeriksaan'")
	}
}

func TestHapusPemeriksaan_Success(t *testing.T) {
	hapusCalled := false
	now := time.Now()
	repo := &mockRepository{
		detailPemeriksaanFunc: func(ctx context.Context, id pemeriksaan.IdPemeriksaan, statusLanjut shared.StatusLanjut) (*pemeriksaan.Pemeriksaan, error) {
			return &pemeriksaan.Pemeriksaan{
				NoRawat:            "2026/04/22/036934",
				TanggalPemeriksaan: now.Format("2006-01-02"),
				JamPemeriksaan:     now.Format("15:04:05"),
				KodeDokterPetugas:  "DK001",
				NamaDokterPetugas:  "dr. Handi",
			}, nil
		},
		hapusPemeriksaanFunc: func(ctx context.Context, id pemeriksaan.IdPemeriksaan, statusLanjut shared.StatusLanjut) error {
			hapusCalled = true
			return nil
		},
	}
	rjRepo := &mockRawatJalanService{}
	log := logger.New()
	svc := pemeriksaan.NewService(repo, rjRepo, 48, log)

	id := pemeriksaan.IdPemeriksaan{
		NoRawat:            "2026/04/22/036934",
		TanggalPemeriksaan: now.Format("2006-01-02"),
		JamPemeriksaan:     now.Format("15:04:05"),
	}

	err := svc.HapusPemeriksaan(context.Background(), "DK001", id, shared.StatusLanjutRawatJalan)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !hapusCalled {
		t.Error("expected HapusPemeriksaan repository to be called")
	}
}

func TestHapusPemeriksaan_ForbiddenDifferentDoctor(t *testing.T) {
	now := time.Now()
	repo := &mockRepository{
		detailPemeriksaanFunc: func(ctx context.Context, id pemeriksaan.IdPemeriksaan, statusLanjut shared.StatusLanjut) (*pemeriksaan.Pemeriksaan, error) {
			return &pemeriksaan.Pemeriksaan{
				NoRawat:            "2026/04/22/036934",
				TanggalPemeriksaan: now.Format("2006-01-02"),
				JamPemeriksaan:     now.Format("15:04:05"),
				KodeDokterPetugas:  "DK999", // Dokter lain
				NamaDokterPetugas:  "dr. Lain",
			}, nil
		},
	}
	rjRepo := &mockRawatJalanService{}
	log := logger.New()
	svc := pemeriksaan.NewService(repo, rjRepo, 48, log)

	id := pemeriksaan.IdPemeriksaan{
		NoRawat:            "2026/04/22/036934",
		TanggalPemeriksaan: now.Format("2006-01-02"),
		JamPemeriksaan:     now.Format("15:04:05"),
	}

	err := svc.HapusPemeriksaan(context.Background(), "DK001", id, shared.StatusLanjutRawatJalan)
	if err == nil {
		t.Fatal("expected ForbiddenError when doctor does not match creator, got nil")
	}

	forbiddenErr, ok := err.(*apperror.ForbiddenError)
	if !ok {
		t.Fatalf("expected *apperror.ForbiddenError, got %T", err)
	}

	if forbiddenErr.Message == "" {
		t.Error("expected non-empty forbidden error message")
	}
}

func TestHapusPemeriksaan_MelebihiBatasWaktu(t *testing.T) {
	// Pemeriksaan 50 jam yang lalu (> 48 jam)
	oldTime := time.Now().Add(-50 * time.Hour)
	repo := &mockRepository{
		detailPemeriksaanFunc: func(ctx context.Context, id pemeriksaan.IdPemeriksaan, statusLanjut shared.StatusLanjut) (*pemeriksaan.Pemeriksaan, error) {
			return &pemeriksaan.Pemeriksaan{
				NoRawat:            "2026/04/22/036934",
				TanggalPemeriksaan: oldTime.Format("2006-01-02"),
				JamPemeriksaan:     oldTime.Format("15:04:05"),
				KodeDokterPetugas:  "DK001",
				NamaDokterPetugas:  "dr. Handi",
			}, nil
		},
	}
	rjRepo := &mockRawatJalanService{}
	log := logger.New()
	svc := pemeriksaan.NewService(repo, rjRepo, 48, log)

	id := pemeriksaan.IdPemeriksaan{
		NoRawat:            "2026/04/22/036934",
		TanggalPemeriksaan: oldTime.Format("2006-01-02"),
		JamPemeriksaan:     oldTime.Format("15:04:05"),
	}

	err := svc.HapusPemeriksaan(context.Background(), "DK001", id, shared.StatusLanjutRawatJalan)
	if err == nil {
		t.Fatal("expected ForbiddenError when examination is older than 48 hours, got nil")
	}

	forbiddenErr, ok := err.(*apperror.ForbiddenError)
	if !ok {
		t.Fatalf("expected *apperror.ForbiddenError, got %T", err)
	}

	if forbiddenErr.Message == "" {
		t.Error("expected non-empty forbidden error message")
	}
}

func TestHapusPemeriksaan_NotFound(t *testing.T) {
	repo := &mockRepository{
		detailPemeriksaanFunc: func(ctx context.Context, id pemeriksaan.IdPemeriksaan, statusLanjut shared.StatusLanjut) (*pemeriksaan.Pemeriksaan, error) {
			return nil, nil
		},
	}
	rjRepo := &mockRawatJalanService{}
	log := logger.New()
	svc := pemeriksaan.NewService(repo, rjRepo, 48, log)

	id := pemeriksaan.IdPemeriksaan{
		NoRawat:            "2026/04/22/036934",
		TanggalPemeriksaan: "2026-04-23",
		JamPemeriksaan:     "12:10:10",
	}

	err := svc.HapusPemeriksaan(context.Background(), "DK001", id, shared.StatusLanjutRawatJalan)
	if err == nil {
		t.Fatal("expected NotFoundError when examination does not exist, got nil")
	}
}

func TestHapusPemeriksaan_InvalidStatusLanjut(t *testing.T) {
	repo := &mockRepository{}
	rjRepo := &mockRawatJalanService{}
	log := logger.New()
	svc := pemeriksaan.NewService(repo, rjRepo, 48, log)

	id := pemeriksaan.IdPemeriksaan{
		NoRawat:            "2026/04/22/036934",
		TanggalPemeriksaan: "2026-04-23",
		JamPemeriksaan:     "12:10:10",
	}

	err := svc.HapusPemeriksaan(context.Background(), "DK001", id, shared.StatusLanjut("Semua"))
	if err == nil {
		t.Fatal("expected error for invalid status lanjut, got nil")
	}
}






