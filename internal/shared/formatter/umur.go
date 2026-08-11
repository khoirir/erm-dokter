package formatter

import (
	"fmt"
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
	y1, m1, d1 := tglLahir.Date()
	y2, m2, d2 := sekarang.Date()

	tahun := y2 - y1
	bulan := int(m2) - int(m1)
	hari := d2 - d1

	if hari < 0 {
		bulan--
		blnSebelumnya := int(m2) - 1
		thnSebelumnya := y2
		if blnSebelumnya == 0 {
			blnSebelumnya = 12
			thnSebelumnya--
		}
		hariDiBulanLalu := time.Date(thnSebelumnya, time.Month(blnSebelumnya+1), 0, 0, 0, 0, 0, time.UTC).Day()
		hari += hariDiBulanLalu
	}

	if bulan < 0 {
		tahun--
		bulan += 12
	}

	if tahun < 0 {
		return Umur{}
	}

	return Umur{
		Tahun: tahun,
		Bulan: bulan,
		Hari:  hari,
	}
}

func FormatUmur(tglLahirStr string) string {
	u := HitungUmur(tglLahirStr)
	return fmt.Sprintf("%d Th %d Bl %d Hr", u.Tahun, u.Bulan, u.Hari)
}
