package laboratorium_test

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"erm-dokter/internal/berkasdigital"
	"erm-dokter/internal/laboratorium"
	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type mockRepository struct {
	daftarHasilLabPKFn     func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratoriumDB, int, error)
	daftarHasilLabPKByRMFn func(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratoriumDB, int, error)
	detailHasilLabPKFn     func(ctx context.Context, noRawat string, kodeTindakan string, tanggalPeriksa string, jamPeriksa string) (*laboratorium.HasilLaboratoriumDB, error)

	daftarHasilLabPAFn     func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratoriumDB, int, error)
	daftarHasilLabPAByRMFn func(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratoriumDB, int, error)
	detailHasilLabPAFn     func(ctx context.Context, noRawat string, kodeTindakan string, tanggalPeriksa string, jamPeriksa string) (*laboratorium.HasilLaboratoriumDB, error)

	daftarHasilLabMBFn     func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratoriumDB, int, error)
	daftarHasilLabMBByRMFn func(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratoriumDB, int, error)
	detailHasilLabMBFn     func(ctx context.Context, noRawat string, kodeTindakan string, tanggalPeriksa string, jamPeriksa string) (*laboratorium.HasilLaboratoriumDB, error)
}

func (m *mockRepository) DaftarHasilLabPK(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratoriumDB, int, error) {
	if m.daftarHasilLabPKFn != nil {
		return m.daftarHasilLabPKFn(ctx, noRawat, statusLanjut, filter)
	}
	return nil, 0, nil
}

func (m *mockRepository) DaftarHasilLabPKByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratoriumDB, int, error) {
	if m.daftarHasilLabPKByRMFn != nil {
		return m.daftarHasilLabPKByRMFn(ctx, noRM, statusLanjut, filter)
	}
	return nil, 0, nil
}

func (m *mockRepository) DetailHasilLabPK(ctx context.Context, noRawat string, kodeTindakan string, tanggalPeriksa string, jamPeriksa string) (*laboratorium.HasilLaboratoriumDB, error) {
	if m.detailHasilLabPKFn != nil {
		return m.detailHasilLabPKFn(ctx, noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa)
	}
	return nil, nil
}

func (m *mockRepository) DaftarHasilLabPA(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratoriumDB, int, error) {
	if m.daftarHasilLabPAFn != nil {
		return m.daftarHasilLabPAFn(ctx, noRawat, statusLanjut, filter)
	}
	return nil, 0, nil
}

func (m *mockRepository) DaftarHasilLabPAByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratoriumDB, int, error) {
	if m.daftarHasilLabPAByRMFn != nil {
		return m.daftarHasilLabPAByRMFn(ctx, noRM, statusLanjut, filter)
	}
	return nil, 0, nil
}

func (m *mockRepository) DetailHasilLabPA(ctx context.Context, noRawat string, kodeTindakan string, tanggalPeriksa string, jamPeriksa string) (*laboratorium.HasilLaboratoriumDB, error) {
	if m.detailHasilLabPAFn != nil {
		return m.detailHasilLabPAFn(ctx, noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa)
	}
	return nil, nil
}

func (m *mockRepository) DaftarHasilLabMB(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratoriumDB, int, error) {
	if m.daftarHasilLabMBFn != nil {
		return m.daftarHasilLabMBFn(ctx, noRawat, statusLanjut, filter)
	}
	return nil, 0, nil
}

func (m *mockRepository) DaftarHasilLabMBByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratoriumDB, int, error) {
	if m.daftarHasilLabMBByRMFn != nil {
		return m.daftarHasilLabMBByRMFn(ctx, noRM, statusLanjut, filter)
	}
	return nil, 0, nil
}

func (m *mockRepository) DetailHasilLabMB(ctx context.Context, noRawat string, kodeTindakan string, tanggalPeriksa string, jamPeriksa string) (*laboratorium.HasilLaboratoriumDB, error) {
	if m.detailHasilLabMBFn != nil {
		return m.detailHasilLabMBFn(ctx, noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa)
	}
	return nil, nil
}

