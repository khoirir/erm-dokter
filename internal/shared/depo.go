package shared

import "strings"

// KodeDepoFarmasi adalah daftar kode bangsal resmi yang bertindak sebagai depo farmasi di Khanza
var KodeDepoFarmasi = []string{
	"DPHD",
	"DPIGD",
	"DPOK",
	"DPRI",
	"DPRJ",
	"GDF",
}

// InClauseDepoFarmasi menghasilkan klausa SQL "column IN ('...')" berdasarkan daftar KodeDepoFarmasi
func InClauseDepoFarmasi(colName string) string {
	quoted := make([]string, len(KodeDepoFarmasi))
	for i, k := range KodeDepoFarmasi {
		quoted[i] = "'" + k + "'"
	}
	return colName + " IN (" + strings.Join(quoted, ",") + ")"
}
