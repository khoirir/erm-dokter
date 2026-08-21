package shared

import (
	"errors"
	"fmt"
	"strconv"
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

// ParseWaktu menggabungkan tanggal (YYYY-MM-DD) dan jam (opsional, contoh: HH:mm:ss atau HH:mm).
// Jika parameter jam kosong, default waktu akan diset ke 00:00:00.
func ParseWaktu(tanggal, jam string) (time.Time, error) {
	tanggal = strings.TrimSpace(tanggal)
	jam = strings.TrimSpace(jam)
	if tanggal == "" {
		return time.Time{}, errors.New("tanggal tidak boleh kosong")
	}

	tgl, err := time.Parse("2006-01-02", tanggal)
	if err != nil {
		return time.Time{}, fmt.Errorf("format tanggal tidak valid (harus YYYY-MM-DD): %w", err)
	}

	hour, min, sec := 0, 0, 0
	if jam != "" {
		jamParts := strings.Split(jam, ":")
		if len(jamParts) >= 1 {
			hour, _ = strconv.Atoi(jamParts[0])
		}
		if len(jamParts) >= 2 {
			min, _ = strconv.Atoi(jamParts[1])
		}
		if len(jamParts) >= 3 {
			sec, _ = strconv.Atoi(jamParts[2])
		}
	}

	return time.Date(tgl.Year(), tgl.Month(), tgl.Day(), hour, min, sec, 0, time.Local), nil
}

// ValidasiBatasWaktuRekamMedis memeriksa apakah data rekam medis masih boleh diubah/dihapus berdasarkan batas waktu (dalam jam).
// Mendukung data dengan tanggal+jam maupun data yang hanya memiliki tanggal saja (jam kosong).
func ValidasiBatasWaktuRekamMedis(tanggal, jam string, maxJam int, aksi string) error {
	if maxJam <= 0 {
		return nil
	}

	waktuData, err := ParseWaktu(tanggal, jam)
	if err != nil {
		return apperror.NewBusinessError(err.Error())
	}

	batasWaktu := waktuData.Add(time.Duration(maxJam) * time.Hour)
	if time.Now().After(batasWaktu) {
		if aksi == "" {
			aksi = "diubah atau dihapus"
		}
		return apperror.NewForbiddenError(fmt.Sprintf("Data tidak dapat %s karena telah melewati batas waktu maksimal %d jam", aksi, maxJam))
	}

	return nil
}