type mockBerkasService struct {
	getMasterBerkasFn     func(ctx context.Context) ([]berkasdigital.MasterBerkasDigital, error)
	getBerkasByNoRawatFn  func(ctx context.Context, noRawat string, kodeList []string) ([]berkasdigital.BerkasDigitalDB, error)
	streamBerkasDigitalFn func(ctx context.Context, encryptedIdBerkas string, w http.ResponseWriter) error
	buildBerkasItemFn     func(kode string, namaBerkas string, lokasiFile string) (*berkasdigital.BerkasDigitalItem, error)
	buildFullURLFn        func(lokasiFile string) string
}

func (m *mockBerkasService) GetMasterBerkas(ctx context.Context) ([]berkasdigital.MasterBerkasDigital, error) {
	if m.getMasterBerkasFn != nil {
		return m.getMasterBerkasFn(ctx)
	}
	return nil, nil
}

func (m *mockBerkasService) GetBerkasByNoRawat(ctx context.Context, noRawat string, kodeList []string) ([]berkasdigital.BerkasDigitalDB, error) {
	if m.getBerkasByNoRawatFn != nil {
		return m.getBerkasByNoRawatFn(ctx, noRawat, kodeList)
	}
	return nil, nil
}

func (m *mockBerkasService) StreamBerkasDigital(ctx context.Context, encryptedIdBerkas string, w http.ResponseWriter) error {
	if m.streamBerkasDigitalFn != nil {
		return m.streamBerkasDigitalFn(ctx, encryptedIdBerkas, w)
	}
	return nil
}

func (m *mockBerkasService) BuildBerkasItem(kode string, namaBerkas string, lokasiFile string) (*berkasdigital.BerkasDigitalItem, error) {
	if m.buildBerkasItemFn != nil {
		return m.buildBerkasItemFn(kode, namaBerkas, lokasiFile)
	}
	return nil, nil
}

func (m *mockBerkasService) BuildFullURL(lokasiFile string) string {
	if m.buildFullURLFn != nil {
		return m.buildFullURLFn(lokasiFile)
	}
	return ""
}

const testJWTSecret = "secret-key-32-bytes-testing-12345"
const testURLBerkas = "http://192.168.30.24/webapps/berkasrawat/"

func TestService_GetRiwayatLabKunjungan_PK_Success(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		daftarHasilLabPKFn: func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratoriumDB, int, error) {
			if noRawat != "2026/04/22/000001" {
				t.Errorf("Expected noRawat 2026/04/22/000001, got %s", noRawat)
			}
			return []laboratorium.HasilLaboratoriumDB{
				{
					NoRawat:           "2026/04/22/000001",
					KodeTindakan:      "LAB001",
					NamaTindakan:      "DARAH LENGKAP",
					Kategori:          "PK",
					Status:            "Ralan",
					TanggalPeriksa:    "2026-04-22",
					JamPeriksa:        "10:00:00",
					KodeDokterPerujuk: "DR001",
					NamaDokterPerujuk: "dr. Handi",
					KodeDokterPJ:      "DR002",
					NamaDokterPJ:      "dr. Budi, Sp.PK",
					NipPetugas:        "P001",
					NamaPetugas:       "Siti Analis",
					DetailPK: []laboratorium.ItemHasilLabPKDB{
						{
							IdTemplate:      101,
							NamaPemeriksaan: "Hemoglobin",
							Nilai:           "14.2",
							Satuan:          "g/dL",
							NilaiRujukan:    "13.5 - 17.5",
							Keterangan:      "Normal",
						},
					},
				},
			}, 1, nil
		},
	}
	mockBerkas := &mockBerkasService{
		getBerkasByNoRawatFn: func(ctx context.Context, noRawat string, kodeList []string) ([]berkasdigital.BerkasDigitalDB, error) {
			return []berkasdigital.BerkasDigitalDB{}, nil
		},
	}

	encNoRawat, _ := crypto.Encrypt("2026/04/22/000001", testJWTSecret)
	svc := laboratorium.NewService(mockRepo, mockBerkas, testJWTSecret, testURLBerkas, []string{"005"}, []string{"015"}, []string{"010"}, log)

	res, meta, err := svc.GetRiwayatLabKunjungan(context.Background(), shared.KategoriLabPK, encNoRawat, shared.StatusLanjut("Semua"), laboratorium.FilterRiwayatLab{Page: 1, Limit: 5})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(res.HasilPemeriksaan) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(res.HasilPemeriksaan))
	}

	if meta.TotalRecords != 1 {
		t.Errorf("Expected total records 1, got %d", meta.TotalRecords)
	}

	item := res.HasilPemeriksaan[0]
	if item.NamaTindakan != "DARAH LENGKAP" || len(item.DetailPK) != 1 {
		t.Errorf("Unexpected mapped item: %+v", item)
	}
	if item.DetailPK[0].NamaPemeriksaan != "Hemoglobin" || item.DetailPK[0].Nilai != "14.2" {
		t.Errorf("Unexpected detail PK: %+v", item.DetailPK[0])
	}
}

