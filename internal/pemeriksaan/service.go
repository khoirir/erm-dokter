package pemeriksaan

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatinap"
	"erm-dokter/internal/rawatjalan"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type Service interface {
	DaftarPemeriksaan(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter FilterDaftarPemeriksaan) ([]Pemeriksaan, shared.PaginationMeta, error)
	DaftarPemeriksaanByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter FilterDaftarPemeriksaan) ([]Pemeriksaan, shared.PaginationMeta, error)
	DetailPemeriksaan(ctx context.Context, id IdPemeriksaan, statusLanjut shared.StatusLanjut) (*Pemeriksaan, error)
	DaftarKesadaran(ctx context.Context) []OpsiReferensi
	SimpanPemeriksaan(ctx context.Context, kodeDokter string, statusLanjut shared.StatusLanjut, req SimpanPemeriksaanRequest) (*Pemeriksaan, error)
	UpdatePemeriksaan(ctx context.Context, kodeDokter string, id IdPemeriksaan, statusLanjut shared.StatusLanjut, req UpdatePemeriksaanRequest) (*Pemeriksaan, error)
	HapusPemeriksaan(ctx context.Context, kodeDokter string, id IdPemeriksaan, statusLanjut shared.StatusLanjut) error
}

type service struct {
	repo              Repository
	rawatJalanService rawatjalan.Service
	rawatInapService  rawatinap.Service
	maxEditJam        int
	log               *logger.Logger
}

func NewService(repo Repository, rawatJalanService rawatjalan.Service, rawatInapService rawatinap.Service, maxEditJam int, log *logger.Logger) Service {
	if maxEditJam <= 0 {
		maxEditJam = 48
	}
	return &service{
		repo:              repo,
		rawatJalanService: rawatJalanService,
		rawatInapService:  rawatInapService,
		maxEditJam:        maxEditJam,
		log:               log,
	}
}

func (s *service) DaftarPemeriksaan(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter FilterDaftarPemeriksaan) ([]Pemeriksaan, shared.PaginationMeta, error) {
	rawParts := strings.Split(noRawat, ",")
	var listNoRawat []string
	for _, p := range rawParts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			listNoRawat = append(listNoRawat, trimmed)
		}
	}

	daftarPemeriksaan, totalData, err := s.repo.DaftarPemeriksaan(ctx, listNoRawat, statusLanjut, filter)
	if err != nil {
		s.log.Error("Gagal query daftar pemeriksaan %s: %v", noRawat, err)
		return nil, shared.PaginationMeta{}, err
	}

	return daftarPemeriksaan, shared.NewPaginationMeta(totalData, filter.Page, filter.Limit), nil
}

func (s *service) DaftarPemeriksaanByRM(ctx context.Context, noRekamMedis string, statusLanjut shared.StatusLanjut, filter FilterDaftarPemeriksaan) ([]Pemeriksaan, shared.PaginationMeta, error) {
	riwayatKunjungan, err := s.rawatJalanService.RiwayatKunjunganPasien(ctx, noRekamMedis)
	if err != nil {
		s.log.Error("Gagal mengambil riwayat kunjungan untuk RM %s: %v", noRekamMedis, err)
		return nil, shared.PaginationMeta{}, err
	}

	if len(riwayatKunjungan) == 0 {
		return make([]Pemeriksaan, 0), shared.NewPaginationMeta(0, filter.Page, filter.Limit), nil
	}

	listNoRawat := make([]string, len(riwayatKunjungan))
	for i, k := range riwayatKunjungan {
		listNoRawat[i] = k.NoRawat
	}

	daftarPemeriksaan, totalData, err := s.repo.DaftarPemeriksaan(ctx, listNoRawat, statusLanjut, filter)
	if err != nil {
		s.log.Error("Gagal query daftar pemeriksaan by RM %s: %v", noRekamMedis, err)
		return nil, shared.PaginationMeta{}, err
	}

	return daftarPemeriksaan, shared.NewPaginationMeta(totalData, filter.Page, filter.Limit), nil
}

