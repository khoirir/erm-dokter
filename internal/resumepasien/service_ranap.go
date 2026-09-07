package resumepasien

import (
	"context"
	"database/sql"
	"errors"

	"erm-dokter/internal/shared/apperror"
)

func (s *service) validasiRegistrasiDanStatusRanap(ctx context.Context, noRawat string) error {
	_, _, exists, err := s.rawatJalanService.GetWaktuRegistrasi(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal mengambil waktu registrasi no_rawat %s: %v", noRawat, err)
		return err
	}
	if !exists {
		return apperror.NewNotFoundError("Data kunjungan pasien tidak ditemukan")
	}

	_, hasRecordKamar, err := s.rawatInapService.CekStatusKamarInap(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal cek status kamar inap untuk no_rawat %s: %v", noRawat, err)
		return err
	}

	if !hasRecordKamar {
		return apperror.NewBusinessError("Pasien belum/tidak terdaftar di rawat inap. Resume medis rawat inap hanya untuk pasien rawat inap.")
	}

	return nil
}

func (s *service) DetailResumePasienRanap(ctx context.Context, noRawat string) (*ResumePasienRanap, error) {
	item, err := s.repo.DetailResumePasienRanap(ctx, noRawat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFoundError("Data resume pasien rawat inap tidak ditemukan")
		}
		s.log.Error("Gagal mengambil detail resume pasien ranap no_rawat %s: %v", noRawat, err)
		return nil, err
	}
	return item, nil
}

func (s *service) RiwayatResumePasienRanapByNoRM(ctx context.Context, noRM string) ([]ResumePasienRanap, error) {
	list, err := s.repo.RiwayatResumePasienRanapByNoRM(ctx, noRM)
	if err != nil {
		s.log.Error("Gagal mengambil riwayat resume pasien ranap no_rkm_medis %s: %v", noRM, err)
		return nil, err
	}
	return list, nil
}

func (s *service) SimpanResumePasienRanap(ctx context.Context, noRawat, kodeDokter string, req SimpanResumePasienRanapRequest) (*ResumePasienRanap, error) {
	if err := s.validasiRegistrasiDanStatusRanap(ctx, noRawat); err != nil {
		return nil, err
	}

	ada, err := s.repo.CekResumePasienRanapAda(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal cek keberadaan resume pasien ranap no_rawat %s: %v", noRawat, err)
		return nil, err
	}
	if ada {
		return nil, apperror.NewBusinessError("Resume medis rawat inap untuk kunjungan ini sudah ada")
	}

	item, err := s.repo.SimpanResumePasienRanap(ctx, noRawat, kodeDokter, req)
	if err != nil {
		s.log.Error("Gagal menyimpan resume pasien ranap no_rawat %s: %v", noRawat, err)
		return nil, err
	}

	s.log.Info("Berhasil menyimpan resume pasien ranap no_rawat %s oleh dokter %s", noRawat, kodeDokter)
	return item, nil
}

func (s *service) UpdateResumePasienRanap(ctx context.Context, kodeDokterLogin, noRawat string, req UpdateResumePasienRanapRequest) (*ResumePasienRanap, error) {
	existing, err := s.repo.DetailResumePasienRanap(ctx, noRawat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFoundError("Data resume pasien rawat inap tidak ditemukan")
		}
		s.log.Error("Gagal mengambil resume pasien ranap no_rawat %s untuk update: %v", noRawat, err)
		return nil, err
	}

	if existing.KodeDokter != kodeDokterLogin {
		s.log.Warn("Dokter %s mencoba mengubah resume pasien ranap milik dokter %s (no_rawat: %s)", kodeDokterLogin, existing.KodeDokter, noRawat)
		return nil, apperror.NewForbiddenError("Anda tidak memiliki akses untuk mengubah resume pasien milik dokter lain")
	}

	if err := s.validasiRegistrasiDanStatusRanap(ctx, noRawat); err != nil {
		return nil, err
	}

	item, err := s.repo.UpdateResumePasienRanap(ctx, noRawat, req)
	if err != nil {
		s.log.Error("Gagal memperbarui resume pasien ranap no_rawat %s: %v", noRawat, err)
		return nil, err
	}

	s.log.Info("Berhasil memperbarui resume pasien ranap no_rawat %s oleh dokter %s", noRawat, kodeDokterLogin)
	return item, nil
}

func (s *service) HapusResumePasienRanap(ctx context.Context, kodeDokterLogin, noRawat string) error {
	existing, err := s.repo.DetailResumePasienRanap(ctx, noRawat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NewNotFoundError("Data resume pasien rawat inap tidak ditemukan")
		}
		s.log.Error("Gagal mengambil resume pasien ranap no_rawat %s untuk hapus: %v", noRawat, err)
		return err
	}

	if existing.KodeDokter != kodeDokterLogin {
		s.log.Warn("Dokter %s mencoba menghapus resume pasien ranap milik dokter %s (no_rawat: %s)", kodeDokterLogin, existing.KodeDokter, noRawat)
		return apperror.NewForbiddenError("Anda tidak memiliki akses untuk menghapus resume pasien milik dokter lain")
	}

	if err := s.validasiRegistrasiDanStatusRanap(ctx, noRawat); err != nil {
		return err
	}

	if err := s.repo.HapusResumePasienRanap(ctx, noRawat); err != nil {
		s.log.Error("Gagal menghapus resume pasien ranap no_rawat %s: %v", noRawat, err)
		return err
	}

	s.log.Info("Berhasil menghapus resume pasien ranap no_rawat %s oleh dokter %s", noRawat, kodeDokterLogin)
	return nil
}

func (s *service) ReferensiRanap(ctx context.Context) ReferensiResumeRanap {
	return ReferensiResumeRanap{
		CaraKeluar:    DaftarOpsiCaraKeluar(),
		KeadaanPulang: DaftarOpsiKeadaanPulang(),
		Dilanjutkan:   DaftarOpsiDilanjutkan(),
	}
}

