package laboratorium_test

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"testing"
	"time"

	"erm-dokter/internal/berkasdigital"
	"erm-dokter/internal/laboratorium"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatjalan"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
	"erm-dokter/internal/tindakan"
)

type mockRepository struct {
	daftarHasilLabPKFn     func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratorium, int, error)
	daftarHasilLabPKByRMFn func(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratorium, int, error)
	detailHasilLabPKFn     func(ctx context.Context, noRawat string, kodeTindakan string, tanggalPeriksa string, jamPeriksa string) (*laboratorium.HasilLaboratorium, error)

	daftarHasilLabPAFn     func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratorium, int, error)
	daftarHasilLabPAByRMFn func(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratorium, int, error)
	detailHasilLabPAFn     func(ctx context.Context, noRawat string, kodeTindakan string, tanggalPeriksa string, jamPeriksa string) (*laboratorium.HasilLaboratorium, error)

	daftarHasilLabMBFn     func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratorium, int, error)
	daftarHasilLabMBByRMFn func(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratorium, int, error)
	detailHasilLabMBFn     func(ctx context.Context, noRawat string, kodeTindakan string, tanggalPeriksa string, jamPeriksa string) (*laboratorium.HasilLaboratorium, error)

	getKunjunganForPermintaanPKFn func(ctx context.Context, noRawat string) (*laboratorium.KunjunganInfoLabPK, error)
	cekStatusKamarInapFn          func(ctx context.Context, noRawat string) (bool, bool, error)
	simpanPermintaanLabPKFn       func(ctx context.Context, noRawat string, kodeDokter string, status string, req laboratorium.SimpanPermintaanLabPKRequest, kodeTindakanList []string, templateMap map[string][]int) (string, error)
	daftarPermintaanLabPKFn       func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) ([]laboratorium.PermintaanLabPK, error)
	daftarPermintaanLabPKByRMFn   func(ctx context.Context, noRkmMedis string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.PermintaanLabPK, int, error)
	detailPermintaanLabPKFn       func(ctx context.Context, noPermintaan string) (*laboratorium.DetailPermintaanLabPK, error)
	hapusPermintaanLabPKFn        func(ctx context.Context, noPermintaan string) error
}

func (m *mockRepository) DaftarHasilLabPK(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratorium, int, error) {
	if m.daftarHasilLabPKFn != nil {
		return m.daftarHasilLabPKFn(ctx, noRawat, statusLanjut, filter)
	}
	return nil, 0, nil
}

func (m *mockRepository) DaftarHasilLabPKByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratorium, int, error) {
	if m.daftarHasilLabPKByRMFn != nil {
		return m.daftarHasilLabPKByRMFn(ctx, noRM, statusLanjut, filter)
	}
	return nil, 0, nil
}

func (m *mockRepository) DetailHasilLabPK(ctx context.Context, noRawat string, kodeTindakan string, tanggalPeriksa string, jamPeriksa string) (*laboratorium.HasilLaboratorium, error) {
	if m.detailHasilLabPKFn != nil {
		return m.detailHasilLabPKFn(ctx, noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa)
	}
	return nil, nil
}

func (m *mockRepository) DaftarHasilLabPA(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratorium, int, error) {
	if m.daftarHasilLabPAFn != nil {
		return m.daftarHasilLabPAFn(ctx, noRawat, statusLanjut, filter)
	}
	return nil, 0, nil
}

func (m *mockRepository) DaftarHasilLabPAByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratorium, int, error) {
	if m.daftarHasilLabPAByRMFn != nil {
		return m.daftarHasilLabPAByRMFn(ctx, noRM, statusLanjut, filter)
	}
	return nil, 0, nil
}

func (m *mockRepository) DetailHasilLabPA(ctx context.Context, noRawat string, kodeTindakan string, tanggalPeriksa string, jamPeriksa string) (*laboratorium.HasilLaboratorium, error) {
	if m.detailHasilLabPAFn != nil {
		return m.detailHasilLabPAFn(ctx, noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa)
	}
	return nil, nil
}

func (m *mockRepository) DaftarHasilLabMB(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratorium, int, error) {
	if m.daftarHasilLabMBFn != nil {
		return m.daftarHasilLabMBFn(ctx, noRawat, statusLanjut, filter)
	}
	return nil, 0, nil
}

