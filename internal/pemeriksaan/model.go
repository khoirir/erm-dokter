package pemeriksaan

import (
	"errors"
	"fmt"
	"strings"
	// "time"

	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type IdPemeriksaan struct {
	NoRawat            string
	TanggalPemeriksaan string
	JamPemeriksaan     string
}

func (id IdPemeriksaan) CompositeKey() string {
	return fmt.Sprintf("%s~%s~%s", id.NoRawat, id.TanggalPemeriksaan, id.JamPemeriksaan)
}

func ParseIdPemeriksaan(decryptedKey string) (IdPemeriksaan, error) {
	parts := strings.Split(decryptedKey, "~")
	if len(parts) != 3 {
		return IdPemeriksaan{}, errors.New("format ID pemeriksaan tidak valid")
	}
	return IdPemeriksaan{
		NoRawat:            parts[0],
		TanggalPemeriksaan: parts[1],
		JamPemeriksaan:     parts[2],
	}, nil
}

func (p *Pemeriksaan) ToIdPemeriksaan() IdPemeriksaan {
	return IdPemeriksaan{
		NoRawat:            p.NoRawat,
		TanggalPemeriksaan: p.TanggalPemeriksaan,
		JamPemeriksaan:     p.JamPemeriksaan,
	}
}

func (p *Pemeriksaan) CompositeKey() string {
	return p.ToIdPemeriksaan().CompositeKey()
}

type Pemeriksaan struct {
	Id                  string              `json:"id"`
	IdKunjungan         string              `json:"id_kunjungan"`
	NoRawat             string              `json:"no_rawat"`
	TanggalPemeriksaan  string              `json:"tanggal_pemeriksaan"`
	JamPemeriksaan      string              `json:"jam_pemeriksaan"`
	SuhuTubuh           string              `json:"suhu_tubuh"`
	Tensi               string              `json:"tensi"`
	Nadi                string              `json:"nadi"`
	Respirasi           string              `json:"respirasi"`
	TinggiBadan         string              `json:"tinggi_badan"`
	BeratBadan          string              `json:"berat_badan"`
	SpO2                string              `json:"spo2"`
	Gcs                 string              `json:"gcs"`
	Kesadaran           Kesadaran           `json:"kesadaran"`
	Keluhan             string              `json:"keluhan"`
	Pemeriksaan         string              `json:"pemeriksaan"`
	Alergi              string              `json:"alergi"`
	LingkarPerut        string              `json:"lingkar_perut,omitempty"`
	RencanaTindakLanjut string              `json:"rencana_tindak_lanjut"`
	Penilaian           string              `json:"penilaian"`
	Instruksi           string              `json:"instruksi"`
	Evaluasi            string              `json:"evaluasi"`
	KodeDokterPetugas   string              `json:"kode_dokter_petugas"`
	NamaDokterPetugas   string              `json:"nama_dokter_petugas"`
	StatusLanjut        shared.StatusLanjut `json:"status_lanjut"`
}

type FilterDaftarPemeriksaan struct {
	Tanggal string `json:"tanggal,omitempty"`
	Halaman int    `json:"halaman,omitempty"`
	Batas   int    `json:"batas,omitempty"`
}

func (f *FilterDaftarPemeriksaan) Sanitize() {
	// if strings.TrimSpace(f.Tanggal) == "" {
	// 	today := time.Now().Format("2006-01-02")
	// 	f.Tanggal = today + "," + today
	// }
	if f.Halaman <= 0 {
		f.Halaman = 1
	}
	if f.Batas <= 0 {
		f.Batas = 20
	} else if f.Batas > 100 {
		f.Batas = 100
	}
}

func (f FilterDaftarPemeriksaan) Offset() int {
	return (f.Halaman - 1) * f.Batas
}

func (f *FilterDaftarPemeriksaan) Validate() apperror.ValidationError {
	f.Sanitize()
	errs := make(apperror.ValidationError)
	if f.Tanggal != "" {
		shared.ValidasiRentangTanggal(f.Tanggal, errs)
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

type SimpanPemeriksaanRequest struct {
	NoRawat             string    `json:"no_rawat" validate:"required,max=17"`
	TanggalPemeriksaan  string    `json:"tanggal_pemeriksaan" validate:"required,datetime=2006-01-02"`
	JamPemeriksaan      string    `json:"jam_pemeriksaan" validate:"required,datetime=15:04:05"`
	SuhuTubuh           string    `json:"suhu_tubuh" validate:"omitempty,max=5"`
	Tensi               string    `json:"tensi" validate:"omitempty,max=8"`
	Nadi                string    `json:"nadi" validate:"omitempty,number,max=3"`
	Respirasi           string    `json:"respirasi" validate:"omitempty,number,max=3"`
	TinggiBadan         string    `json:"tinggi_badan" validate:"omitempty,max=5"`
	BeratBadan          string    `json:"berat_badan" validate:"omitempty,max=5"`
	SpO2                string    `json:"spo2" validate:"omitempty,max=3"`
	Gcs                 string    `json:"gcs" validate:"omitempty,max=10"`
	Kesadaran           Kesadaran `json:"kesadaran" validate:"required"`
	Keluhan             string    `json:"keluhan" validate:"required,max=2000"`
	Pemeriksaan         string    `json:"pemeriksaan" validate:"required,max=2000"`
	Alergi              string    `json:"alergi" validate:"omitempty,max=50"`
	LingkarPerut        string    `json:"lingkar_perut" validate:"omitempty,max=5"`
	RencanaTindakLanjut string    `json:"rencana_tindak_lanjut" validate:"required,max=2000"`
	Penilaian           string    `json:"penilaian" validate:"required,max=2000"`
	Instruksi           string    `json:"instruksi" validate:"required,max=2000"`
	Evaluasi            string    `json:"evaluasi" validate:"required,max=2000"`
}
