package laboratorium

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"erm-dokter/internal/berkasdigital"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatinap"
	"erm-dokter/internal/rawatjalan"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
	"erm-dokter/internal/tindakan"
)

type Service interface {
	GetRiwayatLabKunjungan(ctx context.Context, kategori shared.KategoriLab, noRawat string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) (*HasilLaboratoriumKunjungan, shared.PaginationMeta, error)
	GetRiwayatLabPasien(ctx context.Context, kategori shared.KategoriLab, noRM string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]HasilLaboratorium, shared.PaginationMeta, error)
	GetDetailHasilLab(ctx context.Context, kategori shared.KategoriLab, noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa string) (*HasilLaboratorium, error)

	SimpanPermintaanLabPK(ctx context.Context, kodeDokterLogin string, statusLanjut shared.StatusLanjut, req SimpanPermintaanLabPKRequest) (*DetailPermintaanLabPK, error)
	GetDaftarPermintaanLabPK(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) ([]PermintaanLabPK, error)
	GetRiwayatPermintaanLabPKByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]PermintaanLabPK, shared.PaginationMeta, error)
	GetDetailPermintaanLabPK(ctx context.Context, noRawat string, noPermintaan string, statusLanjut shared.StatusLanjut) (*DetailPermintaanLabPK, error)
	UpdatePermintaanLabPK(ctx context.Context, kodeDokterLogin, noRawat, noPermintaan string, statusLanjut shared.StatusLanjut, req SimpanPermintaanLabPKRequest) (*DetailPermintaanLabPK, error)
	HapusPermintaanLabPK(ctx context.Context, noRawat string, noPermintaan string, statusLanjut shared.StatusLanjut, kodeDokterLogin string) error

	SimpanPermintaanLabPA(ctx context.Context, kodeDokterLogin string, statusLanjut shared.StatusLanjut, req SimpanPermintaanLabPARequest) (*DetailPermintaanLabPA, error)
	GetDaftarPermintaanLabPA(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) ([]PermintaanLabPA, error)
	GetRiwayatPermintaanLabPAByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]PermintaanLabPA, shared.PaginationMeta, error)
	GetDetailPermintaanLabPA(ctx context.Context, noRawat string, noPermintaan string, statusLanjut shared.StatusLanjut) (*DetailPermintaanLabPA, error)
	UpdatePermintaanLabPA(ctx context.Context, kodeDokterLogin, noRawat, noPermintaan string, statusLanjut shared.StatusLanjut, req SimpanPermintaanLabPARequest) (*DetailPermintaanLabPA, error)
	HapusPermintaanLabPA(ctx context.Context, noRawat string, noPermintaan string, statusLanjut shared.StatusLanjut, kodeDokterLogin string) error

	SimpanPermintaanLabMB(ctx context.Context, kodeDokterLogin string, statusLanjut shared.StatusLanjut, req SimpanPermintaanLabMBRequest) (*DetailPermintaanLabMB, error)
	GetDaftarPermintaanLabMB(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) ([]PermintaanLabMB, error)
	GetRiwayatPermintaanLabMBByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]PermintaanLabMB, shared.PaginationMeta, error)
	GetDetailPermintaanLabMB(ctx context.Context, noRawat string, noPermintaan string, statusLanjut shared.StatusLanjut) (*DetailPermintaanLabMB, error)
	UpdatePermintaanLabMB(ctx context.Context, kodeDokterLogin, noRawat, noPermintaan string, statusLanjut shared.StatusLanjut, req SimpanPermintaanLabMBRequest) (*DetailPermintaanLabMB, error)
	HapusPermintaanLabMB(ctx context.Context, noRawat string, noPermintaan string, statusLanjut shared.StatusLanjut, kodeDokterLogin string) error
}

type service struct {
	repo              Repository
	berkasSvc         berkasdigital.Service
	rawatJalanService rawatjalan.Service
	rawatInapService  rawatinap.Service
	tindakanService   tindakan.Service
	maxEditJam        int
	urlBerkasDigital  string
	kodeBerkasPK      []string
	kodeBerkasPA      []string
	kodeBerkasMB      []string
	log               *logger.Logger
}