func (m *mockRepository) DaftarHasilLabMBByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratorium, int, error) {
	if m.daftarHasilLabMBByRMFn != nil {
		return m.daftarHasilLabMBByRMFn(ctx, noRM, statusLanjut, filter)
	}
	return nil, 0, nil
}

func (m *mockRepository) DetailHasilLabMB(ctx context.Context, noRawat string, kodeTindakan string, tanggalPeriksa string, jamPeriksa string) (*laboratorium.HasilLaboratorium, error) {
	if m.detailHasilLabMBFn != nil {
		return m.detailHasilLabMBFn(ctx, noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa)
	}
	return nil, nil
}

func (m *mockRepository) GetKunjunganForPermintaanPK(ctx context.Context, noRawat string) (*laboratorium.KunjunganInfoLabPK, error) {
	if m.getKunjunganForPermintaanPKFn != nil {
		return m.getKunjunganForPermintaanPKFn(ctx, noRawat)
	}
	return nil, nil
}

func (m *mockRepository) CekStatusKamarInap(ctx context.Context, noRawat string) (bool, bool, error) {
	if m.cekStatusKamarInapFn != nil {
		return m.cekStatusKamarInapFn(ctx, noRawat)
	}
	return false, false, nil
}

func (m *mockRepository) SimpanPermintaanLabPK(ctx context.Context, noRawat string, kodeDokter string, status string, req laboratorium.SimpanPermintaanLabPKRequest, kodeTindakanList []string, templateMap map[string][]int) (string, error) {
	if m.simpanPermintaanLabPKFn != nil {
		return m.simpanPermintaanLabPKFn(ctx, noRawat, kodeDokter, status, req, kodeTindakanList, templateMap)
	}
	return "", nil
}

func (m *mockRepository) DaftarPermintaanLabPK(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) ([]laboratorium.PermintaanLabPK, error) {
	if m.daftarPermintaanLabPKFn != nil {
		return m.daftarPermintaanLabPKFn(ctx, noRawat, statusLanjut)
	}
	return nil, nil
}

func (m *mockRepository) DaftarPermintaanLabPKByRM(ctx context.Context, noRkmMedis string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.PermintaanLabPK, int, error) {
	if m.daftarPermintaanLabPKByRMFn != nil {
		return m.daftarPermintaanLabPKByRMFn(ctx, noRkmMedis, statusLanjut, filter)
	}
	return nil, 0, nil
}

func (m *mockRepository) DetailPermintaanLabPK(ctx context.Context, noPermintaan string) (*laboratorium.DetailPermintaanLabPK, error) {
	if m.detailPermintaanLabPKFn != nil {
		return m.detailPermintaanLabPKFn(ctx, noPermintaan)
	}
	return nil, nil
}

func (m *mockRepository) HapusPermintaanLabPK(ctx context.Context, noPermintaan string) error {
	if m.hapusPermintaanLabPKFn != nil {
		return m.hapusPermintaanLabPKFn(ctx, noPermintaan)
	}
	return nil
}

type mockBerkasService struct {
	getBerkasByNoRawatFn func(ctx context.Context, noRawat string, kodeList []string) ([]berkasdigital.BerkasDigitalDB, error)
	buildBerkasItemFn    func(kode string, namaBerkas string, lokasiFile string) (*berkasdigital.BerkasDigitalItem, error)
}

func (m *mockBerkasService) GetMasterBerkas(ctx context.Context) ([]berkasdigital.MasterBerkasDigital, error) {
	return nil, nil
}

func (m *mockBerkasService) GetBerkasByNoRawat(ctx context.Context, noRawat string, kodeList []string) ([]berkasdigital.BerkasDigitalDB, error) {
	if m.getBerkasByNoRawatFn != nil {
		return m.getBerkasByNoRawatFn(ctx, noRawat, kodeList)
	}
	return nil, nil
}

func (m *mockBerkasService) StreamBerkasDigital(ctx context.Context, encryptedIdBerkas string, w http.ResponseWriter) error {
	return nil
}

