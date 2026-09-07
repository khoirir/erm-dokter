package resumepasien

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

func (s *service) validasiRegistrasiDanStatusRalan(ctx context.Context, noRawat string) error {
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

	_, hasRecordKamar, err := s.rawatInapService.CekStatusKamarInap(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal cek status kamar inap untuk no_rawat %s: %v", noRawat, err)
		return err
	}
	if hasRecordKamar {
		return apperror.NewBusinessError("Pasien sudah terdaftar di rawat inap. Resume medis rawat jalan hanya untuk kunjungan rawat jalan.")
	}

	batasWaktu := waktuRegistrasi.Add(time.Duration(s.maxEditJam) * time.Hour)
	if time.Now().After(batasWaktu) {
		errMsg := fmt.Sprintf("Batas waktu resume pasien untuk kunjungan rawat jalan ini telah berakhir (maksimal %d jam dari waktu registrasi: %s %s)", s.maxEditJam, tanggalRegistrasiStr, jamRegistrasiStr)
		s.log.Warn("Resume pasien ditolak karena lewat batas %d jam untuk no_rawat %s: %s", s.maxEditJam, noRawat, errMsg)
		return apperror.NewForbiddenError(errMsg)
	}

	return nil
}

func (s *service) DetailResumePasienRalan(ctx context.Context, noRawat string) (*ResumePasienRalan, error) {
	item, err := s.repo.DetailResumePasienRalan(ctx, noRawat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFoundError("Data resume pasien rawat jalan tidak ditemukan")
		}
		s.log.Error("Gagal mengambil detail resume pasien ralan no_rawat %s: %v", noRawat, err)
		return nil, err
	}
	return item, nil
}

func (s *service) RiwayatResumePasienRalanByNoRM(ctx context.Context, noRM string) ([]ResumePasienRalan, error) {
	list, err := s.repo.RiwayatResumePasienRalanByNoRM(ctx, noRM)
	if err != nil {
		s.log.Error("Gagal mengambil riwayat resume pasien ralan no_rkm_medis %s: %v", noRM, err)
		return nil, err
	}
	return list, nil
}

func (s *service) SimpanResumePasienRalan(ctx context.Context, noRawat, kodeDokter string, req SimpanResumePasienRalanRequest) (*ResumePasienRalan, error) {
	if err := s.validasiRegistrasiDanStatusRalan(ctx, noRawat); err != nil {
		return nil, err
	}

	ada, err := s.repo.CekResumePasienRalanAda(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal cek keberadaan resume pasien ralan no_rawat %s: %v", noRawat, err)
		return nil, err
	}
	if ada {
		return nil, apperror.NewBusinessError("Resume medis rawat jalan untuk kunjungan ini sudah ada")
	}

	item, err := s.repo.SimpanResumePasienRalan(ctx, noRawat, kodeDokter, req)
	if err != nil {
		s.log.Error("Gagal menyimpan resume pasien ralan no_rawat %s: %v", noRawat, err)
		return nil, err
	}

	s.log.Info("Berhasil menyimpan resume pasien ralan no_rawat %s oleh dokter %s", noRawat, kodeDokter)
	return item, nil
}

func (s *service) UpdateResumePasienRalan(ctx context.Context, kodeDokterLogin, noRawat string, req UpdateResumePasienRalanRequest) (*ResumePasienRalan, error) {
	existing, err := s.repo.DetailResumePasienRalan(ctx, noRawat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFoundError("Data resume pasien rawat jalan tidak ditemukan")
		}
		s.log.Error("Gagal mengambil resume pasien ralan no_rawat %s untuk update: %v", noRawat, err)
		return nil, err
	}

	if existing.KodeDokter != kodeDokterLogin {
		s.log.Warn("Dokter %s mencoba mengubah resume pasien ralan milik dokter %s (no_rawat: %s)", kodeDokterLogin, existing.KodeDokter, noRawat)
		return nil, apperror.NewForbiddenError("Anda tidak memiliki akses untuk mengubah resume pasien milik dokter lain")
	}

	if err := s.validasiRegistrasiDanStatusRalan(ctx, noRawat); err != nil {
		return nil, err
	}

	item, err := s.repo.UpdateResumePasienRalan(ctx, noRawat, req)
	if err != nil {
		s.log.Error("Gagal memperbarui resume pasien ralan no_rawat %s: %v", noRawat, err)
		return nil, err
	}

	s.log.Info("Berhasil memperbarui resume pasien ralan no_rawat %s oleh dokter %s", noRawat, kodeDokterLogin)
	return item, nil
}

func (s *service) HapusResumePasienRalan(ctx context.Context, kodeDokterLogin, noRawat string) error {
	existing, err := s.repo.DetailResumePasienRalan(ctx, noRawat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NewNotFoundError("Data resume pasien rawat jalan tidak ditemukan")
		}
		s.log.Error("Gagal mengambil resume pasien ralan no_rawat %s untuk hapus: %v", noRawat, err)
		return err
	}

	if existing.KodeDokter != kodeDokterLogin {
		s.log.Warn("Dokter %s mencoba menghapus resume pasien ralan milik dokter %s (no_rawat: %s)", kodeDokterLogin, existing.KodeDokter, noRawat)
		return apperror.NewForbiddenError("Anda tidak memiliki akses untuk menghapus resume pasien milik dokter lain")
	}

	if err := s.validasiRegistrasiDanStatusRalan(ctx, noRawat); err != nil {
		return err
	}

	if err := s.repo.HapusResumePasienRalan(ctx, noRawat); err != nil {
		s.log.Error("Gagal menghapus resume pasien ralan no_rawat %s: %v", noRawat, err)
		return err
	}

	s.log.Info("Berhasil menghapus resume pasien ralan no_rawat %s oleh dokter %s", noRawat, kodeDokterLogin)
	return nil
}

func (s *service) ReferensiRalan(ctx context.Context) ReferensiResumeRalan {
	return ReferensiResumeRalan{
		KondisiPulang: DaftarOpsiKondisiPulang(),
	}
}
