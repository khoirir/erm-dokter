package usecase

import (
	"context"
	"math"

	"erm-dokter/internal/domain"
	"erm-dokter/internal/dto"
	"erm-dokter/pkg/logger"
)

type rawatJalanUsecase struct {
	rawatJalanRepo domain.RawatJalanRepository
	Log            *logger.Logger
}

func NewRawatJalanUsecase(repo domain.RawatJalanRepository, log *logger.Logger) domain.RawatJalanUsecase {
	return &rawatJalanUsecase{
		rawatJalanRepo: repo,
		Log:            log,
	}
}

func (u *rawatJalanUsecase) DaftarAntreanDokter(ctx context.Context, filter dto.FilterAntreanDokter) ([]domain.KunjunganRawatJalan, dto.MetaPaginasi, error) {
	if errs := filter.Validate(); errs != nil {
		u.Log.Warn("Filter validasi gagal: %+v", errs)
		return nil, dto.MetaPaginasi{}, errs
	}

	if len(filter.KataKunci) > 0 && len(filter.KataKunci) < 3 {
		return nil, dto.MetaPaginasi{}, domain.NewBusinessError("kata kunci pencarian minimal 3 karakter")
	}

	daftarAntrean, totalData, err := u.rawatJalanRepo.DaftarAntreanDokter(ctx, filter)
	if err != nil {
		u.Log.Error("Gagal query antrean dokter %s: %v", filter.KodeDokter, err)
		return nil, dto.MetaPaginasi{}, err
	}

	totalHalaman := int(math.Ceil(float64(totalData) / float64(filter.Batas)))

	meta := dto.MetaPaginasi{
		TotalData:    totalData,
		TotalHalaman: totalHalaman,
		HalamanAktif: filter.Halaman,
		BatasData:    filter.Batas,
	}

	return daftarAntrean, meta, nil
}

func (u *rawatJalanUsecase) DetailKunjungan(ctx context.Context, noRawat string, kodeDokter string) (*domain.KunjunganRawatJalan, error) {
	if noRawat == "" {
		return nil, domain.NewBusinessError("nomor rawat tidak boleh kosong")
	}

	kunjungan, err := u.rawatJalanRepo.DetailKunjungan(ctx, noRawat, kodeDokter)
	if err != nil {
		u.Log.Error("Gagal query detail kunjungan %s: %v", noRawat, err)
		return nil, err
	}

	return kunjungan, nil
}

func (u *rawatJalanUsecase) GetReferensiFilter(ctx context.Context) dto.ReferensiFilterRawatJalan {
	return dto.ReferensiFilterRawatJalan{
		StatusPemeriksaan: []dto.OpsiReferensi{
			{Value: string(domain.StatusBelum), Label: "Belum Periksa"},
			{Value: string(domain.StatusSudah), Label: "Sudah Periksa"},
			{Value: string(domain.StatusBatal), Label: "Batal Periksa"},
			{Value: string(domain.StatusBerkasDiterima), Label: "Berkas Diterima"},
			{Value: string(domain.StatusDirujuk), Label: "Dirujuk"},
			{Value: string(domain.StatusMeninggal), Label: "Meninggal"},
			{Value: string(domain.StatusDirawat), Label: "Dirawat"},
			{Value: string(domain.StatusPulangPaksa), Label: "Pulang Paksa"},
		},
		StatusLanjut: []dto.OpsiReferensi{
			{Value: string(domain.StatusLanjutRawatJalan), Label: "Rawat Jalan"},
			{Value: string(domain.StatusLanjutRawatInap), Label: "Rawat Inap"},
		},
		StatusBayar: []dto.OpsiReferensi{
			{Value: string(domain.StatusBayarSudah), Label: "Sudah Bayar"},
			{Value: string(domain.StatusBayarBelum), Label: "Belum Bayar"},
		},
		JenisAntrean: []dto.OpsiReferensi{
			{Value: string(domain.JenisAntreanRujukan), Label: "Rujukan"},
			{Value: string(domain.JenisAntreanTidakRujukan), Label: "Bukan Rujukan"},
		},
	}
}
