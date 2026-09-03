package laboratorium

type HasilLabPA struct {
	DiagnosaKlinik string `json:"diagnosa_klinik"`
	Makroskopis    string `json:"makroskopis"`
	Mikroskopis    string `json:"mikroskopis"`
	Kesimpulan     string `json:"kesimpulan"`
	Kesan          string `json:"kesan"`
}
