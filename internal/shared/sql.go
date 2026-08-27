package shared

import "strings"

// CreateInPlaceholders membuat deretan placeholder (?,?,?) sesuai jumlah parameter untuk query SQL IN (...)
func CreateInPlaceholders(count int) string {
	if count <= 0 {
		return "?"
	}
	placeholders := make([]string, count)
	for i := range count {
		placeholders[i] = "?"
	}
	return strings.Join(placeholders, ",")
}
