package penilaianmedis

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatjalan"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type Service interface {
	Referensi(ctx context.Context) ReferensiPenilaianMedis

	DetailPenilaianMedisRalan(ctx context.Context, noRawat string) (*PenilaianMedisRalan, error)
	RiwayatPenilaianMedisRalanByNoRM(ctx context.Context, noRM string) ([]PenilaianMedisRalan, error)
	SimpanPenilaianMedisRalan(ctx context.Context, kodeDokterLogin, noRawat string, req SimpanPenilaianMedisRalanRequest) (*PenilaianMedisRalan, error)
	UpdatePenilaianMedisRalan(ctx context.Context, kodeDokterLogin, noRawat string, req UpdatePenilaianMedisRalanRequest) (*PenilaianMedisRalan, error)
	HapusPenilaianMedisRalan(ctx context.Context, kodeDokterLogin, noRawat string) error

	DetailPenilaianMedisIGD(ctx context.Context, noRawat string) (*PenilaianMedisIGD, error)
	RiwayatPenilaianMedisIGDByNoRM(ctx context.Context, noRM string) ([]PenilaianMedisIGD, error)
	SimpanPenilaianMedisIGD(ctx context.Context, kodeDokterLogin, noRawat string, req SimpanPenilaianMedisIGDRequest) (*PenilaianMedisIGD, error)
	UpdatePenilaianMedisIGD(ctx context.Context, kodeDokterLogin, noRawat string, req UpdatePenilaianMedisIGDRequest) (*PenilaianMedisIGD, error)
	HapusPenilaianMedisIGD(ctx context.Context, kodeDokterLogin, noRawat string) error
}

type service struct {
	repo              Repository
	rawatJalanService rawatjalan.Service
	log               *logger.Logger
	encryptionKey     string
}

func NewService(repo Repository, rawatJalanService rawatjalan.Service, log *logger.Logger, encryptionKey string) Service {
	return &service{
		repo:              repo,
		rawatJalanService: rawatJalanService,
		log:               log,
		encryptionKey:     encryptionKey,
	}
}

func (s *service) Referensi(ctx context.Context) ReferensiPenilaianMedis {
	return GetReferensiPenilaianMedis()
}

func (s *service) DetailPenilaianMedisRalan(ctx context.Context, noRawat string) (*PenilaianMedisRalan, error) {
	item, err := s.repo.DetailPenilaianMedisRalan(ctx, noRawat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFoundError("Data penilaian awal medis rawat jalan tidak ditemukan")
		}
		s.log.Error("Gagal mengambil detail penilaian medis ralan no_rawat %s: %v", noRawat, err)
		return nil, err
	}

	if encryptedIdKunjungan, err := crypto.Encrypt(item.NoRawat, s.encryptionKey); err == nil {
		item.IdKunjungan = encryptedIdKunjungan
	}
	if encryptedIdPasien, err := crypto.Encrypt(item.NoRM, s.encryptionKey); err == nil {
		item.IdPasien = encryptedIdPasien
	}

	return item, nil
}

func (s *service) RiwayatPenilaianMedisRalanByNoRM(ctx context.Context, noRM string) ([]PenilaianMedisRalan, error) {
	list, err := s.repo.RiwayatPenilaianMedisRalanByNoRM(ctx, noRM)
	if err != nil {
		s.log.Error("Gagal mengambil riwayat penilaian medis ralan no_rkm_medis %s: %v", noRM, err)
		return nil, err
	}

	for i := range list {
		item := &list[i]
		if encryptedIdKunjungan, err := crypto.Encrypt(item.NoRawat, s.encryptionKey); err == nil {
			item.IdKunjungan = encryptedIdKunjungan
		}
		if encryptedIdPasien, err := crypto.Encrypt(item.NoRM, s.encryptionKey); err == nil {
			item.IdPasien = encryptedIdPasien
		}
	}

	return list, nil
}

