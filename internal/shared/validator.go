package shared

import (
	"strings"
	"time"

	"erm-dokter/internal/shared/apperror"
)

func ValidasiRentangTanggal(tanggal string, errs apperror.ValidationError) {
	const layout = "2006-01-02"
	parts := strings.Split(tanggal, ",")

	if len(parts) != 2 {
		errs["tanggal"] = "format tanggal wajib menggunakan rentang 2 tanggal (contoh: YYYY-MM-DD,YYYY-MM-DD)"
		return
	}

	tglAwal, err := time.Parse(layout, strings.TrimSpace(parts[0]))
	if err != nil {
		errs["tanggal"] = "format tanggal awal tidak valid (YYYY-MM-DD)"
		return
	}

	tglAkhir, err := time.Parse(layout, strings.TrimSpace(parts[1]))
	if err != nil {
		errs["tanggal"] = "format tanggal akhir tidak valid (YYYY-MM-DD)"
		return
	}

	if tglAwal.After(tglAkhir) {
		errs["tanggal"] = "tanggal awal harus lebih kecil atau sama dengan tanggal akhir"
	}
}
