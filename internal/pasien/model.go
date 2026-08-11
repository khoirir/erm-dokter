package pasien

import (
	"context"

	"erm-dokter/internal/shared/formatter"
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
	return formatter.FormatNoRekamMedis(p.NoRM)
}

func (p *Pasien) HitungUmurLengkap() formatter.Umur {
	return formatter.HitungUmur(p.TanggalLahir)
}

func (p *Pasien) FormatUmur() string {
	return formatter.FormatUmur(p.TanggalLahir)
}

func (p *Pasien) FormatJenisKelamin() string {
	return formatter.FormatJenisKelamin(p.JenisKelamin)
}

type Repository interface {
	CariByNoRM(ctx context.Context, noRM string) (*Pasien, error)
	CariByNIK(ctx context.Context, nik string) (*Pasien, error)
	CariPasien(ctx context.Context, kataKunci string, limit int) ([]Pasien, error)
}

type Service interface {
	DetailPasien(ctx context.Context, noRM string) (*Pasien, error)
	CariPasien(ctx context.Context, kataKunci string) ([]Pasien, error)
}
