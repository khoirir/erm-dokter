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

func (r *rawatJalanRepository) DaftarAntreanDokter(ctx context.Context, filter domain.FilterAntreanDokter) ([]domain.KunjunganRawatJalan, int, error) {
	var (
		conditions []string
		args       []interface{}
	)

	conditions = append(conditions, "r.stts <> 'Batal'")

	if filter.KodeDokter != "" {
		conditions = append(conditions, "(r.kd_dokter = ? OR rip.kd_dokter = ?)")
		args = append(args, filter.KodeDokter, filter.KodeDokter)
	}

	if filter.Tanggal != "" {
		conditions = append(conditions, "r.tgl_registrasi = ?")
		args = append(args, filter.Tanggal)
	}

	if filter.KodePenjamin != "" {
		conditions = append(conditions, "r.kd_pj = ?")
		args = append(args, filter.KodePenjamin)
	}

	if filter.StatusPemeriksaan != "" {
		conditions = append(conditions, "r.stts = ?")
		args = append(args, filter.StatusPemeriksaan)
	}

	if filter.JenisAntrean == domain.JenisAntreanRujukan {
		conditions = append(conditions, "rip.no_rawat IS NOT NULL")
	}
	if filter.JenisAntrean == domain.JenisAntreanTidakRujukan {
		conditions = append(conditions, "rip.no_rawat IS NULL")
	}

	if filter.KataKunci != "" {
		conditions = append(conditions, "(p.nm_pasien LIKE ? OR r.no_rkm_medis LIKE ? OR r.no_rawat LIKE ?)")
		keywordPattern := "%" + filter.KataKunci + "%"
		args = append(args, keywordPattern, keywordPattern, keywordPattern)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := `
		SELECT COUNT(*) 
		FROM reg_periksa r
		INNER JOIN pasien p ON r.no_rkm_medis = p.no_rkm_medis
		LEFT JOIN rujukan_internal_poli rip ON r.no_rawat = rip.no_rawat
		INNER JOIN poliklinik pol ON r.kd_poli = pol.kd_poli
		INNER JOIN dokter d ON r.kd_dokter = d.kd_dokter
		INNER JOIN penjab pj ON r.kd_pj = pj.kd_pj
	` + whereClause

	var totalData int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalData)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal menghitung total antrean: %w", err)
	}

	if totalData == 0 {
		return []domain.KunjunganRawatJalan{}, 0, nil
	}

	offset := (filter.Halaman - 1) * filter.Batas

	dataQuery := `
		SELECT 
			r.no_rawat,
			r.no_reg AS no_registrasi,
			r.tgl_registrasi AS tanggal_registrasi,
			r.jam_reg AS jam_registrasi,
			r.no_rkm_medis AS no_rekam_medis,
			p.nm_pasien AS nama_pasien,
			p.jk AS jenis_kelamin,
			p.tgl_lahir AS tanggal_lahir,
			r.almt_pj AS alamat,
			r.kd_poli AS kode_poli,
			pol.nm_poli AS nama_poli,
			r.kd_dokter AS kode_dokter,
			d.nm_dokter AS nama_dokter,
			r.kd_pj AS kode_penjamin,
			pj.png_jawab AS nama_penjamin,
			r.stts AS status_pemeriksaan,
			r.status_lanjut,
			r.status_bayar,
			CASE 
				WHEN rip.no_rawat IS NOT NULL THEN 'Rujukan'
				ELSE 'Bukan Rujukan'
			END AS jenis_antrean
		FROM reg_periksa r
		INNER JOIN pasien p ON r.no_rkm_medis = p.no_rkm_medis
		LEFT JOIN rujukan_internal_poli rip ON r.no_rawat = rip.no_rawat
		INNER JOIN poliklinik pol ON r.kd_poli = pol.kd_poli
		INNER JOIN dokter d ON r.kd_dokter = d.kd_dokter
		INNER JOIN penjab pj ON r.kd_pj = pj.kd_pj
	` + whereClause

	sortDir := "ASC"
	if strings.ToUpper(filter.SortOrder) == "DESC" {
		sortDir = "DESC"
	}
	switch filter.OrderBy {
	case "waktu_registrasi":
		dataQuery += fmt.Sprintf(" ORDER BY r.tgl_registrasi %s, r.jam_reg %s", sortDir, sortDir)
	case "nama_pasien":
		dataQuery += fmt.Sprintf(" ORDER BY p.nm_pasien %s", sortDir)
	default:
		dataQuery += fmt.Sprintf(" ORDER BY r.tgl_registrasi %s, r.jam_reg %s", sortDir, sortDir)
	}

	dataQuery += " LIMIT ? OFFSET ?"
	args = append(args, filter.Batas, offset)
	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal query data antrean: %w", err)
	}
	defer rows.Close()

	var daftarAntrean []domain.KunjunganRawatJalan
	for rows.Next() {
		var kunjungan domain.KunjunganRawatJalan
		err := rows.Scan(
			&kunjungan.NoRawat,
			&kunjungan.NoRegistrasi,
			&kunjungan.TanggalRegistrasi,
			&kunjungan.JamRegistrasi,
			&kunjungan.NoRekamMedis,
			&kunjungan.NamaPasien,
			&kunjungan.JenisKelamin,
			&kunjungan.TanggalLahir,
			&kunjungan.Alamat,
			&kunjungan.KodePoli,
			&kunjungan.NamaPoli,
			&kunjungan.KodeDokter,
			&kunjungan.NamaDokter,
			&kunjungan.KodePenjamin,
			&kunjungan.NamaPenjamin,
			&kunjungan.StatusPemeriksaan,
			&kunjungan.StatusLanjut,
			&kunjungan.StatusBayar,
			&kunjungan.JenisAntrean,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("gagal scan data antrean: %w", err)
		}
		kunjungan.JenisKelamin = kunjungan.FormatJenisKelamin()
		kunjungan.Umur = kunjungan.FormatUmur()
		kunjungan.NoRekamMedis = kunjungan.FormatNoRekamMedis()
		daftarAntrean = append(daftarAntrean, kunjungan)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("gagal iterasi data antrean: %w", err)
	}
	return daftarAntrean, totalData, nil
}

