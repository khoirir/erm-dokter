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

func (s StatusPemeriksaan) IsValid() bool {
	switch s {
	case StatusBelum, StatusSudah, StatusBatal, StatusBerkasDiterima,
		StatusDirujuk, StatusMeninggal, StatusDirawat, StatusPulangPaksa:
		return true
	default:
		return false
	}
}

type StatusBayar string

const (
	StatusBayarSudah StatusBayar = "Sudah Bayar"
	StatusBayarBelum StatusBayar = "Belum Bayar"
)

func (s StatusBayar) IsValid() bool {
	switch s {
	case StatusBayarSudah, StatusBayarBelum:
		return true
	default:
		return false
	}
}

type JenisAntrean string

const (
	JenisAntreanRujukan      JenisAntrean = "Rujukan"
	JenisAntreanTidakRujukan JenisAntrean = "Bukan Rujukan"
)

func (j JenisAntrean) IsValid() bool {
	switch j {
	case JenisAntreanRujukan, JenisAntreanTidakRujukan:
		return true
	default:
		return false
	}
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

