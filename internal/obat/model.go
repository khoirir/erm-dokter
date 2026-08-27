package obat

import (
	"strings"

	"erm-dokter/internal/shared/apperror"
)

type Obat struct {
	Id           string     `json:"id"`
	KodeObat     string     `json:"kode_obat"`
	NamaObat     string     `json:"nama_obat"`
	Komposisi    string     `json:"komposisi"`
	Harga        string     `json:"harga"`
	Satuan       string     `json:"satuan"`
	Stok         float64    `json:"stok"`
	Kapasitas    float64    `json:"kapasitas"`
	KodeDepo     string     `json:"kode_depo,omitempty"`
	NamaDepo     string     `json:"nama_depo,omitempty"`
	KodeJenis    string     `json:"kode_jenis"`
	NamaJenis    string     `json:"nama_jenis"`
	KodeGolongan string     `json:"kode_golongan"`
	NamaGolongan string     `json:"nama_golongan"`
	KodeKategori string     `json:"kode_kategori"`
	NamaKategori string     `json:"nama_kategori"`
	StokDepo     []StokDepo `json:"stok_depo,omitempty"`
}

type StokDepo struct {
	KodeDepo string  `json:"kode_depo"`
	NamaDepo string  `json:"nama_depo"`
	Stok     float64 `json:"stok"`
}

type JenisObat struct {
	Kode string `json:"kode"`
	Nama string `json:"nama"`
}

type GolonganObat struct {
	Kode string `json:"kode"`
	Nama string `json:"nama"`
}

type KategoriObat struct {
	Kode string `json:"kode"`
	Nama string `json:"nama"`
}

type FilterDaftarObat struct {
	Depo      string `json:"depo,omitempty"`
	Jenis     string `json:"jenis,omitempty"`
	Golongan  string `json:"golongan,omitempty"`
	Kategori  string `json:"kategori,omitempty"`
	Keyword   string `json:"keyword,omitempty"`
	OrderBy   string `json:"order_by,omitempty"`
	SortOrder string `json:"sort_order,omitempty"`
	Page      int    `json:"page,omitempty"`
	Limit     int    `json:"limit,omitempty"`
}

func (f *FilterDaftarObat) Sanitize() {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Limit <= 0 {
		f.Limit = 20
	} else if f.Limit > 100 {
		f.Limit = 100
	}

	validSortColumns := map[string]string{
		"nama_obat":     "dtb.nama_brng",
		"stok":          "stok",
		"nama_depo":     "bg.nm_bangsal",
		"nama_jenis":    "jn.nama",
		"nama_golongan": "gb.nama",
		"nama_kategori": "kb.nama",
	}

	if col, ok := validSortColumns[strings.ToLower(f.OrderBy)]; ok {
		f.OrderBy = col
	} else {
		f.OrderBy = "dtb.nama_brng"
	}

	if strings.ToUpper(f.SortOrder) == "DESC" {
		f.SortOrder = "DESC"
	} else {
		f.SortOrder = "ASC"
	}
}

func (f FilterDaftarObat) Offset() int {
	return (f.Page - 1) * f.Limit
}

func (f *FilterDaftarObat) Validate() apperror.ValidationError {
	f.Sanitize()
	errs := make(apperror.ValidationError)

	keyword := strings.TrimSpace(f.Keyword)
	if len(keyword) > 0 && len(keyword) < 3 {
		errs["keyword"] = "Kata kunci pencarian minimal 3 karakter"
	}

	validSortOrders := map[string]bool{
		"ASC":  true,
		"DESC": true,
	}
	if f.SortOrder != "" && !validSortOrders[strings.ToUpper(strings.TrimSpace(f.SortOrder))] {
		errs["sort_order"] = "Jenis pengurutan tidak valid"
	}

	validSortColumns := map[string]bool{
		"nama_obat":     true,
		"stok":          true,
		"nama_depo":     true,
		"nama_jenis":    true,
		"nama_golongan": true,
		"nama_kategori": true,
	}
	if f.OrderBy != "" && !validSortColumns[strings.ToLower(strings.TrimSpace(f.OrderBy))] {
		errs["order_by"] = "Jenis pengurutan tidak valid"
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}