func (s *service) DetailPemeriksaan(ctx context.Context, id IdPemeriksaan, statusLanjut shared.StatusLanjut) (*Pemeriksaan, error) {
	pemeriksaan, err := s.repo.DetailPemeriksaan(ctx, id, statusLanjut)
	if err != nil {
		s.log.Error("Gagal query detail pemeriksaan %+v (%s): %v", id, statusLanjut, err)
		return nil, err
	}
	if pemeriksaan == nil {
		return nil, apperror.NewNotFoundError("Detail pemeriksaan tidak ditemukan")
	}

	return pemeriksaan, nil
}

func (s *service) DaftarKesadaran(ctx context.Context) []OpsiReferensi {
	opsi := make([]OpsiReferensi, len(ListKesadaran))
	for i, k := range ListKesadaran {
		opsi[i] = OpsiReferensi{
			Value: string(k),
			Label: string(k),
		}
	}
	return opsi
}

func (s *service) SimpanPemeriksaan(ctx context.Context, kodeDokter string, statusLanjut shared.StatusLanjut, req SimpanPemeriksaanRequest) (*Pemeriksaan, error) {
	if err := s.validasiRegistrasiDanStatus(ctx, req.NoRawat, req.TanggalPemeriksaan, req.JamPemeriksaan, statusLanjut, "membuat"); err != nil {
		return nil, err
	}

	if err := s.repo.SimpanPemeriksaan(ctx, kodeDokter, statusLanjut, req); err != nil {
		var mysqlErr *mysql.MySQLError
		if (errors.As(err, &mysqlErr) && mysqlErr.Number == 1062) || strings.Contains(err.Error(), "1062") || strings.Contains(err.Error(), "Duplicate entry") {
			errMsg := fmt.Sprintf("Data pemeriksaan pada tanggal %s jam %s sudah pernah disimpan sebelumnya. Silakan sesuaikan jam pemeriksaan.", req.TanggalPemeriksaan, req.JamPemeriksaan)
			s.log.Warn("Penyimpanan pemeriksaan duplikat untuk no_rawat %s (%s %s): %s", req.NoRawat, req.TanggalPemeriksaan, req.JamPemeriksaan, errMsg)
			return nil, apperror.ValidationError{
				"jam_pemeriksaan": errMsg,
			}
		}
		s.log.Error("Gagal menyimpan pemeriksaan no_rawat %s (%s) oleh dokter %s: %v", req.NoRawat, statusLanjut, kodeDokter, err)
		return nil, err
	}

	idPemeriksaan := IdPemeriksaan{
		NoRawat:            req.NoRawat,
		TanggalPemeriksaan: req.TanggalPemeriksaan,
		JamPemeriksaan:     req.JamPemeriksaan,
	}

	detail, err := s.repo.DetailPemeriksaan(ctx, idPemeriksaan, statusLanjut)
	if err != nil {
		s.log.Error("Gagal mengambil detail pemeriksaan setelah simpan no_rawat %s (%s): %v", req.NoRawat, statusLanjut, err)
		return nil, err
	}
	if detail == nil {
		return nil, apperror.NewNotFoundError("Data pemeriksaan yang baru disimpan tidak ditemukan")
	}

	s.log.Info("Berhasil menyimpan pemeriksaan no_rawat %s (%s) oleh dokter %s", req.NoRawat, statusLanjut, kodeDokter)
	return detail, nil
}

