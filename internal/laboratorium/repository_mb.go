package laboratorium

import (
	"context"
	"fmt"
	"strings"

	"erm-dokter/internal/shared"
)

func (r *repository) DaftarHasilLabMB(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]HasilLaboratoriumDB, int, error) {
	return r.queryRiwayatLabMB(ctx, "pl.no_rawat = ?", noRawat, statusLanjut, filter)
}

func (r *repository) DaftarHasilLabMBByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]HasilLaboratoriumDB, int, error) {
	return r.queryRiwayatLabMB(ctx, "rp.no_rkm_medis = ?", noRM, statusLanjut, filter)
}

func (r *repository) queryRiwayatLabMB(ctx context.Context, whereClause string, paramValue string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]HasilLaboratoriumDB, int, error) {
	var conditions []string
	var args []any

	conditions = append(conditions, whereClause)
	args = append(args, paramValue)

	conditions = append(conditions, "pl.kategori = 'MB'")

	switch statusLanjut {
	case shared.StatusLanjutRawatJalan:
		conditions = append(conditions, "pl.status = 'Ralan'")
	case shared.StatusLanjutRawatInap:
		conditions = append(conditions, "pl.status = 'Ranap'")
	}

	if filter.Tanggal != "" {
		tglParts := strings.Split(filter.Tanggal, ",")
		if len(tglParts) == 2 {
			conditions = append(conditions, "pl.tgl_periksa BETWEEN ? AND ?")
			args = append(args, strings.TrimSpace(tglParts[0]), strings.TrimSpace(tglParts[1]))
		} else {
			conditions = append(conditions, "pl.tgl_periksa = ?")
			args = append(args, strings.TrimSpace(tglParts[0]))
		}
	}

	whereSQL := "WHERE " + strings.Join(conditions, " AND ")

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM periksa_lab pl
		INNER JOIN reg_periksa rp ON rp.no_rawat = pl.no_rawat
		%s
	`, whereSQL)

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []HasilLaboratoriumDB{}, 0, nil
	}

	filter.Sanitize()
	dataQuery := fmt.Sprintf(`
		SELECT 
			pl.no_rawat,
			pl.kd_jenis_prw AS kode_tindakan,
			COALESCE(jpl.nm_perawatan, pl.kd_jenis_prw) AS nama_tindakan,
			pl.kategori,
			pl.status,
			pl.tgl_periksa,
			pl.jam AS jam_periksa,
			pl.dokter_perujuk AS kd_dokter_perujuk,
			COALESCE(d_perujuk.nm_dokter, '-') AS nm_dokter_perujuk,
			pl.kd_dokter AS kd_dokter_pj,
			COALESCE(d_pj.nm_dokter, '-') AS nm_dokter_pj,
			pl.nip AS nip_petugas,
			COALESCE(p.nama, '-') AS nama_petugas
		FROM periksa_lab pl
		INNER JOIN reg_periksa rp ON rp.no_rawat = pl.no_rawat
		LEFT JOIN jns_perawatan_lab jpl ON jpl.kd_jenis_prw = pl.kd_jenis_prw
		LEFT JOIN dokter d_perujuk ON d_perujuk.kd_dokter = pl.dokter_perujuk
		LEFT JOIN dokter d_pj ON d_pj.kd_dokter = pl.kd_dokter
		LEFT JOIN petugas p ON p.nip = pl.nip
		%s
		ORDER BY pl.tgl_periksa DESC, pl.jam DESC
		LIMIT ? OFFSET ?
	`, whereSQL)

	argsWithPaging := append(args, filter.Limit, filter.Offset())
	rows, err := r.db.QueryContext(ctx, dataQuery, argsWithPaging...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []HasilLaboratoriumDB
	for rows.Next() {
		var item HasilLaboratoriumDB
		err = rows.Scan(
			&item.NoRawat,
			&item.KodeTindakan,
			&item.NamaTindakan,
			&item.Kategori,
			&item.Status,
			&item.TanggalPeriksa,
			&item.JamPeriksa,
			&item.KodeDokterPerujuk,
			&item.NamaDokterPerujuk,
			&item.KodeDokterPJ,
			&item.NamaDokterPJ,
			&item.NipPetugas,
			&item.NamaPetugas,
		)
		if err != nil {
			return nil, 0, err
		}

		item.DetailPK = []ItemHasilLabPKDB{}
		list = append(list, item)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *repository) DetailHasilLabMB(ctx context.Context, noRawat string, kodeTindakan string, tanggalPeriksa string, jamPeriksa string) (*HasilLaboratoriumDB, error) {
	query := `
		SELECT 
			pl.no_rawat,
			pl.kd_jenis_prw AS kode_tindakan,
			COALESCE(jpl.nm_perawatan, pl.kd_jenis_prw) AS nama_tindakan,
			pl.kategori,
			pl.status,
			pl.tgl_periksa,
			pl.jam AS jam_periksa,
			pl.dokter_perujuk AS kd_dokter_perujuk,
			COALESCE(d_perujuk.nm_dokter, '-') AS nm_dokter_perujuk,
			pl.kd_dokter AS kd_dokter_pj,
			COALESCE(d_pj.nm_dokter, '-') AS nm_dokter_pj,
			pl.nip AS nip_petugas,
			COALESCE(p.nama, '-') AS nama_petugas
		FROM periksa_lab pl
		INNER JOIN reg_periksa rp ON rp.no_rawat = pl.no_rawat
		LEFT JOIN jns_perawatan_lab jpl ON jpl.kd_jenis_prw = pl.kd_jenis_prw
		LEFT JOIN dokter d_perujuk ON d_perujuk.kd_dokter = pl.dokter_perujuk
		LEFT JOIN dokter d_pj ON d_pj.kd_dokter = pl.kd_dokter
		LEFT JOIN petugas p ON p.nip = pl.nip
		WHERE pl.no_rawat = ? AND pl.kd_jenis_prw = ? AND pl.tgl_periksa = ? AND pl.jam = ? AND pl.kategori = 'MB'
	`

	var item HasilLaboratoriumDB
	err := r.db.QueryRowContext(ctx, query, noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa).Scan(
		&item.NoRawat,
		&item.KodeTindakan,
		&item.NamaTindakan,
		&item.Kategori,
		&item.Status,
		&item.TanggalPeriksa,
		&item.JamPeriksa,
		&item.KodeDokterPerujuk,
		&item.NamaDokterPerujuk,
		&item.KodeDokterPJ,
		&item.NamaDokterPJ,
		&item.NipPetugas,
		&item.NamaPetugas,
	)
	if err != nil {
		return nil, err
	}

	item.DetailPK = []ItemHasilLabPKDB{}
	return &item, nil
}
