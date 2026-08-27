package master

import (
	"context"
	"database/sql"
	"fmt"

	"erm-dokter/internal/shared"
)

type Repository interface {
	DaftarPenjamin(ctx context.Context) ([]Penjamin, error)
	DaftarDepo(ctx context.Context) ([]Depo, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) DaftarPenjamin(ctx context.Context) ([]Penjamin, error) {
	query := `SELECT kd_pj AS kode, png_jawab AS nama FROM penjab WHERE status = '1' ORDER BY png_jawab ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("gagal query daftar penjamin: %w", err)
	}
	defer rows.Close()

	var list []Penjamin
	for rows.Next() {
		var p Penjamin
		if err := rows.Scan(&p.Kode, &p.Nama); err != nil {
			return nil, fmt.Errorf("gagal scan data penjamin: %w", err)
		}
		list = append(list, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error saat iterasi penjamin: %w", err)
	}

	return list, nil
}

func (r *repository) DaftarDepo(ctx context.Context) ([]Depo, error) {
	query := fmt.Sprintf(`SELECT kd_bangsal AS kode, nm_bangsal AS nama 
		FROM bangsal 
		WHERE %s 
		ORDER BY nm_bangsal DESC`, shared.InClauseDepoFarmasi("kd_bangsal"))
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("gagal query daftar depo: %w", err)
	}
	defer rows.Close()

	var list []Depo
	for rows.Next() {
		var d Depo
		if err := rows.Scan(&d.Kode, &d.Nama); err != nil {
			return nil, fmt.Errorf("gagal scan data depo: %w", err)
		}
		list = append(list, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error saat iterasi depo: %w", err)
	}

	return list, nil
}
