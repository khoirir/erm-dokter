package domain

// StatusPemeriksaan menentukan status pemeriksaan pasien
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

// StatusLanjut menentukan jenis perawatan lanjutan
type StatusLanjut string

const (
	StatusLanjutRawatJalan StatusLanjut = "Ralan"
	StatusLanjutRawatInap  StatusLanjut = "Ranap"
)

// StatusBayar menentukan status pembayaran
type StatusBayar string

const (
	StatusBayarSudah StatusBayar = "Sudah Bayar"
	StatusBayarBelum StatusBayar = "Belum Bayar"
)

// JenisAntrean menentukan jenis antrean pasien
type JenisAntrean string

const (
	JenisAntreanRujukan      JenisAntrean = "Rujukan"
	JenisAntreanTidakRujukan JenisAntrean = "Bukan Rujukan"
)

// Kesadaran menentukan tingkat kesadaran pasien
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
