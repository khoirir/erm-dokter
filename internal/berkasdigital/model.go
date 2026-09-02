package berkasdigital

type MasterBerkasDigital struct {
	Kode string `json:"kode"`
	Nama string `json:"nama"`
}

type BerkasDigitalItem struct {
	Kode       string `json:"kode"`
	NamaBerkas string `json:"nama_berkas"`
	IdBerkas   string `json:"id_berkas"`
	UrlBerkas  string `json:"url_berkas"`
}

type BerkasDigitalDB struct {
	NoRawat       string
	Kode          string
	NamaBerkas    string
	LokasiFile    string
	StatusLanjut  string
	TglRegistrasi string
	JamReg        string
	KodeDokter    string
	NamaDokter    string
}
