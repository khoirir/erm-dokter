package shared

import "strings"

var KodeDepoFarmasi = []string{
	"DPHD",
	"DPIGD",
	"DPOK",
	"DPRI",
	"DPRJ",
	"GDF",
}

func InClauseDepoFarmasi(colName string) string {
	quoted := make([]string, len(KodeDepoFarmasi))
	for i, k := range KodeDepoFarmasi {
		quoted[i] = "'" + k + "'"
	}
	return colName + " IN (" + strings.Join(quoted, ",") + ")"
}
