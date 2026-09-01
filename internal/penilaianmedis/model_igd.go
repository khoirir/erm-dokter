package penilaianmedis

import (
	"fmt"
	"strings"
	"time"

	"erm-dokter/internal/shared/apperror"
)

type DataPenilaianMedisIGD struct {
	TanggalPenilaian        string    `json:"tanggal_penilaian"`
	Anamnesis               Anamnesis `json:"anamnesis"`
	Hubungan                string    `json:"hubungan"`
	KeluhanUtama            string    `json:"keluhan_utama"`
	RiwayatPenyakitSekarang string    `json:"riwayat_penyakit_sekarang"`
	RiwayatPenyakitDahulu   string    `json:"riwayat_penyakit_dahulu"`
	RiwayatPenyakitKeluarga string    `json:"riwayat_penyakit_keluarga"`
	RiwayatPengobatan       string    `json:"riwayat_pengobatan"`
	Alergi                  string    `json:"alergi"`

	Keadaan     Keadaan       `json:"keadaan"`
	Gcs         string        `json:"gcs"`
	Kesadaran   KesadaranAwal `json:"kesadaran"`
	Tensi       string        `json:"tensi"`
	Nadi        string        `json:"nadi"`
	Respirasi   string        `json:"respirasi"`
	SuhuTubuh   string        `json:"suhu_tubuh"`
	SpO2        string        `json:"spo2"`
	BeratBadan  string        `json:"berat_badan"`
	TinggiBadan string        `json:"tinggi_badan"`

	Kepala          StatusFisik `json:"kepala"`
	Mata            StatusFisik `json:"mata"`
	Gigi            StatusFisik `json:"gigi"`
	Leher           StatusFisik `json:"leher"`
	Thoraks         StatusFisik `json:"thoraks"`
	Abdomen         StatusFisik `json:"abdomen"`
	Genital         StatusFisik `json:"genital"`
	Ekstremitas     StatusFisik `json:"ekstremitas"`
	KeteranganFisik string      `json:"keterangan_fisik"`

	KeteranganLokalis string `json:"keterangan_lokalis"`
	EKG               string `json:"ekg"`
	Radiologi         string `json:"radiologi"`
	Laboratorium      string `json:"laboratorium"`

	Diagnosis   string `json:"diagnosis"`
	TataLaksana string `json:"tata_laksana"`
}

func (d *DataPenilaianMedisIGD) Sanitize() {
	d.TanggalPenilaian = strings.TrimSpace(d.TanggalPenilaian)
	if d.TanggalPenilaian == "" {
		d.TanggalPenilaian = time.Now().Format("2006-01-02 15:04:05")
	} else if len(d.TanggalPenilaian) == 10 {
		d.TanggalPenilaian = fmt.Sprintf("%s %s", d.TanggalPenilaian, time.Now().Format("15:04:05"))
	}

	d.Anamnesis = Anamnesis(strings.TrimSpace(string(d.Anamnesis)))
	if d.Anamnesis == "" {
		d.Anamnesis = Autoanamnesis
	}
	d.Hubungan = strings.TrimSpace(d.Hubungan)
	d.KeluhanUtama = strings.TrimSpace(d.KeluhanUtama)
	d.RiwayatPenyakitSekarang = strings.TrimSpace(d.RiwayatPenyakitSekarang)
	d.RiwayatPenyakitDahulu = strings.TrimSpace(d.RiwayatPenyakitDahulu)
	d.RiwayatPenyakitKeluarga = strings.TrimSpace(d.RiwayatPenyakitKeluarga)
	d.RiwayatPengobatan = strings.TrimSpace(d.RiwayatPengobatan)
	d.Alergi = strings.TrimSpace(d.Alergi)

	d.Keadaan = Keadaan(strings.TrimSpace(string(d.Keadaan)))
	if d.Keadaan == "" {
		d.Keadaan = KeadaanSehat
	}
	d.Gcs = strings.TrimSpace(d.Gcs)
	d.Kesadaran = KesadaranAwal(strings.TrimSpace(string(d.Kesadaran)))
	if d.Kesadaran == "" {
		d.Kesadaran = KesadaranComposMentis
	}
	d.Tensi = strings.ReplaceAll(strings.TrimSpace(d.Tensi), " ", "")
	d.Nadi = strings.TrimSpace(d.Nadi)
	d.Respirasi = strings.TrimSpace(d.Respirasi)
	d.SuhuTubuh = strings.ReplaceAll(strings.TrimSpace(d.SuhuTubuh), ",", ".")
	d.SpO2 = strings.TrimSpace(d.SpO2)
	d.BeratBadan = strings.ReplaceAll(strings.TrimSpace(d.BeratBadan), ",", ".")
	d.TinggiBadan = strings.ReplaceAll(strings.TrimSpace(d.TinggiBadan), ",", ".")

	d.Kepala = StatusFisik(strings.TrimSpace(string(d.Kepala)))
	if d.Kepala == "" {
		d.Kepala = StatusFisikNormal
	}
	d.Mata = StatusFisik(strings.TrimSpace(string(d.Mata)))
	if d.Mata == "" {
		d.Mata = StatusFisikNormal
	}
	d.Gigi = StatusFisik(strings.TrimSpace(string(d.Gigi)))
	if d.Gigi == "" {
		d.Gigi = StatusFisikNormal
	}
	d.Leher = StatusFisik(strings.TrimSpace(string(d.Leher)))
	if d.Leher == "" {
		d.Leher = StatusFisikNormal
	}
	d.Thoraks = StatusFisik(strings.TrimSpace(string(d.Thoraks)))
	if d.Thoraks == "" {
		d.Thoraks = StatusFisikNormal
	}
	d.Abdomen = StatusFisik(strings.TrimSpace(string(d.Abdomen)))
	if d.Abdomen == "" {
		d.Abdomen = StatusFisikNormal
	}
	d.Genital = StatusFisik(strings.TrimSpace(string(d.Genital)))
	if d.Genital == "" {
		d.Genital = StatusFisikNormal
	}
	d.Ekstremitas = StatusFisik(strings.TrimSpace(string(d.Ekstremitas)))
	if d.Ekstremitas == "" {
		d.Ekstremitas = StatusFisikNormal
	}
	d.KeteranganFisik = strings.TrimSpace(d.KeteranganFisik)

	d.KeteranganLokalis = strings.TrimSpace(d.KeteranganLokalis)
	d.EKG = strings.TrimSpace(d.EKG)
	d.Radiologi = strings.TrimSpace(d.Radiologi)
	d.Laboratorium = strings.TrimSpace(d.Laboratorium)
	d.Diagnosis = strings.TrimSpace(d.Diagnosis)
	d.TataLaksana = strings.TrimSpace(d.TataLaksana)
}

