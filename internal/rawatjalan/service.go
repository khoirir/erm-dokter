package rawatjalan

import (
	"context"
	"math"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type Service interface {
	DaftarAntreanDokter(ctx context.Context, filter FilterAntreanDokter) ([]KunjunganRawatJalan, MetaPaginasi, error)
	DetailKunjungan(ctx context.Context, noRawat string, kodeDokter string) (*KunjunganRawatJalan, error)
	GetReferensiFilter(ctx context.Context) ReferensiFilterRawatJalan
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

func (s *service) DaftarAntreanDokter(ctx context.Context, filter FilterAntreanDokter) ([]KunjunganRawatJalan, MetaPaginasi, error) {
	if errs := filter.Validate(); errs != nil {
		s.log.Warn("Filter validasi gagal: %+v", errs)
		return nil, MetaPaginasi{}, errs
	}

	if len(filter.KataKunci) > 0 && len(filter.KataKunci) < 3 {
		return nil, MetaPaginasi{}, apperror.NewBusinessError("kata kunci pencarian minimal 3 karakter")
	}

	daftarAntrean, totalData, err := s.repo.DaftarAntreanDokter(ctx, filter)
	if err != nil {
		s.log.Error("Gagal query antrean dokter %s: %v", filter.KodeDokter, err)
		return nil, MetaPaginasi{}, err
	}

	totalHalaman := int(math.Ceil(float64(totalData) / float64(filter.Batas)))

	meta := MetaPaginasi{
		TotalData:    totalData,
		TotalHalaman: totalHalaman,
		HalamanAktif: filter.Halaman,
		BatasData:    filter.Batas,
	}

	return daftarAntrean, meta, nil
}

func (s *service) DetailKunjungan(ctx context.Context, noRawat string, kodeDokter string) (*KunjunganRawatJalan, error) {
	if noRawat == "" {
		return nil, apperror.NewBusinessError("nomor rawat tidak boleh kosong")
	}

	kunjungan, err := s.repo.DetailKunjungan(ctx, noRawat, kodeDokter)
	if err != nil {
		s.log.Error("Gagal query detail kunjungan %s: %v", noRawat, err)
		return nil, err
	}

	return kunjungan, nil
}

func (s *service) GetReferensiFilter(ctx context.Context) ReferensiFilterRawatJalan {
	return ReferensiFilterRawatJalan{
		StatusPemeriksaan: []OpsiReferensi{
			{Value: string(shared.StatusBelum), Label: "Belum Periksa"},
			{Value: string(shared.StatusSudah), Label: "Sudah Periksa"},
			{Value: string(shared.StatusBatal), Label: "Batal Periksa"},
			{Value: string(shared.StatusBerkasDiterima), Label: "Berkas Diterima"},
			{Value: string(shared.StatusDirujuk), Label: "Dirujuk"},
			{Value: string(shared.StatusMeninggal), Label: "Meninggal"},
			{Value: string(shared.StatusDirawat), Label: "Dirawat"},
			{Value: string(shared.StatusPulangPaksa), Label: "Pulang Paksa"},
		},
		StatusLanjut: []OpsiReferensi{
			{Value: string(shared.StatusLanjutRawatJalan), Label: "Rawat Jalan"},
			{Value: string(shared.StatusLanjutRawatInap), Label: "Rawat Inap"},
		},
		StatusBayar: []OpsiReferensi{
			{Value: string(shared.StatusBayarSudah), Label: "Sudah Bayar"},
			{Value: string(shared.StatusBayarBelum), Label: "Belum Bayar"},
		},
		JenisAntrean: []OpsiReferensi{
			{Value: string(shared.JenisAntreanRujukan), Label: "Rujukan"},
			{Value: string(shared.JenisAntreanTidakRujukan), Label: "Bukan Rujukan"},
		},
	}
}
