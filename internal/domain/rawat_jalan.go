package domain

import (
	"context"

	"erm-dokter/internal/dto"
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

type KunjunganRawatJalan struct {
	NoRawat           string            `json:"no_rawat"`
	DetailURL         string            `json:"detail_url,omitempty"`
	NoRegistrasi      string            `json:"no_registrasi"`
	TanggalRegistrasi string            `json:"tanggal_registrasi"`
	JamRegistrasi     string            `json:"jam_registrasi"`
	NoRekamMedis      string            `json:"no_rekam_medis"`
	NamaPasien        string            `json:"nama_pasien"`
	JenisKelamin      string            `json:"jenis_kelamin"`
	TanggalLahir      string            `json:"tanggal_lahir"`
	Umur              string            `json:"umur"`
	Alamat            string            `json:"alamat"`
	KodePoliAsal      string            `json:"kode_poli_asal"`
	NamaPoliAsal      string            `json:"nama_poli_asal"`
	KodeDokterAsal    string            `json:"kode_dokter_asal"`
	NamaDokterAsal    string            `json:"nama_dokter_asal"`
	KodePoliRujukan   string            `json:"kode_poli_rujukan,omitempty"`
	NamaPoliRujukan   string            `json:"nama_poli_rujukan,omitempty"`
	KodeDokterRujukan string            `json:"kode_dokter_rujukan,omitempty"`
	NamaDokterRujukan string            `json:"nama_dokter_rujukan,omitempty"`
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
	DaftarAntreanDokter(ctx context.Context, filter dto.FilterAntreanDokter) ([]KunjunganRawatJalan, int, error)
	DetailKunjungan(ctx context.Context, noRawat string, kodeDokter string) (*KunjunganRawatJalan, error)
}

type RawatJalanUsecase interface {
	DaftarAntreanDokter(ctx context.Context, filter dto.FilterAntreanDokter) ([]KunjunganRawatJalan, dto.MetaPaginasi, error)
	DetailKunjungan(ctx context.Context, noRawat string, kodeDokter string) (*KunjunganRawatJalan, error)
	GetReferensiFilter(ctx context.Context) dto.ReferensiFilterRawatJalan
}

