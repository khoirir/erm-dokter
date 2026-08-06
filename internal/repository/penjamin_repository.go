package repository

import (
	"context"
	"database/sql"
	"fmt"

	"erm-dokter/internal/domain"
)

type penjaminRepository struct {
	db *sql.DB
}

func NewPenjaminRepository(db *sql.DB) domain.PenjaminRepository {
	return &penjaminRepository{
		db: db,
	}
}

func (r *penjaminRepository) DaftarPenjamin(ctx context.Context) ([]domain.Penjamin, error) {
	query := `SELECT kd_pj, png_jawab FROM penjab WHERE status = '1' ORDER BY png_jawab ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("gagal query daftar penjamin: %w", err)
	}
	defer rows.Close()

	var penjamin []domain.Penjamin
	for rows.Next() {
		var pj domain.Penjamin
		if err := rows.Scan(&pj.KodePenjamin, &pj.Nama); err != nil {
			return nil, fmt.Errorf("gagal scan data penjamin: %w", err)
		}
		penjamin = append(penjamin, pj)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error saat iterasi data penjamin: %w", err)
	}

	return penjamin, nil
}
