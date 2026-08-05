package domain

import (
	"fmt"
	"strings"
	"time"
)

type Umur struct {
	Tahun int `json:"tahun"`
	Bulan int `json:"bulan"`
	Hari  int `json:"hari"`
}

func HitungUmur(tglLahirStr string) Umur {
	tglLahir, err := time.Parse("2006-01-02", tglLahirStr)
	if err != nil {
		return Umur{}
	}

	sekarang := time.Now()

	tahun := sekarang.Year() - tglLahir.Year()
	bulan := int(sekarang.Month()) - int(tglLahir.Month())
	hari := sekarang.Day() - tglLahir.Day()

	if hari < 0 {
		bulan--
		bulanLalu := sekarang.AddDate(0, -1, 0)
		hariDiBulanLalu := time.Date(bulanLalu.Year(), bulanLalu.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
		hari += hariDiBulanLalu
	}

	if bulan < 0 {
		tahun--
		bulan += 12
	}

	return Umur{
		Tahun: tahun,
		Bulan: bulan,
		Hari:  hari,
	}
}

func FormatUmur(tglLahirStr string) string {
	u := HitungUmur(tglLahirStr)
	return fmt.Sprintf("%d Th %d Bln %d Hari", u.Tahun, u.Bulan, u.Hari)
}

func FormatJenisKelamin(jk string) string {
	switch strings.ToUpper(jk) {
	case "L":
		return "Laki-Laki"
	case "P":
		return "Perempuan"
	default:
		return jk
	}
}

func FormatNoRekamMedis(noRM string) string {
	noRM = strings.TrimSpace(noRM)
	if len(noRM) == 8 && strings.HasPrefix(noRM, "00") {
		return noRM[2:]
	}
	return noRM
}

type MetaPaginasi struct {
	TotalData    int `json:"total_data"`
	TotalHalaman int `json:"total_halaman"`
	HalamanAktif int `json:"halaman_aktif"`
	BatasData    int `json:"batas_data"`
}
