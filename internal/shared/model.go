package shared

type MetaPaginasi struct {
	TotalData    int `json:"total_data"`
	TotalHalaman int `json:"total_halaman"`
	HalamanAktif int `json:"halaman_aktif"`
	BatasData    int `json:"batas_data"`
}