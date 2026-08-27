package rawatjalan

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

var ListStatusPemeriksaan = []struct {
	Value StatusPemeriksaan
	Label string
}{
	{Value: StatusBelum, Label: "Belum Periksa"},
	{Value: StatusSudah, Label: "Sudah Periksa"},
	{Value: StatusBatal, Label: "Batal Periksa"},
	{Value: StatusBerkasDiterima, Label: "Berkas Diterima"},
	{Value: StatusDirujuk, Label: "Dirujuk"},
	{Value: StatusMeninggal, Label: "Meninggal"},
	{Value: StatusDirawat, Label: "Dirawat"},
	{Value: StatusPulangPaksa, Label: "Pulang Paksa"},
}

func (s StatusPemeriksaan) IsValid() bool {
	for _, item := range ListStatusPemeriksaan {
		if s == item.Value {
			return true
		}
	}
	return false
}

type StatusBayar string

const (
	StatusBayarSudah StatusBayar = "Sudah Bayar"
	StatusBayarBelum StatusBayar = "Belum Bayar"
)

var ListStatusBayar = []struct {
	Value StatusBayar
	Label string
}{
	{Value: StatusBayarSudah, Label: "Sudah Bayar"},
	{Value: StatusBayarBelum, Label: "Belum Bayar"},
}

func (s StatusBayar) IsValid() bool {
	for _, item := range ListStatusBayar {
		if s == item.Value {
			return true
		}
	}
	return false
}

type JenisAntrean string

const (
	JenisAntreanRujukan      JenisAntrean = "Rujukan"
	JenisAntreanTidakRujukan JenisAntrean = "Bukan Rujukan"
)

var ListJenisAntrean = []struct {
	Value JenisAntrean
	Label string
}{
	{Value: JenisAntreanRujukan, Label: "Rujukan"},
	{Value: JenisAntreanTidakRujukan, Label: "Bukan Rujukan"},
}

func (j JenisAntrean) IsValid() bool {
	for _, item := range ListJenisAntrean {
		if j == item.Value {
			return true
		}
	}
	return false
}

type OrderBy string

const (
	OrderByWaktuRegistrasi OrderBy = "waktu_registrasi"
	OrderByNamaPasien     OrderBy = "nama_pasien"
)

func (o OrderBy) IsValid() bool {
	switch o {
	case OrderByWaktuRegistrasi, OrderByNamaPasien:
		return true
	default:
		return false
	}
}