func NewService(
	repo Repository,
	berkasSvc berkasdigital.Service,
	rawatJalanService rawatjalan.Service,
	rawatInapService rawatinap.Service,
	tindakanService tindakan.Service,
	maxEditJam int,
	urlBerkasDigital string,
	kodeBerkasPK, kodeBerkasPA, kodeBerkasMB []string,
	log *logger.Logger,
) Service {
	baseURL := ""
	if urlBerkasDigital != "" {
		baseURL = strings.TrimRight(urlBerkasDigital, "/") + "/"
	}
	return &service{
		repo:              repo,
		berkasSvc:         berkasSvc,
		rawatJalanService: rawatJalanService,
		rawatInapService:  rawatInapService,
		tindakanService:   tindakanService,
		maxEditJam:        maxEditJam,
		urlBerkasDigital:  baseURL,
		kodeBerkasPK:      kodeBerkasPK,
		kodeBerkasPA:      kodeBerkasPA,
		kodeBerkasMB:      kodeBerkasMB,
		log:               log,
	}
}

func (s *service) GetRiwayatLabKunjungan(ctx context.Context, kategori shared.KategoriLab, noRawat string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) (*HasilLaboratoriumKunjungan, shared.PaginationMeta, error) {
	var items []HasilLaboratorium
	var total int
	var kodeBerkas []string
	var err error

	switch kategori {
	case shared.KategoriLabPK:
		items, total, err = s.repo.DaftarHasilLabPK(ctx, noRawat, statusLanjut, filter)
		kodeBerkas = s.kodeBerkasPK
	case shared.KategoriLabPA:
		items, total, err = s.repo.DaftarHasilLabPA(ctx, noRawat, statusLanjut, filter)
		kodeBerkas = s.kodeBerkasPA
	case shared.KategoriLabMB:
		items, total, err = s.repo.DaftarHasilLabMB(ctx, noRawat, statusLanjut, filter)
		kodeBerkas = s.kodeBerkasMB
	default:
		return nil, shared.PaginationMeta{}, apperror.NewBusinessError("Kategori laboratorium tidak valid")
	}

	if err != nil {
		s.log.Error("gagal mengambil riwayat lab kunjungan %s kategori %s: %v", noRawat, kategori, err)
		return nil, shared.PaginationMeta{}, err
	}

	berkasList := s.fetchBerkasDigitalKunjungan(ctx, noRawat, kodeBerkas)

	kunjunganData := &HasilLaboratoriumKunjungan{
		HasilPemeriksaan: items,
		BerkasDigital:    berkasList,
	}

	meta := shared.NewPaginationMeta(total, filter.Page, filter.Limit)
	return kunjunganData, meta, nil
}

func (s *service) GetRiwayatLabPasien(ctx context.Context, kategori shared.KategoriLab, noRM string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]HasilLaboratorium, shared.PaginationMeta, error) {
	var items []HasilLaboratorium
	var total int
	var err error

	switch kategori {
	case shared.KategoriLabPK:
		items, total, err = s.repo.DaftarHasilLabPKByRM(ctx, noRM, statusLanjut, filter)
	case shared.KategoriLabPA:
		items, total, err = s.repo.DaftarHasilLabPAByRM(ctx, noRM, statusLanjut, filter)
	case shared.KategoriLabMB:
		items, total, err = s.repo.DaftarHasilLabMBByRM(ctx, noRM, statusLanjut, filter)
	default:
		return nil, shared.PaginationMeta{}, apperror.NewBusinessError("Kategori laboratorium tidak valid")
	}

	if err != nil {
		s.log.Error("gagal mengambil riwayat lab pasien %s kategori %s: %v", noRM, kategori, err)
		return nil, shared.PaginationMeta{}, err
	}

	meta := shared.NewPaginationMeta(total, filter.Page, filter.Limit)
	return items, meta, nil
}

func (s *service) GetDetailHasilLab(ctx context.Context, kategori shared.KategoriLab, noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa string) (*HasilLaboratorium, error) {
	var item *HasilLaboratorium
	var err error

	switch kategori {
	case shared.KategoriLabPK:
		item, err = s.repo.DetailHasilLabPK(ctx, noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa)
	case shared.KategoriLabPA:
		item, err = s.repo.DetailHasilLabPA(ctx, noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa)
	case shared.KategoriLabMB:
		item, err = s.repo.DetailHasilLabMB(ctx, noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa)
	default:
		return nil, apperror.NewBusinessError("Kategori laboratorium tidak valid")
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFoundError("Data hasil laboratorium tidak ditemukan")
		}
		s.log.Error("gagal mengambil detail hasil lab %s %s: %v", noRawat, kodeTindakan, err)
		return nil, err
	}

	return item, nil
}

