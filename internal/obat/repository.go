package obat

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type Repository interface {
	DaftarObat(ctx context.Context, filter FilterDaftarObat) ([]Obat, int64, error)
	DetailObat(ctx context.Context, kodeObat string) (*Obat, error)
	DaftarJenis(ctx context.Context) ([]JenisObat, error)
	DaftarGolongan(ctx context.Context) ([]GolonganObat, error)
	DaftarKategori(ctx context.Context) ([]KategoriObat, error)
	CekKeberadaanObat(ctx context.Context, listKodeObat []string) (map[string]bool, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

var orderByMapping = map[string]func(string) string{
	"nama_obat":     func(dir string) string { return fmt.Sprintf("dtb.nama_brng %s", dir) },
	"stok":          func(dir string) string { return fmt.Sprintf("stok %s", dir) },
	"nama_depo":     func(dir string) string { return fmt.Sprintf("bg.nm_bangsal %s", dir) },
	"nama_jenis":    func(dir string) string { return fmt.Sprintf("jn.nama %s", dir) },
	"nama_golongan": func(dir string) string { return fmt.Sprintf("gb.nama %s", dir) },
	"nama_kategori": func(dir string) string { return fmt.Sprintf("kb.nama %s", dir) },
}

func (r *repository) DaftarObat(ctx context.Context, filter FilterDaftarObat) ([]Obat, int64, error) {

	var baseFrom string
	var selectCols string
	var conditions []string
	var args []any

	if filter.Depo != "" {
		baseFrom = `FROM databarang dtb 
			INNER JOIN gudangbarang gdb ON dtb.kode_brng = gdb.kode_brng 
			INNER JOIN bangsal bg ON gdb.kd_bangsal = bg.kd_bangsal
			INNER JOIN kategori_barang kb ON dtb.kode_kategori = kb.kode 
			INNER JOIN golongan_barang gb ON dtb.kode_golongan = gb.kode 
			INNER JOIN jenis jn ON dtb.kdjns = jn.kdjns
			WHERE dtb.status = '1'`

		selectCols = `SELECT 
			dtb.kode_brng AS kode_obat,
			dtb.nama_brng AS nama_obat,
			dtb.letak_barang AS komposisi,
			dtb.ralan AS harga,
			dtb.kode_sat AS satuan,
			gdb.stok,
			dtb.kapasitas,
			gdb.kd_bangsal AS kode_depo,
			bg.nm_bangsal AS nama_depo,
			dtb.kdjns AS kode_jenis,
			jn.nama AS nama_jenis,
			dtb.kode_golongan AS kode_golongan,
			gb.nama AS nama_golongan,
			dtb.kode_kategori AS kode_kategori,
			kb.nama AS nama_kategori `

		conditions = append(conditions, "gdb.kd_bangsal = ?")
		args = append(args, filter.Depo)
	} else {
		baseFrom = `FROM databarang dtb 
			INNER JOIN kategori_barang kb ON dtb.kode_kategori = kb.kode 
			INNER JOIN golongan_barang gb ON dtb.kode_golongan = gb.kode 
			INNER JOIN jenis jn ON dtb.kdjns = jn.kdjns
			WHERE dtb.status = '1'`

		selectCols = `SELECT 
			dtb.kode_brng AS kode_obat,
			dtb.nama_brng AS nama_obat,
			dtb.letak_barang AS komposisi,
			dtb.ralan AS harga,
			dtb.kode_sat AS satuan,
			COALESCE((SELECT SUM(gdb.stok) FROM gudangbarang gdb WHERE gdb.kode_brng = dtb.kode_brng AND ` + shared.InClauseDepoFarmasi("gdb.kd_bangsal") + `), 0) AS stok,
			dtb.kapasitas,
			'' AS kode_depo,
			'Semua Depo' AS nama_depo,
			dtb.kdjns AS kode_jenis,
			jn.nama AS nama_jenis,
			dtb.kode_golongan AS kode_golongan,
			gb.nama AS nama_golongan,
			dtb.kode_kategori AS kode_kategori,
			kb.nama AS nama_kategori `
	}

	if filter.Jenis != "" {
		conditions = append(conditions, "dtb.kdjns = ?")
		args = append(args, filter.Jenis)
	}

	if filter.Golongan != "" {
		conditions = append(conditions, "dtb.kode_golongan = ?")
		args = append(args, filter.Golongan)
	}

	if filter.Kategori != "" {
		conditions = append(conditions, "dtb.kode_kategori = ?")
		args = append(args, filter.Kategori)
	}

	if filter.Keyword != "" {
		keywordPattern := "%" + filter.Keyword + "%"
		conditions = append(conditions, "(dtb.nama_brng LIKE ? OR dtb.letak_barang LIKE ?)")
		args = append(args, keywordPattern, keywordPattern)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " AND " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) " + baseFrom + whereClause
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("gagal menghitung total data obat: %w", err)
	}

	if total == 0 {
		return []Obat{}, 0, nil
	}

	builder, exists := orderByMapping[filter.OrderBy]
	if !exists {
		builder = orderByMapping["nama_obat"]
	}
	orderClause := " ORDER BY " + builder(filter.SortOrder)

	selectQuery := selectCols + baseFrom + whereClause + orderClause + " LIMIT ? OFFSET ?"


	dataArgs := append(args, filter.Limit, filter.Offset())
	rows, err := r.db.QueryContext(ctx, selectQuery, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal query daftar obat: %w", err)
	}
	defer rows.Close()

	var listObat []Obat
	for rows.Next() {
		var o Obat
		if err := rows.Scan(
			&o.KodeObat,
			&o.NamaObat,
			&o.Komposisi,
			&o.Harga,
			&o.Satuan,
			&o.Stok,
			&o.Kapasitas,
			&o.KodeDepo,
			&o.NamaDepo,
			&o.KodeJenis,
			&o.NamaJenis,
			&o.KodeGolongan,
			&o.NamaGolongan,
			&o.KodeKategori,
			&o.NamaKategori,
		); err != nil {
			return nil, 0, fmt.Errorf("gagal scan data obat: %w", err)
		}
		listObat = append(listObat, o)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error saat iterasi data obat: %w", err)
	}

	return listObat, total, nil
}

func (r *repository) DetailObat(ctx context.Context, kodeObat string) (*Obat, error) {
	query := `SELECT 
		dtb.kode_brng AS kode_obat,
		dtb.nama_brng AS nama_obat,
		dtb.letak_barang AS komposisi,
		dtb.ralan AS harga,
		dtb.kode_sat AS satuan,
		COALESCE((SELECT SUM(gdb.stok) FROM gudangbarang gdb WHERE gdb.kode_brng = dtb.kode_brng AND ` + shared.InClauseDepoFarmasi("gdb.kd_bangsal") + `), 0) AS stok,
		dtb.kapasitas,
		dtb.kdjns AS kode_jenis,
		jn.nama AS nama_jenis,
		dtb.kode_golongan AS kode_golongan,
		gb.nama AS nama_golongan,
		dtb.kode_kategori AS kode_kategori,
		kb.nama AS nama_kategori
	FROM databarang dtb
	INNER JOIN kategori_barang kb ON dtb.kode_kategori = kb.kode 
	INNER JOIN golongan_barang gb ON dtb.kode_golongan = gb.kode 
	INNER JOIN jenis jn ON dtb.kdjns = jn.kdjns
	WHERE dtb.kode_brng = ? AND dtb.status = '1'`

	var o Obat
	err := r.db.QueryRowContext(ctx, query, kodeObat).Scan(
		&o.KodeObat,
		&o.NamaObat,
		&o.Komposisi,
		&o.Harga,
		&o.Satuan,
		&o.Stok,
		&o.Kapasitas,
		&o.KodeJenis,
		&o.NamaJenis,
		&o.KodeGolongan,
		&o.NamaGolongan,
		&o.KodeKategori,
		&o.NamaKategori,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFoundError("Data obat tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal query detail obat: %w", err)
	}

	queryStok := fmt.Sprintf(`SELECT 
		gdb.kd_bangsal AS kode_depo,
		bg.nm_bangsal AS nama_depo,
		gdb.stok
	FROM gudangbarang gdb
	INNER JOIN bangsal bg ON gdb.kd_bangsal = bg.kd_bangsal
	WHERE gdb.kode_brng = ? AND %s
	ORDER BY bg.nm_bangsal DESC`, shared.InClauseDepoFarmasi("gdb.kd_bangsal"))

	rows, err := r.db.QueryContext(ctx, queryStok, kodeObat)
	if err != nil {
		return nil, fmt.Errorf("gagal query stok depo obat: %w", err)
	}
	defer rows.Close()

	o.StokDepo = make([]StokDepo, 0)
	for rows.Next() {
		var s StokDepo
		if err := rows.Scan(&s.KodeDepo, &s.NamaDepo, &s.Stok); err != nil {
			return nil, fmt.Errorf("gagal scan stok depo: %w", err)
		}
		o.StokDepo = append(o.StokDepo, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error saat iterasi stok depo: %w", err)
	}

	return &o, nil
}

func (r *repository) DaftarJenis(ctx context.Context) ([]JenisObat, error) {
	query := `SELECT kdjns AS kode, nama FROM jenis ORDER BY nama ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("gagal query daftar jenis obat: %w", err)
	}
	defer rows.Close()

	var list []JenisObat
	for rows.Next() {
		var j JenisObat
		if err := rows.Scan(&j.Kode, &j.Nama); err != nil {
			return nil, fmt.Errorf("gagal scan data jenis obat: %w", err)
		}
		list = append(list, j)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error saat iterasi jenis obat: %w", err)
	}

	return list, nil
}

func (r *repository) DaftarGolongan(ctx context.Context) ([]GolonganObat, error) {
	query := `SELECT kode, nama FROM golongan_barang ORDER BY nama ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("gagal query daftar golongan obat: %w", err)
	}
	defer rows.Close()

	var list []GolonganObat
	for rows.Next() {
		var g GolonganObat
		if err := rows.Scan(&g.Kode, &g.Nama); err != nil {
			return nil, fmt.Errorf("gagal scan data golongan obat: %w", err)
		}
		list = append(list, g)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error saat iterasi golongan obat: %w", err)
	}

	return list, nil
}

func (r *repository) DaftarKategori(ctx context.Context) ([]KategoriObat, error) {
	query := `SELECT kode, nama FROM kategori_barang ORDER BY nama ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("gagal query daftar kategori obat: %w", err)
	}
	defer rows.Close()

	var list []KategoriObat
	for rows.Next() {
		var k KategoriObat
		if err := rows.Scan(&k.Kode, &k.Nama); err != nil {
			return nil, fmt.Errorf("gagal scan data kategori obat: %w", err)
		}
		list = append(list, k)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error saat iterasi kategori obat: %w", err)
	}

	return list, nil
}

func (r *repository) CekKeberadaanObat(ctx context.Context, listKodeObat []string) (map[string]bool, error) {
	if len(listKodeObat) == 0 {
		return map[string]bool{}, nil
	}

	uniqueCodes := make([]string, 0, len(listKodeObat))
	seen := make(map[string]bool)
	for _, code := range listKodeObat {
		trimmed := strings.TrimSpace(code)
		if trimmed != "" && !seen[trimmed] {
			seen[trimmed] = true
			uniqueCodes = append(uniqueCodes, trimmed)
		}
	}

	if len(uniqueCodes) == 0 {
		return map[string]bool{}, nil
	}

	inPlaceholders := shared.CreateInPlaceholders(len(uniqueCodes))
	args := make([]any, len(uniqueCodes))
	for i, c := range uniqueCodes {
		args[i] = c
	}

	query := fmt.Sprintf("SELECT kode_brng FROM databarang WHERE kode_brng IN (%s)", inPlaceholders)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	foundMap := make(map[string]bool)
	for rows.Next() {
		var kode string
		if err := rows.Scan(&kode); err != nil {
			return nil, err
		}
		foundMap[kode] = true
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return foundMap, nil
}

