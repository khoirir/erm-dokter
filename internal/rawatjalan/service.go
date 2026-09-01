package rawatjalan

import (
	"context"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type Service interface {
	DaftarAntreanDokter(ctx context.Context, kodeDokter string, filter FilterAntreanDokter) ([]KunjunganRawatJalan, shared.PaginationMeta, error)
	DetailKunjungan(ctx context.Context, noRawat string, kodeDokter string) (*KunjunganRawatJalan, error)
	RiwayatKunjunganPasien(ctx context.Context, noRM string) ([]KunjunganRawatJalan, error)
	GetWaktuRegistrasi(ctx context.Context, noRawat string) (tanggal string, jam string, exists bool, err error)

	DaftarStatusPemeriksaan(ctx context.Context) []OpsiReferensi
	DaftarStatusLanjut(ctx context.Context) []OpsiReferensi
	DaftarStatusBayar(ctx context.Context) []OpsiReferensi
	DaftarJenisAntrean(ctx context.Context) []OpsiReferensi
}

type service struct {
	repo Repository
	log  *logger.Logger
}

func NewService(repo Repository, log *logger.Logger) Service {
	return &service{
		repo: repo,
		log:  log,
	}
}

func (s *service) DaftarAntreanDokter(ctx context.Context, kodeDokter string, filter FilterAntreanDokter) ([]KunjunganRawatJalan, shared.PaginationMeta, error) {
	daftarAntrean, totalData, err := s.repo.DaftarAntreanDokter(ctx, kodeDokter, filter)
	if err != nil {
		s.log.Error("Gagal query antrean dokter %s: %v", kodeDokter, err)
		return nil, shared.PaginationMeta{}, err
	}

	return daftarAntrean, shared.NewPaginationMeta(totalData, filter.Page, filter.Limit), nil
}

func (s *service) DetailKunjungan(ctx context.Context, noRawat string, kodeDokter string) (*KunjunganRawatJalan, error) {
	if noRawat == "" {
		return nil, apperror.NewBusinessError("nomor rawat tidak boleh kosong")
	}

	kunjungan, err := s.repo.DetailKunjungan(ctx, noRawat, kodeDokter)
	if err != nil {
		s.log.Error("Gagal query detail kunjungan %s oleh dokter %s: %v", noRawat, kodeDokter, err)
		return nil, err
	}

	return kunjungan, nil
}

func (s *service) RiwayatKunjunganPasien(ctx context.Context, noRM string) ([]KunjunganRawatJalan, error) {
	if noRM == "" {
		return nil, apperror.NewBusinessError("nomor rekam medis tidak boleh kosong")
	}
	listKunjungan, err := s.repo.RiwayatKunjunganPasien(ctx, noRM)
	if err != nil {
		s.log.Error("Gagal query riwayat kunjungan pasien %s: %v", noRM, err)
		return nil, err
	}
	return listKunjungan, nil
}

func (s *service) GetWaktuRegistrasi(ctx context.Context, noRawat string) (string, string, bool, error) {
	if noRawat == "" {
		return "", "", false, apperror.NewBusinessError("nomor rawat tidak boleh kosong")
	}
	tglReg, jamReg, exists, err := s.repo.GetWaktuRegistrasi(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal query waktu registrasi %s: %v", noRawat, err)
		return "", "", false, err
	}
	return tglReg, jamReg, exists, nil
}

func (s *service) DaftarStatusPemeriksaan(ctx context.Context) []OpsiReferensi {
	opsi := make([]OpsiReferensi, len(ListStatusPemeriksaan))
	for i, item := range ListStatusPemeriksaan {
		opsi[i] = OpsiReferensi{
			Value: string(item.Value),
			Label: item.Label,
		}
	}
	return opsi
}

func (s *service) DaftarStatusLanjut(ctx context.Context) []OpsiReferensi {
	return []OpsiReferensi{
		{Value: string(shared.StatusLanjutRawatJalan), Label: "Rawat Jalan"},
		{Value: string(shared.StatusLanjutRawatInap), Label: "Rawat Inap"},
	}
}

func (s *service) DaftarStatusBayar(ctx context.Context) []OpsiReferensi {
	opsi := make([]OpsiReferensi, len(ListStatusBayar))
	for i, item := range ListStatusBayar {
		opsi[i] = OpsiReferensi{
			Value: string(item.Value),
			Label: item.Label,
		}
	}
	return opsi
}

func (s *service) DaftarJenisAntrean(ctx context.Context) []OpsiReferensi {
	opsi := make([]OpsiReferensi, len(ListJenisAntrean))
	for i, item := range ListJenisAntrean {
		opsi[i] = OpsiReferensi{
			Value: string(item.Value),
			Label: item.Label,
		}
	}
	return opsi
}