func (d *DataPenilaianMedisIGD) Validate(errs apperror.ValidationError) {
	d.Sanitize()

	if d.KeluhanUtama == "" {
		errs["keluhan_utama"] = "Keluhan utama wajib diisi"
	} else if len(d.KeluhanUtama) > 2000 {
		errs["keluhan_utama"] = "Keluhan utama maksimal 2000 karakter"
	}

	if len(d.RiwayatPenyakitSekarang) > 2000 {
		errs["riwayat_penyakit_sekarang"] = "Riwayat penyakit sekarang maksimal 2000 karakter"
	}
	if len(d.RiwayatPenyakitDahulu) > 1000 {
		errs["riwayat_penyakit_dahulu"] = "Riwayat penyakit dahulu maksimal 1000 karakter"
	}
	if len(d.RiwayatPenyakitKeluarga) > 1000 {
		errs["riwayat_penyakit_keluarga"] = "Riwayat penyakit keluarga maksimal 1000 karakter"
	}
	if len(d.RiwayatPengobatan) > 1000 {
		errs["riwayat_pengobatan"] = "Riwayat pengobatan maksimal 1000 karakter"
	}
	if len(d.Alergi) > 100 {
		errs["alergi"] = "Alergi maksimal 100 karakter"
	}

	if !d.Anamnesis.IsValid() {
		errs["anamnesis"] = "Anamnesis tidak valid (pilih Autoanamnesis atau Alloanamnesis)"
	} else if d.Anamnesis == Alloanamnesis && d.Hubungan == "" {
		errs["hubungan"] = "Hubungan keluarga/pengantar wajib diisi untuk Alloanamnesis"
	}
	if len(d.Hubungan) > 100 {
		errs["hubungan"] = "Hubungan keluarga/pengantar maksimal 100 karakter"
	}

	if !d.Keadaan.IsValid() {
		errs["keadaan"] = "Keadaan umum tidak valid"
	}

	if !d.Kesadaran.IsValid() {
		errs["kesadaran"] = "Tingkat kesadaran tidak valid"
	}

	if d.SuhuTubuh != "" {
		if msg := validateSuhuTubuh(d.SuhuTubuh); msg != "" {
			errs["suhu_tubuh"] = msg
		}
	}
	if d.Tensi != "" {
		if msg := validateTensi(d.Tensi); msg != "" {
			errs["tensi"] = msg
		}
	}
	if d.Nadi != "" {
		if msg := validateNadi(d.Nadi); msg != "" {
			errs["nadi"] = msg
		}
	}
	if d.Respirasi != "" {
		if msg := validateRespirasi(d.Respirasi); msg != "" {
			errs["respirasi"] = msg
		}
	}
	if d.TinggiBadan != "" {
		if msg := validateTinggiBadan(d.TinggiBadan); msg != "" {
			errs["tinggi_badan"] = msg
		}
	}
	if d.BeratBadan != "" {
		if msg := validateBeratBadan(d.BeratBadan); msg != "" {
			errs["berat_badan"] = msg
		}
	}
	if d.SpO2 != "" {
		if msg := validateSpO2(d.SpO2); msg != "" {
			errs["spo2"] = msg
		}
	}
	if len(d.Gcs) > 10 {
		errs["gcs"] = "GCS maksimal 10 karakter"
	}

	if !d.Kepala.IsValid() {
		errs["kepala"] = "Status pemeriksaan kepala tidak valid"
	}
	if !d.Mata.IsValid() {
		errs["mata"] = "Status pemeriksaan mata tidak valid"
	}
	if !d.Gigi.IsValid() {
		errs["gigi"] = "Status pemeriksaan gigi tidak valid"
	}
	if !d.Leher.IsValid() {
		errs["leher"] = "Status pemeriksaan leher tidak valid"
	}
	if !d.Thoraks.IsValid() {
		errs["thoraks"] = "Status pemeriksaan thoraks tidak valid"
	}
	if !d.Abdomen.IsValid() {
		errs["abdomen"] = "Status pemeriksaan abdomen tidak valid"
	}
	if !d.Genital.IsValid() {
		errs["genital"] = "Status pemeriksaan genital tidak valid"
	}
	if !d.Ekstremitas.IsValid() {
		errs["ekstremitas"] = "Status pemeriksaan ekstremitas tidak valid"
	}

	if t, err := time.ParseInLocation("2006-01-02 15:04:05", d.TanggalPenilaian, time.Local); err != nil {
		errs["tanggal_penilaian"] = "Format tanggal penilaian harus YYYY-MM-DD HH:mm:ss (contoh: 2026-08-31 09:30:00) atau YYYY-MM-DD"
	} else if t.After(time.Now()) {
		errs["tanggal_penilaian"] = "Waktu penilaian medis tidak boleh melebihi waktu saat ini"
	}

	if d.Diagnosis == "" {
		errs["diagnosis"] = "Diagnosis wajib diisi"
	} else if len(d.Diagnosis) > 500 {
		errs["diagnosis"] = "Diagnosis maksimal 500 karakter"
	}

	if d.TataLaksana == "" {
		errs["tata_laksana"] = "Tata laksana wajib diisi"
	}
}