func TestService_GetRiwayatLabKunjungan_PA_Success_WithBerkas(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		daftarHasilLabPAFn: func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratoriumDB, int, error) {
			return []laboratorium.HasilLaboratoriumDB{
				{
					NoRawat:           "2026/04/22/000002",
					KodeTindakan:      "PA001",
					NamaTindakan:      "BIOPSI JARINGAN",
					Kategori:          "PA",
					Status:            "Ranap",
					TanggalPeriksa:    "2026-04-22",
					JamPeriksa:        "11:00:00",
					KodeDokterPerujuk: "DR001",
					NamaDokterPerujuk: "dr. Handi",
					KodeDokterPJ:      "DR003",
					NamaDokterPJ:      "dr. Ratna, Sp.PA",
					NipPetugas:        "P002",
					NamaPetugas:       "Agus Analis",
					DetailPA: &laboratorium.HasilLabPADB{
						DiagnosaKlinik: "Tumor Coli",
						Makroskopis:    "Jaringan kenyal...",
						Mikroskopis:    "Tampak sel atipik...",
						Kesimpulan:     "Adenocarcinoma",
						Kesan:          "Maligna",
					},
				},
			}, 1, nil
		},
	}
	mockBerkas := &mockBerkasService{
		getBerkasByNoRawatFn: func(ctx context.Context, noRawat string, kodeList []string) ([]berkasdigital.BerkasDigitalDB, error) {
			return []berkasdigital.BerkasDigitalDB{
				{
					Kode:       "015",
					NamaBerkas: "Hasil PA FNAB",
					LokasiFile: "pages/upload/fnab_hasil.pdf",
				},
			}, nil
		},
	}

	encNoRawat, _ := crypto.Encrypt("2026/04/22/000002", testJWTSecret)
	svc := laboratorium.NewService(mockRepo, mockBerkas, testJWTSecret, testURLBerkas, []string{"005"}, []string{"015"}, []string{"010"}, log)

	res, _, err := svc.GetRiwayatLabKunjungan(context.Background(), shared.KategoriLabPA, encNoRawat, shared.StatusLanjutRawatInap, laboratorium.FilterRiwayatLab{Page: 1, Limit: 5})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(res.HasilPemeriksaan) != 1 || res.HasilPemeriksaan[0].DetailPA == nil {
		t.Fatalf("Expected 1 PA result with DetailPA, got %+v", res)
	}

	if res.HasilPemeriksaan[0].DetailPA.Kesimpulan != "Adenocarcinoma" {
		t.Errorf("Expected Kesimpulan 'Adenocarcinoma', got %s", res.HasilPemeriksaan[0].DetailPA.Kesimpulan)
	}

	if len(res.BerkasDigital) != 1 {
		t.Fatalf("Expected 1 berkas digital, got %d", len(res.BerkasDigital))
	}

	// Check URL berkas digital points to /api/v1/berkas-digital/{id}
	if res.BerkasDigital[0].UrlBerkas != "/api/v1/berkas-digital/"+res.BerkasDigital[0].IdBerkas {
		t.Errorf("Expected url_berkas to start with /api/v1/berkas-digital/, got %s", res.BerkasDigital[0].UrlBerkas)
	}

	decURL, errDec := crypto.Decrypt(res.BerkasDigital[0].IdBerkas, testJWTSecret)
	if errDec != nil || decURL != "http://192.168.30.24/webapps/berkasrawat/pages/upload/fnab_hasil.pdf" {
		t.Errorf("Expected decrypted URL http://192.168.30.24/webapps/berkasrawat/pages/upload/fnab_hasil.pdf, got %s, err: %v", decURL, errDec)
	}
}

