package resumepasien

type KondisiPulang string

const (
	KondisiPulangHidup     KondisiPulang = "Hidup"
	KondisiPulangMeninggal KondisiPulang = "Meninggal"
)

func (k KondisiPulang) IsValid() bool {
	switch k {
	case KondisiPulangHidup, KondisiPulangMeninggal:
		return true
	default:
		return false
	}
}

type CaraKeluar string

const (
	CaraKeluarAtasIzinDokter CaraKeluar = "Atas Izin Dokter"
	CaraKeluarPindahRS       CaraKeluar = "Pindah RS"
	CaraKeluarPulangSendiri  CaraKeluar = "Pulang Atas Permintaan Sendiri"
	CaraKeluarLainnya        CaraKeluar = "Lainnya"
)

func (c CaraKeluar) IsValid() bool {
	switch c {
	case CaraKeluarAtasIzinDokter, CaraKeluarPindahRS, CaraKeluarPulangSendiri, CaraKeluarLainnya:
		return true
	default:
		return false
	}
}

type KeadaanPulang string

const (
	KeadaanPulangMembaik       KeadaanPulang = "Membaik"
	KeadaanPulangSembuh        KeadaanPulang = "Sembuh"
	KeadaanPulangRujuk         KeadaanPulang = "Rujuk"
	KeadaanPulangKeadaanKhusus KeadaanPulang = "Keadaan Khusus"
	KeadaanPulangMeninggal     KeadaanPulang = "Meninggal"
)

func (k KeadaanPulang) IsValid() bool {
	switch k {
	case KeadaanPulangMembaik, KeadaanPulangSembuh, KeadaanPulangRujuk, KeadaanPulangKeadaanKhusus, KeadaanPulangMeninggal:
		return true
	default:
		return false
	}
}

type Dilanjutkan string

const (
	DilanjutkanKembaliKeRS Dilanjutkan = "Kembali Ke RS"
	DilanjutkanRSLain      Dilanjutkan = "RS Lain"
	DilanjutkanDokterLuar  Dilanjutkan = "Dokter Luar"
	DilanjutkanPuskesmas   Dilanjutkan = "Puskesmes"
	DilanjutkanLainnya     Dilanjutkan = "Lainnya"
)

func (d Dilanjutkan) IsValid() bool {
	switch d {
	case DilanjutkanKembaliKeRS, DilanjutkanRSLain, DilanjutkanDokterLuar, DilanjutkanPuskesmas, DilanjutkanLainnya:
		return true
	default:
		return false
	}
}

type OpsiReferensi struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type ReferensiResumeRalan struct {
	KondisiPulang []OpsiReferensi `json:"kondisi_pulang"`
}

type ReferensiResumeRanap struct {
	CaraKeluar    []OpsiReferensi `json:"cara_keluar"`
	KeadaanPulang []OpsiReferensi `json:"keadaan_pulang"`
	Dilanjutkan   []OpsiReferensi `json:"dilanjutkan"`
}

func DaftarOpsiKondisiPulang() []OpsiReferensi {
	return []OpsiReferensi{
		{Value: string(KondisiPulangHidup), Label: "Hidup"},
		{Value: string(KondisiPulangMeninggal), Label: "Meninggal"},
	}
}

func DaftarOpsiCaraKeluar() []OpsiReferensi {
	return []OpsiReferensi{
		{Value: string(CaraKeluarAtasIzinDokter), Label: "Atas Izin Dokter"},
		{Value: string(CaraKeluarPindahRS), Label: "Pindah RS"},
		{Value: string(CaraKeluarPulangSendiri), Label: "Pulang Atas Permintaan Sendiri"},
		{Value: string(CaraKeluarLainnya), Label: "Lainnya"},
	}
}

func DaftarOpsiKeadaanPulang() []OpsiReferensi {
	return []OpsiReferensi{
		{Value: string(KeadaanPulangMembaik), Label: "Membaik"},
		{Value: string(KeadaanPulangSembuh), Label: "Sembuh"},
		{Value: string(KeadaanPulangRujuk), Label: "Rujuk"},
		{Value: string(KeadaanPulangKeadaanKhusus), Label: "Keadaan Khusus"},
		{Value: string(KeadaanPulangMeninggal), Label: "Meninggal"},
	}
}

func DaftarOpsiDilanjutkan() []OpsiReferensi {
	return []OpsiReferensi{
		{Value: string(DilanjutkanKembaliKeRS), Label: "Kembali Ke RS"},
		{Value: string(DilanjutkanRSLain), Label: "RS Lain"},
		{Value: string(DilanjutkanDokterLuar), Label: "Dokter Luar"},
		{Value: string(DilanjutkanPuskesmas), Label: "Puskesmas"},
		{Value: string(DilanjutkanLainnya), Label: "Lainnya"},
	}
}
