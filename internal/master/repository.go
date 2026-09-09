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
	DaftarPoliklinik(ctx context.Context) ([]Poliklinik, error)
	DaftarBangsal(ctx context.Context) ([]Bangsal, error)
	DaftarKelas(ctx context.Context) ([]KelasKamar, error)
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
		return nil, err
	}
	defer rows.Close()

	list := make([]Penjamin, 0)
	for rows.Next() {
		var p Penjamin
		if err := rows.Scan(&p.Kode, &p.Nama); err != nil {
			return nil, err
		}
		list = append(list, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *repository) DaftarDepo(ctx context.Context) ([]Depo, error) {
	query := fmt.Sprintf(`SELECT kd_bangsal AS kode, nm_bangsal AS nama 
		FROM bangsal 
		WHERE %s 
		ORDER BY nm_bangsal ASC`, shared.InClauseDepoFarmasi("kd_bangsal"))
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]Depo, 0)
	for rows.Next() {
		var d Depo
		if err := rows.Scan(&d.Kode, &d.Nama); err != nil {
			return nil, err
		}
		list = append(list, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *repository) DaftarPoliklinik(ctx context.Context) ([]Poliklinik, error) {
	query := `SELECT kd_poli AS kode, nm_poli AS nama FROM poliklinik WHERE status = '1' ORDER BY nm_poli ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]Poliklinik, 0)
	for rows.Next() {
		var p Poliklinik
		if err := rows.Scan(&p.Kode, &p.Nama); err != nil {
			return nil, err
		}
		list = append(list, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *repository) DaftarBangsal(ctx context.Context) ([]Bangsal, error) {
	query := `
		SELECT DISTINCT b.kd_bangsal AS kode, b.nm_bangsal AS nama 
		FROM bangsal b 
		INNER JOIN kamar k ON b.kd_bangsal = k.kd_bangsal 
		WHERE b.status = '1' AND k.statusdata = '1' 
		ORDER BY b.nm_bangsal ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]Bangsal, 0)
	for rows.Next() {
		var b Bangsal
		if err := rows.Scan(&b.Kode, &b.Nama); err != nil {
			return nil, err
		}
		list = append(list, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *repository) DaftarKelas(ctx context.Context) ([]KelasKamar, error) {
	query := `
		SELECT DISTINCT k.kelas AS kode, k.kelas AS nama 
		FROM kamar k 
		INNER JOIN bangsal b ON k.kd_bangsal = b.kd_bangsal 
		WHERE k.statusdata = '1' AND b.status = '1' AND k.kelas <> '' 
		ORDER BY k.kelas ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]KelasKamar, 0)
	for rows.Next() {
		var k KelasKamar
		if err := rows.Scan(&k.Kode, &k.Nama); err != nil {
			return nil, err
		}
		list = append(list, k)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}
