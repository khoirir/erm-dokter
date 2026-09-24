package penilaianmedis

import (
	"strconv"
	"strings"
)

type Anamnesis string

const (
	Autoanamnesis Anamnesis = "Autoanamnesis"
	Alloanamnesis Anamnesis = "Alloanamnesis"
)

func (a Anamnesis) IsValid() bool {
	switch a {
	case Autoanamnesis, Alloanamnesis:
		return true
	default:
		return false
	}
}

type Keadaan string

const (
	KeadaanSehat       Keadaan = "Sehat"
	KeadaanSakitRingan Keadaan = "Sakit Ringan"
	KeadaanSakitSedang Keadaan = "Sakit Sedang"
	KeadaanSakitBerat  Keadaan = "Sakit Berat"
)

func (k Keadaan) IsValid() bool {
	switch k {
	case KeadaanSehat, KeadaanSakitRingan, KeadaanSakitSedang, KeadaanSakitBerat:
		return true
	default:
		return false
	}
}

type KesadaranAwal string

const (
	KesadaranComposMentis KesadaranAwal = "Compos Mentis"
	KesadaranApatis       KesadaranAwal = "Apatis"
	KesadaranSomnolen     KesadaranAwal = "Somnolen"
	KesadaranSopor        KesadaranAwal = "Sopor"
	KesadaranKoma         KesadaranAwal = "Koma"
)

func (k KesadaranAwal) IsValid() bool {
	switch k {
	case KesadaranComposMentis, KesadaranApatis, KesadaranSomnolen, KesadaranSopor, KesadaranKoma:
		return true
	default:
		return false
	}
}

type StatusFisik string

const (
	StatusFisikNormal         StatusFisik = "Normal"
	StatusFisikAbnormal       StatusFisik = "Abnormal"
	StatusFisikTidakDiperiksa StatusFisik = "Tidak Diperiksa"
)

func (s StatusFisik) IsValid() bool {
	switch s {
	case StatusFisikNormal, StatusFisikAbnormal, StatusFisikTidakDiperiksa:
		return true
	default:
		return false
	}
}

type OpsiReferensi struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

func DaftarOpsiAnamnesis() []OpsiReferensi {
	return []OpsiReferensi{
		{Value: string(Autoanamnesis), Label: "Autoanamnesis (Pasien Sendiri)"},
		{Value: string(Alloanamnesis), Label: "Alloanamnesis (Keluarga / Pengantar)"},
	}
}

func DaftarOpsiKeadaan() []OpsiReferensi {
	return []OpsiReferensi{
		{Value: string(KeadaanSehat), Label: "Sehat"},
		{Value: string(KeadaanSakitRingan), Label: "Sakit Ringan"},
		{Value: string(KeadaanSakitSedang), Label: "Sakit Sedang"},
		{Value: string(KeadaanSakitBerat), Label: "Sakit Berat"},
	}
}

func DaftarOpsiKesadaran() []OpsiReferensi {
	return []OpsiReferensi{
		{Value: string(KesadaranComposMentis), Label: "Compos Mentis"},
		{Value: string(KesadaranApatis), Label: "Apatis"},
		{Value: string(KesadaranSomnolen), Label: "Somnolen"},
		{Value: string(KesadaranSopor), Label: "Sopor"},
		{Value: string(KesadaranKoma), Label: "Koma"},
	}
}

func DaftarOpsiStatusFisik() []OpsiReferensi {
	return []OpsiReferensi{
		{Value: string(StatusFisikNormal), Label: "Normal"},
		{Value: string(StatusFisikAbnormal), Label: "Abnormal"},
		{Value: string(StatusFisikTidakDiperiksa), Label: "Tidak Diperiksa"},
	}
}

type ReferensiPenilaianMedis struct {
	Anamnesis   []OpsiReferensi `json:"anamnesis"`
	Keadaan     []OpsiReferensi `json:"keadaan"`
	Kesadaran   []OpsiReferensi `json:"kesadaran"`
	StatusFisik []OpsiReferensi `json:"status_fisik"`
}

func GetReferensiPenilaianMedis() ReferensiPenilaianMedis {
	return ReferensiPenilaianMedis{
		Anamnesis:   DaftarOpsiAnamnesis(),
		Keadaan:     DaftarOpsiKeadaan(),
		Kesadaran:   DaftarOpsiKesadaran(),
		StatusFisik: DaftarOpsiStatusFisik(),
	}
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
