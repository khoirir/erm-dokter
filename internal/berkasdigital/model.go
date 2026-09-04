package berkasdigital

import "io"

type MasterBerkasDigital struct {
	Kode string `json:"kode"`
	Nama string `json:"nama"`
}

type BerkasDigital struct {
	Kode       string `json:"kode"`
	NamaBerkas string `json:"nama_berkas"`
	IdBerkas   string `json:"id_berkas"`
	UrlBerkas  string `json:"url_berkas"`
}

type BerkasDigitalPerawatan struct {
	NoRawat    string `json:"no_rawat"`
	Kode       string `json:"kode"`
	NamaBerkas string `json:"nama_berkas"`
	LokasiFile string `json:"lokasi_file"`
}

type BerkasStream struct {
	Body          io.ReadCloser
	ContentType   string
	ContentLength string
	Filename      string
}

