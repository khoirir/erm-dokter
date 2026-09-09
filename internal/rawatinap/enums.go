package rawatinap

type StatusPulang string

const (
	StatusPulangBelumPulang           StatusPulang = "-"
	StatusPulangSehat                 StatusPulang = "Sehat"
	StatusPulangRujuk                 StatusPulang = "Rujuk"
	StatusPulangAPS                   StatusPulang = "APS"
	StatusPulangPlus                  StatusPulang = "+"
	StatusPulangMeninggal             StatusPulang = "Meninggal"
	StatusPulangSembuh                StatusPulang = "Sembuh"
	StatusPulangMembaik               StatusPulang = "Membaik"
	StatusPulangPulangPaksa           StatusPulang = "Pulang Paksa"
	StatusPulangPindahKamar           StatusPulang = "Pindah Kamar"
	StatusPulangAtasPersetujuanDokter StatusPulang = "Atas Persetujuan Dokter"
	StatusPulangAtasPermintaanSendiri StatusPulang = "Atas Permintaan Sendiri"
	StatusPulangIsoman                StatusPulang = "Isoman"
	StatusPulangLainLain              StatusPulang = "Lain-lain"
)

var ListStatusPulang = []struct {
	Value StatusPulang
	Label string
}{
	{Value: StatusPulangBelumPulang, Label: "Belum Pulang (-)"},
	{Value: StatusPulangSehat, Label: "Sehat"},
	{Value: StatusPulangRujuk, Label: "Rujuk"},
	{Value: StatusPulangAPS, Label: "APS"},
	{Value: StatusPulangPlus, Label: "+"},
	{Value: StatusPulangMeninggal, Label: "Meninggal"},
	{Value: StatusPulangSembuh, Label: "Sembuh"},
	{Value: StatusPulangMembaik, Label: "Membaik"},
	{Value: StatusPulangPulangPaksa, Label: "Pulang Paksa"},
	{Value: StatusPulangPindahKamar, Label: "Pindah Kamar"},
	{Value: StatusPulangAtasPersetujuanDokter, Label: "Atas Persetujuan Dokter"},
	{Value: StatusPulangAtasPermintaanSendiri, Label: "Atas Permintaan Sendiri"},
	{Value: StatusPulangIsoman, Label: "Isoman"},
	{Value: StatusPulangLainLain, Label: "Lain-lain"},
}

func (s StatusPulang) IsValid() bool {
	for _, item := range ListStatusPulang {
		if s == item.Value {
			return true
		}
	}
	return false
}
