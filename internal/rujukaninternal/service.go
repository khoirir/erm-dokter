package rujukaninternal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatjalan"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type Service interface {
	DaftarOpsiPoliDokter(ctx context.Context, kodeDokterLogin, keyword string) ([]OpsiPoliDokter, error)
	DaftarRujukanInternal(ctx context.Context, noRawat string) ([]RujukanInternal, error)
	SimpanRujukanInternal(ctx context.Context, kodeDokterLogin, noRawat, targetKodePoli, targetKodeDokter string) (*RujukanInternal, error)
	HapusRujukanInternal(ctx context.Context, kodeDokterLogin, noRawat, targetKodeDokter string) error
}

type service struct {
	repo              Repository
	rawatJalanService rawatjalan.Service
	maxEditJam        int
	log               *logger.Logger
}

func NewService(repo Repository, rawatJalanService rawatjalan.Service, maxEditJam int, log *logger.Logger) Service {
	if maxEditJam <= 0 {
		maxEditJam = 48
	}
	return &service{
		repo:              repo,
		rawatJalanService: rawatJalanService,
		maxEditJam:        maxEditJam,
		log:               log,
	}
}

func (s *service) DaftarOpsiPoliDokter(ctx context.Context, kodeDokterLogin, keyword string) ([]OpsiPoliDokter, error) {
	cleanKeyword := strings.TrimSpace(keyword)
	list, err := s.repo.DaftarOpsiPoliDokter(ctx, kodeDokterLogin, cleanKeyword)
	if err != nil {
		s.log.Error("Gagal mengambil opsi poli dokter untuk dokter %s (keyword: %s): %v", kodeDokterLogin, cleanKeyword, err)
		return nil, err
	}
	return list, nil
}

func (s *service) DaftarRujukanInternal(ctx context.Context, noRawat string) ([]RujukanInternal, error) {
	list, err := s.repo.DaftarRujukanInternalByNoRawat(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal mengambil daftar rujukan internal no_rawat %s: %v", noRawat, err)
		return nil, err
	}
	return list, nil
}

func (s *service) SimpanRujukanInternal(ctx context.Context, kodeDokterLogin, noRawat, targetKodePoli, targetKodeDokter string) (*RujukanInternal, error) {
	if targetKodeDokter == kodeDokterLogin {
		s.log.Warn("Dokter %s mencoba membuat rujukan internal ke dokter pemeriksa pada no_rawat %s", kodeDokterLogin, noRawat)
		return nil, apperror.ValidationError{"id_tujuan": "Tujuan rujukan tidak dapat dipilih"}
	}

	if err := s.validasiRegistrasiDanStatus(ctx, noRawat, "disimpan"); err != nil {
		return nil, err
	}

	ada, err := s.repo.CekRujukanInternalAda(ctx, noRawat, targetKodeDokter)
	if err != nil {
		s.log.Error("Gagal memeriksa duplikasi rujukan internal no_rawat %s dokter %s: %v", noRawat, targetKodeDokter, err)
		return nil, err
	}
	if ada {
		s.log.Warn("Percobaan duplikasi rujukan internal no_rawat %s ke dokter %s", noRawat, targetKodeDokter)
		return nil, apperror.ValidationError{"id_tujuan": "Tujuan rujukan sudah dipilih"}
	}

	if err := s.repo.SimpanRujukanInternal(ctx, noRawat, targetKodePoli, targetKodeDokter); err != nil {
		s.log.Error("Gagal menyimpan data rujukan internal no_rawat '%s' poli '%s' dokter '%s': %v", noRawat, targetKodePoli, targetKodeDokter, err)
		return nil, err
	}

	s.log.Info("Berhasil menyimpan data rujukan internal no_rawat '%s' ke poli '%s' dokter '%s' oleh dokter '%s'", noRawat, targetKodePoli, targetKodeDokter, kodeDokterLogin)

	list, err := s.repo.DaftarRujukanInternalByNoRawat(ctx, noRawat)
	if err != nil {
		return nil, nil
	}

	for i := range list {
		if list[i].KodeDokter == targetKodeDokter && list[i].KodePoli == targetKodePoli {
			return &list[i], nil
		}
	}

	return nil, nil
}

func (s *service) HapusRujukanInternal(ctx context.Context, kodeDokterLogin, noRawat, targetKodeDokter string) error {
	if err := s.validasiRegistrasiDanStatus(ctx, noRawat, "dihapus"); err != nil {
		return err
	}

	if err := s.repo.HapusRujukanInternal(ctx, noRawat, targetKodeDokter); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NewNotFoundError(
				"Data rujukan internal tidak ditemukan",
				fmt.Sprintf("Data rujukan internal tidak ditemukan untuk no_rawat '%s' dokter '%s'", noRawat, targetKodeDokter),
			)
		}
		s.log.Error("Gagal menghapus data rujukan internal no_rawat '%s' dokter '%s' oleh dokter '%s': %v", noRawat, targetKodeDokter, kodeDokterLogin, err)
		return err
	}

	s.log.Info("Berhasil menghapus data rujukan internal no_rawat '%s' dokter '%s' oleh dokter '%s'", noRawat, targetKodeDokter, kodeDokterLogin)
	return nil
}

func (s *service) validasiRegistrasiDanStatus(ctx context.Context, noRawat, aksi string) error {
	infoRegistrasi, err := s.rawatJalanService.GetInfoRegistrasi(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal mengambil data registrasi no_rawat '%s': %v", noRawat, err)
		return err
	}
	if infoRegistrasi == nil {
		return apperror.NewNotFoundError(
			"Data registrasi kunjungan pasien tidak ditemukan",
			fmt.Sprintf("Data registrasi kunjungan pasien tidak ditemukan untuk no_rawat '%s'", noRawat),
		)
	}

	if infoRegistrasi.StatusBayar == "Sudah Bayar" && infoRegistrasi.KodePenjamin == "BPJ" {
		return apperror.NewBusinessError(fmt.Sprintf("Pasien BPJS sudah bayar, rujukan tidak dapat %s", aksi))
	}

	tanggalRegistrasi := infoRegistrasi.TanggalRegistrasi
	jamRegistrasi := infoRegistrasi.JamRegistrasi

	waktuRegistrasi, err := shared.ParseWaktu(tanggalRegistrasi, jamRegistrasi)
	if err != nil {
		s.log.Error("Gagal parse waktu registrasi no_rawat '%s' (%s %s): %v", noRawat, tanggalRegistrasi, jamRegistrasi, err)
		return err
	}

	batasWaktu := waktuRegistrasi.Add(time.Duration(s.maxEditJam) * time.Hour)
	if time.Now().After(batasWaktu) {
		return apperror.NewForbiddenError(
			fmt.Sprintf("Rujukan internal tidak dapat %s, melebihi batas %d jam", aksi, s.maxEditJam),
			fmt.Sprintf("Rujukan internal no_rawat '%s' ditolak untuk %s karena melewati batas %d jam dari registrasi (%s)", noRawat, aksi, s.maxEditJam, waktuRegistrasi.Format("2006-01-02 15:04:05")),
		)
	}

	return nil
}