func (s *service) UpdatePemeriksaan(ctx context.Context, kodeDokter string, id IdPemeriksaan, statusLanjut shared.StatusLanjut, req UpdatePemeriksaanRequest) (*Pemeriksaan, error) {
	pemeriksaan, err := s.repo.DetailPemeriksaan(ctx, id, statusLanjut)
	if err != nil {
		s.log.Error("Gagal mengambil detail pemeriksaan untuk update %+v (%s): %v", id, statusLanjut, err)
		return nil, err
	}
	if pemeriksaan == nil {
		return nil, apperror.NewNotFoundError("Data pemeriksaan tidak ditemukan")
	}

	if pemeriksaan.KodeDokterPetugas != kodeDokter {
		s.log.Warn("Percobaan mengubah pemeriksaan no_rawat %s (%s %s) oleh dokter %s ditolak: diinput oleh %s (%s)", id.NoRawat, id.TanggalPemeriksaan, id.JamPemeriksaan, kodeDokter, pemeriksaan.KodeDokterPetugas, pemeriksaan.NamaDokterPetugas)
		return nil, apperror.NewForbiddenError(fmt.Sprintf("Anda tidak memiliki hak akses untuk mengubah data pemeriksaan ini karena diinput oleh dokter/petugas lain (%s)", pemeriksaan.NamaDokterPetugas))
	}

	if err := s.validasiRegistrasiDanStatus(ctx, id.NoRawat, req.TanggalPemeriksaan, req.JamPemeriksaan, statusLanjut, "mengubah"); err != nil {
		return nil, err
	}

	if err := s.repo.UpdatePemeriksaan(ctx, id, statusLanjut, req); err != nil {
		var mysqlErr *mysql.MySQLError
		if (errors.As(err, &mysqlErr) && mysqlErr.Number == 1062) || strings.Contains(err.Error(), "1062") || strings.Contains(err.Error(), "Duplicate entry") {
			errMsg := fmt.Sprintf("Data pemeriksaan pada tanggal %s jam %s sudah pernah disimpan sebelumnya. Silakan sesuaikan jam pemeriksaan.", req.TanggalPemeriksaan, req.JamPemeriksaan)
			s.log.Warn("Pembaruan pemeriksaan duplikat untuk no_rawat %s (%s %s): %s", id.NoRawat, req.TanggalPemeriksaan, req.JamPemeriksaan, errMsg)
			return nil, apperror.ValidationError{
				"jam_pemeriksaan": errMsg,
			}
		}
		s.log.Error("Gagal memperbarui pemeriksaan no_rawat %s (%s) oleh dokter %s: %v", id.NoRawat, statusLanjut, kodeDokter, err)
		return nil, err
	}

	updatedId := IdPemeriksaan{
		NoRawat:            id.NoRawat,
		TanggalPemeriksaan: req.TanggalPemeriksaan,
		JamPemeriksaan:     req.JamPemeriksaan,
	}

	detail, err := s.repo.DetailPemeriksaan(ctx, updatedId, statusLanjut)
	if err != nil {
		s.log.Error("Gagal mengambil detail pemeriksaan setelah update no_rawat %s (%s): %v", id.NoRawat, statusLanjut, err)
		return nil, err
	}
	if detail == nil {
		return nil, apperror.NewNotFoundError("Data pemeriksaan yang baru diperbarui tidak ditemukan")
	}

	s.log.Info("Berhasil memperbarui data pemeriksaan no_rawat %s (%s %s -> %s %s) oleh dokter %s", id.NoRawat, id.TanggalPemeriksaan, id.JamPemeriksaan, req.TanggalPemeriksaan, req.JamPemeriksaan, kodeDokter)
	return detail, nil
}

func (s *service) HapusPemeriksaan(ctx context.Context, kodeDokter string, id IdPemeriksaan, statusLanjut shared.StatusLanjut) error {
	pemeriksaan, err := s.repo.DetailPemeriksaan(ctx, id, statusLanjut)
	if err != nil {
		s.log.Error("Gagal mengambil detail pemeriksaan untuk hapus %+v (%s): %v", id, statusLanjut, err)
		return err
	}
	if pemeriksaan == nil {
		return apperror.NewNotFoundError("Data pemeriksaan tidak ditemukan")
	}

	if pemeriksaan.KodeDokterPetugas != kodeDokter {
		s.log.Warn("Percobaan menghapus pemeriksaan no_rawat %s (%s %s) oleh dokter %s ditolak: diinput oleh %s (%s)", id.NoRawat, id.TanggalPemeriksaan, id.JamPemeriksaan, kodeDokter, pemeriksaan.KodeDokterPetugas, pemeriksaan.NamaDokterPetugas)
		return apperror.NewForbiddenError(fmt.Sprintf("Anda tidak memiliki hak akses untuk menghapus data pemeriksaan ini karena diinput oleh dokter/petugas lain (%s)", pemeriksaan.NamaDokterPetugas))
	}

	if err := s.validasiRegistrasiDanStatus(ctx, id.NoRawat, "", "", statusLanjut, "menghapus"); err != nil {
		return err
	}

	if err := s.repo.HapusPemeriksaan(ctx, id, statusLanjut); err != nil {
		s.log.Error("Gagal menghapus pemeriksaan no_rawat %s (%s) oleh dokter %s: %v", id.NoRawat, statusLanjut, kodeDokter, err)
		return err
	}

	s.log.Info("Berhasil menghapus data pemeriksaan no_rawat %s (%s %s) oleh dokter %s", id.NoRawat, id.TanggalPemeriksaan, id.JamPemeriksaan, kodeDokter)
	return nil
}