func (m *mockBerkasService) BuildBerkasItem(kode string, namaBerkas string, lokasiFile string) (*berkasdigital.BerkasDigitalItem, error) {
	if m.buildBerkasItemFn != nil {
		return m.buildBerkasItemFn(kode, namaBerkas, lokasiFile)
	}
	return &berkasdigital.BerkasDigitalItem{
		Kode:       kode,
		NamaBerkas: namaBerkas,
		IdBerkas:   "enc-" + kode,
		UrlBerkas:  "/api/v1/berkas-digital/enc-" + kode,
	}, nil
}

func (m *mockBerkasService) BuildFullURL(lokasiFile string) string {
	return lokasiFile
}

type mockRawatJalanService struct {
	getWaktuRegistrasiFn func(ctx context.Context, noRawat string) (tanggal string, jam string, exists bool, err error)
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

func (m *mockRawatJalanService) GetWaktuRegistrasi(ctx context.Context, noRawat string) (tanggal string, jam string, exists bool, err error) {
	if m.getWaktuRegistrasiFn != nil {
		return m.getWaktuRegistrasiFn(ctx, noRawat)
	}
	return "2026-09-03", "08:00:00", true, nil
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

type mockTindakanService struct {
	cekKeberadaanTindakanLabFn func(ctx context.Context, kategori shared.KategoriLab, listKodeTindakan []string) (map[string]bool, error)
	cekKeberadaanTemplateLabFn func(ctx context.Context, listKodeTindakan []string, templateMap map[string][]int) (map[string]map[int]bool, error)
}

func (m *mockTindakanService) GetDaftarTindakanLab(ctx context.Context, kategori shared.KategoriLab, filter tindakan.FilterDaftarTindakanLab) ([]tindakan.TindakanLab, shared.PaginationMeta, error) {
	return nil, shared.PaginationMeta{}, nil
}

func (m *mockTindakanService) GetDetailTindakanLab(ctx context.Context, kategori shared.KategoriLab, encryptedId string) (*tindakan.DetailTindakanLab, error) {
	return nil, nil
}

func (m *mockTindakanService) CekKeberadaanTindakanLab(ctx context.Context, kategori shared.KategoriLab, listKodeTindakan []string) (map[string]bool, error) {
	if m.cekKeberadaanTindakanLabFn != nil {
		return m.cekKeberadaanTindakanLabFn(ctx, kategori, listKodeTindakan)
	}
	res := make(map[string]bool)
	for _, k := range listKodeTindakan {
		res[k] = true
	}
	return res, nil
}

func (m *mockTindakanService) CekKeberadaanTemplateLab(ctx context.Context, listKodeTindakan []string, templateMap map[string][]int) (map[string]map[int]bool, error) {
	if m.cekKeberadaanTemplateLabFn != nil {
		return m.cekKeberadaanTemplateLabFn(ctx, listKodeTindakan, templateMap)
	}
	res := make(map[string]map[int]bool)
	for _, k := range listKodeTindakan {
		res[k] = make(map[int]bool)
		for _, id := range templateMap[k] {
			res[k][id] = true
		}
	}
	return res, nil
}

const testURLBerkas = "http://192.168.30.24/webapps/berkasrawat/"

func TestService_GetRiwayatLabKunjungan_PK_Success(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		daftarHasilLabPKFn: func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratorium, int, error) {
			if noRawat != "2026/04/22/000001" {
				t.Errorf("Expected noRawat 2026/04/22/000001, got %s", noRawat)
			}
			return []laboratorium.HasilLaboratorium{
				{
					NoRawat:        "2026/04/22/000001",
					KodeTindakan:   "LAB001",
					NamaTindakan:   "DARAH LENGKAP",
					Kategori:       "PK",
					Status:         "Ralan",
					TanggalPeriksa: "2026-04-22",
					JamPeriksa:     "10:00:00",
					DokterPerujuk: laboratorium.DokterInfo{
						KodeDokter: "DR001",
						NamaDokter: "dr. Handi",
					},
					DokterPJ: laboratorium.DokterInfo{
						KodeDokter: "DR002",
						NamaDokter: "dr. Budi, Sp.PK",
					},
					Petugas: laboratorium.PetugasInfo{
						Nip:  "P001",
						Nama: "Siti Analis",
					},
					DetailPK: []laboratorium.ItemHasilLabPK{
						{
							IdTemplate:      "101",
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
	mockRawatJalan := &mockRawatJalanService{}
	mockTindakan := &mockTindakanService{}

	svc := laboratorium.NewService(mockRepo, mockBerkas, mockRawatJalan, mockTindakan, 48, testURLBerkas, []string{"005"}, []string{"015"}, []string{"010"}, log)

	res, meta, err := svc.GetRiwayatLabKunjungan(context.Background(), shared.KategoriLabPK, "2026/04/22/000001", shared.StatusLanjut("Semua"), laboratorium.FilterRiwayatLab{Page: 1, Limit: 5})
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
		daftarHasilLabPAFn: func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.HasilLaboratorium, int, error) {
			return []laboratorium.HasilLaboratorium{
				{
					NoRawat:        "2026/04/22/000002",
					KodeTindakan:   "PA001",
					NamaTindakan:   "BIOPSI JARINGAN",
					Kategori:       "PA",
					Status:         "Ranap",
					TanggalPeriksa: "2026-04-22",
					JamPeriksa:     "11:00:00",
					DokterPerujuk: laboratorium.DokterInfo{
						KodeDokter: "DR001",
						NamaDokter: "dr. Handi",
					},
					DokterPJ: laboratorium.DokterInfo{
						KodeDokter: "DR003",
						NamaDokter: "dr. Ratna, Sp.PA",
					},
					Petugas: laboratorium.PetugasInfo{
						Nip:  "P002",
						Nama: "Agus Analis",
					},
					DetailPA: &laboratorium.HasilLabPA{
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
	mockRawatJalan := &mockRawatJalanService{}
	mockTindakan := &mockTindakanService{}

	svc := laboratorium.NewService(mockRepo, mockBerkas, mockRawatJalan, mockTindakan, 48, testURLBerkas, []string{"005"}, []string{"015"}, []string{"010"}, log)

	res, _, err := svc.GetRiwayatLabKunjungan(context.Background(), shared.KategoriLabPA, "2026/04/22/000002", shared.StatusLanjutRawatInap, laboratorium.FilterRiwayatLab{Page: 1, Limit: 5})
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
}

func TestService_GetDetailHasilLab_Success(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		detailHasilLabPKFn: func(ctx context.Context, noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa string) (*laboratorium.HasilLaboratorium, error) {
			return &laboratorium.HasilLaboratorium{
				NoRawat:        noRawat,
				KodeTindakan:   kodeTindakan,
				NamaTindakan:   "DARAH LENGKAP",
				Kategori:       "PK",
				TanggalPeriksa: tanggalPeriksa,
				JamPeriksa:     jamPeriksa,
			}, nil
		},
	}
	mockBerkas := &mockBerkasService{}
	mockRawatJalan := &mockRawatJalanService{}
	mockTindakan := &mockTindakanService{}

	svc := laboratorium.NewService(mockRepo, mockBerkas, mockRawatJalan, mockTindakan, 48, testURLBerkas, []string{"005"}, []string{"015"}, []string{"010"}, log)
	detail, err := svc.GetDetailHasilLab(context.Background(), shared.KategoriLabPK, "2026/04/22/000001", "LAB001", "2026-04-22", "10:00:00")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if detail.NamaTindakan != "DARAH LENGKAP" {
		t.Errorf("Expected NamaTindakan 'DARAH LENGKAP', got %s", detail.NamaTindakan)
	}
}

func TestService_GetDetailHasilLab_NotFound(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		detailHasilLabPKFn: func(ctx context.Context, noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa string) (*laboratorium.HasilLaboratorium, error) {
			return nil, sql.ErrNoRows
		},
	}
	mockBerkas := &mockBerkasService{}
	mockRawatJalan := &mockRawatJalanService{}
	mockTindakan := &mockTindakanService{}

	svc := laboratorium.NewService(mockRepo, mockBerkas, mockRawatJalan, mockTindakan, 48, testURLBerkas, []string{"005"}, []string{"015"}, []string{"010"}, log)
	_, err := svc.GetDetailHasilLab(context.Background(), shared.KategoriLabPK, "2026/04/22/000001", "LAB001", "2026-04-22", "10:00:00")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	var notFound *apperror.NotFoundError
	if !errors.As(err, &notFound) {
		t.Errorf("Expected NotFoundError, got %T: %v", err, err)
	}
}

func TestService_PermintaanLabPK_BusinessScenarios(t *testing.T) {
	log := logger.New()
	mockBerkas := &mockBerkasService{}

	nowStr := time.Now().Format("2006-01-02")
	jamNowStr := time.Now().Format("15:04:05")
	tglKemarinStr := time.Now().Add(-2 * time.Hour).Format("2006-01-02")
	jamKemarinStr := time.Now().Add(-2 * time.Hour).Format("15:04:05")

	t.Run("Sukses Simpan Permintaan Lab PK Ralan dengan Template", func(t *testing.T) {
		mockRepo := &mockRepository{
			getKunjunganForPermintaanPKFn: func(ctx context.Context, noRawat string) (*laboratorium.KunjunganInfoLabPK, error) {
				return &laboratorium.KunjunganInfoLabPK{
					NoRawat:           noRawat,
					NoRkmMedis:        "123456",
					TanggalRegistrasi: tglKemarinStr,
					JamRegistrasi:     jamKemarinStr,
					KodePenjamin:      "UMU",
					StatusBayar:       "Belum Bayar",
					StatusLanjut:      "Ralan",
				}, nil
			},
			cekStatusKamarInapFn: func(ctx context.Context, noRawat string) (bool, bool, error) {
				return false, false, nil
			},
			simpanPermintaanLabPKFn: func(ctx context.Context, noRawat, kodeDokter, status string, req laboratorium.SimpanPermintaanLabPKRequest, kodeTindakanList []string, templateMap map[string][]int) (string, error) {
				return "PK202609030001", nil
			},
			detailPermintaanLabPKFn: func(ctx context.Context, noPermintaan string) (*laboratorium.DetailPermintaanLabPK, error) {
				return &laboratorium.DetailPermintaanLabPK{
					PermintaanLabPK: laboratorium.PermintaanLabPK{
						NoPermintaan:      noPermintaan,
						NoRawat:           "2026/09/03/000001",
						TanggalPermintaan: nowStr,
						JamPermintaan:     jamNowStr,
						KodeDokterPerujuk: "DR01",
						NamaDokterPerujuk: "dr. Sp.PK",
						Status:            "ralan",
						InformasiTambahan: "-",
						DiagnosaKlinis:    "Febris H-3",
					},
					Pemeriksaan: []laboratorium.PemeriksaanLabPKItem{
						{
							KodeTindakan: "TND001",
							NamaTindakan: "Darah Lengkap",
							StatusBayar:  "Belum",
							DetailTemplate: []laboratorium.DetailTemplateLabPKItem{
								{
									IdTemplate:      "101",
									NamaPemeriksaan: "Hemoglobin",
									Satuan:          "g/dL",
									NilaiRujukan:    "12-16",
								},
							},
						},
					},
				}, nil
			},
		}
		mockRawatJalan := &mockRawatJalanService{
			getWaktuRegistrasiFn: func(ctx context.Context, noRawat string) (string, string, bool, error) {
				return tglKemarinStr, jamKemarinStr, true, nil
			},
		}
		mockTindakan := &mockTindakanService{
			cekKeberadaanTindakanLabFn: func(ctx context.Context, kategori shared.KategoriLab, listKodeTindakan []string) (map[string]bool, error) {
				return map[string]bool{"TND001": true}, nil
			},
			cekKeberadaanTemplateLabFn: func(ctx context.Context, listKodeTindakan []string, templateMap map[string][]int) (map[string]map[int]bool, error) {
				return map[string]map[int]bool{"TND001": {101: true}}, nil
			},
		}

		svc := laboratorium.NewService(mockRepo, mockBerkas, mockRawatJalan, mockTindakan, 48, testURLBerkas, nil, nil, nil, log)
		res, err := svc.SimpanPermintaanLabPK(context.Background(), "DR01", shared.StatusLanjutRawatJalan, laboratorium.SimpanPermintaanLabPKRequest{
			NoRawat:           "2026/09/03/000001",
			TanggalPermintaan: nowStr,
			JamPermintaan:     jamNowStr,
			DiagnosaKlinis:    "Febris H-3",
			Pemeriksaan: []laboratorium.ItemPemeriksaanLabPKRequest{
				{
					KodeTindakan: "TND001",
					KodeTemplate: []int{101},
				},
			},
		})

		if err != nil {
			t.Fatalf("Expected success, got error: %v", err)
		}
		if res.NoPermintaan != "PK202609030001" {
			t.Errorf("Expected NoPermintaan PK202609030001, got %s", res.NoPermintaan)
		}
		if len(res.Pemeriksaan) != 1 || len(res.Pemeriksaan[0].DetailTemplate) != 1 {
			t.Errorf("Expected 1 pemeriksaan with 1 template, got %v", res.Pemeriksaan)
		}
	})

	t.Run("Gagal Simpan jika Tindakan Tidak Ditemukan di Master DB", func(t *testing.T) {
		mockRepo := &mockRepository{
			getKunjunganForPermintaanPKFn: func(ctx context.Context, noRawat string) (*laboratorium.KunjunganInfoLabPK, error) {
				return &laboratorium.KunjunganInfoLabPK{
					NoRawat:      noRawat,
					StatusLanjut: "Ralan",
				}, nil
			},
			cekStatusKamarInapFn: func(ctx context.Context, noRawat string) (bool, bool, error) {
				return false, false, nil
			},
		}
		mockRawatJalan := &mockRawatJalanService{
			getWaktuRegistrasiFn: func(ctx context.Context, noRawat string) (string, string, bool, error) {
				return tglKemarinStr, jamKemarinStr, true, nil
			},
		}
		mockTindakan := &mockTindakanService{
			cekKeberadaanTindakanLabFn: func(ctx context.Context, kategori shared.KategoriLab, listKodeTindakan []string) (map[string]bool, error) {
				return map[string]bool{"TND001": false}, nil
			},
		}

		svc := laboratorium.NewService(mockRepo, mockBerkas, mockRawatJalan, mockTindakan, 48, testURLBerkas, nil, nil, nil, log)
		_, err := svc.SimpanPermintaanLabPK(context.Background(), "DR01", shared.StatusLanjutRawatJalan, laboratorium.SimpanPermintaanLabPKRequest{
			NoRawat:           "2026/09/03/000001",
			TanggalPermintaan: nowStr,
			JamPermintaan:     jamNowStr,
			DiagnosaKlinis:    "Febris H-3",
			Pemeriksaan: []laboratorium.ItemPemeriksaanLabPKRequest{
				{KodeTindakan: "TND001"},
			},
		})

		if err == nil {
			t.Fatal("Expected validation error for non-existent tindakan, got nil")
		}
		var valErr apperror.ValidationError
		if !errors.As(err, &valErr) {
			t.Errorf("Expected ValidationError, got %T: %v", err, err)
		}
	})

	t.Run("Gagal Simpan jika Template Tidak Cocok dengan Tindakan", func(t *testing.T) {
		mockRepo := &mockRepository{
			getKunjunganForPermintaanPKFn: func(ctx context.Context, noRawat string) (*laboratorium.KunjunganInfoLabPK, error) {
				return &laboratorium.KunjunganInfoLabPK{
					NoRawat:      noRawat,
					StatusLanjut: "Ralan",
				}, nil
			},
			cekStatusKamarInapFn: func(ctx context.Context, noRawat string) (bool, bool, error) {
				return false, false, nil
			},
		}
		mockRawatJalan := &mockRawatJalanService{
			getWaktuRegistrasiFn: func(ctx context.Context, noRawat string) (string, string, bool, error) {
				return tglKemarinStr, jamKemarinStr, true, nil
			},
		}
		mockTindakan := &mockTindakanService{
			cekKeberadaanTindakanLabFn: func(ctx context.Context, kategori shared.KategoriLab, listKodeTindakan []string) (map[string]bool, error) {
				return map[string]bool{"TND001": true}, nil
			},
			cekKeberadaanTemplateLabFn: func(ctx context.Context, listKodeTindakan []string, templateMap map[string][]int) (map[string]map[int]bool, error) {
				return map[string]map[int]bool{"TND001": {999: false}}, nil
			},
		}

		svc := laboratorium.NewService(mockRepo, mockBerkas, mockRawatJalan, mockTindakan, 48, testURLBerkas, nil, nil, nil, log)
		_, err := svc.SimpanPermintaanLabPK(context.Background(), "DR01", shared.StatusLanjutRawatJalan, laboratorium.SimpanPermintaanLabPKRequest{
			NoRawat:           "2026/09/03/000001",
			TanggalPermintaan: nowStr,
			JamPermintaan:     jamNowStr,
			DiagnosaKlinis:    "Febris H-3",
			Pemeriksaan: []laboratorium.ItemPemeriksaanLabPKRequest{
				{
					KodeTindakan: "TND001",
					KodeTemplate: []int{999},
				},
			},
		})

		if err == nil {
			t.Fatal("Expected validation error for template mismatch, got nil")
		}
		var valErr apperror.ValidationError
		if !errors.As(err, &valErr) {
			t.Errorf("Expected ValidationError, got %T: %v", err, err)
		}
	})

	t.Run("Gagal Simpan jika Pasien BPJS Sudah Bayar", func(t *testing.T) {
		mockRepo := &mockRepository{
			getKunjunganForPermintaanPKFn: func(ctx context.Context, noRawat string) (*laboratorium.KunjunganInfoLabPK, error) {
				return &laboratorium.KunjunganInfoLabPK{
					NoRawat:           noRawat,
					TanggalRegistrasi: nowStr,
					JamRegistrasi:     jamKemarinStr,
					KodePenjamin:      "BPJ",
					StatusBayar:       "Sudah Bayar",
					StatusLanjut:      "Ralan",
				}, nil
			},
		}
		mockRawatJalan := &mockRawatJalanService{}
		mockTindakan := &mockTindakanService{}

		svc := laboratorium.NewService(mockRepo, mockBerkas, mockRawatJalan, mockTindakan, 48, testURLBerkas, nil, nil, nil, log)
		_, err := svc.SimpanPermintaanLabPK(context.Background(), "DR01", shared.StatusLanjutRawatJalan, laboratorium.SimpanPermintaanLabPKRequest{
			NoRawat:           "2026/09/03/000001",
			TanggalPermintaan: nowStr,
			JamPermintaan:     jamNowStr,
			DiagnosaKlinis:    "Febris H-3",
			Pemeriksaan: []laboratorium.ItemPemeriksaanLabPKRequest{
				{KodeTindakan: "TND001"},
			},
		})

		if err == nil {
			t.Fatal("Expected error for paid BPJS patient, got nil")
		}
	})

	t.Run("Gagal Simpan jika Melewati 48 Jam Ralan", func(t *testing.T) {
		tgl3HariLaluStr := time.Now().Add(-72 * time.Hour).Format("2006-01-02")
		mockRepo := &mockRepository{
			getKunjunganForPermintaanPKFn: func(ctx context.Context, noRawat string) (*laboratorium.KunjunganInfoLabPK, error) {
				return &laboratorium.KunjunganInfoLabPK{
					NoRawat:           noRawat,
					TanggalRegistrasi: tgl3HariLaluStr,
					JamRegistrasi:     "08:00:00",
					KodePenjamin:      "UMU",
					StatusBayar:       "Belum Bayar",
					StatusLanjut:      "Ralan",
				}, nil
			},
			cekStatusKamarInapFn: func(ctx context.Context, noRawat string) (bool, bool, error) {
				return false, false, nil
			},
		}
		mockRawatJalan := &mockRawatJalanService{
			getWaktuRegistrasiFn: func(ctx context.Context, noRawat string) (string, string, bool, error) {
				return tgl3HariLaluStr, "08:00:00", true, nil
			},
		}
		mockTindakan := &mockTindakanService{}

		svc := laboratorium.NewService(mockRepo, mockBerkas, mockRawatJalan, mockTindakan, 48, testURLBerkas, nil, nil, nil, log)
		_, err := svc.SimpanPermintaanLabPK(context.Background(), "DR01", shared.StatusLanjutRawatJalan, laboratorium.SimpanPermintaanLabPKRequest{
			NoRawat:           "2026/09/03/000001",
			TanggalPermintaan: nowStr,
			JamPermintaan:     jamNowStr,
			DiagnosaKlinis:    "Febris H-3",
			Pemeriksaan: []laboratorium.ItemPemeriksaanLabPKRequest{
				{KodeTindakan: "TND001"},
			},
		})

		if err == nil {
			t.Fatal("Expected error for 48h expiration, got nil")
		}
	})

	t.Run("Gagal Simpan jika Pasien Ranap Sudah Checkout", func(t *testing.T) {
		mockRepo := &mockRepository{
			getKunjunganForPermintaanPKFn: func(ctx context.Context, noRawat string) (*laboratorium.KunjunganInfoLabPK, error) {
				return &laboratorium.KunjunganInfoLabPK{
					NoRawat:           noRawat,
					TanggalRegistrasi: nowStr,
					JamRegistrasi:     jamKemarinStr,
					KodePenjamin:      "UMU",
					StatusBayar:       "Belum Bayar",
					StatusLanjut:      "Ranap",
				}, nil
			},
			cekStatusKamarInapFn: func(ctx context.Context, noRawat string) (bool, bool, error) {
				return false, true, nil // has record but not active (checked out)
			},
		}
		mockRawatJalan := &mockRawatJalanService{
			getWaktuRegistrasiFn: func(ctx context.Context, noRawat string) (string, string, bool, error) {
				return nowStr, jamKemarinStr, true, nil
			},
		}
		mockTindakan := &mockTindakanService{}

		svc := laboratorium.NewService(mockRepo, mockBerkas, mockRawatJalan, mockTindakan, 48, testURLBerkas, nil, nil, nil, log)
		_, err := svc.SimpanPermintaanLabPK(context.Background(), "DR01", shared.StatusLanjutRawatInap, laboratorium.SimpanPermintaanLabPKRequest{
			NoRawat:           "2026/09/03/000001",
			TanggalPermintaan: nowStr,
			JamPermintaan:     jamNowStr,
			DiagnosaKlinis:    "Febris H-3",
			Pemeriksaan: []laboratorium.ItemPemeriksaanLabPKRequest{
				{KodeTindakan: "TND001"},
			},
		})

		if err == nil {
			t.Fatal("Expected error for checked out inpatient, got nil")
		}
	})

	t.Run("Hapus Permintaan Lab PK Sukses dan Gagal", func(t *testing.T) {
		mockRepo := &mockRepository{
			detailPermintaanLabPKFn: func(ctx context.Context, noPermintaan string) (*laboratorium.DetailPermintaanLabPK, error) {
				if noPermintaan == "PK-SAMPEL" {
					return &laboratorium.DetailPermintaanLabPK{
						PermintaanLabPK: laboratorium.PermintaanLabPK{
							NoPermintaan:      noPermintaan,
							NoRawat:           "2026/09/03/000001",
							Status:            "ralan",
							KodeDokterPerujuk: "DR01",
							TanggalSampel:     "2026-09-03",
						},
					}, nil
				}
				if noPermintaan == "PK-OTHER-DOC" {
					return &laboratorium.DetailPermintaanLabPK{
						PermintaanLabPK: laboratorium.PermintaanLabPK{
							NoPermintaan:      noPermintaan,
							NoRawat:           "2026/09/03/000001",
							Status:            "ralan",
							KodeDokterPerujuk: "DR99",
						},
					}, nil
				}
				return &laboratorium.DetailPermintaanLabPK{
					PermintaanLabPK: laboratorium.PermintaanLabPK{
						NoPermintaan:      noPermintaan,
						NoRawat:           "2026/09/03/000001",
						Status:            "ralan",
						KodeDokterPerujuk: "DR01",
						TanggalSampel:     "0000-00-00",
						TanggalHasil:      "0000-00-00",
					},
				}, nil
			},
			hapusPermintaanLabPKFn: func(ctx context.Context, noPermintaan string) error {
				return nil
			},
		}
		mockRawatJalan := &mockRawatJalanService{}
		mockTindakan := &mockTindakanService{}

		svc := laboratorium.NewService(mockRepo, mockBerkas, mockRawatJalan, mockTindakan, 48, testURLBerkas, nil, nil, nil, log)

		// Sukses
		err := svc.HapusPermintaanLabPK(context.Background(), "2026/09/03/000001", "PK-VALID", shared.StatusLanjutRawatJalan, "DR01")
		if err != nil {
			t.Errorf("Expected success, got %v", err)
		}

		// Gagal dokter lain
		errDoc := svc.HapusPermintaanLabPK(context.Background(), "2026/09/03/000001", "PK-OTHER-DOC", shared.StatusLanjutRawatJalan, "DR01")
		if errDoc == nil {
			t.Error("Expected ForbiddenError for other doctor")
		}

		// Gagal sampel sudah diambil
		errSampel := svc.HapusPermintaanLabPK(context.Background(), "2026/09/03/000001", "PK-SAMPEL", shared.StatusLanjutRawatJalan, "DR01")
		if errSampel == nil {
			t.Error("Expected BusinessError when sample is already taken")
		}
	})
}
