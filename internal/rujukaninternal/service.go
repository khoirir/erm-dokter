package rujukaninternal

import (
	"context"
	"database/sql"
	"errors"
	"strings"

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
	log               *logger.Logger
}

func NewService(repo Repository, rawatJalanService rawatjalan.Service, log *logger.Logger) Service {
	return &service{
		repo:              repo,
		rawatJalanService: rawatJalanService,
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
		s.log.Warn("Dokter %s mencoba membuat rujukan internal ke dirinya sendiri pada no_rawat %s", kodeDokterLogin, noRawat)
		return nil, apperror.ValidationError{
			"id_tujuan": "Tidak dapat membuat rujukan internal ke diri sendiri",
		}
	}

	if err := s.validasiBatasWaktu(ctx, noRawat, "dibuat"); err != nil {
		return nil, err
	}

	ada, err := s.repo.CekRujukanInternalAda(ctx, noRawat, targetKodeDokter)
	if err != nil {
		s.log.Error("Gagal memeriksa duplikasi rujukan internal no_rawat %s dokter %s: %v", noRawat, targetKodeDokter, err)
		return nil, err
	}
	if ada {
		s.log.Warn("Percobaan duplikasi rujukan internal no_rawat %s ke dokter %s", noRawat, targetKodeDokter)
		return nil, apperror.ValidationError{
			"id_tujuan": "Pasien sudah pernah dirujuk ke dokter ini pada kunjungan yang sama",
		}
	}

	if err := s.repo.SimpanRujukanInternal(ctx, noRawat, targetKodePoli, targetKodeDokter); err != nil {
		s.log.Error("Gagal menyimpan rujukan internal no_rawat %s poli %s dokter %s: %v", noRawat, targetKodePoli, targetKodeDokter, err)
		return nil, err
	}

	s.log.Info("Berhasil menyimpan rujukan internal no_rawat %s ke poli %s dokter %s oleh dokter %s", noRawat, targetKodePoli, targetKodeDokter, kodeDokterLogin)

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
	if err := s.validasiBatasWaktu(ctx, noRawat, "menghapus"); err != nil {
		return err
	}

	if err := s.repo.HapusRujukanInternal(ctx, noRawat, targetKodeDokter); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NewNotFoundError("Data rujukan internal tidak ditemukan")
		}
		s.log.Error("Gagal menghapus rujukan internal no_rawat %s dokter %s: %v", noRawat, targetKodeDokter, err)
		return err
	}

	s.log.Info("Berhasil menghapus rujukan internal no_rawat %s ke dokter %s oleh dokter %s", noRawat, targetKodeDokter, kodeDokterLogin)
	return nil
}

func (s *service) validasiBatasWaktu(ctx context.Context, noRawat, aksi string) error {
	tanggalRegistrasiStr, jamRegistrasiStr, exists, err := s.rawatJalanService.GetWaktuRegistrasi(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal mengambil waktu registrasi no_rawat %s: %v", noRawat, err)
		return err
	}
	if !exists {
		return apperror.NewNotFoundError("Data kunjungan pasien tidak ditemukan")
	}

	if err := shared.ValidasiBatasWaktuRekamMedis(tanggalRegistrasiStr, jamRegistrasiStr, 48, aksi); err != nil {
		s.log.Warn("Validasi batas waktu 48 jam rujukan internal gagal untuk no_rawat %s (%s %s): %v", noRawat, tanggalRegistrasiStr, jamRegistrasiStr, err)
		return err
	}

	return nil
}
