package laboratorium_test

import (
	"context"
	"database/sql"
	"errors"
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

	cekStatusKamarInapFn          func(ctx context.Context, noRawat string) (bool, bool, error)
	simpanPermintaanLabPKFn       func(ctx context.Context, noRawat string, kodeDokter string, status string, req laboratorium.SimpanPermintaanLabPKRequest, kodeTindakanList []string, templateMap map[string][]int) (string, error)
	daftarPermintaanLabPKFn       func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) ([]laboratorium.PermintaanLabPK, error)
	daftarPermintaanLabPKByRMFn   func(ctx context.Context, noRkmMedis string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.PermintaanLabPK, int, error)
	detailPermintaanLabPKFn       func(ctx context.Context, noPermintaan string) (*laboratorium.DetailPermintaanLabPK, error)
	hapusPermintaanLabPKFn        func(ctx context.Context, noPermintaan string) error

	simpanPermintaanLabPAFn     func(ctx context.Context, noRawat string, kodeDokter string, status string, req laboratorium.SimpanPermintaanLabPARequest, kodeTindakanList []string) (string, error)
	daftarPermintaanLabPAFn     func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) ([]laboratorium.PermintaanLabPA, error)
	daftarPermintaanLabPAByRMFn func(ctx context.Context, noRkmMedis string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.PermintaanLabPA, int, error)
	detailPermintaanLabPAFn     func(ctx context.Context, noPermintaan string) (*laboratorium.DetailPermintaanLabPA, error)
	hapusPermintaanLabPAFn      func(ctx context.Context, noPermintaan string) error
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

func (m *mockRepository) SimpanPermintaanLabPA(ctx context.Context, noRawat string, kodeDokter string, status string, req laboratorium.SimpanPermintaanLabPARequest, kodeTindakanList []string) (string, error) {
	if m.simpanPermintaanLabPAFn != nil {
		return m.simpanPermintaanLabPAFn(ctx, noRawat, kodeDokter, status, req, kodeTindakanList)
	}
	return "", nil
}

func (m *mockRepository) DaftarPermintaanLabPA(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) ([]laboratorium.PermintaanLabPA, error) {
	if m.daftarPermintaanLabPAFn != nil {
		return m.daftarPermintaanLabPAFn(ctx, noRawat, statusLanjut)
	}
	return nil, nil
}

func (m *mockRepository) DaftarPermintaanLabPAByRM(ctx context.Context, noRkmMedis string, statusLanjut shared.StatusLanjut, filter laboratorium.FilterRiwayatLab) ([]laboratorium.PermintaanLabPA, int, error) {
	if m.daftarPermintaanLabPAByRMFn != nil {
		return m.daftarPermintaanLabPAByRMFn(ctx, noRkmMedis, statusLanjut, filter)
	}
	return nil, 0, nil
}

func (m *mockRepository) DetailPermintaanLabPA(ctx context.Context, noPermintaan string) (*laboratorium.DetailPermintaanLabPA, error) {
	if m.detailPermintaanLabPAFn != nil {
		return m.detailPermintaanLabPAFn(ctx, noPermintaan)
	}
	return nil, nil
}

func (m *mockRepository) HapusPermintaanLabPA(ctx context.Context, noPermintaan string) error {
	if m.hapusPermintaanLabPAFn != nil {
		return m.hapusPermintaanLabPAFn(ctx, noPermintaan)
	}
	return nil
}

type mockBerkasService struct {
	getBerkasByNoRawatFn func(ctx context.Context, noRawat string, kodeList []string) ([]berkasdigital.BerkasDigitalPerawatan, error)
	buildFullURLFn       func(lokasiFile string) string
}

func (m *mockBerkasService) GetMasterBerkas(ctx context.Context) ([]berkasdigital.MasterBerkasDigital, error) {
	return nil, nil
}

func (m *mockBerkasService) GetBerkasByNoRawat(ctx context.Context, noRawat string, kodeList []string) ([]berkasdigital.BerkasDigitalPerawatan, error) {
	if m.getBerkasByNoRawatFn != nil {
		return m.getBerkasByNoRawatFn(ctx, noRawat, kodeList)
	}
	return nil, nil
}

func (m *mockBerkasService) GetBerkasStream(ctx context.Context, targetURL string) (*berkasdigital.BerkasStream, error) {
	return nil, nil
}

