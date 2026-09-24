package obat

type OrderBy string

const (
	OrderByNamaObat     OrderBy = "nama_obat"
	OrderByStok         OrderBy = "stok"
	OrderByNamaDepo     OrderBy = "nama_depo"
	OrderByNamaJenis    OrderBy = "nama_jenis"
	OrderByNamaGolongan OrderBy = "nama_golongan"
	OrderByNamaKategori OrderBy = "nama_kategori"
)

func (o OrderBy) IsValid() bool {
	switch o {
	case OrderByNamaObat,
		OrderByStok,
		OrderByNamaDepo,
		OrderByNamaJenis,
		OrderByNamaGolongan,
		OrderByNamaKategori:
		return true
	default:
		return false
	}
}
