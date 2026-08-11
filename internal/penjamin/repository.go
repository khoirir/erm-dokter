package penjamin

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	DaftarPenjamin(ctx context.Context) ([]Penjamin, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) DaftarPenjamin(ctx context.Context) ([]Penjamin, error) {
	query := `SELECT kd_pj, png_jawab FROM penjab WHERE status = '1' ORDER BY png_jawab ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("gagal query daftar penjamin: %w", err)
	}
	defer rows.Close()

	var penjamin []Penjamin
	for rows.Next() {
		var pj Penjamin
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
