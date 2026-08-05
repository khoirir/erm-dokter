package domain

import (
	"context"
)

type Pasien struct {
	NoRM         string `json:"no_rkm_medis"`
	Nama         string `json:"nm_pasien"`
	NIK          string `json:"no_ktp"`
	NoBPJS       string `json:"no_peserta"`
	JenisKelamin string `json:"jk"`
	TanggalLahir string `json:"tgl_lahir"`
	TempatLahir  string `json:"tmp_lahir"`
	Alamat       string `json:"alamat"`
	NoTelp       string `json:"no_tlp"`
}

func (p *Pasien) FormatNoRekamMedis() string {
	return FormatNoRekamMedis(p.NoRM)
}

func (p *Pasien) HitungUmurLengkap() Umur {
	return HitungUmur(p.TanggalLahir)
}

func (p *Pasien) FormatUmur() string {
	return FormatUmur(p.TanggalLahir)
}

func (p *Pasien) FormatJenisKelamin() string {
	return FormatJenisKelamin(p.JenisKelamin)
}

type PasienRepository interface {
	CariByNoRM(ctx context.Context, noRM string) (*Pasien, error)
	CariByNIK(ctx context.Context, nik string) (*Pasien, error)
	CariPasien(ctx context.Context, kataKunci string, limit int) ([]Pasien, error)
}

type PasienUsecase interface {
	DetailPasien(ctx context.Context, noRM string) (*Pasien, error)
	CariPasien(ctx context.Context, kataKunci string) ([]Pasien, error)
}
