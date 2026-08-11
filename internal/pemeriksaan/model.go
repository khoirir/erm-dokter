package pemeriksaan

import (
	"context"

	"erm-dokter/internal/shared"
)

type Pemeriksaan struct {
	NoRawat             string               `json:"no_rawat"`
	TanggalPemeriksaan  string               `json:"tanggal_pemeriksaan"`
	JamPemeriksaan      string               `json:"jam_pemeriksaan"`
	SuhuTubuh           string               `json:"suhu_tubuh"`
	Tensi               string               `json:"tensi"`
	Nadi                string               `json:"nadi"`
	Respirasi           string               `json:"respirasi"`
	TinggiBadan         string               `json:"tinggi_badan"`
	BeratBadan          string               `json:"berat_badan"`
	SpO2                string               `json:"spo2"`
	Gcs                 string               `json:"gcs"`
	Kesadaran           shared.Kesadaran     `json:"kesadaran"`
	Keluhan             string               `json:"keluhan"`
	Pemeriksaan         string               `json:"pemeriksaan"`
	Alergi              string               `json:"alergi"`
	LingkarPerut        string               `json:"lingkar_perut,omitempty"`
	RencanaTindakLanjut string               `json:"rencana_tindak_lanjut"`
	Penilaian           string               `json:"penilaian"`
	Instruksi           string               `json:"instruksi"`
	Evaluasi            string               `json:"evaluasi"`
	KodeDokterPetugas   string               `json:"kode_dokter_petugas"`
	NamaDokterPetugas   string               `json:"nama_dokter_petugas"`
	StatusLanjut        shared.StatusLanjut  `json:"status_lanjut"`
}

type Repository interface {
	DaftarPemeriksaan(ctx context.Context) ([]Pemeriksaan, error)
}

type Service interface {
	DaftarPemeriksaan(ctx context.Context) ([]Pemeriksaan, error)
}
