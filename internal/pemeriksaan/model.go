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

type OpsiReferensi struct {
	Value string `json:"value"`
	Label string `json:"label"`
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
	NoRawat             string    `json:"no_rawat"`
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

func (r *SimpanPemeriksaanRequest) Sanitize() {
	r.NoRawat = strings.TrimSpace(r.NoRawat)
	r.TanggalPemeriksaan = strings.TrimSpace(r.TanggalPemeriksaan)
	r.JamPemeriksaan = strings.TrimSpace(r.JamPemeriksaan)
	r.SuhuTubuh = strings.ReplaceAll(strings.TrimSpace(r.SuhuTubuh), ",", ".")
	r.Tensi = strings.ReplaceAll(strings.TrimSpace(r.Tensi), " ", "")
	r.Nadi = strings.TrimSpace(r.Nadi)
	r.Respirasi = strings.TrimSpace(r.Respirasi)
	r.TinggiBadan = strings.ReplaceAll(strings.TrimSpace(r.TinggiBadan), ",", ".")
	r.BeratBadan = strings.ReplaceAll(strings.TrimSpace(r.BeratBadan), ",", ".")
	r.SpO2 = strings.TrimSpace(r.SpO2)
	r.Gcs = strings.TrimSpace(r.Gcs)
	r.Kesadaran = Kesadaran(strings.TrimSpace(string(r.Kesadaran)))
	r.Keluhan = strings.TrimSpace(r.Keluhan)
	r.Pemeriksaan = strings.TrimSpace(r.Pemeriksaan)
	r.Alergi = strings.TrimSpace(r.Alergi)
	r.LingkarPerut = strings.TrimSpace(r.LingkarPerut)
	r.RencanaTindakLanjut = strings.TrimSpace(r.RencanaTindakLanjut)
	r.Penilaian = strings.TrimSpace(r.Penilaian)
	r.Instruksi = strings.TrimSpace(r.Instruksi)
	r.Evaluasi = strings.TrimSpace(r.Evaluasi)

	if r.TanggalPemeriksaan == "" {
		r.TanggalPemeriksaan = time.Now().Format("2006-01-02")
	}
	if r.JamPemeriksaan == "" {
		r.JamPemeriksaan = time.Now().Format("15:04:05")
	}
}

func (r *SimpanPemeriksaanRequest) Validate() apperror.ValidationError {
	r.Sanitize()
	errs := make(apperror.ValidationError)

	if r.NoRawat == "" {
		errs["no_rawat"] = "Nomor rawat wajib diisi"
	} else if len(r.NoRawat) > 17 {
		errs["no_rawat"] = "Nomor rawat maksimal 17 karakter"
	}

	tgl, errTgl := time.Parse("2006-01-02", r.TanggalPemeriksaan)
	if errTgl != nil {
		errs["tanggal_pemeriksaan"] = "Format tanggal pemeriksaan harus YYYY-MM-DD"
	}

	jam, errJam := time.Parse("15:04:05", r.JamPemeriksaan)
	if errJam != nil {
		errs["jam_pemeriksaan"] = "Format jam pemeriksaan harus HH:mm:ss (contoh: 12:10:00)"
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

	if r.Kesadaran == "" {
		errs["kesadaran"] = "Tingkat kesadaran wajib diisi"
	} else if !r.Kesadaran.IsValid() {
		errs["kesadaran"] = "Tingkat kesadaran tidak valid"
	}

	if r.Keluhan == "" {
		errs["keluhan"] = "Keluhan wajib diisi"
	} else if len(r.Keluhan) > 2000 {
		errs["keluhan"] = "Keluhan maksimal 2000 karakter"
	}

	if r.Pemeriksaan == "" {
		errs["pemeriksaan"] = "Pemeriksaan fisik/objektif wajib diisi"
	} else if len(r.Pemeriksaan) > 2000 {
		errs["pemeriksaan"] = "Pemeriksaan maksimal 2000 karakter"
	}

	if r.Penilaian == "" {
		errs["penilaian"] = "Penilaian/Assessment wajib diisi"
	} else if len(r.Penilaian) > 2000 {
		errs["penilaian"] = "Penilaian maksimal 2000 karakter"
	}

	if r.RencanaTindakLanjut == "" {
		errs["rencana_tindak_lanjut"] = "Rencana tindak lanjut wajib diisi"
	} else if len(r.RencanaTindakLanjut) > 2000 {
		errs["rencana_tindak_lanjut"] = "Rencana tindak lanjut maksimal 2000 karakter"
	}

	if r.Instruksi == "" {
		errs["instruksi"] = "Instruksi medis wajib diisi"
	} else if len(r.Instruksi) > 2000 {
		errs["instruksi"] = "Instruksi maksimal 2000 karakter"
	}

	if r.Evaluasi == "" {
		errs["evaluasi"] = "Evaluasi wajib diisi"
	} else if len(r.Evaluasi) > 2000 {
		errs["evaluasi"] = "Evaluasi maksimal 2000 karakter"
	}

	if r.SuhuTubuh != "" {
		if msg := validateSuhuTubuh(r.SuhuTubuh); msg != "" {
			errs["suhu_tubuh"] = msg
		}
	}
	if r.Tensi != "" {
		if msg := validateTensi(r.Tensi); msg != "" {
			errs["tensi"] = msg
		}
	}
	if r.Nadi != "" {
		if msg := validateNadi(r.Nadi); msg != "" {
			errs["nadi"] = msg
		}
	}
	if r.Respirasi != "" {
		if msg := validateRespirasi(r.Respirasi); msg != "" {
			errs["respirasi"] = msg
		}
	}
	if r.TinggiBadan != "" {
		if msg := validateTinggiBadan(r.TinggiBadan); msg != "" {
			errs["tinggi_badan"] = msg
		}
	}
	if r.BeratBadan != "" {
		if msg := validateBeratBadan(r.BeratBadan); msg != "" {
			errs["berat_badan"] = msg
		}
	}
	if r.SpO2 != "" {
		if msg := validateSpO2(r.SpO2); msg != "" {
			errs["spo2"] = msg
		}
	}
	if len(r.Gcs) > 10 {
		errs["gcs"] = "GCS maksimal 10 karakter"
	}
	if len(r.Alergi) > 50 {
		errs["alergi"] = "Alergi maksimal 50 karakter"
	}
	if len(r.LingkarPerut) > 5 {
		errs["lingkar_perut"] = "Lingkar perut maksimal 5 karakter"
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func validateSuhuTubuh(suhuStr string) string {
	if len(suhuStr) > 5 {
		return "Suhu tubuh maksimal 5 karakter (contoh: 36.5)"
	}
	suhu, err := strconv.ParseFloat(suhuStr, 64)
	if err != nil {
		return "Suhu tubuh harus berupa angka (contoh: 36.5)"
	}
	if suhu < 25.0 || suhu > 45.0 {
		return "Suhu tubuh tidak wajar untuk manusia (rentang wajar: 25.0 - 45.0 °C)"
	}
	return ""
}

func validateTensi(tensi string) string {
	if len(tensi) > 8 {
		return "Tensi maksimal 8 karakter (contoh: 120/80)"
	}

	parts := strings.Split(tensi, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "Format tensi harus Sistolik/Diastolik (contoh: 120/80)"
	}

	sistolik, errSis := strconv.Atoi(parts[0])
	diastolik, errDia := strconv.Atoi(parts[1])
	if errSis != nil || errDia != nil {
		return "Nilai sistolik dan diastolik tensi harus berupa angka (contoh: 120/80)"
	}

	if sistolik < 50 || sistolik > 300 || diastolik < 30 || diastolik > 200 {
		return "Nilai tensi tidak wajar untuk manusia (sistolik: 50-300, diastolik: 30-200 mmHg)"
	}

	if sistolik <= diastolik {
		return "Nilai sistolik harus lebih besar dari diastolik (contoh: 120/80)"
	}

	return ""
}

func validateNadi(nadiStr string) string {
	if len(nadiStr) > 3 {
		return "Nadi maksimal 3 karakter"
	}
	nadi, err := strconv.Atoi(nadiStr)
	if err != nil {
		return "Nadi harus berupa angka bulat (contoh: 80)"
	}
	if nadi < 20 || nadi > 300 {
		return "Nilai nadi tidak wajar untuk manusia (rentang wajar: 20 - 300 x/menit)"
	}
	return ""
}

func validateRespirasi(respStr string) string {
	if len(respStr) > 3 {
		return "Respirasi maksimal 3 karakter"
	}
	resp, err := strconv.Atoi(respStr)
	if err != nil {
		return "Respirasi harus berupa angka bulat (contoh: 20)"
	}
	if resp < 5 || resp > 100 {
		return "Nilai respirasi tidak wajar untuk manusia (rentang wajar: 5 - 100 x/menit)"
	}
	return ""
}

func validateTinggiBadan(tbStr string) string {
	if len(tbStr) > 5 {
		return "Tinggi badan maksimal 5 karakter"
	}
	tb, err := strconv.ParseFloat(tbStr, 64)
	if err != nil {
		return "Tinggi badan harus berupa angka (contoh: 170)"
	}
	if tb < 20.0 || tb > 250.0 {
		return "Tinggi badan tidak wajar untuk manusia (rentang wajar: 20 - 250 cm)"
	}
	return ""
}

func validateBeratBadan(bbStr string) string {
	if len(bbStr) > 5 {
		return "Berat badan maksimal 5 karakter"
	}
	bb, err := strconv.ParseFloat(bbStr, 64)
	if err != nil {
		return "Berat badan harus berupa angka (contoh: 65)"
	}
	if bb < 0.5 || bb > 500.0 {
		return "Berat badan tidak wajar untuk manusia (rentang wajar: 0.5 - 500 kg)"
	}
	return ""
}

func validateSpO2(spo2Str string) string {
	if len(spo2Str) > 3 {
		return "SpO2 maksimal 3 karakter"
	}
	spo2, err := strconv.Atoi(spo2Str)
	if err != nil {
		return "SpO2 harus berupa angka bulat (contoh: 98)"
	}
	if spo2 < 0 || spo2 > 100 {
		return "SpO2 harus berada dalam rentang 0 - 100%"
	}
	return ""
}
