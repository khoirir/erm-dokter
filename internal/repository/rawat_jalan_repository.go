package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"erm-dokter/internal/domain"
)

type rawatJalanRepository struct {
	db *sql.DB
}

func NewRawatJalanRepository(db *sql.DB) domain.RawatJalanRepository {
	return &rawatJalanRepository{
		db: db,
	}
}

const baseSelectKunjunganQuery = `
	SELECT 
		r.no_rawat,
		r.no_reg AS no_registrasi,
		DATE_FORMAT(r.tgl_registrasi, '%Y-%m-%d') AS tanggal_registrasi,
		r.jam_reg AS jam_registrasi,
		r.no_rkm_medis AS no_rekam_medis,
		p.nm_pasien AS nama_pasien,
		p.jk AS jenis_kelamin,
		DATE_FORMAT(p.tgl_lahir, '%Y-%m-%d') AS tanggal_lahir,
		r.almt_pj AS alamat,
		r.kd_poli AS kode_poli_asal,
		pol.nm_poli AS nama_poli_asal,
		r.kd_dokter AS kode_dokter_asal,
		d.nm_dokter AS nama_dokter_asal,
		'' AS kode_poli_rujukan,
		'' AS nama_poli_rujukan,
		'' AS kode_dokter_rujukan,
		'' AS nama_dokter_rujukan,
		r.kd_pj AS kode_penjamin,
		pj.png_jawab AS nama_penjamin,
		r.stts AS status_pemeriksaan,
		r.status_lanjut,
		r.status_bayar,
		'Bukan Rujukan' AS jenis_antrean
	FROM reg_periksa r
	INNER JOIN pasien p ON r.no_rkm_medis = p.no_rkm_medis
	INNER JOIN poliklinik pol ON r.kd_poli = pol.kd_poli
	INNER JOIN dokter d ON r.kd_dokter = d.kd_dokter
	INNER JOIN penjab pj ON r.kd_pj = pj.kd_pj
	WHERE r.kd_dokter = ?
	UNION ALL
	SELECT 
		r.no_rawat,
		r.no_reg AS no_registrasi,
		DATE_FORMAT(r.tgl_registrasi, '%Y-%m-%d') AS tanggal_registrasi,
		r.jam_reg AS jam_registrasi,
		r.no_rkm_medis AS no_rekam_medis,
		p.nm_pasien AS nama_pasien,
		p.jk AS jenis_kelamin,
		DATE_FORMAT(p.tgl_lahir, '%Y-%m-%d') AS tanggal_lahir,
		r.almt_pj AS alamat,
		r.kd_poli AS kode_poli_asal,
		pol.nm_poli AS nama_poli_asal,
		r.kd_dokter AS kode_dokter_asal,
		d.nm_dokter AS nama_dokter_asal,
		rip.kd_poli AS kode_poli_rujukan,
		pol_rip.nm_poli AS nama_poli_rujukan,
		rip.kd_dokter AS kode_dokter_rujukan,
		d_rip.nm_dokter AS nama_dokter_rujukan,
		r.kd_pj AS kode_penjamin,
		pj.png_jawab AS nama_penjamin,
		r.stts AS status_pemeriksaan,
		r.status_lanjut,
		r.status_bayar,
		'Rujukan' AS jenis_antrean
	FROM rujukan_internal_poli rip
	INNER JOIN reg_periksa r ON rip.no_rawat = r.no_rawat
	INNER JOIN pasien p ON r.no_rkm_medis = p.no_rkm_medis
	INNER JOIN poliklinik pol ON r.kd_poli = pol.kd_poli
	INNER JOIN dokter d ON r.kd_dokter = d.kd_dokter
	INNER JOIN penjab pj ON r.kd_pj = pj.kd_pj
	INNER JOIN poliklinik pol_rip ON rip.kd_poli = pol_rip.kd_poli
	INNER JOIN dokter d_rip ON rip.kd_dokter = d_rip.kd_dokter
	WHERE rip.kd_dokter = ?
`

type scanner interface {
	Scan(dest ...any) error
}

func scanKunjungan(s scanner) (*domain.KunjunganRawatJalan, error) {
	var k domain.KunjunganRawatJalan
	err := s.Scan(
		&k.NoRawat,
		&k.NoRegistrasi,
		&k.TanggalRegistrasi,
		&k.JamRegistrasi,
		&k.NoRekamMedis,
		&k.NamaPasien,
		&k.JenisKelamin,
		&k.TanggalLahir,
		&k.Alamat,
		&k.KodePoliAsal,
		&k.NamaPoliAsal,
		&k.KodeDokterAsal,
		&k.NamaDokterAsal,
		&k.KodePoliRujukan,
		&k.NamaPoliRujukan,
		&k.KodeDokterRujukan,
		&k.NamaDokterRujukan,
		&k.KodePenjamin,
		&k.NamaPenjamin,
		&k.StatusPemeriksaan,
		&k.StatusLanjut,
		&k.StatusBayar,
		&k.JenisAntrean,
	)
	if err != nil {
		return nil, err
	}
	k.JenisKelamin = k.FormatJenisKelamin()
	k.Umur = k.FormatUmur()
	k.NoRekamMedis = k.FormatNoRekamMedis()
	return &k, nil
}