func (s *service) SimpanPenilaianMedisRalan(ctx context.Context, kodeDokterLogin, noRawat string, req SimpanPenilaianMedisRalanRequest) (*PenilaianMedisRalan, error) {
	if errs := req.Validate(); errs != nil {
		return nil, errs
	}

	if err := s.validasiWaktuPenilaianMedis(ctx, noRawat, req.TanggalPenilaian, "dibuat"); err != nil {
		return nil, err
	}

	ada, err := s.repo.CekPenilaianMedisRalanAda(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal memeriksa keberadaan penilaian medis ralan no_rawat %s: %v", noRawat, err)
		return nil, err
	}
	if ada {
		s.log.Warn("Percobaan duplikasi penilaian medis ralan pada no_rawat %s", noRawat)
		return nil, apperror.NewBusinessError("Penilaian awal medis rawat jalan untuk kunjungan ini sudah ada")
	}

	if err := s.repo.SimpanPenilaianMedisRalan(ctx, noRawat, kodeDokterLogin, req); err != nil {
		s.log.Error("Gagal menyimpan penilaian medis ralan no_rawat %s oleh dokter %s: %v", noRawat, kodeDokterLogin, err)
		return nil, err
	}

	s.log.Info("Berhasil menyimpan penilaian medis ralan no_rawat %s oleh dokter %s", noRawat, kodeDokterLogin)
	return s.DetailPenilaianMedisRalan(ctx, noRawat)
}

func (s *service) UpdatePenilaianMedisRalan(ctx context.Context, kodeDokterLogin, noRawat string, req UpdatePenilaianMedisRalanRequest) (*PenilaianMedisRalan, error) {
	if errs := req.Validate(); errs != nil {
		return nil, errs
	}

	existing, err := s.repo.DetailPenilaianMedisRalan(ctx, noRawat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFoundError("Data penilaian awal medis rawat jalan tidak ditemukan")
		}
		s.log.Error("Gagal memeriksa data penilaian medis ralan sebelum update no_rawat %s: %v", noRawat, err)
		return nil, err
	}

	if existing.KodeDokter != kodeDokterLogin {
		s.log.Warn("Percobaan mengubah penilaian medis ralan no_rawat %s ditolak: dibuat oleh dokter %s (%s), dicoba oleh %s",
			noRawat, existing.KodeDokter, existing.NamaDokter, kodeDokterLogin)
		return nil, apperror.NewForbiddenError(fmt.Sprintf("Anda tidak memiliki hak akses untuk mengubah penilaian medis ini karena dibuat oleh dokter lain (%s)", existing.NamaDokter))
	}

	if err := s.validasiWaktuPenilaianMedis(ctx, noRawat, req.TanggalPenilaian, "diubah"); err != nil {
		return nil, err
	}

	if err := s.repo.UpdatePenilaianMedisRalan(ctx, noRawat, req); err != nil {
		s.log.Error("Gagal memperbarui penilaian medis ralan no_rawat %s oleh dokter %s: %v", noRawat, kodeDokterLogin, err)
		return nil, err
	}

	s.log.Info("Berhasil memperbarui penilaian medis ralan no_rawat %s oleh dokter %s", noRawat, kodeDokterLogin)
	return s.DetailPenilaianMedisRalan(ctx, noRawat)
}

func (s *service) HapusPenilaianMedisRalan(ctx context.Context, kodeDokterLogin, noRawat string) error {
	existing, err := s.repo.DetailPenilaianMedisRalan(ctx, noRawat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NewNotFoundError("Data penilaian awal medis rawat jalan tidak ditemukan")
		}
		s.log.Error("Gagal memeriksa data penilaian medis ralan sebelum hapus no_rawat %s: %v", noRawat, err)
		return err
	}

	if existing.KodeDokter != kodeDokterLogin {
		s.log.Warn("Percobaan menghapus penilaian medis ralan no_rawat %s ditolak: dibuat oleh dokter %s (%s), dicoba oleh %s",
			noRawat, existing.KodeDokter, existing.NamaDokter, kodeDokterLogin)
		return apperror.NewForbiddenError(fmt.Sprintf("Anda tidak memiliki hak akses untuk menghapus penilaian medis ini karena dibuat oleh dokter lain (%s)", existing.NamaDokter))
	}

	if err := s.validasiWaktuPenilaianMedis(ctx, noRawat, "", "dihapus"); err != nil {
		return err
	}

	if err := s.repo.HapusPenilaianMedisRalan(ctx, noRawat); err != nil {
		s.log.Error("Gagal menghapus penilaian medis ralan no_rawat %s oleh dokter %s: %v", noRawat, kodeDokterLogin, err)
		return err
	}

	s.log.Info("Berhasil menghapus penilaian medis ralan no_rawat %s oleh dokter %s", noRawat, kodeDokterLogin)
	return nil
}

