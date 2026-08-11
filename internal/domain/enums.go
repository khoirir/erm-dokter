package domain

type StatusPemeriksaan string

const (
	StatusBelum          StatusPemeriksaan = "Belum"
	StatusSudah          StatusPemeriksaan = "Sudah"
	StatusBatal          StatusPemeriksaan = "Batal"
	StatusBerkasDiterima StatusPemeriksaan = "Berkas Diterima"
	StatusDirujuk        StatusPemeriksaan = "Dirujuk"
	StatusMeninggal      StatusPemeriksaan = "Meninggal"
	StatusDirawat        StatusPemeriksaan = "Dirawat"
	StatusPulangPaksa    StatusPemeriksaan = "Pulang Paksa"
)

type StatusLanjut string

const (
	StatusLanjutRawatJalan StatusLanjut = "Ralan"
	StatusLanjutRawatInap  StatusLanjut = "Ranap"
)

type StatusBayar string

const (
	StatusBayarSudah StatusBayar = "Sudah Bayar"
	StatusBayarBelum StatusBayar = "Belum Bayar"
)

type JenisAntrean string

const (
	JenisAntreanRujukan      JenisAntrean = "Rujukan"
	JenisAntreanTidakRujukan JenisAntrean = "Bukan Rujukan"
)

type Kesadaran string

const (
	KesadaranComposMentis Kesadaran = "Compos Mentis"
	KesadaranSomnolen     Kesadaran = "Somnolen"
	KesadaranSopor        Kesadaran = "Sopor"
	KesadaranKoma         Kesadaran = "Koma"
	KesadaranAlert        Kesadaran = "Alert"
	KesadaranConfusion    Kesadaran = "Confusion"
	KesadaranVoice        Kesadaran = "Voice"
	KesadaranPain         Kesadaran = "Pain"
	KesadaranUnresponsive Kesadaran = "Unresponsive"
	KesadaranApatis       Kesadaran = "Apatis"
	KesadaranDelirium     Kesadaran = "Delirium"
)
