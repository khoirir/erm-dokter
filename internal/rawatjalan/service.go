package rawatjalan

import (
	"context"
	"math"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type Service interface {
	DaftarAntreanDokter(ctx context.Context, kodeDokter string, filter FilterAntreanDokter) ([]KunjunganRawatJalan, shared.MetaPaginasi, error)
	DetailKunjungan(ctx context.Context, noRawat string, kodeDokter string) (*KunjunganRawatJalan, error)
	RiwayatKunjunganPasien(ctx context.Context, noRM string) ([]KunjunganRawatJalan, error)
	GetWaktuRegistrasi(ctx context.Context, noRawat string) (tanggal string, jam string, exists bool, err error)
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

func (s *service) DaftarAntreanDokter(ctx context.Context, kodeDokter string, filter FilterAntreanDokter) ([]KunjunganRawatJalan, shared.MetaPaginasi, error) {
	if errs := filter.Validate(); errs != nil {
		s.log.Warn("Filter validasi gagal: %+v", errs)
		return nil, shared.MetaPaginasi{}, errs
	}

	if len(filter.KataKunci) > 0 && len(filter.KataKunci) < 3 {
		return nil, shared.MetaPaginasi{}, apperror.NewBusinessError("kata kunci pencarian minimal 3 karakter")
	}

	daftarAntrean, totalData, err := s.repo.DaftarAntreanDokter(ctx, kodeDokter, filter)
	if err != nil {
		s.log.Error("Gagal query antrean dokter %s: %v", kodeDokter, err)
		return nil, shared.MetaPaginasi{}, err
	}

	totalHalaman := int(math.Ceil(float64(totalData) / float64(filter.Batas)))

	meta := shared.MetaPaginasi{
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
			{Value: string(StatusBelum), Label: "Belum Periksa"},
			{Value: string(StatusSudah), Label: "Sudah Periksa"},
			{Value: string(StatusBatal), Label: "Batal Periksa"},
			{Value: string(StatusBerkasDiterima), Label: "Berkas Diterima"},
			{Value: string(StatusDirujuk), Label: "Dirujuk"},
			{Value: string(StatusMeninggal), Label: "Meninggal"},
			{Value: string(StatusDirawat), Label: "Dirawat"},
			{Value: string(StatusPulangPaksa), Label: "Pulang Paksa"},
		},
		StatusLanjut: []OpsiReferensi{
			{Value: string(shared.StatusLanjutRawatJalan), Label: "Rawat Jalan"},
			{Value: string(shared.StatusLanjutRawatInap), Label: "Rawat Inap"},
		},
		StatusBayar: []OpsiReferensi{
			{Value: string(StatusBayarSudah), Label: "Sudah Bayar"},
			{Value: string(StatusBayarBelum), Label: "Belum Bayar"},
		},
		JenisAntrean: []OpsiReferensi{
			{Value: string(JenisAntreanRujukan), Label: "Rujukan"},
			{Value: string(JenisAntreanTidakRujukan), Label: "Bukan Rujukan"},
		},
	}
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

