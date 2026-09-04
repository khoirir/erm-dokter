package pemeriksaan

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

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
	Id          string `json:"id"`
	IdKunjungan string `json:"id_kunjungan"`
	NoRawat     string `json:"no_rawat"`
	DataPemeriksaan
	KodeDokterPetugas string              `json:"kode_dokter_petugas"`
	NamaDokterPetugas string              `json:"nama_dokter_petugas"`
	StatusLanjut      shared.StatusLanjut `json:"status_lanjut"`
}

type OpsiReferensi struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type FilterDaftarPemeriksaan struct {
	Tanggal string `json:"tanggal,omitempty"`
	Page    int    `json:"page,omitempty"`
	Limit   int    `json:"limit,omitempty"`
}

func (f *FilterDaftarPemeriksaan) Sanitize() {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Limit <= 0 {
		f.Limit = 20
	} else if f.Limit > 100 {
		f.Limit = 100
	}
}

func (f FilterDaftarPemeriksaan) Offset() int {
	return (f.Page - 1) * f.Limit
}

func (f *FilterDaftarPemeriksaan) Validate() apperror.ValidationError {
	errs := make(apperror.ValidationError)
	if f.Tanggal != "" {
		shared.ValidasiRentangTanggal(f.Tanggal, errs)
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

type DataPemeriksaan struct {
	TanggalPemeriksaan  string    `json:"tanggal_pemeriksaan"`
	JamPemeriksaan      string    `json:"jam_pemeriksaan"`
	SuhuTubuh           string    `json:"suhu_tubuh"`
	Tensi               string    `json:"tensi"`
	Nadi                string    `json:"nadi"`
	Respirasi           string    `json:"respirasi"`
	TinggiBadan         string    `json:"tinggi_badan"`
	BeratBadan          string    `json:"berat_badan"`
	SpO2                string    `json:"spo2"`
	Gcs                 string    `json:"gcs"`
	Kesadaran           Kesadaran `json:"kesadaran"`
	Keluhan             string    `json:"keluhan"`
	Pemeriksaan         string    `json:"pemeriksaan"`
	Alergi              string    `json:"alergi"`
	LingkarPerut        string    `json:"lingkar_perut,omitempty"`
	RencanaTindakLanjut string    `json:"rencana_tindak_lanjut"`
	Penilaian           string    `json:"penilaian"`
	Instruksi           string    `json:"instruksi"`
	Evaluasi            string    `json:"evaluasi"`
}

func (d *DataPemeriksaan) Sanitize() {
	d.TanggalPemeriksaan = strings.TrimSpace(d.TanggalPemeriksaan)
	d.JamPemeriksaan = strings.TrimSpace(d.JamPemeriksaan)
	if len(d.JamPemeriksaan) == 5 && strings.Count(d.JamPemeriksaan, ":") == 1 {
		d.JamPemeriksaan += ":00"
	}
	d.SuhuTubuh = strings.ReplaceAll(strings.TrimSpace(d.SuhuTubuh), ",", ".")
	d.Tensi = strings.ReplaceAll(strings.TrimSpace(d.Tensi), " ", "")
	d.Nadi = strings.TrimSpace(d.Nadi)
	d.Respirasi = strings.TrimSpace(d.Respirasi)
	d.TinggiBadan = strings.ReplaceAll(strings.TrimSpace(d.TinggiBadan), ",", ".")
	d.BeratBadan = strings.ReplaceAll(strings.TrimSpace(d.BeratBadan), ",", ".")
	d.SpO2 = strings.TrimSpace(d.SpO2)
	d.Gcs = strings.TrimSpace(d.Gcs)
	d.Kesadaran = Kesadaran(strings.TrimSpace(string(d.Kesadaran)))
	d.Keluhan = strings.TrimSpace(d.Keluhan)
	d.Pemeriksaan = strings.TrimSpace(d.Pemeriksaan)
	d.Alergi = strings.TrimSpace(d.Alergi)
	d.LingkarPerut = strings.TrimSpace(d.LingkarPerut)
	d.RencanaTindakLanjut = strings.TrimSpace(d.RencanaTindakLanjut)
	d.Penilaian = strings.TrimSpace(d.Penilaian)
	d.Instruksi = strings.TrimSpace(d.Instruksi)
	d.Evaluasi = strings.TrimSpace(d.Evaluasi)
}

func (d *DataPemeriksaan) Validate(errs apperror.ValidationError) {
	tgl, errTgl := time.Parse("2006-01-02", d.TanggalPemeriksaan)
	if errTgl != nil {
		errs["tanggal_pemeriksaan"] = "Format tanggal pemeriksaan harus YYYY-MM-DD (2026-01-01)"
	}

	jam, errJam := time.Parse("15:04:05", d.JamPemeriksaan)
	if errJam != nil {
		errs["jam_pemeriksaan"] = "Format jam pemeriksaan harus HH:mm:ss (12:10:00)"
	}

	if errTgl == nil && errJam == nil {
		waktuPemeriksaan := time.Date(
			tgl.Year(), tgl.Month(), tgl.Day(),
			jam.Hour(), jam.Minute(), jam.Second(), 0,
			time.Local,
		)
		if waktuPemeriksaan.After(time.Now()) {
			errs["tanggal_pemeriksaan"] = "Waktu pemeriksaan tidak boleh melebihi waktu saat ini"
		}
	}

	if d.Kesadaran == "" {
		errs["kesadaran"] = "Tingkat kesadaran wajib diisi"
	} else if !d.Kesadaran.IsValid() {
		errs["kesadaran"] = "Tingkat kesadaran tidak valid"
	}

	if d.Keluhan == "" {
		errs["keluhan"] = "Keluhan wajib diisi (Subjective)"
	} else if len(d.Keluhan) > 2000 {
		errs["keluhan"] = "Keluhan maksimal 2000 karakter (Subjective)"
	}

	if d.Pemeriksaan == "" {
		errs["pemeriksaan"] = "Pemeriksaan wajib diisi (Objective)"
	} else if len(d.Pemeriksaan) > 2000 {
		errs["pemeriksaan"] = "Pemeriksaan maksimal 2000 karakter (Objective)"
	}

	if d.Penilaian == "" {
		errs["penilaian"] = "Penilaian wajib diisi (Assessment)"
	} else if len(d.Penilaian) > 2000 {
		errs["penilaian"] = "Penilaian maksimal 2000 karakter (Assessment)"
	}

	if d.RencanaTindakLanjut == "" {
		errs["rencana_tindak_lanjut"] = "Rencana tindak lanjut wajib diisi (Plan)"
	} else if len(d.RencanaTindakLanjut) > 2000 {
		errs["rencana_tindak_lanjut"] = "Rencana tindak lanjut maksimal 2000 karakter (Plan)"
	}

	if d.Instruksi == "" {
		errs["instruksi"] = "Instruksi medis wajib diisi (Instruction)"
	} else if len(d.Instruksi) > 2000 {
		errs["instruksi"] = "Instruksi maksimal 2000 karakter (Instruction)"
	}

	if d.Evaluasi == "" {
		errs["evaluasi"] = "Evaluasi wajib diisi (Evaluation)"
	} else if len(d.Evaluasi) > 2000 {
		errs["evaluasi"] = "Evaluasi maksimal 2000 karakter (Evaluation)"
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
	if len(d.Alergi) > 50 {
		errs["alergi"] = "Alergi maksimal 50 karakter"
	}
	if len(d.LingkarPerut) > 5 {
		errs["lingkar_perut"] = "Lingkar perut maksimal 5 karakter"
	}
}

type SimpanPemeriksaanRequest struct {
	NoRawat string `json:"no_rawat"`
	DataPemeriksaan
}

func (r *SimpanPemeriksaanRequest) Sanitize() {
	r.NoRawat = strings.TrimSpace(r.NoRawat)
	if strings.TrimSpace(r.TanggalPemeriksaan) == "" {
		r.TanggalPemeriksaan = time.Now().Format("2006-01-02")
	}
	if strings.TrimSpace(r.JamPemeriksaan) == "" {
		r.JamPemeriksaan = time.Now().Format("15:04:05")
	}
	r.DataPemeriksaan.Sanitize()
}

func (r *SimpanPemeriksaanRequest) Validate() apperror.ValidationError {
	errs := make(apperror.ValidationError)

	if r.NoRawat == "" {
		errs["no_rawat"] = "Nomor rawat wajib diisi"
	}

	r.DataPemeriksaan.Validate(errs)

	if len(errs) > 0 {
		return errs
	}
	return nil
}

type UpdatePemeriksaanRequest struct {
	DataPemeriksaan
}

func (r *UpdatePemeriksaanRequest) Sanitize() {
	r.DataPemeriksaan.Sanitize()
}

func (r *UpdatePemeriksaanRequest) Validate() apperror.ValidationError {
	errs := make(apperror.ValidationError)

	r.DataPemeriksaan.Validate(errs)

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func validateSuhuTubuh(suhuStr string) string {
	if len(suhuStr) > 5 {
		return "Suhu tubuh maksimal 5 karakter (36.5)"
	}
	suhu, err := strconv.ParseFloat(suhuStr, 64)
	if err != nil {
		return "Suhu tubuh harus berupa angka (36.5)"
	}
	if suhu < 25.0 || suhu > 45.0 {
		return "Suhu tubuh harus berada dalam rentang 25.0 - 45.0 °C"
	}
	return ""
}

func validateTensi(tensi string) string {
	if len(tensi) > 8 {
		return "Tensi maksimal 8 karakter (120/80)"
	}

	parts := strings.Split(tensi, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "Format tensi harus Sistolik/Diastolik (120/80)"
	}

	sistolik, errSis := strconv.Atoi(parts[0])
	diastolik, errDia := strconv.Atoi(parts[1])
	if errSis != nil || errDia != nil {
		return "Nilai sistolik dan diastolik tensi harus berupa angka (120/80)"
	}

	if sistolik < 50 || sistolik > 300 || diastolik < 30 || diastolik > 200 {
		return "Tensi harus berada dalam rentang sistolik: 50-300 dan diastolik: 30-200 mmHg"
	}

	if sistolik <= diastolik {
		return "Nilai sistolik harus lebih besar dari diastolik (120/80)"
	}

	return ""
}

func validateNadi(nadiStr string) string {
	if len(nadiStr) > 3 {
		return "Nadi maksimal 3 karakter"
	}
	nadi, err := strconv.Atoi(nadiStr)
	if err != nil {
		return "Nadi harus berupa angka bulat (80)"
	}
	if nadi < 20 || nadi > 300 {
		return "Nadi harus berada dalam rentang 20 - 300 x/menit"
	}
	return ""
}

func validateRespirasi(respStr string) string {
	if len(respStr) > 3 {
		return "Respirasi maksimal 3 karakter"
	}
	resp, err := strconv.Atoi(respStr)
	if err != nil {
		return "Respirasi harus berupa angka bulat (20)"
	}
	if resp < 5 || resp > 100 {
		return "Respirasi harus berada dalam rentang 5 - 100 x/menit"
	}
	return ""
}

func validateTinggiBadan(tbStr string) string {
	if len(tbStr) > 5 {
		return "Tinggi badan maksimal 5 karakter"
	}
	tb, err := strconv.ParseFloat(tbStr, 64)
	if err != nil {
		return "Tinggi badan harus berupa angka (170)"
	}
	if tb < 20.0 || tb > 250.0 {
		return "Tinggi badan harus berada dalam rentang 20 - 250 cm"
	}
	return ""
}

func validateBeratBadan(bbStr string) string {
	if len(bbStr) > 5 {
		return "Berat badan maksimal 5 karakter"
	}
	bb, err := strconv.ParseFloat(bbStr, 64)
	if err != nil {
		return "Berat badan harus berupa angka (65)"
	}
	if bb < 0.5 || bb > 500.0 {
		return "Berat badan harus berada dalam rentang 0.5 - 500 kg"
	}
	return ""
}

func validateSpO2(spo2Str string) string {
	if len(spo2Str) > 3 {
		return "SpO2 maksimal 3 karakter"
	}
	spo2, err := strconv.Atoi(spo2Str)
	if err != nil {
		return "SpO2 harus berupa angka bulat (98)"
	}
	if spo2 < 0 || spo2 > 100 {
		return "SpO2 harus berada dalam rentang 0 - 100%"
	}
	return ""
}
