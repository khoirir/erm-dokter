package laboratorium

import (
	"context"
	"database/sql"

	"erm-dokter/internal/shared"
)

type Repository interface {
	DaftarHasilLabPK(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]HasilLaboratoriumDB, int, error)
	DaftarHasilLabPKByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]HasilLaboratoriumDB, int, error)
	DetailHasilLabPK(ctx context.Context, noRawat string, kodeTindakan string, tanggalPeriksa string, jamPeriksa string) (*HasilLaboratoriumDB, error)

	DaftarHasilLabPA(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]HasilLaboratoriumDB, int, error)
	DaftarHasilLabPAByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]HasilLaboratoriumDB, int, error)
	DetailHasilLabPA(ctx context.Context, noRawat string, kodeTindakan string, tanggalPeriksa string, jamPeriksa string) (*HasilLaboratoriumDB, error)

	DaftarHasilLabMB(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]HasilLaboratoriumDB, int, error)
	DaftarHasilLabMBByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]HasilLaboratoriumDB, int, error)
	DetailHasilLabMB(ctx context.Context, noRawat string, kodeTindakan string, tanggalPeriksa string, jamPeriksa string) (*HasilLaboratoriumDB, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}