func (m *mockBerkasService) BuildFullURL(lokasiFile string) string {
	if m.buildFullURLFn != nil {
		return m.buildFullURLFn(lokasiFile)
	}
	return testURLBerkas + lokasiFile
}

type mockRawatJalanService struct {
	getWaktuRegistrasiFn func(ctx context.Context, noRawat string) (tanggal string, jam string, exists bool, err error)
	getInfoRegistrasiFn  func(ctx context.Context, noRawat string) (*rawatjalan.InfoRegistrasiPasien, error)
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

func (m *mockRawatJalanService) GetInfoRegistrasi(ctx context.Context, noRawat string) (*rawatjalan.InfoRegistrasiPasien, error) {
	if m.getInfoRegistrasiFn != nil {
		return m.getInfoRegistrasiFn(ctx, noRawat)
	}
	tgl := "2026-09-03"
	jam := "08:00:00"
	if m.getWaktuRegistrasiFn != nil {
		t, j, exists, err := m.getWaktuRegistrasiFn(ctx, noRawat)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, apperror.NewNotFoundError("Data kunjungan pasien tidak ditemukan")
		}
		tgl = t
		jam = j
	}
	return &rawatjalan.InfoRegistrasiPasien{
		TanggalRegistrasi: tgl,
		JamRegistrasi:     jam,
		KodePenjamin:      "UMU",
		StatusBayar:       "Belum Bayar",
	}, nil
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
		getBerkasByNoRawatFn: func(ctx context.Context, noRawat string, kodeList []string) ([]berkasdigital.BerkasDigitalPerawatan, error) {
			return []berkasdigital.BerkasDigitalPerawatan{}, nil
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
		getBerkasByNoRawatFn: func(ctx context.Context, noRawat string, kodeList []string) ([]berkasdigital.BerkasDigitalPerawatan, error) {
			return []berkasdigital.BerkasDigitalPerawatan{
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
			cekStatusKamarInapFn: func(ctx context.Context, noRawat string) (bool, bool, error) {
				return false, false, nil
			},
			simpanPermintaanLabPKFn: func(ctx context.Context, noRawat, kodeDokter, status string, req laboratorium.SimpanPermintaanLabPKRequest, kodeTindakanList []string, templateMap map[string][]int) (string, error) {
				return "PK202609030001", nil
			},
			detailPermintaanLabPKFn: func(ctx context.Context, noPermintaan string) (*laboratorium.DetailPermintaanLabPK, error) {
				return &laboratorium.DetailPermintaanLabPK{
					PermintaanLabPK: laboratorium.PermintaanLabPK{
						PermintaanLabHeader: laboratorium.PermintaanLabHeader{
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
			PermintaanLabHeaderRequest: laboratorium.PermintaanLabHeaderRequest{
				NoRawat:           "2026/09/03/000001",
				TanggalPermintaan: nowStr,
				JamPermintaan:     jamNowStr,
				DiagnosaKlinis:    "Febris H-3",
			},
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
			PermintaanLabHeaderRequest: laboratorium.PermintaanLabHeaderRequest{
				NoRawat:           "2026/09/03/000001",
				TanggalPermintaan: nowStr,
				JamPermintaan:     jamNowStr,
				DiagnosaKlinis:    "Febris H-3",
			},
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
			PermintaanLabHeaderRequest: laboratorium.PermintaanLabHeaderRequest{
				NoRawat:           "2026/09/03/000001",
				TanggalPermintaan: nowStr,
				JamPermintaan:     jamNowStr,
				DiagnosaKlinis:    "Febris H-3",
			},
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
		mockRepo := &mockRepository{}
		mockRawatJalan := &mockRawatJalanService{
			getInfoRegistrasiFn: func(ctx context.Context, noRawat string) (*rawatjalan.InfoRegistrasiPasien, error) {
				return &rawatjalan.InfoRegistrasiPasien{
					TanggalRegistrasi: nowStr,
					JamRegistrasi:     jamKemarinStr,
					KodePenjamin:      "BPJ",
					StatusBayar:       "Sudah Bayar",
				}, nil
			},
		}
		mockTindakan := &mockTindakanService{}

		svc := laboratorium.NewService(mockRepo, mockBerkas, mockRawatJalan, mockTindakan, 48, testURLBerkas, nil, nil, nil, log)
		_, err := svc.SimpanPermintaanLabPK(context.Background(), "DR01", shared.StatusLanjutRawatJalan, laboratorium.SimpanPermintaanLabPKRequest{
			PermintaanLabHeaderRequest: laboratorium.PermintaanLabHeaderRequest{
				NoRawat:           "2026/09/03/000001",
				TanggalPermintaan: nowStr,
				JamPermintaan:     jamNowStr,
				DiagnosaKlinis:    "Febris H-3",
			},
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
			PermintaanLabHeaderRequest: laboratorium.PermintaanLabHeaderRequest{
				NoRawat:           "2026/09/03/000001",
				TanggalPermintaan: nowStr,
				JamPermintaan:     jamNowStr,
				DiagnosaKlinis:    "Febris H-3",
			},
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
			PermintaanLabHeaderRequest: laboratorium.PermintaanLabHeaderRequest{
				NoRawat:           "2026/09/03/000001",
				TanggalPermintaan: nowStr,
				JamPermintaan:     jamNowStr,
				DiagnosaKlinis:    "Febris H-3",
			},
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
							PermintaanLabHeader: laboratorium.PermintaanLabHeader{
								NoPermintaan:      noPermintaan,
								NoRawat:           "2026/09/03/000001",
								Status:            "ralan",
								KodeDokterPerujuk: "DR01",
								TanggalSampel:     "2026-09-03",
							},
						},
					}, nil
				}
				if noPermintaan == "PK-OTHER-DOC" {
					return &laboratorium.DetailPermintaanLabPK{
						PermintaanLabPK: laboratorium.PermintaanLabPK{
							PermintaanLabHeader: laboratorium.PermintaanLabHeader{
								NoPermintaan:      noPermintaan,
								NoRawat:           "2026/09/03/000001",
								Status:            "ralan",
								KodeDokterPerujuk: "DR99",
							},
						},
					}, nil
				}
				return &laboratorium.DetailPermintaanLabPK{
					PermintaanLabPK: laboratorium.PermintaanLabPK{
						PermintaanLabHeader: laboratorium.PermintaanLabHeader{
							NoPermintaan:      noPermintaan,
							NoRawat:           "2026/09/03/000001",
							Status:            "ralan",
							KodeDokterPerujuk: "DR01",
							TanggalSampel:     "0000-00-00",
							TanggalHasil:      "0000-00-00",
						},
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

func TestService_PermintaanLabPA(t *testing.T) {
	log := logger.New()
	now := time.Now()
	nowDate := now.Format("2006-01-02")
	nowTime := now.Format("15:04:05")

	mockBerkas := &mockBerkasService{}
	mockRawatJalan := &mockRawatJalanService{
		getWaktuRegistrasiFn: func(ctx context.Context, noRawat string) (string, string, bool, error) {
			if noRawat == "2026/09/04/000001" {
				return nowDate, "07:00:00", true, nil
			}
			return "", "", false, nil
		},
	}
	mockTindakan := &mockTindakanService{
		cekKeberadaanTindakanLabFn: func(ctx context.Context, kategori shared.KategoriLab, listKodeTindakan []string) (map[string]bool, error) {
			res := make(map[string]bool)
			for _, kd := range listKodeTindakan {
				if kd == "PA00001" {
					res[kd] = true
				}
			}
			return res, nil
		},
	}

	t.Run("SimpanPermintaanLabPA_Success", func(t *testing.T) {
		mockRepo := &mockRepository{
			cekStatusKamarInapFn: func(ctx context.Context, noRawat string) (bool, bool, error) {
				return false, false, nil
			},
			simpanPermintaanLabPAFn: func(ctx context.Context, noRawat, kodeDokter, status string, req laboratorium.SimpanPermintaanLabPARequest, kodeTindakanList []string) (string, error) {
				return "PA202609040001", nil
			},
			detailPermintaanLabPAFn: func(ctx context.Context, noPermintaan string) (*laboratorium.DetailPermintaanLabPA, error) {
				return &laboratorium.DetailPermintaanLabPA{
					PermintaanLabPA: laboratorium.PermintaanLabPA{
						PermintaanLabHeader: laboratorium.PermintaanLabHeader{
							NoPermintaan:      noPermintaan,
							NoRawat:           "2026/09/04/000001",
							TanggalPermintaan: nowDate,
							JamPermintaan:     nowTime,
							Status:            "ralan",
						},
					},
					Pemeriksaan: []laboratorium.PemeriksaanLabPAItem{
						{KodeTindakan: "PA00001", NamaTindakan: "Pemeriksaan PA Sediaan Kecil"},
					},
				}, nil
			},
		}

		svc := laboratorium.NewService(mockRepo, mockBerkas, mockRawatJalan, mockTindakan, 48, testURLBerkas, nil, nil, nil, log)
		req := laboratorium.SimpanPermintaanLabPARequest{
			PermintaanLabHeaderRequest: laboratorium.PermintaanLabHeaderRequest{
				NoRawat:           "2026/09/04/000001",
				TanggalPermintaan: nowDate,
				JamPermintaan:     nowTime,
				DiagnosaKlinis:    "Tumor Mammae",
			},
			Pemeriksaan: []laboratorium.ItemPemeriksaanLabPARequest{
				{KodeTindakan: "PA00001"},
			},
		}

		detail, err := svc.SimpanPermintaanLabPA(context.Background(), "DR01", shared.StatusLanjutRawatJalan, req)
		if err != nil {
			t.Fatalf("Expected success, got %v", err)
		}
		if detail.NoPermintaan != "PA202609040001" {
			t.Errorf("Expected NoPermintaan PA202609040001, got %s", detail.NoPermintaan)
		}
	})

	t.Run("SimpanPermintaanLabPA_BPJSSudahBayar", func(t *testing.T) {
		mockRepo := &mockRepository{}
		mockRawatJalanBPJS := &mockRawatJalanService{
			getInfoRegistrasiFn: func(ctx context.Context, noRawat string) (*rawatjalan.InfoRegistrasiPasien, error) {
				return &rawatjalan.InfoRegistrasiPasien{
					TanggalRegistrasi: nowDate,
					JamRegistrasi:     nowTime,
					KodePenjamin:      "BPJ",
					StatusBayar:       "Sudah Bayar",
				}, nil
			},
		}

		svc := laboratorium.NewService(mockRepo, mockBerkas, mockRawatJalanBPJS, mockTindakan, 48, testURLBerkas, nil, nil, nil, log)
		req := laboratorium.SimpanPermintaanLabPARequest{
			PermintaanLabHeaderRequest: laboratorium.PermintaanLabHeaderRequest{
				NoRawat:           "2026/09/04/000001",
				TanggalPermintaan: nowDate,
				JamPermintaan:     nowTime,
				DiagnosaKlinis:    "Tumor Mammae",
			},
			Pemeriksaan: []laboratorium.ItemPemeriksaanLabPARequest{
				{KodeTindakan: "PA00001"},
			},
		}

		_, err := svc.SimpanPermintaanLabPA(context.Background(), "DR01", shared.StatusLanjutRawatJalan, req)
		if err == nil {
			t.Fatal("Expected error for BPJS already paid, got nil")
		}
	})

	t.Run("SimpanPermintaanLabPA_TindakanNotFound", func(t *testing.T) {
		mockRepo := &mockRepository{
			cekStatusKamarInapFn: func(ctx context.Context, noRawat string) (bool, bool, error) {
				return false, false, nil
			},
		}

		svc := laboratorium.NewService(mockRepo, mockBerkas, mockRawatJalan, mockTindakan, 48, testURLBerkas, nil, nil, nil, log)
		req := laboratorium.SimpanPermintaanLabPARequest{
			PermintaanLabHeaderRequest: laboratorium.PermintaanLabHeaderRequest{
				NoRawat:           "2026/09/04/000001",
				TanggalPermintaan: nowDate,
				JamPermintaan:     nowTime,
				DiagnosaKlinis:    "Tumor Mammae",
			},
			Pemeriksaan: []laboratorium.ItemPemeriksaanLabPARequest{
				{KodeTindakan: "PA-TIDAK-ADA"},
			},
		}

		_, err := svc.SimpanPermintaanLabPA(context.Background(), "DR01", shared.StatusLanjutRawatJalan, req)
		if err == nil {
			t.Fatal("Expected error for invalid tindakan, got nil")
		}
	})

	t.Run("GetDaftarPermintaanLabPA_Success", func(t *testing.T) {
		mockRepo := &mockRepository{
			daftarPermintaanLabPAFn: func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) ([]laboratorium.PermintaanLabPA, error) {
				return []laboratorium.PermintaanLabPA{
					{
						PermintaanLabHeader: laboratorium.PermintaanLabHeader{
							NoPermintaan: "PA001",
							NoRawat:      noRawat,
							TanggalHasil: "2026-09-04",
							JamHasil:     "10:00:00",
						},
					},
				}, nil
			},
		}
		svc := laboratorium.NewService(mockRepo, mockBerkas, mockRawatJalan, mockTindakan, 48, testURLBerkas, nil, nil, nil, log)
		list, err := svc.GetDaftarPermintaanLabPA(context.Background(), "2026/09/04/000001", shared.StatusLanjutRawatJalan)
		if err != nil {
			t.Fatalf("Expected success, got %v", err)
		}
		if len(list) != 1 || list[0].StatusProses != "Selesai" {
			t.Errorf("Expected 1 item with StatusProses Selesai, got %+v", list)
		}
	})

	t.Run("HapusPermintaanLabPA_SuccessAndGuards", func(t *testing.T) {
		mockRepo := &mockRepository{
			detailPermintaanLabPAFn: func(ctx context.Context, noPermintaan string) (*laboratorium.DetailPermintaanLabPA, error) {
				if noPermintaan == "PA-VALID" {
					return &laboratorium.DetailPermintaanLabPA{
						PermintaanLabPA: laboratorium.PermintaanLabPA{
							PermintaanLabHeader: laboratorium.PermintaanLabHeader{
								NoPermintaan:      noPermintaan,
								NoRawat:           "2026/09/04/000001",
								KodeDokterPerujuk: "DR01",
								TanggalSampel:     "0000-00-00",
								TanggalHasil:      "0000-00-00",
								Status:            "ralan",
							},
						},
					}, nil
				}
				if noPermintaan == "PA-OTHER-DOC" {
					return &laboratorium.DetailPermintaanLabPA{
						PermintaanLabPA: laboratorium.PermintaanLabPA{
							PermintaanLabHeader: laboratorium.PermintaanLabHeader{
								NoPermintaan:      noPermintaan,
								NoRawat:           "2026/09/04/000001",
								KodeDokterPerujuk: "DR02",
								Status:            "ralan",
							},
						},
					}, nil
				}
				if noPermintaan == "PA-SAMPEL" {
					return &laboratorium.DetailPermintaanLabPA{
						PermintaanLabPA: laboratorium.PermintaanLabPA{
							PermintaanLabHeader: laboratorium.PermintaanLabHeader{
								NoPermintaan:      noPermintaan,
								NoRawat:           "2026/09/04/000001",
								KodeDokterPerujuk: "DR01",
								TanggalSampel:     "2026-09-04",
								Status:            "ralan",
							},
						},
					}, nil
				}
				return nil, sql.ErrNoRows
			},
			hapusPermintaanLabPAFn: func(ctx context.Context, noPermintaan string) error {
				return nil
			},
		}

		svc := laboratorium.NewService(mockRepo, mockBerkas, mockRawatJalan, mockTindakan, 48, testURLBerkas, nil, nil, nil, log)

		// Sukses
		err := svc.HapusPermintaanLabPA(context.Background(), "2026/09/04/000001", "PA-VALID", shared.StatusLanjutRawatJalan, "DR01")
		if err != nil {
			t.Errorf("Expected success, got %v", err)
		}

		// Gagal dokter lain
		errDoc := svc.HapusPermintaanLabPA(context.Background(), "2026/09/04/000001", "PA-OTHER-DOC", shared.StatusLanjutRawatJalan, "DR01")
		if errDoc == nil {
			t.Error("Expected ForbiddenError for other doctor")
		}

		// Gagal sampel sudah diambil
		errSampel := svc.HapusPermintaanLabPA(context.Background(), "2026/09/04/000001", "PA-SAMPEL", shared.StatusLanjutRawatJalan, "DR01")
		if errSampel == nil {
			t.Error("Expected BusinessError when sample is already taken")
		}
	})
}