func (r *rawatJalanRepository) DetailKunjungan(ctx context.Context, noRawat string) (*domain.KunjunganRawatJalan, error) {
	query := `
		SELECT 
			r.no_rawat,
			r.no_reg AS no_registrasi,
			r.tgl_registrasi AS tanggal_registrasi,
			r.jam_reg AS jam_registrasi,
			r.no_rkm_medis AS no_rekam_medis,
			p.nm_pasien AS nama_pasien,
			p.jk AS jenis_kelamin,
			p.tgl_lahir AS tanggal_lahir,
			r.almt_pj AS alamat,
			r.kd_poli AS kode_poli,
			pol.nm_poli AS nama_poli,
			r.kd_dokter AS kode_dokter,
			d.nm_dokter AS nama_dokter,
			r.kd_pj AS kode_penjamin,
			pj.png_jawab AS nama_penjamin,
			r.stts AS status_pemeriksaan,
			r.status_lanjut,
			r.status_bayar,
			CASE 
				WHEN rip.no_rawat IS NOT NULL THEN 'Rujukan'
				ELSE 'Bukan Rujukan'
			END AS jenis_antrean
		FROM reg_periksa r
		INNER JOIN pasien p ON r.no_rkm_medis = p.no_rkm_medis
		LEFT JOIN rujukan_internal_poli rip ON r.no_rawat = rip.no_rawat
		INNER JOIN poliklinik pol ON r.kd_poli = pol.kd_poli
		INNER JOIN dokter d ON r.kd_dokter = d.kd_dokter
		INNER JOIN penjab pj ON r.kd_pj = pj.kd_pj
		WHERE r.stts <> 'Batal' AND r.no_rawat = ?
		LIMIT 1
	`

	var kunjungan domain.KunjunganRawatJalan

	err := r.db.QueryRowContext(ctx, query, noRawat).Scan(
		&kunjungan.NoRawat,
		&kunjungan.NoRegistrasi,
		&kunjungan.TanggalRegistrasi,
		&kunjungan.JamRegistrasi,
		&kunjungan.NoRekamMedis,
		&kunjungan.NamaPasien,
		&kunjungan.JenisKelamin,
		&kunjungan.TanggalLahir,
		&kunjungan.Alamat,
		&kunjungan.KodePoli,
		&kunjungan.NamaPoli,
		&kunjungan.KodeDokter,
		&kunjungan.NamaDokter,
		&kunjungan.KodePenjamin,
		&kunjungan.NamaPenjamin,
		&kunjungan.StatusPemeriksaan,
		&kunjungan.StatusLanjut,
		&kunjungan.StatusBayar,
		&kunjungan.JenisAntrean,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("gagal query detail kunjungan: %w", err)
	}

	kunjungan.Umur = kunjungan.FormatUmur()
	kunjungan.JenisKelamin = kunjungan.FormatJenisKelamin()
	kunjungan.NoRekamMedis = kunjungan.FormatNoRekamMedis()
	return &kunjungan, nil
}
