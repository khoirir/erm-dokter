package formatter

import "strings"

func ParseRentangTanggal(tanggal string) (string, string) {
	parts := strings.Split(tanggal, ",")
	if len(parts) != 2 {
		return "", ""
	}
	awal := strings.TrimSpace(parts[0])
	akhir := strings.TrimSpace(parts[1])
	if awal == "" || akhir == "" {
		return "", ""
	}
	return awal, akhir
}
