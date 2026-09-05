package radiologi

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"erm-dokter/internal/shared"
)

type Repository interface {
	DaftarHasilRadiologiKunjungan(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter FilterRiwayatRadiologi) ([]HasilRadiologi, int, error)
	DaftarHasilRadiologiPasien(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter FilterRiwayatRadiologi) ([]HasilRadiologi, int, error)
	DetailHasilRadiologi(ctx context.Context, idHasil IdHasilRadiologi, statusLanjut shared.StatusLanjut) (*HasilRadiologi, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) DaftarHasilRadiologiKunjungan(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter FilterRiwayatRadiologi) ([]HasilRadiologi, int, error) {
	return r.queryRiwayatRadiologi(ctx, "pr.no_rawat = ?", noRawat, statusLanjut, filter)
}

func (r *repository) DaftarHasilRadiologiPasien(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter FilterRiwayatRadiologi) ([]HasilRadiologi, int, error) {
	return r.queryRiwayatRadiologi(ctx, "rp.no_rkm_medis = ?", noRM, statusLanjut, filter)
}

func (r *repository) queryRiwayatRadiologi(ctx context.Context, whereClause string, paramValue string, statusLanjut shared.StatusLanjut, filter FilterRiwayatRadiologi) ([]HasilRadiologi, int, error) {
	var conditions []string
	var args []any

	conditions = append(conditions, whereClause)
	args = append(args, paramValue)

	switch statusLanjut {
	case shared.StatusLanjutRawatJalan:
		conditions = append(conditions, "pr.status = 'Ralan'")
	case shared.StatusLanjutRawatInap:
		conditions = append(conditions, "pr.status = 'Ranap'")
	}

	if filter.Tanggal != "" {
		tanggalParts := strings.Split(filter.Tanggal, ",")
		if len(tanggalParts) == 2 {
			conditions = append(conditions, "pr.tgl_periksa BETWEEN ? AND ?")
			args = append(args, strings.TrimSpace(tanggalParts[0]), strings.TrimSpace(tanggalParts[1]))
		} else {
			conditions = append(conditions, "pr.tgl_periksa = ?")
			args = append(args, strings.TrimSpace(tanggalParts[0]))
		}
	}

	whereSQL := "WHERE " + strings.Join(conditions, " AND ")

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM periksa_radiologi pr
		INNER JOIN reg_periksa rp ON rp.no_rawat = pr.no_rawat
		%s
	`, whereSQL)

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return make([]HasilRadiologi, 0), 0, nil
	}

	dataQuery := fmt.Sprintf(`
		SELECT 
			pr.no_rawat,
			pr.kd_jenis_prw AS kode_tindakan,
			COALESCE(jpr.nm_perawatan, pr.kd_jenis_prw) AS nama_tindakan,
			pr.status,
			DATE_FORMAT(pr.tgl_periksa, '%%Y-%%m-%%d') AS tanggal_periksa,
			pr.jam AS jam_periksa,
			pr.dokter_perujuk AS kd_dokter_perujuk,
			COALESCE(d_perujuk.nm_dokter, '-') AS nm_dokter_perujuk,
			pr.kd_dokter AS kd_dokter_radiologi,
			COALESCE(d_rad.nm_dokter, '-') AS nm_dokter_radiologi,
			pr.nip AS nip_petugas,
			COALESCE(p.nama, '-') AS nama_petugas,
			COALESCE(hr.hasil, '') AS hasil
		FROM periksa_radiologi pr
		INNER JOIN reg_periksa rp ON rp.no_rawat = pr.no_rawat
		LEFT JOIN jns_perawatan_radiologi jpr ON jpr.kd_jenis_prw = pr.kd_jenis_prw
		LEFT JOIN dokter d_perujuk ON d_perujuk.kd_dokter = pr.dokter_perujuk
		LEFT JOIN dokter d_rad ON d_rad.kd_dokter = pr.kd_dokter
		LEFT JOIN petugas p ON p.nip = pr.nip
		LEFT JOIN hasil_radiologi hr ON hr.no_rawat = pr.no_rawat AND hr.tgl_periksa = pr.tgl_periksa AND hr.jam = pr.jam
		%s
		ORDER BY pr.tgl_periksa DESC, pr.jam DESC
		LIMIT ? OFFSET ?
	`, whereSQL)

	argsWithPaging := append(args, filter.Limit, filter.Offset())
	rows, err := r.db.QueryContext(ctx, dataQuery, argsWithPaging...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []HasilRadiologi
	for rows.Next() {
		var item HasilRadiologi
		err = rows.Scan(
			&item.NoRawat,
			&item.KodeTindakan,
			&item.NamaTindakan,
			&item.Status,
			&item.TanggalPeriksa,
			&item.JamPeriksa,
			&item.DokterPerujuk.KodeDokter,
			&item.DokterPerujuk.NamaDokter,
			&item.DokterRadiologi.KodeDokter,
			&item.DokterRadiologi.NamaDokter,
			&item.Petugas.Nip,
			&item.Petugas.Nama,
			&item.Hasil,
		)
		if err != nil {
			return nil, 0, err
		}

		list = append(list, item)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	if err = r.attachGambarPACS(ctx, list); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *repository) DetailHasilRadiologi(ctx context.Context, idHasil IdHasilRadiologi, statusLanjut shared.StatusLanjut) (*HasilRadiologi, error) {
	query := `
		SELECT 
			pr.no_rawat,
			pr.kd_jenis_prw AS kode_tindakan,
			COALESCE(jpr.nm_perawatan, pr.kd_jenis_prw) AS nama_tindakan,
			pr.status,
			DATE_FORMAT(pr.tgl_periksa, '%Y-%m-%d') AS tanggal_periksa,
			pr.jam AS jam_periksa,
			pr.dokter_perujuk AS kd_dokter_perujuk,
			COALESCE(d_perujuk.nm_dokter, '-') AS nm_dokter_perujuk,
			pr.kd_dokter AS kd_dokter_radiologi,
			COALESCE(d_rad.nm_dokter, '-') AS nm_dokter_radiologi,
			pr.nip AS nip_petugas,
			COALESCE(p.nama, '-') AS nama_petugas,
			COALESCE(hr.hasil, '') AS hasil
		FROM periksa_radiologi pr
		INNER JOIN reg_periksa rp ON rp.no_rawat = pr.no_rawat
		LEFT JOIN jns_perawatan_radiologi jpr ON jpr.kd_jenis_prw = pr.kd_jenis_prw
		LEFT JOIN dokter d_perujuk ON d_perujuk.kd_dokter = pr.dokter_perujuk
		LEFT JOIN dokter d_rad ON d_rad.kd_dokter = pr.kd_dokter
		LEFT JOIN petugas p ON p.nip = pr.nip
		LEFT JOIN hasil_radiologi hr ON hr.no_rawat = pr.no_rawat AND hr.tgl_periksa = pr.tgl_periksa AND hr.jam = pr.jam
		WHERE pr.no_rawat = ? AND pr.kd_jenis_prw = ? AND pr.tgl_periksa = ? AND pr.jam = ?
	`
	args := []any{idHasil.NoRawat, idHasil.KodeTindakan, idHasil.TanggalPeriksa, idHasil.JamPeriksa}

	switch statusLanjut {
	case shared.StatusLanjutRawatJalan:
		query += " AND pr.status = 'Ralan'"
	case shared.StatusLanjutRawatInap:
		query += " AND pr.status = 'Ranap'"
	}

	var item HasilRadiologi
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&item.NoRawat,
		&item.KodeTindakan,
		&item.NamaTindakan,
		&item.Status,
		&item.TanggalPeriksa,
		&item.JamPeriksa,
		&item.DokterPerujuk.KodeDokter,
		&item.DokterPerujuk.NamaDokter,
		&item.DokterRadiologi.KodeDokter,
		&item.DokterRadiologi.NamaDokter,
		&item.Petugas.Nip,
		&item.Petugas.Nama,
		&item.Hasil,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	items := []HasilRadiologi{item}
	if err := r.attachGambarPACS(ctx, items); err != nil {
		return nil, err
	}

	return &items[0], nil
}

func (r *repository) attachGambarPACS(ctx context.Context, list []HasilRadiologi) error {
	if len(list) == 0 {
		return nil
	}

	var placeholders []string
	var args []any
	visited := make(map[string]bool)

	for _, item := range list {
		key := item.NoRawat + "|" + item.TanggalPeriksa + "|" + item.JamPeriksa
		if visited[key] {
			continue
		}
		visited[key] = true
		placeholders = append(placeholders, "(?, ?, ?)")
		args = append(args, item.NoRawat, item.TanggalPeriksa, item.JamPeriksa)
	}

	if len(placeholders) == 0 {
		return nil
	}

	inClause := strings.Join(placeholders, ", ")

	queryGambar := fmt.Sprintf(`
		SELECT 
			no_rawat,
			DATE_FORMAT(tgl_periksa, '%%Y-%%m-%%d') AS tanggal_periksa,
			jam AS jam_periksa,
			lokasi_gambar
		FROM gambar_radiologi_pacs
		WHERE (no_rawat, tgl_periksa, jam) IN (%s)
		ORDER BY lokasi_gambar ASC
	`, inClause)

	rowsGambar, err := r.db.QueryContext(ctx, queryGambar, args...)
	if err != nil {
		return err
	}
	defer rowsGambar.Close()

	mapGambar := make(map[string][]string)
	for rowsGambar.Next() {
		var noRawat, tanggalPeriksa, jamPeriksa, url string
		if err := rowsGambar.Scan(&noRawat, &tanggalPeriksa, &jamPeriksa, &url); err != nil {
			return err
		}
		key := noRawat + "|" + tanggalPeriksa + "|" + jamPeriksa
		mapGambar[key] = append(mapGambar[key], url)
	}
	if err := rowsGambar.Err(); err != nil {
		return err
	}

	for i := range list {
		key := list[i].NoRawat + "|" + list[i].TanggalPeriksa + "|" + list[i].JamPeriksa
		if urls, exists := mapGambar[key]; exists {
			list[i].GambarPACS = urls
		} else {
			list[i].GambarPACS = make([]string, 0)
		}
	}

	return nil
}