type PenilaianMedisIGD struct {
	IdKunjungan       string `json:"id_kunjungan"`
	IdPasien          string `json:"id_pasien"`
	NoRawat           string `json:"no_rawat"`
	NoRM              string `json:"no_rkm_medis"`
	KodePoli          string `json:"kode_poli"`
	NamaPoli          string `json:"nama_poli"`
	TanggalRegistrasi string `json:"tanggal_registrasi"`
	KodeDokter        string `json:"kode_dokter"`
	NamaDokter        string `json:"nama_dokter"`

	DataPenilaianMedisIGD
}

type SimpanPenilaianMedisIGDRequest struct {
	NoRawat string `json:"no_rawat"`
	DataPenilaianMedisIGD
}

func (req *SimpanPenilaianMedisIGDRequest) Sanitize() {
	req.NoRawat = strings.TrimSpace(req.NoRawat)
	req.DataPenilaianMedisIGD.Sanitize()
}

func (req *SimpanPenilaianMedisIGDRequest) Validate() apperror.ValidationError {
	req.Sanitize()
	errs := make(apperror.ValidationError)

	if req.NoRawat == "" {
		errs["no_rawat"] = "Nomor rawat wajib diisi"
	}

	req.DataPenilaianMedisIGD.Validate(errs)
	if len(errs) > 0 {
		return errs
	}
	return nil
}

type UpdatePenilaianMedisIGDRequest struct {
	DataPenilaianMedisIGD
}

func (req *UpdatePenilaianMedisIGDRequest) Sanitize() {
	req.DataPenilaianMedisIGD.Sanitize()
}

func (req *UpdatePenilaianMedisIGDRequest) Validate() apperror.ValidationError {
	req.Sanitize()
	errs := make(apperror.ValidationError)

	req.DataPenilaianMedisIGD.Validate(errs)
	if len(errs) > 0 {
		return errs
	}
	return nil
}
