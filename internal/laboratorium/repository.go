package laboratorium

import (
	"context"
	"database/sql"

	"erm-dokter/internal/shared"
)

type Repository interface {
	DaftarHasilLabPK(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]HasilLaboratorium, int, error)
	DaftarHasilLabPKByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]HasilLaboratorium, int, error)
	DetailHasilLabPK(ctx context.Context, noRawat string, kodeTindakan string, tanggalPeriksa string, jamPeriksa string) (*HasilLaboratorium, error)

	DaftarHasilLabPA(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]HasilLaboratorium, int, error)
	DaftarHasilLabPAByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]HasilLaboratorium, int, error)
	DetailHasilLabPA(ctx context.Context, noRawat string, kodeTindakan string, tanggalPeriksa string, jamPeriksa string) (*HasilLaboratorium, error)

	DaftarHasilLabMB(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]HasilLaboratorium, int, error)
	DaftarHasilLabMBByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]HasilLaboratorium, int, error)
	DetailHasilLabMB(ctx context.Context, noRawat string, kodeTindakan string, tanggalPeriksa string, jamPeriksa string) (*HasilLaboratorium, error)

	SimpanPermintaanLabPK(ctx context.Context, noRawat string, kodeDokter string, statusLanjut shared.StatusLanjut, req SimpanPermintaanLabPKRequest, kodeTindakanList []string, templateMap map[string][]int) (string, error)
	DaftarPermintaanLabPK(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) ([]PermintaanLabPK, error)
	DaftarPermintaanLabPKByRM(ctx context.Context, noRkmMedis string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]PermintaanLabPK, int, error)
	DetailPermintaanLabPK(ctx context.Context, noPermintaan string) (*DetailPermintaanLabPK, error)
	UpdatePermintaanLabPK(ctx context.Context, noPermintaan string, req SimpanPermintaanLabPKRequest, kodeTindakanList []string, templateMap map[string][]int) error
	HapusPermintaanLabPK(ctx context.Context, noPermintaan string) error

	SimpanPermintaanLabPA(ctx context.Context, noRawat string, kodeDokter string, statusLanjut shared.StatusLanjut, req SimpanPermintaanLabPARequest, kodeTindakanList []string) (string, error)
	DaftarPermintaanLabPA(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) ([]PermintaanLabPA, error)
	DaftarPermintaanLabPAByRM(ctx context.Context, noRkmMedis string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]PermintaanLabPA, int, error)
	DetailPermintaanLabPA(ctx context.Context, noPermintaan string) (*DetailPermintaanLabPA, error)
	UpdatePermintaanLabPA(ctx context.Context, noPermintaan string, req SimpanPermintaanLabPARequest, kodeTindakanList []string) error
	HapusPermintaanLabPA(ctx context.Context, noPermintaan string) error

	SimpanPermintaanLabMB(ctx context.Context, noRawat string, kodeDokter string, statusLanjut shared.StatusLanjut, req SimpanPermintaanLabMBRequest, kodeTindakanList []string, templateMap map[string][]int) (string, error)
	DaftarPermintaanLabMB(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) ([]PermintaanLabMB, error)
	DaftarPermintaanLabMBByRM(ctx context.Context, noRkmMedis string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]PermintaanLabMB, int, error)
	DetailPermintaanLabMB(ctx context.Context, noPermintaan string) (*DetailPermintaanLabMB, error)
	UpdatePermintaanLabMB(ctx context.Context, noPermintaan string, req SimpanPermintaanLabMBRequest, kodeTindakanList []string, templateMap map[string][]int) error
	HapusPermintaanLabMB(ctx context.Context, noPermintaan string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}