func TestService_GetRiwayatLabKunjungan_MB_MultiTindakan_NoDuplicateBerkas(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		daftarHasilLabMBFn: func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratoriumDB, int, error) {
			return []laboratorium.HasilLaboratoriumDB{
				{
					NoRawat:           "2026/04/22/000003",
					KodeTindakan:      "MB001",
					NamaTindakan:      "KULTUR URIN DAN SENSITIVITAS",
					Kategori:          "MB",
					Status:            "Ralan",
					TanggalPeriksa:    "2026-04-22",
					JamPeriksa:        "09:30:00",
					KodeDokterPerujuk: "DR001",
					NamaDokterPerujuk: "dr. Handi",
					KodeDokterPJ:      "DR004",
					NamaDokterPJ:      "dr. Maya, Sp.MK",
					NipPetugas:        "P003",
					NamaPetugas:       "Rani Analis",
					DetailPK:          []laboratorium.ItemHasilLabPKDB{},
				},
				{
					NoRawat:           "2026/04/22/000003",
					KodeTindakan:      "MB002",
					NamaTindakan:      "PEWARNAAN GRAM",
					Kategori:          "MB",
					Status:            "Ralan",
					TanggalPeriksa:    "2026-04-22",
					JamPeriksa:        "10:00:00",
					KodeDokterPerujuk: "DR001",
					NamaDokterPerujuk: "dr. Handi",
					KodeDokterPJ:      "DR004",
					NamaDokterPJ:      "dr. Maya, Sp.MK",
					NipPetugas:        "P003",
					NamaPetugas:       "Rani Analis",
					DetailPK:          []laboratorium.ItemHasilLabPKDB{},
				},
			}, 2, nil
		},
	}
	mockBerkas := &mockBerkasService{
		getBerkasByNoRawatFn: func(ctx context.Context, noRawat string, kodeList []string) ([]berkasdigital.BerkasDigitalDB, error) {
			return []berkasdigital.BerkasDigitalDB{
				{
					Kode:       "010",
					NamaBerkas: "HASIL KULTUR MIKROBIOLOGI",
					LokasiFile: "pages/upload/kultur_urin.pdf",
				},
				{
					Kode:       "010",
					NamaBerkas: "HASIL RESISTENSI MIKROBIOLOGI",
					LokasiFile: "pages/upload/resistensi.pdf",
				},
			}, nil
		},
	}

	encNoRawat, _ := crypto.Encrypt("2026/04/22/000003", testJWTSecret)
	svc := laboratorium.NewService(mockRepo, mockBerkas, testJWTSecret, testURLBerkas, []string{"005"}, []string{"015"}, []string{"010"}, log)

	res, meta, err := svc.GetRiwayatLabKunjungan(context.Background(), shared.KategoriLabMB, encNoRawat, shared.StatusLanjut("Semua"), laboratorium.FilterRiwayatLab{Page: 1, Limit: 5})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(res.HasilPemeriksaan) != 2 || meta.TotalRecords != 2 {
		t.Fatalf("Expected 2 MB records, got len=%d total=%d", len(res.HasilPemeriksaan), meta.TotalRecords)
	}

	// 2 distinct tests
	if res.HasilPemeriksaan[0].KodeTindakan != "MB001" || res.HasilPemeriksaan[1].KodeTindakan != "MB002" {
		t.Errorf("Unexpected MB tindakan data: %+v", res.HasilPemeriksaan)
	}

	// Exactly 2 berkas digital in root level (NOT 4 / 2x2 duplicate)
	if len(res.BerkasDigital) != 2 {
		t.Fatalf("Expected exactly 2 berkas digital at root level without duplicate, got %d", len(res.BerkasDigital))
	}
}

