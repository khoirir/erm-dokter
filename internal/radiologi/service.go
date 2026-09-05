package radiologi

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatinap"
	"erm-dokter/internal/rawatjalan"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
	"erm-dokter/internal/tindakan"
)

type Service interface {
	GetRiwayatRadiologiKunjungan(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter FilterRiwayatRadiologi) ([]HasilRadiologi, shared.PaginationMeta, error)
	GetRiwayatRadiologiPasien(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter FilterRiwayatRadiologi) ([]HasilRadiologi, shared.PaginationMeta, error)
	GetDetailHasilRadiologi(ctx context.Context, idHasil IdHasilRadiologi, statusLanjut shared.StatusLanjut) (*HasilRadiologi, error)

	SimpanPermintaanRadiologi(ctx context.Context, kodeDokterLogin string, statusLanjut shared.StatusLanjut, req SimpanPermintaanRadiologiRequest) (*DetailPermintaanRadiologi, error)
	GetDaftarPermintaanRadiologi(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) ([]DetailPermintaanRadiologi, error)
	GetRiwayatPermintaanRadiologiByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter FilterRiwayatPermintaanRadiologi) ([]DetailPermintaanRadiologi, shared.PaginationMeta, error)
	GetDetailPermintaanRadiologi(ctx context.Context, noRawat string, noPermintaan string, statusLanjut shared.StatusLanjut) (*DetailPermintaanRadiologi, error)
	UpdatePermintaanRadiologi(ctx context.Context, kodeDokterLogin, noRawat, noPermintaan string, statusLanjut shared.StatusLanjut, req SimpanPermintaanRadiologiRequest) (*DetailPermintaanRadiologi, error)
	HapusPermintaanRadiologi(ctx context.Context, noRawat string, noPermintaan string, statusLanjut shared.StatusLanjut, kodeDokterLogin string) error
}

type service struct {
	repo              Repository
	rawatJalanService rawatjalan.Service
	rawatInapService  rawatinap.Service
	tindakanService   tindakan.Service
	maxEditJam        int
	log               *logger.Logger
}

func NewService(
	repo Repository,
	rawatJalanService rawatjalan.Service,
	rawatInapService rawatinap.Service,
	tindakanService tindakan.Service,
	maxEditJam int,
	log *logger.Logger,
) Service {
	return &service{
		repo:              repo,
		rawatJalanService: rawatJalanService,
		rawatInapService:  rawatInapService,
		tindakanService:   tindakanService,
		maxEditJam:        maxEditJam,
		log:               log,
	}
}

func (s *service) GetRiwayatRadiologiKunjungan(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter FilterRiwayatRadiologi) ([]HasilRadiologi, shared.PaginationMeta, error) {
	list, total, err := s.repo.DaftarHasilRadiologiKunjungan(ctx, noRawat, statusLanjut, filter)
	if err != nil {
		s.log.Error("Gagal mengambil riwayat radiologi kunjungan %s: %v", noRawat, err)
		return nil, shared.PaginationMeta{}, err
	}

	meta := shared.NewPaginationMeta(total, filter.Page, filter.Limit)
	return list, meta, nil
}

func (s *service) GetRiwayatRadiologiPasien(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter FilterRiwayatRadiologi) ([]HasilRadiologi, shared.PaginationMeta, error) {
	list, total, err := s.repo.DaftarHasilRadiologiPasien(ctx, noRM, statusLanjut, filter)
	if err != nil {
		s.log.Error("Gagal mengambil riwayat radiologi pasien %s: %v", noRM, err)
		return nil, shared.PaginationMeta{}, err
	}

	meta := shared.NewPaginationMeta(total, filter.Page, filter.Limit)
	return list, meta, nil
}

func (s *service) GetDetailHasilRadiologi(ctx context.Context, idHasil IdHasilRadiologi, statusLanjut shared.StatusLanjut) (*HasilRadiologi, error) {
	item, err := s.repo.DetailHasilRadiologi(ctx, idHasil, statusLanjut)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFoundError("Data hasil pemeriksaan radiologi tidak ditemukan")
		}
		s.log.Error("Gagal mengambil detail hasil radiologi no_rawat %s kode_tindakan %s: %v", idHasil.NoRawat, idHasil.KodeTindakan, err)
		return nil, err
	}

	return item, nil
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
		return apperror.NewBusinessError(fmt.Sprintf("Pasien BPJS yang sudah menyelesaikan pembayaran / administrasi tidak dapat %s permintaan radiologi", passive))
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
				"tanggal_permintaan": fmt.Sprintf("Waktu permintaan radiologi (%s %s) tidak boleh lebih awal dari waktu registrasi pasien (%s %s)", tglPermintaan, jamPermintaan, tglRegStr, jamRegStr),
			}
			s.log.Warn("Validasi waktu permintaan radiologi gagal untuk no_rawat %s: %+v", noRawat, errs)
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

	if hasRecordKamar {
		if !isKamarAktif {
			return apperror.NewBusinessError("Pasien rawat inap sudah keluar / checkout dari kamar inap")
		}
		return nil
	}

	if strings.EqualFold(string(statusLanjut), string(shared.StatusLanjutRawatInap)) {
		return apperror.NewBusinessError("Pasien belum/tidak terdaftar di kamar inap. Permintaan radiologi harus menggunakan status 'Ralan'.")
	}

	batasWaktu := waktuRegistrasi.Add(time.Duration(s.maxEditJam) * time.Hour)
	if time.Now().After(batasWaktu) {
		errMsg := fmt.Sprintf("Batas waktu permintaan radiologi untuk kunjungan rawat jalan ini telah berakhir (maksimal %d jam dari waktu registrasi: %s %s)", s.maxEditJam, tglRegStr, jamRegStr)
		s.log.Warn("Permintaan radiologi ditolak karena lewat batas %d jam untuk no_rawat %s: %s", s.maxEditJam, noRawat, errMsg)
		return apperror.NewForbiddenError(errMsg)
	}

	return nil
}