func (s *service) fetchBerkasDigitalKunjungan(ctx context.Context, noRawat string, kodeBerkas []string) []berkasdigital.BerkasDigital {
	if len(kodeBerkas) == 0 {
		return make([]berkasdigital.BerkasDigital, 0)
	}

	berkasList, err := s.berkasSvc.GetBerkasByNoRawat(ctx, noRawat, kodeBerkas)
	if err != nil || len(berkasList) == 0 {
		return make([]berkasdigital.BerkasDigital, 0)
	}

	berkasDigital := make([]berkasdigital.BerkasDigital, 0, len(berkasList))
	for _, b := range berkasList {
		if b.LokasiFile == "" {
			continue
		}
		fullURL := s.berkasSvc.BuildFullURL(b.LokasiFile)
		berkasDigital = append(berkasDigital, berkasdigital.BerkasDigital{
			Kode:       b.Kode,
			NamaBerkas: b.NamaBerkas,
			IdBerkas:   fullURL,
		})
	}
	return berkasDigital
}

func (s *service) validasiRegistrasiDanStatus(ctx context.Context, noRawat, tglPermintaan, jamPermintaan string, statusLanjut shared.StatusLanjut, action string) error {
	infoReg, err := s.rawatJalanService.GetInfoRegistrasi(ctx, noRawat)
	if err != nil {
		return err
	}

	if infoReg.StatusBayar == "Sudah Bayar" && infoReg.KodePenjamin == "BPJ" {
		passive := "membuat atau mengubah"
		if action == "menghapus" || action == "membatalkan" {
			passive = "membatalkan"
		}
		return apperror.NewBusinessError(fmt.Sprintf("Pasien BPJS yang sudah menyelesaikan pembayaran / administrasi tidak dapat %s permintaan laboratorium", passive))
	}

	tglRegStr := infoReg.TanggalRegistrasi
	jamRegStr := infoReg.JamRegistrasi

	waktuRegistrasi, err := shared.ParseWaktu(tglRegStr, jamRegStr)
	if err != nil {
		s.log.Error("Gagal parse waktu registrasi no_rawat %s (%s %s): %v", noRawat, tglRegStr, jamRegStr, err)
		return err
	}

	if tglPermintaan != "" && jamPermintaan != "" {
		waktuPermintaan, err := shared.ParseWaktu(tglPermintaan, jamPermintaan)
		if err != nil {
			return apperror.NewBusinessError(err.Error())
		}

		if waktuPermintaan.Before(waktuRegistrasi) {
			errs := apperror.ValidationError{
				"tanggal_permintaan": fmt.Sprintf("Waktu permintaan laboratorium (%s %s) tidak boleh lebih awal dari waktu registrasi pasien (%s %s)", tglPermintaan, jamPermintaan, tglRegStr, jamRegStr),
			}
			s.log.Warn("Validasi waktu permintaan lab gagal untuk no_rawat %s: %+v", noRawat, errs)
			return errs
		}
	}

	return s.validasiStatusKamarDanBatasWaktu(ctx, noRawat, statusLanjut, waktuRegistrasi, tglRegStr, jamRegStr)
}

func (s *service) validasiStatusKamarDanBatasWaktu(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, waktuRegistrasi time.Time, tglRegStr, jamRegStr string) error {
	isKamarAktif, hasRecordKamar, err := s.rawatInapService.CekStatusKamarInap(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal memeriksa status kamar inap pasien %s: %v", noRawat, err)
		return err
	}

	// 1. Jika pasien pernah / sedang masuk rawat inap:
	if hasRecordKamar {
		if !isKamarAktif {
			return apperror.NewBusinessError("Pasien rawat inap sudah keluar / checkout dari kamar inap")
		}
		// Selama pasien masih dirawat di kamar inap (belum checkout), transaksi Ranap maupun Ralan tetap diizinkan
		return nil
	}

	// 2. Jika pasien murni rawat jalan (tidak pernah masuk rawat inap):
	if strings.EqualFold(string(statusLanjut), string(shared.StatusLanjutRawatInap)) {
		return apperror.NewBusinessError("Pasien belum/tidak terdaftar di kamar inap. Permintaan laboratorium harus menggunakan status 'Ralan'.")
	}

	batasWaktu := waktuRegistrasi.Add(time.Duration(s.maxEditJam) * time.Hour)
	if time.Now().After(batasWaktu) {
		errMsg := fmt.Sprintf("Batas waktu permintaan laboratorium untuk kunjungan rawat jalan ini telah berakhir (maksimal %d jam dari waktu registrasi: %s %s)", s.maxEditJam, tglRegStr, jamRegStr)
		s.log.Warn("Permintaan laboratorium ditolak karena lewat batas %d jam untuk no_rawat %s: %s", s.maxEditJam, noRawat, errMsg)
		return apperror.NewForbiddenError(errMsg)
	}

	return nil
}
