package penilaianmedis

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/shared/apperror"
)

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