func (s *service) DetailPenilaianMedisIGD(ctx context.Context, noRawat string) (*PenilaianMedisIGD, error) {
	item, err := s.repo.DetailPenilaianMedisIGD(ctx, noRawat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFoundError("Data penilaian awal medis IGD tidak ditemukan")
		}
		s.log.Error("Gagal mengambil detail penilaian medis IGD no_rawat %s: %v", noRawat, err)
		return nil, err
	}

	if encryptedIdKunjungan, err := crypto.Encrypt(item.NoRawat, s.encryptionKey); err == nil {
		item.IdKunjungan = encryptedIdKunjungan
	}
	if encryptedIdPasien, err := crypto.Encrypt(item.NoRM, s.encryptionKey); err == nil {
		item.IdPasien = encryptedIdPasien
	}

	return item, nil
}

func (s *service) RiwayatPenilaianMedisIGDByNoRM(ctx context.Context, noRM string) ([]PenilaianMedisIGD, error) {
	list, err := s.repo.RiwayatPenilaianMedisIGDByNoRM(ctx, noRM)
	if err != nil {
		s.log.Error("Gagal mengambil riwayat penilaian medis IGD no_rkm_medis %s: %v", noRM, err)
		return nil, err
	}

	for i := range list {
		item := &list[i]
		if encryptedIdKunjungan, err := crypto.Encrypt(item.NoRawat, s.encryptionKey); err == nil {
			item.IdKunjungan = encryptedIdKunjungan
		}
		if encryptedIdPasien, err := crypto.Encrypt(item.NoRM, s.encryptionKey); err == nil {
			item.IdPasien = encryptedIdPasien
		}
	}

	return list, nil
}

func (s *service) SimpanPenilaianMedisIGD(ctx context.Context, kodeDokterLogin, noRawat string, req SimpanPenilaianMedisIGDRequest) (*PenilaianMedisIGD, error) {
	if errs := req.Validate(); errs != nil {
		return nil, errs
	}

	if err := s.validasiWaktuPenilaianMedis(ctx, noRawat, req.TanggalPenilaian, "dibuat"); err != nil {
		return nil, err
	}

	ada, err := s.repo.CekPenilaianMedisIGDAda(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal memeriksa keberadaan penilaian medis IGD no_rawat %s: %v", noRawat, err)
		return nil, err
	}
	if ada {
		s.log.Warn("Percobaan duplikasi penilaian medis IGD pada no_rawat %s", noRawat)
		return nil, apperror.NewBusinessError("Penilaian awal medis IGD untuk kunjungan ini sudah ada")
	}

	if err := s.repo.SimpanPenilaianMedisIGD(ctx, noRawat, kodeDokterLogin, req); err != nil {
		s.log.Error("Gagal menyimpan penilaian medis IGD no_rawat %s oleh dokter %s: %v", noRawat, kodeDokterLogin, err)
		return nil, err
	}

	s.log.Info("Berhasil menyimpan penilaian medis IGD no_rawat %s oleh dokter %s", noRawat, kodeDokterLogin)
	return s.DetailPenilaianMedisIGD(ctx, noRawat)
}

