package domain

import (
	"context"
)

type StatusPemeriksaan string

const (
	StatusBelum          StatusPemeriksaan = "Belum"
	StatusSudah          StatusPemeriksaan = "Sudah"
	StatusBatal          StatusPemeriksaan = "Batal"
	StatusBerkasDiterima StatusPemeriksaan = "Berkas Diterima"
	StatusDirujuk        StatusPemeriksaan = "Dirujuk"
	StatusMeninggal      StatusPemeriksaan = "Meninggal"
	StatusDirawat        StatusPemeriksaan = "Dirawat"
	StatusPulangPaksa    StatusPemeriksaan = "Pulang Paksa"
)

type StatusLanjut string

const (
	StatusLanjutRawatJalan StatusLanjut = "Ralan"
	StatusLanjutRawatInap  StatusLanjut = "Ranap"
)

type StatusBayar string

const (
	StatusBayarSudah StatusBayar = "Sudah Bayar"
	StatusBayarBelum StatusBayar = "Belum Bayar"
)

type JenisAntrean string

const (
	JenisAntreanRujukan      JenisAntrean = "Rujukan"
	JenisAntreanTidakRujukan JenisAntrean = "Bukan Rujukan"
)

type FilterAntreanDokter struct {
	KodeDokter        string            `json:"kode_dokter"`
	Tanggal           string            `json:"tanggal"`
	KodePenjamin      string            `json:"kode_penjamin"`
	StatusPemeriksaan StatusPemeriksaan `json:"status_pemeriksaan"`
	JenisAntrean      JenisAntrean      `json:"jenis_antrean"`
	KataKunci         string            `json:"keyword"`
	OrderBy           string            `json:"order_by"`
	SortOrder         string            `json:"sort_order"`
	Halaman           int               `json:"halaman"`
	Batas             int               `json:"batas"`
}

type KunjunganRawatJalan struct {
	NoRawat           string            `json:"no_rawat"`
	NoRegistrasi      string            `json:"no_registrasi"`
	TanggalRegistrasi string            `json:"tanggal_registrasi"`
	JamRegistrasi     string            `json:"jam_registrasi"`
	NoRekamMedis      string            `json:"no_rekam_medis"`
	NamaPasien        string            `json:"nama_pasien"`
	JenisKelamin      string            `json:"jenis_kelamin"`
	TanggalLahir      string            `json:"tanggal_lahir"`
	Umur              string            `json:"umur"`
	Alamat            string            `json:"alamat"`
	KodePoli          string            `json:"kode_poli"`
	NamaPoli          string            `json:"nama_poli"`
	KodeDokter        string            `json:"kode_dokter"`
	NamaDokter        string            `json:"nama_dokter"`
	KodePenjamin      string            `json:"kode_penjamin"`
	NamaPenjamin      string            `json:"nama_penjamin"`
	StatusPemeriksaan StatusPemeriksaan `json:"status_pemeriksaan"`
	StatusLanjut      StatusLanjut      `json:"status_lanjut"`
	StatusBayar       StatusBayar       `json:"status_bayar"`
	JenisAntrean      JenisAntrean      `json:"jenis_antrean"`
}

func (k *KunjunganRawatJalan) FormatJenisKelamin() string {
	return FormatJenisKelamin(k.JenisKelamin)
}

func (k *KunjunganRawatJalan) FormatUmur() string {
	return FormatUmur(k.TanggalLahir)
}

func (k *KunjunganRawatJalan) FormatNoRekamMedis() string {
	return FormatNoRekamMedis(k.NoRekamMedis)
}

type RawatJalanRepository interface {
	DaftarAntreanDokter(ctx context.Context, filter FilterAntreanDokter) ([]KunjunganRawatJalan, int, error)
	DetailKunjungan(ctx context.Context, noRawat string) (*KunjunganRawatJalan, error)
}

type RawatJalanUsecase interface {
	DaftarAntreanDokter(ctx context.Context, filter FilterAntreanDokter) ([]KunjunganRawatJalan, MetaPaginasi, error)
	DetailKunjungan(ctx context.Context, noRawat string) (*KunjunganRawatJalan, error)
}
