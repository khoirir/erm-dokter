package pemeriksaan

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"erm-dokter/internal/shared"
)

type Repository interface {
	DaftarPemeriksaan(ctx context.Context, listNoRawat []string, statusLanjut shared.StatusLanjut, filter FilterDaftarPemeriksaan) ([]Pemeriksaan, int, error)
	DetailPemeriksaan(ctx context.Context, idPemeriksaan IdPemeriksaan, statusLanjut shared.StatusLanjut) (*Pemeriksaan, error)
	SimpanPemeriksaan(ctx context.Context, kodeDokter string, statusLanjut shared.StatusLanjut, req SimpanPemeriksaanRequest) error
	HapusPemeriksaan(ctx context.Context, id IdPemeriksaan, statusLanjut shared.StatusLanjut) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

type scanner interface {
	Scan(dest ...any) error
}

func scanPemeriksaan(s scanner) (*Pemeriksaan, error) {
	var p Pemeriksaan
	err := s.Scan(
		&p.NoRawat,
		&p.TanggalPemeriksaan,
		&p.JamPemeriksaan,
		&p.SuhuTubuh,
		&p.Tensi,
		&p.Nadi,
		&p.Respirasi,
		&p.TinggiBadan,
		&p.BeratBadan,
		&p.SpO2,
		&p.Gcs,
		&p.Kesadaran,
		&p.Keluhan,
		&p.Pemeriksaan,
		&p.Alergi,
		&p.LingkarPerut,
		&p.RencanaTindakLanjut,
		&p.Penilaian,
		&p.Instruksi,
		&p.Evaluasi,
		&p.KodeDokterPetugas,
		&p.NamaDokterPetugas,
		&p.StatusLanjut,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

const selectPemeriksaanRalan = `
	SELECT 
		pr.no_rawat,
		DATE_FORMAT(pr.tgl_perawatan, '%%Y-%%m-%%d') AS tanggal_pemeriksaan,
		pr.jam_rawat AS jam_pemeriksaan,
		pr.suhu_tubuh,
		pr.tensi,
		pr.nadi,
		pr.respirasi,
		pr.tinggi AS tinggi_badan,
		pr.berat AS berat_badan,
		pr.spo2,
		pr.gcs,
		pr.kesadaran,
		pr.keluhan,
		pr.pemeriksaan,
		pr.alergi,
		pr.lingkar_perut,
		pr.rtl AS rencana_tindak_lanjut,
		pr.penilaian,
		pr.instruksi,
		pr.evaluasi,
		pr.nip AS kode_dokter_petugas,
		COALESCE(d.nm_dokter, p.nama, pr.nip) AS nama_dokter_petugas,
		'Ralan' AS status_lanjut
	FROM pemeriksaan_ralan pr
	LEFT JOIN dokter d ON pr.nip = d.kd_dokter
	LEFT JOIN petugas p ON pr.nip = p.nip
	WHERE pr.no_rawat IN (%s)
`

const selectPemeriksaanRanap = `
	SELECT 
		pr.no_rawat,
		DATE_FORMAT(pr.tgl_perawatan, '%%Y-%%m-%%d') AS tanggal_pemeriksaan,
		pr.jam_rawat AS jam_pemeriksaan,
		pr.suhu_tubuh,
		pr.tensi,
		pr.nadi,
		pr.respirasi,
		pr.tinggi AS tinggi_badan,
		pr.berat AS berat_badan,
		pr.spo2,
		pr.gcs,
		pr.kesadaran,
		pr.keluhan,
		pr.pemeriksaan,
		pr.alergi,
		'' AS lingkar_perut,
		pr.rtl AS rencana_tindak_lanjut,
		pr.penilaian,
		pr.instruksi,
		pr.evaluasi,
		pr.nip AS kode_dokter_petugas,
		COALESCE(d.nm_dokter, p.nama, pr.nip) AS nama_dokter_petugas,
		'Ranap' AS status_lanjut
	FROM pemeriksaan_ranap pr
	LEFT JOIN dokter d ON pr.nip = d.kd_dokter
	LEFT JOIN petugas p ON pr.nip = p.nip
	WHERE pr.no_rawat IN (%s)
`

const insertPemeriksaanRalan = `
	INSERT INTO pemeriksaan_ralan (
		no_rawat, tgl_perawatan, jam_rawat, suhu_tubuh, tensi, nadi, respirasi,
		tinggi, berat, spo2, gcs, kesadaran, keluhan, pemeriksaan, alergi,
		lingkar_perut, rtl, penilaian, instruksi, evaluasi, nip
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

const insertPemeriksaanRanap = `
	INSERT INTO pemeriksaan_ranap (
		no_rawat, tgl_perawatan, jam_rawat, suhu_tubuh, tensi, nadi, respirasi,
		tinggi, berat, spo2, gcs, kesadaran, keluhan, pemeriksaan, alergi,
		rtl, penilaian, instruksi, evaluasi, nip
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

func createInPlaceholders(count int) string {
	if count <= 0 {
		return "?"
	}
	placeholders := make([]string, count)
	for i := range count {
		placeholders[i] = "?"
	}
	return strings.Join(placeholders, ",")
}

func buildBaseQuery(listNoRawat []string, statusLanjut shared.StatusLanjut, filter FilterDaftarPemeriksaan) (string, []any) {
	var tglAwal, tglAkhir string
	useTglFilter := false

	tglParts := strings.Split(filter.Tanggal, ",")
	if len(tglParts) == 2 {
		tglAwal = strings.TrimSpace(tglParts[0])
		tglAkhir = strings.TrimSpace(tglParts[1])
		useTglFilter = true
	}

	inClause := createInPlaceholders(len(listNoRawat))
	sqlRalan := fmt.Sprintf(selectPemeriksaanRalan, inClause)
	sqlRanap := fmt.Sprintf(selectPemeriksaanRanap, inClause)

	queryWithFilter := func(baseSQL string, args *[]any) string {
		q := baseSQL
		for _, nr := range listNoRawat {
			*args = append(*args, nr)
		}
		if useTglFilter {
			q += " AND pr.tgl_perawatan BETWEEN ? AND ?"
			*args = append(*args, tglAwal, tglAkhir)
		}
		return q
	}

	var args []any
	var baseQuery string

	switch statusLanjut {
	case shared.StatusLanjutRawatJalan:
		baseQuery = queryWithFilter(sqlRalan, &args)
	case shared.StatusLanjutRawatInap:
		baseQuery = queryWithFilter(sqlRanap, &args)
	default:
		q1 := queryWithFilter(sqlRalan, &args)
		q2 := queryWithFilter(sqlRanap, &args)
		baseQuery = fmt.Sprintf("%s UNION ALL %s", q1, q2)
	}

	return baseQuery, args
}

func (r *repository) DaftarPemeriksaan(ctx context.Context, listNoRawat []string, statusLanjut shared.StatusLanjut, filter FilterDaftarPemeriksaan) ([]Pemeriksaan, int, error) {
	baseQuery, baseArgs := buildBaseQuery(listNoRawat, statusLanjut, filter)

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS t", baseQuery)
	var totalData int
	if err := r.db.QueryRowContext(ctx, countQuery, baseArgs...).Scan(&totalData); err != nil || totalData == 0 {
		return []Pemeriksaan{}, 0, err
	}

	dataQuery := fmt.Sprintf("SELECT * FROM (%s) AS t ORDER BY t.tanggal_pemeriksaan DESC, t.jam_pemeriksaan DESC LIMIT ? OFFSET ?", baseQuery)

	dataArgs := append(baseArgs, filter.Batas, filter.Offset())

	rows, err := r.db.QueryContext(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal query data pemeriksaan: %w", err)
	}
	defer rows.Close()

	var daftarPemeriksaan []Pemeriksaan
	for rows.Next() {
		pemeriksaan, err := scanPemeriksaan(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("gagal scan data pemeriksaan: %w", err)
		}
		daftarPemeriksaan = append(daftarPemeriksaan, *pemeriksaan)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error saat iterasi data pemeriksaan: %w", err)
	}

	return daftarPemeriksaan, totalData, nil
}

func (r *repository) DetailPemeriksaan(ctx context.Context, id IdPemeriksaan, statusLanjut shared.StatusLanjut) (*Pemeriksaan, error) {
	var query string
	var args []interface{}

	switch statusLanjut {
	case shared.StatusLanjutRawatJalan:
		query = fmt.Sprintf(selectPemeriksaanRalan, "?") + " AND pr.tgl_perawatan = ? AND pr.jam_rawat = ? LIMIT 1"
		args = append(args, id.NoRawat, id.TanggalPemeriksaan, id.JamPemeriksaan)

	case shared.StatusLanjutRawatInap:
		query = fmt.Sprintf(selectPemeriksaanRanap, "?") + " AND pr.tgl_perawatan = ? AND pr.jam_rawat = ? LIMIT 1"
		args = append(args, id.NoRawat, id.TanggalPemeriksaan, id.JamPemeriksaan)

	default:
		return nil, fmt.Errorf("status lanjut tidak valid (Ralan atau Ranap)")
	}

	row := r.db.QueryRowContext(ctx, query, args...)
	pemeriksaan, err := scanPemeriksaan(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("gagal query detail pemeriksaan: %w", err)
	}

	return pemeriksaan, nil
}

func (r *repository) SimpanPemeriksaan(ctx context.Context, kodeDokter string, statusLanjut shared.StatusLanjut, req SimpanPemeriksaanRequest) error {
	var query string
	var args []any

	switch statusLanjut {
	case shared.StatusLanjutRawatJalan:
		query = insertPemeriksaanRalan
		args = []any{
			req.NoRawat, req.TanggalPemeriksaan, req.JamPemeriksaan, req.SuhuTubuh, req.Tensi, req.Nadi, req.Respirasi,
			req.TinggiBadan, req.BeratBadan, req.SpO2, req.Gcs, string(req.Kesadaran), req.Keluhan, req.Pemeriksaan, req.Alergi,
			req.LingkarPerut, req.RencanaTindakLanjut, req.Penilaian, req.Instruksi, req.Evaluasi, kodeDokter,
		}
	case shared.StatusLanjutRawatInap:
		query = insertPemeriksaanRanap
		args = []any{
			req.NoRawat, req.TanggalPemeriksaan, req.JamPemeriksaan, req.SuhuTubuh, req.Tensi, req.Nadi, req.Respirasi,
			req.TinggiBadan, req.BeratBadan, req.SpO2, req.Gcs, string(req.Kesadaran), req.Keluhan, req.Pemeriksaan, req.Alergi,
			req.RencanaTindakLanjut, req.Penilaian, req.Instruksi, req.Evaluasi, kodeDokter,
		}
	default:
		return fmt.Errorf("status lanjut tidak valid (harus Ralan atau Ranap)")
	}

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("gagal menyimpan data pemeriksaan: %w", err)
	}

	return nil
}

const deletePemeriksaanRalan = `
	DELETE FROM pemeriksaan_ralan 
	WHERE no_rawat = ? AND tgl_perawatan = ? AND jam_rawat = ?
`

const deletePemeriksaanRanap = `
	DELETE FROM pemeriksaan_ranap 
	WHERE no_rawat = ? AND tgl_perawatan = ? AND jam_rawat = ?
`

func (r *repository) HapusPemeriksaan(ctx context.Context, id IdPemeriksaan, statusLanjut shared.StatusLanjut) error {
	var query string
	switch statusLanjut {
	case shared.StatusLanjutRawatJalan:
		query = deletePemeriksaanRalan
	case shared.StatusLanjutRawatInap:
		query = deletePemeriksaanRanap
	default:
		return fmt.Errorf("status lanjut tidak valid (harus Ralan atau Ranap)")
	}

	_, err := r.db.ExecContext(ctx, query, id.NoRawat, id.TanggalPemeriksaan, id.JamPemeriksaan)
	if err != nil {
		return fmt.Errorf("gagal menghapus data pemeriksaan: %w", err)
	}

	return nil
}



