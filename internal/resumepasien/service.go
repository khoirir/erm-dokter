package resumepasien

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatinap"
	"erm-dokter/internal/rawatjalan"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type Service interface {
	DetailResumePasien(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) (*ResumePasien, error)
	RiwayatResumePasienByNoRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut) ([]ResumePasien, error)
	SimpanResumePasien(ctx context.Context, noRawat, kodeDokter string, statusLanjut shared.StatusLanjut, req SimpanResumePasienRequest) (*ResumePasien, error)
	UpdateResumePasien(ctx context.Context, kodeDokterLogin, noRawat string, statusLanjut shared.StatusLanjut, req UpdateResumePasienRequest) (*ResumePasien, error)
	HapusResumePasien(ctx context.Context, kodeDokterLogin, noRawat string, statusLanjut shared.StatusLanjut) error
}

type service struct {
	repo              Repository
	rawatJalanService rawatjalan.Service
	rawatInapService  rawatinap.Service
	maxEditJam        int
	log               *logger.Logger
}

func NewService(
	repo Repository,
	rawatJalanService rawatjalan.Service,
	rawatInapService rawatinap.Service,
	maxEditJam int,
	log *logger.Logger,
) Service {
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

func (s *service) validasiStatusLanjut(statusLanjut shared.StatusLanjut) error {
	if !statusLanjut.IsValid() {
		return apperror.NewBusinessError("Status lanjut tidak valid")
	}
	if statusLanjut != shared.StatusLanjutRawatJalan {
		return apperror.NewBusinessError("Resume pasien ini khusus untuk rawat jalan (Ralan)")
	}
	return nil
}

func (s *service) validasiRegistrasiDanStatus(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) error {
	if err := s.validasiStatusLanjut(statusLanjut); err != nil {
		return err
	}

	tanggalRegistrasiStr, jamRegistrasiStr, exists, err := s.rawatJalanService.GetWaktuRegistrasi(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal mengambil waktu registrasi no_rawat %s: %v", noRawat, err)
		return err
	}
	if !exists {
		return apperror.NewNotFoundError("Data kunjungan pasien tidak ditemukan")
	}

	waktuRegistrasi, err := shared.ParseWaktu(tanggalRegistrasiStr, jamRegistrasiStr)
	if err != nil {
		s.log.Error("Gagal parse waktu registrasi no_rawat %s (%s %s): %v", noRawat, tanggalRegistrasiStr, jamRegistrasiStr, err)
		return err
	}

	return s.validasiStatusKamarDanBatasWaktu(ctx, noRawat, waktuRegistrasi, tanggalRegistrasiStr, jamRegistrasiStr)
}

func (s *service) validasiStatusKamarDanBatasWaktu(ctx context.Context, noRawat string, waktuRegistrasi time.Time, tglRegStr, jamRegStr string) error {
	isAktifRanap, hasRecordKamar, err := s.rawatInapService.CekStatusKamarInap(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal cek status kamar inap untuk no_rawat %s: %v", noRawat, err)
		return err
	}

	if hasRecordKamar || isAktifRanap {
		return apperror.NewBusinessError("Pasien sudah terdaftar di rawat inap. Pengisian resume pasien ini hanya untuk kunjungan rawat jalan (Ralan).")
	}

	batasWaktu := waktuRegistrasi.Add(time.Duration(s.maxEditJam) * time.Hour)
	if time.Now().After(batasWaktu) {
		errMsg := fmt.Sprintf("Batas waktu resume pasien untuk kunjungan rawat jalan ini telah berakhir (maksimal %d jam dari waktu registrasi: %s %s)", s.maxEditJam, tglRegStr, jamRegStr)
		s.log.Warn("Resume pasien ditolak karena lewat batas %d jam untuk no_rawat %s: %s", s.maxEditJam, noRawat, errMsg)
		return apperror.NewForbiddenError(errMsg)
	}

	return nil
}

func (s *service) DetailResumePasien(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) (*ResumePasien, error) {
	if err := s.validasiStatusLanjut(statusLanjut); err != nil {
		return nil, err
	}

	item, err := s.repo.DetailResumePasien(ctx, noRawat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFoundError("Data resume pasien tidak ditemukan")
		}
		s.log.Error("Gagal mengambil detail resume pasien no_rawat %s: %v", noRawat, err)
		return nil, err
	}
	return item, nil
}

func (s *service) RiwayatResumePasienByNoRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut) ([]ResumePasien, error) {
	if err := s.validasiStatusLanjut(statusLanjut); err != nil {
		return nil, err
	}

	list, err := s.repo.RiwayatResumePasienByNoRM(ctx, noRM)
	if err != nil {
		s.log.Error("Gagal mengambil riwayat resume pasien no_rkm_medis %s: %v", noRM, err)
		return nil, err
	}
	return list, nil
}

func (s *service) SimpanResumePasien(ctx context.Context, noRawat, kodeDokter string, statusLanjut shared.StatusLanjut, req SimpanResumePasienRequest) (*ResumePasien, error) {
	if err := s.validasiRegistrasiDanStatus(ctx, noRawat, statusLanjut); err != nil {
		return nil, err
	}

	ada, err := s.repo.CekResumePasienAda(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal cek keberadaan resume pasien no_rawat %s: %v", noRawat, err)
		return nil, err
	}
	if ada {
		return nil, apperror.NewBusinessError("Resume pasien untuk kunjungan ini sudah ada")
	}

	item, err := s.repo.SimpanResumePasien(ctx, noRawat, kodeDokter, req)
	if err != nil {
		s.log.Error("Gagal menyimpan resume pasien no_rawat %s: %v", noRawat, err)
		return nil, err
	}

	s.log.Info("Berhasil menyimpan resume pasien no_rawat %s oleh dokter %s", noRawat, kodeDokter)
	return item, nil
}

func (s *service) UpdateResumePasien(ctx context.Context, kodeDokterLogin, noRawat string, statusLanjut shared.StatusLanjut, req UpdateResumePasienRequest) (*ResumePasien, error) {
	existing, err := s.repo.DetailResumePasien(ctx, noRawat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFoundError("Data resume pasien tidak ditemukan")
		}
		s.log.Error("Gagal mengambil resume pasien no_rawat %s untuk update: %v", noRawat, err)
		return nil, err
	}

	if existing.KodeDokter != kodeDokterLogin {
		s.log.Warn("Dokter %s mencoba mengubah resume pasien milik dokter %s (no_rawat: %s)", kodeDokterLogin, existing.KodeDokter, noRawat)
		return nil, apperror.NewForbiddenError("Anda tidak memiliki akses untuk mengubah resume pasien milik dokter lain")
	}

	if err := s.validasiRegistrasiDanStatus(ctx, noRawat, statusLanjut); err != nil {
		return nil, err
	}

	item, err := s.repo.UpdateResumePasien(ctx, noRawat, req)
	if err != nil {
		s.log.Error("Gagal memperbarui resume pasien no_rawat %s: %v", noRawat, err)
		return nil, err
	}

	s.log.Info("Berhasil memperbarui resume pasien no_rawat %s oleh dokter %s", noRawat, kodeDokterLogin)
	return item, nil
}

func (s *service) HapusResumePasien(ctx context.Context, kodeDokterLogin, noRawat string, statusLanjut shared.StatusLanjut) error {
	existing, err := s.repo.DetailResumePasien(ctx, noRawat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NewNotFoundError("Data resume pasien tidak ditemukan")
		}
		s.log.Error("Gagal mengambil resume pasien no_rawat %s untuk hapus: %v", noRawat, err)
		return err
	}

	if existing.KodeDokter != kodeDokterLogin {
		s.log.Warn("Dokter %s mencoba menghapus resume pasien milik dokter %s (no_rawat: %s)", kodeDokterLogin, existing.KodeDokter, noRawat)
		return apperror.NewForbiddenError("Anda tidak memiliki akses untuk menghapus resume pasien milik dokter lain")
	}

	if err := s.validasiRegistrasiDanStatus(ctx, noRawat, statusLanjut); err != nil {
		return err
	}

	if err := s.repo.HapusResumePasien(ctx, noRawat); err != nil {
		s.log.Error("Gagal menghapus resume pasien no_rawat %s: %v", noRawat, err)
		return err
	}

	s.log.Info("Berhasil menghapus resume pasien no_rawat %s oleh dokter %s", noRawat, kodeDokterLogin)
	return nil
}
