package pemeriksaan

import (
	"erm-dokter/internal/shared"
)

type Pemeriksaan struct {
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