func TestService_GetRiwayatLabPasien_Success(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		daftarHasilLabPKByRMFn: func(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratoriumDB, int, error) {
			if noRM != "123456" {
				t.Errorf("Expected noRM 123456, got %s", noRM)
			}
			return []laboratorium.HasilLaboratoriumDB{
				{
					NoRawat:        "2026/04/22/000001",
					KodeTindakan:   "LAB001",
					NamaTindakan:   "DARAH LENGKAP",
					Kategori:       "PK",
					TanggalPeriksa: "2026-04-22",
					JamPeriksa:     "10:00:00",
				},
			}, 1, nil
		},
	}
	mockBerkas := &mockBerkasService{}

	encNoRM, _ := crypto.Encrypt("123456", testJWTSecret)
	svc := laboratorium.NewService(mockRepo, mockBerkas, testJWTSecret, testURLBerkas, []string{"005"}, []string{"015"}, []string{"010"}, log)

	res, meta, err := svc.GetRiwayatLabPasien(context.Background(), shared.KategoriLabPK, encNoRM, shared.StatusLanjut("Semua"), laboratorium.FilterRiwayatLab{Page: 1, Limit: 5})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(res) != 1 || meta.TotalRecords != 1 {
		t.Errorf("Expected 1 record, got %+v", res)
	}
}

func TestService_GetDetailHasilLab_Success(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		detailHasilLabPKFn: func(ctx context.Context, noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa string) (*laboratorium.HasilLaboratoriumDB, error) {
			if noRawat != "2026/04/22/000001" || kodeTindakan != "LAB001" {
				t.Errorf("Unexpected params: %s %s", noRawat, kodeTindakan)
			}
			return &laboratorium.HasilLaboratoriumDB{
				NoRawat:        "2026/04/22/000001",
				KodeTindakan:   "LAB001",
				NamaTindakan:   "DARAH LENGKAP",
				Kategori:       "PK",
				TanggalPeriksa: "2026-04-22",
				JamPeriksa:     "10:00:00",
			}, nil
		},
	}
	mockBerkas := &mockBerkasService{}

	encNoRawat, _ := crypto.Encrypt("2026/04/22/000001", testJWTSecret)
	encHasil, _ := crypto.Encrypt("2026/04/22/000001~LAB001~2026-04-22~10:00:00", testJWTSecret)

	svc := laboratorium.NewService(mockRepo, mockBerkas, testJWTSecret, testURLBerkas, []string{"005"}, []string{"015"}, []string{"010"}, log)
	detail, err := svc.GetDetailHasilLab(context.Background(), shared.KategoriLabPK, encNoRawat, encHasil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if detail.NamaTindakan != "DARAH LENGKAP" {
		t.Errorf("Expected NamaTindakan 'DARAH LENGKAP', got %s", detail.NamaTindakan)
	}
}

func TestService_GetDetailHasilLab_MismatchNoRawat(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{}
	mockBerkas := &mockBerkasService{}

	encNoRawat, _ := crypto.Encrypt("2026/04/22/000001", testJWTSecret)
	encHasil, _ := crypto.Encrypt("2026/04/22/999999~LAB001~2026-04-22~10:00:00", testJWTSecret)

	svc := laboratorium.NewService(mockRepo, mockBerkas, testJWTSecret, testURLBerkas, []string{"005"}, []string{"015"}, []string{"010"}, log)
	_, err := svc.GetDetailHasilLab(context.Background(), shared.KategoriLabPK, encNoRawat, encHasil)
	if err == nil {
		t.Fatalf("Expected error for mismatched no_rawat, got nil")
	}

	var businessErr *apperror.BusinessError
	if !errors.As(err, &businessErr) {
		t.Errorf("Expected BusinessError, got %T: %v", err, err)
	}
}

func TestService_GetDetailHasilLab_NotFound(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		detailHasilLabPKFn: func(ctx context.Context, noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa string) (*laboratorium.HasilLaboratoriumDB, error) {
			return nil, sql.ErrNoRows
		},
	}
	mockBerkas := &mockBerkasService{}

	encNoRawat, _ := crypto.Encrypt("2026/04/22/000001", testJWTSecret)
	encHasil, _ := crypto.Encrypt("2026/04/22/000001~LAB001~2026-04-22~10:00:00", testJWTSecret)

	svc := laboratorium.NewService(mockRepo, mockBerkas, testJWTSecret, testURLBerkas, []string{"005"}, []string{"015"}, []string{"010"}, log)
	_, err := svc.GetDetailHasilLab(context.Background(), shared.KategoriLabPK, encNoRawat, encHasil)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	var notFound *apperror.NotFoundError
	if !errors.As(err, &notFound) {
		t.Errorf("Expected NotFoundError, got %T: %v", err, err)
	}
}