func (r *rawatJalanRepository) DaftarAntreanDokter(ctx context.Context, filter domain.FilterAntreanDokter) ([]domain.KunjunganRawatJalan, int, error) {
	var (
		conditions []string
		args       []interface{}
	)
	args = append(args, filter.KodeDokter, filter.KodeDokter)

	if filter.JenisAntrean == domain.JenisAntreanRujukan {
		conditions = append(conditions, "t.jenis_antrean = 'Rujukan'")
	} else if filter.JenisAntrean == domain.JenisAntreanTidakRujukan {
		conditions = append(conditions, "t.jenis_antrean = 'Bukan Rujukan'")
	}

	if filter.Tanggal != "" {
		tglAwal := strings.TrimSpace(filter.Tanggal)
		tglAkhir := tglAwal
		tglParts := strings.Split(filter.Tanggal, ",")
		if len(tglParts) == 2 && strings.TrimSpace(tglParts[0]) != "" && strings.TrimSpace(tglParts[1]) != "" {
			tglAwal = strings.TrimSpace(tglParts[0])
			tglAkhir = strings.TrimSpace(tglParts[1])
		}
		conditions = append(conditions, "t.tanggal_registrasi BETWEEN ? AND ?")
		args = append(args, tglAwal, tglAkhir)
	}

	if filter.KodePenjamin != "" {
		conditions = append(conditions, "t.kode_penjamin = ?")
		args = append(args, filter.KodePenjamin)
	}

	if filter.StatusPemeriksaan != "" {
		conditions = append(conditions, "t.status_pemeriksaan = ?")
		args = append(args, filter.StatusPemeriksaan)
	}

	if filter.KataKunci != "" {
		conditions = append(conditions, "(t.nama_pasien LIKE ? OR t.no_rekam_medis LIKE ? OR t.no_rawat LIKE ?)")
		keywordPattern := "%" + filter.KataKunci + "%"
		args = append(args, keywordPattern, keywordPattern, keywordPattern)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM (%s) AS t
		%s
	`, baseSelectKunjunganQuery, whereClause)

	var totalData int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalData)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal menghitung total antrean: %w", err)
	}

	if totalData == 0 {
		return []domain.KunjunganRawatJalan{}, 0, nil
	}

	offset := (filter.Halaman - 1) * filter.Batas

	dataQuery := fmt.Sprintf(`
		SELECT * 
		FROM (%s) AS t
		%s
	`, baseSelectKunjunganQuery, whereClause)

	sortDir := "ASC"
	if strings.ToUpper(filter.SortOrder) == "DESC" {
		sortDir = "DESC"
	}
	switch filter.OrderBy {
	case "waktu_registrasi":
		dataQuery += fmt.Sprintf(" ORDER BY t.tanggal_registrasi %s, t.jam_registrasi %s", sortDir, sortDir)
	case "nama_pasien":
		dataQuery += fmt.Sprintf(" ORDER BY t.nama_pasien %s", sortDir)
	default:
		dataQuery += fmt.Sprintf(" ORDER BY t.tanggal_registrasi %s, t.jam_registrasi %s", sortDir, sortDir)
	}

	dataQuery += " LIMIT ? OFFSET ?"
	dataArgs := make([]interface{}, len(args))
	copy(dataArgs, args)
	dataArgs = append(dataArgs, filter.Batas, offset)
	rows, err := r.db.QueryContext(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal query data antrean: %w", err)
	}
	defer rows.Close()

	var daftarAntrean []domain.KunjunganRawatJalan
	for rows.Next() {
		kunjungan, err := scanKunjungan(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("gagal scan data antrean: %w", err)
		}
		daftarAntrean = append(daftarAntrean, *kunjungan)
	}

	return daftarAntrean, totalData, nil
}

func (r *rawatJalanRepository) DetailKunjungan(ctx context.Context, noRawat string, kodeDokter string) (*domain.KunjunganRawatJalan, error) {
	query := fmt.Sprintf(`
		SELECT * 
		FROM (%s) AS t 
		WHERE t.status_pemeriksaan <> 'Batal' AND t.no_rawat = ? 
		LIMIT 1
	`, baseSelectKunjunganQuery)
	row := r.db.QueryRowContext(ctx, query, kodeDokter, kodeDokter, noRawat)
	kunjungan, err := scanKunjungan(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("gagal query detail kunjungan: %w", err)
	}
	return kunjungan, nil
}
