package formatter

import "strings"

func FormatJenisKelamin(jk string) string {
	switch strings.ToUpper(strings.TrimSpace(jk)) {
	case "L":
		return "Laki-Laki"
	case "P":
		return "Perempuan"
	default:
		return jk
	}
}