func (s *service) UpdatePenilaianMedisIGD(ctx context.Context, kodeDokterLogin, noRawat string, req UpdatePenilaianMedisIGDRequest) (*PenilaianMedisIGD, error) {
	if errs := req.Validate(); errs != nil {
		return nil, errs
	}

	existing, err := s.repo.DetailPenilaianMedisIGD(ctx, noRawat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFoundError("Data penilaian awal medis IGD tidak ditemukan")
		}
		s.log.Error("Gagal memeriksa data penilaian medis IGD sebelum update no_rawat %s: %v", noRawat, err)
		return nil, err
	}

	if existing.KodeDokter != kodeDokterLogin {
		s.log.Warn("Percobaan mengubah penilaian medis IGD no_rawat %s ditolak: dibuat oleh dokter %s (%s), dicoba oleh %s",
			noRawat, existing.KodeDokter, existing.NamaDokter, kodeDokterLogin)
		return nil, apperror.NewForbiddenError(fmt.Sprintf("Anda tidak memiliki hak akses untuk mengubah penilaian medis ini karena dibuat oleh dokter lain (%s)", existing.NamaDokter))
	}

	if err := s.validasiWaktuPenilaianMedis(ctx, noRawat, req.TanggalPenilaian, "diubah"); err != nil {
		return nil, err
	}

	if err := s.repo.UpdatePenilaianMedisIGD(ctx, noRawat, req); err != nil {
		s.log.Error("Gagal memperbarui penilaian medis IGD no_rawat %s oleh dokter %s: %v", noRawat, kodeDokterLogin, err)
		return nil, err
	}

	s.log.Info("Berhasil memperbarui penilaian medis IGD no_rawat %s oleh dokter %s", noRawat, kodeDokterLogin)
	return s.DetailPenilaianMedisIGD(ctx, noRawat)
}

func (s *service) HapusPenilaianMedisIGD(ctx context.Context, kodeDokterLogin, noRawat string) error {
	existing, err := s.repo.DetailPenilaianMedisIGD(ctx, noRawat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NewNotFoundError("Data penilaian awal medis IGD tidak ditemukan")
		}
		s.log.Error("Gagal memeriksa data penilaian medis IGD sebelum hapus no_rawat %s: %v", noRawat, err)
		return err
	}

	if existing.KodeDokter != kodeDokterLogin {
		s.log.Warn("Percobaan menghapus penilaian medis IGD no_rawat %s ditolak: dibuat oleh dokter %s (%s), dicoba oleh %s",
			noRawat, existing.KodeDokter, existing.NamaDokter, kodeDokterLogin)
		return apperror.NewForbiddenError(fmt.Sprintf("Anda tidak memiliki hak akses untuk menghapus penilaian medis ini karena dibuat oleh dokter lain (%s)", existing.NamaDokter))
	}

	if err := s.validasiWaktuPenilaianMedis(ctx, noRawat, "", "dihapus"); err != nil {
		return err
	}

	if err := s.repo.HapusPenilaianMedisIGD(ctx, noRawat); err != nil {
		s.log.Error("Gagal menghapus penilaian medis IGD no_rawat %s oleh dokter %s: %v", noRawat, kodeDokterLogin, err)
		return err
	}

	s.log.Info("Berhasil menghapus penilaian medis IGD no_rawat %s oleh dokter %s", noRawat, kodeDokterLogin)
	return nil
}

func (s *service) validasiWaktuPenilaianMedis(ctx context.Context, noRawat, tglPeriksa, aksi string) error {
	tglRegStr, jamRegStr, exists, err := s.rawatJalanService.GetWaktuRegistrasi(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal mengambil waktu registrasi no_rawat %s: %v", noRawat, err)
		return err
	}
	if !exists {
		return apperror.NewNotFoundError("Data kunjungan pasien tidak ditemukan")
	}

	if err := shared.ValidasiBatasWaktuRekamMedis(tglRegStr, jamRegStr, 48, aksi); err != nil {
		s.log.Warn("Validasi batas waktu 48 jam penilaian medis gagal untuk no_rawat %s (%s %s): %v", noRawat, tglRegStr, jamRegStr, err)
		return err
	}

	if tglPeriksa != "" {
		waktuRegistrasi, errReg := shared.ParseWaktu(tglRegStr, jamRegStr)
		waktuPemeriksaan, errPer := time.ParseInLocation("2006-01-02 15:04:05", tglPeriksa, time.Local)
		if errReg == nil && errPer == nil && waktuPemeriksaan.Before(waktuRegistrasi) {
			errs := apperror.ValidationError{
				"tanggal_penilaian": fmt.Sprintf("Waktu penilaian medis (%s) tidak boleh lebih awal dari waktu registrasi pasien (%s %s)", tglPeriksa, tglRegStr, jamRegStr),
			}
			s.log.Warn("Validasi waktu penilaian medis gagal untuk no_rawat %s: %+v", noRawat, errs)
			return errs
		}
	}

	return nil
}