func (s *service) validasiRegistrasiDanStatus(ctx context.Context, noRawat, tanggalPeriksa, jamPeriksa string, statusLanjut shared.StatusLanjut, action string) error {
	tanggalRegistrasiStr, jamRegistrasiStr, exists, err := s.rawatJalanService.GetWaktuRegistrasi(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal mengambil data registrasi no_rawat %s: %v", noRawat, err)
		return err
	}
	if !exists {
		return apperror.NewNotFoundError("Data registrasi kunjungan pasien tidak ditemukan")
	}

	waktuRegistrasi, err := shared.ParseWaktu(tanggalRegistrasiStr, jamRegistrasiStr)
	if err != nil {
		s.log.Error("Gagal parse waktu registrasi no_rawat %s (%s %s): %v", noRawat, tanggalRegistrasiStr, jamRegistrasiStr, err)
		return err
	}

	if tanggalPeriksa != "" && jamPeriksa != "" {
		waktuPemeriksaan, err := shared.ParseWaktu(tanggalPeriksa, jamPeriksa)
		if err != nil {
			return apperror.NewBusinessError(err.Error())
		}

		if waktuPemeriksaan.Before(waktuRegistrasi) {
			errs := apperror.ValidationError{
				"tanggal_pemeriksaan": fmt.Sprintf("Waktu pemeriksaan (%s %s) tidak boleh lebih awal dari waktu registrasi pasien (%s %s)", tanggalPeriksa, jamPeriksa, tanggalRegistrasiStr, jamRegistrasiStr),
			}
			s.log.Warn("Validasi waktu pemeriksaan gagal untuk no_rawat %s: %+v", noRawat, errs)
			return errs
		}
	}

	return s.validasiStatusKamarDanBatasWaktu(ctx, noRawat, statusLanjut, waktuRegistrasi, tanggalRegistrasiStr, jamRegistrasiStr)
}

func (s *service) validasiStatusKamarDanBatasWaktu(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, waktuRegistrasi time.Time, tglRegStr, jamRegStr string) error {
	isAktifRanap, hasRecordKamar, err := s.rawatInapService.CekStatusKamarInap(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal cek status kamar inap untuk no_rawat %s: %v", noRawat, err)
		return err
	}

	if hasRecordKamar {
		if !isAktifRanap {
			return apperror.NewBusinessError("Pasien rawat inap sudah keluar / checkout dari kamar inap")
		}
		return nil
	}

	if strings.EqualFold(string(statusLanjut), string(shared.StatusLanjutRawatInap)) {
		return apperror.NewBusinessError("Pasien belum/tidak terdaftar di kamar inap. Pemeriksaan harus menggunakan status 'Ralan'.")
	}

	batasWaktu := waktuRegistrasi.Add(time.Duration(s.maxEditJam) * time.Hour)
	if time.Now().After(batasWaktu) {
		errMsg := fmt.Sprintf("Batas waktu pemeriksaan medis untuk kunjungan rawat jalan ini telah berakhir (maksimal %d jam dari waktu registrasi: %s %s)", s.maxEditJam, tglRegStr, jamRegStr)
		s.log.Warn("Pemeriksaan ditolak karena lewat batas %d jam untuk no_rawat %s: %s", s.maxEditJam, noRawat, errMsg)
		return apperror.NewForbiddenError(errMsg)
	}

	return nil
}