func (s *service) SimpanPermintaanRadiologi(ctx context.Context, kodeDokterLogin string, statusLanjut shared.StatusLanjut, req SimpanPermintaanRadiologiRequest) (*DetailPermintaanRadiologi, error) {
	if err := s.validasiRegistrasiDanStatus(ctx, req.NoRawat, req.TanggalPermintaan, req.JamPermintaan, statusLanjut, "membuat"); err != nil {
		return nil, err
	}

	var kodeTindakanList []string
	for _, item := range req.Pemeriksaan {
		kodeTindakanList = append(kodeTindakanList, item.KodeTindakan)
	}

	adaMap, err := s.tindakanService.CekKeberadaanTindakanRadiologi(ctx, kodeTindakanList)
	if err != nil {
		s.log.Error("Gagal memeriksa keberadaan tindakan radiologi: %v", err)
		return nil, err
	}

	for _, kode := range kodeTindakanList {
		if !adaMap[kode] {
			return nil, apperror.NewBusinessError(fmt.Sprintf("Pemeriksaan radiologi '%s' tidak ditemukan atau tidak aktif", kode))
		}
	}

	noPermintaan, err := s.repo.SimpanPermintaanRadiologi(ctx, req.NoRawat, kodeDokterLogin, statusLanjut, req, kodeTindakanList)
	if err != nil {
		s.log.Error("Gagal menyimpan permintaan radiologi untuk no_rawat %s: %v", req.NoRawat, err)
		return nil, err
	}

	detail, err := s.repo.DetailPermintaanRadiologi(ctx, noPermintaan)
	if err != nil {
		s.log.Error("Gagal mengambil detail setelah simpan permintaan radiologi no_permintaan %s: %v", noPermintaan, err)
		return nil, err
	}

	return detail, nil
}

func (s *service) GetDaftarPermintaanRadiologi(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) ([]DetailPermintaanRadiologi, error) {
	list, err := s.repo.DaftarPermintaanRadiologiKunjungan(ctx, noRawat, statusLanjut)
	if err != nil {
		s.log.Error("Gagal mengambil daftar permintaan radiologi no_rawat %s: %v", noRawat, err)
		return nil, err
	}
	return list, nil
}

func (s *service) GetRiwayatPermintaanRadiologiByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter FilterRiwayatPermintaanRadiologi) ([]DetailPermintaanRadiologi, shared.PaginationMeta, error) {
	list, total, err := s.repo.DaftarPermintaanRadiologiPasien(ctx, noRM, statusLanjut, filter)
	if err != nil {
		s.log.Error("Gagal mengambil riwayat permintaan radiologi pasien no_rm %s: %v", noRM, err)
		return nil, shared.PaginationMeta{}, err
	}

	meta := shared.NewPaginationMeta(total, filter.Page, filter.Limit)
	return list, meta, nil
}

func (s *service) GetDetailPermintaanRadiologi(ctx context.Context, noRawat string, noPermintaan string, statusLanjut shared.StatusLanjut) (*DetailPermintaanRadiologi, error) {
	detail, err := s.repo.DetailPermintaanRadiologi(ctx, noPermintaan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFoundError("Data permintaan radiologi tidak ditemukan")
		}
		s.log.Error("Gagal mengambil detail permintaan radiologi no_permintaan %s: %v", noPermintaan, err)
		return nil, err
	}

	if detail.NoRawat != noRawat {
		return nil, apperror.NewBusinessError("Permintaan radiologi tidak sesuai dengan kunjungan pasien")
	}

	if statusLanjut != "" && !strings.EqualFold(detail.Status, string(statusLanjut)) {
		return nil, apperror.NewNotFoundError("Data permintaan radiologi tidak ditemukan")
	}

	return detail, nil
}

