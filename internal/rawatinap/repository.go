package rawatinap

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

type Repository interface {
	CekStatusKamarInap(ctx context.Context, noRawat string) (isKamarAktif bool, hasRecordKamar bool, err error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CekStatusKamarInap(ctx context.Context, noRawat string) (bool, bool, error) {
	query := `
		SELECT stts_pulang, tgl_keluar, jam_keluar 
		FROM kamar_inap 
		WHERE no_rawat = ? 
		ORDER BY tgl_masuk DESC, jam_masuk DESC 
		LIMIT 1
	`
	var sttsPulang, tglKeluar, jamKeluar string
	err := r.db.QueryRowContext(ctx, query, noRawat).Scan(&sttsPulang, &tglKeluar, &jamKeluar)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, false, nil
		}
		return false, false, err
	}

	isBelumPulangStts := sttsPulang == "-" || strings.TrimSpace(sttsPulang) == ""
	isBelumKeluarTgl := tglKeluar == "0000-00-00" || strings.TrimSpace(tglKeluar) == ""
	isBelumKeluarJam := jamKeluar == "00:00:00" || strings.TrimSpace(jamKeluar) == ""

	if isBelumPulangStts && isBelumKeluarTgl && isBelumKeluarJam {
		return true, true, nil
	}

	return false, true, nil
}