func (s *service) UpdatePermintaanRadiologi(ctx context.Context, kodeDokterLogin, noRawat, noPermintaan string, statusLanjut shared.StatusLanjut, req SimpanPermintaanRadiologiRequest) (*DetailPermintaanRadiologi, error) {
	detail, err := s.repo.DetailPermintaanRadiologi(ctx, noPermintaan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFoundError("Data permintaan radiologi tidak ditemukan")
		}
		s.log.Error("Gagal memeriksa detail permintaan radiologi sebelum update no_permintaan %s: %v", noPermintaan, err)
		return nil, err
	}

	if detail.NoRawat != noRawat {
		return nil, apperror.NewBusinessError("Permintaan radiologi tidak sesuai dengan kunjungan pasien")
	}

	if statusLanjut != "" && !strings.EqualFold(detail.Status, string(statusLanjut)) {
		return nil, apperror.NewNotFoundError("Data permintaan radiologi tidak ditemukan")
	}

	if detail.DokterPerujuk.KodeDokter != kodeDokterLogin {
		return nil, apperror.NewForbiddenError("Hanya dokter pembuat order yang dapat memperbarui permintaan radiologi ini")
	}

	if (detail.TanggalSampel != "0000-00-00" && detail.TanggalSampel != "") || (detail.TanggalHasil != "0000-00-00" && detail.TanggalHasil != "") {
		return nil, apperror.NewBusinessError("Permintaan radiologi sudah diproses oleh petugas dan tidak dapat diubah")
	}

	for _, p := range detail.Pemeriksaan {
		if strings.EqualFold(p.StatusBayar, "Sudah") {
			return nil, apperror.NewBusinessError("Permintaan radiologi yang pemeriksaannya telah dibayar tidak dapat diubah")
		}
	}

	if err := s.validasiRegistrasiDanStatus(ctx, noRawat, req.TanggalPermintaan, req.JamPermintaan, statusLanjut, "mengubah"); err != nil {
		return nil, err
	}

	var kodeTindakanList []string
	for _, item := range req.Pemeriksaan {
		kodeTindakanList = append(kodeTindakanList, item.KodeTindakan)
	}

	adaMap, err := s.tindakanService.CekKeberadaanTindakanRadiologi(ctx, kodeTindakanList)
	if err != nil {
		s.log.Error("Gagal memeriksa keberadaan tindakan radiologi: %v", err)
		return nil, err
	}

	for _, kode := range kodeTindakanList {
		if !adaMap[kode] {
			return nil, apperror.NewBusinessError(fmt.Sprintf("Pemeriksaan radiologi '%s' tidak ditemukan atau tidak aktif", kode))
		}
	}

	if err := s.repo.UpdatePermintaanRadiologi(ctx, noPermintaan, req, kodeTindakanList); err != nil {
		s.log.Error("Gagal memperbarui permintaan radiologi no_permintaan %s: %v", noPermintaan, err)
		return nil, err
	}

	return s.repo.DetailPermintaanRadiologi(ctx, noPermintaan)
}

func (s *service) HapusPermintaanRadiologi(ctx context.Context, noRawat string, noPermintaan string, statusLanjut shared.StatusLanjut, kodeDokterLogin string) error {
	detail, err := s.repo.DetailPermintaanRadiologi(ctx, noPermintaan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NewNotFoundError("Data permintaan radiologi tidak ditemukan")
		}
		s.log.Error("Gagal memeriksa detail permintaan radiologi sebelum hapus no_permintaan %s: %v", noPermintaan, err)
		return err
	}

	if detail.NoRawat != noRawat {
		return apperror.NewBusinessError("Permintaan radiologi tidak sesuai dengan kunjungan pasien")
	}

	if statusLanjut != "" && !strings.EqualFold(detail.Status, string(statusLanjut)) {
		return apperror.NewNotFoundError("Data permintaan radiologi tidak ditemukan")
	}

	if detail.DokterPerujuk.KodeDokter != kodeDokterLogin {
		return apperror.NewForbiddenError("Hanya dokter pembuat order yang dapat membatalkan permintaan radiologi ini")
	}

	if (detail.TanggalSampel != "0000-00-00" && detail.TanggalSampel != "") || (detail.TanggalHasil != "0000-00-00" && detail.TanggalHasil != "") {
		return apperror.NewBusinessError("Permintaan radiologi sudah diproses oleh petugas dan tidak dapat dibatalkan")
	}

	for _, p := range detail.Pemeriksaan {
		if strings.EqualFold(p.StatusBayar, "Sudah") {
			return apperror.NewBusinessError("Permintaan radiologi yang pemeriksaannya telah dibayar tidak dapat dibatalkan")
		}
	}

	if err := s.validasiRegistrasiDanStatus(ctx, noRawat, "", "", statusLanjut, "membatalkan"); err != nil {
		return err
	}

	if err := s.repo.HapusPermintaanRadiologi(ctx, noPermintaan); err != nil {
		s.log.Error("Gagal menghapus permintaan radiologi no_permintaan %s: %v", noPermintaan, err)
		return err
	}

	return nil
}
