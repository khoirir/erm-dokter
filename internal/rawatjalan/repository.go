package rawatjalan

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type Repository interface {
	DaftarAntreanDokter(ctx context.Context, kodeDokter string, filter FilterAntreanDokter) ([]KunjunganRawatJalan, int, error)
	DetailKunjungan(ctx context.Context, noRawat string, kodeDokter string) (*KunjunganRawatJalan, error)
	RiwayatKunjunganPasien(ctx context.Context, noRM string) ([]KunjunganRawatJalan, error)
	GetWaktuRegistrasi(ctx context.Context, noRawat string) (tanggal string, jam string, exists bool, err error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

const selectKunjunganBukanRujukan = `
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
`

const selectKunjunganRujukan = `
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
`

var orderByMapping = map[string]func(string) string{
	"nama_pasien":      func(dir string) string { return fmt.Sprintf("t.nama_pasien %s", dir) },
	"waktu_registrasi": func(dir string) string { return fmt.Sprintf("t.tanggal_registrasi %s, t.jam_registrasi %s", dir, dir) },
}

type scanner interface {
	Scan(dest ...any) error
}

func scanKunjungan(s scanner) (*KunjunganRawatJalan, error) {
	var k KunjunganRawatJalan
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
	return &k, nil
}

func buildBranchConditions(dokterCol string, kodeDokter string, filter FilterAntreanDokter) (string, []any) {
	var conditions []string
	var args []any

	if kodeDokter != "" {
		conditions = append(conditions, dokterCol+" = ?")
		args = append(args, kodeDokter)
	}

	if filter.Tanggal != "" {
		tglParts := strings.Split(filter.Tanggal, ",")
		tglAwal := strings.TrimSpace(tglParts[0])
		tglAkhir := strings.TrimSpace(tglParts[1])
		conditions = append(conditions, "r.tgl_registrasi BETWEEN ? AND ?")
		args = append(args, tglAwal, tglAkhir)
	}

	if filter.KodePenjamin != "" {
		conditions = append(conditions, "r.kd_pj = ?")
		args = append(args, filter.KodePenjamin)
	}

	if filter.StatusPemeriksaan != "" {
		conditions = append(conditions, "r.stts = ?")
		args = append(args, filter.StatusPemeriksaan)
	}

	if filter.StatusLanjut != "" {
		conditions = append(conditions, "r.status_lanjut = ?")
		args = append(args, filter.StatusLanjut)
	}

	if filter.KataKunci != "" {
		conditions = append(conditions, "(p.nm_pasien LIKE ? OR r.no_rkm_medis LIKE ? OR r.no_rawat LIKE ?)")
		keywordPattern := "%" + filter.KataKunci + "%"
		args = append(args, keywordPattern, keywordPattern, keywordPattern)
	}

	if len(conditions) == 0 {
		return "", args
	}

	return " WHERE " + strings.Join(conditions, " AND "), args
}

func buildBaseQuery(kodeDokter string, filter FilterAntreanDokter) (string, []any) {
	switch filter.JenisAntrean {
	case JenisAntreanTidakRujukan:
		where, args := buildBranchConditions("r.kd_dokter", kodeDokter, filter)
		return selectKunjunganBukanRujukan + where, args

	case JenisAntreanRujukan:
		where, args := buildBranchConditions("rip.kd_dokter", kodeDokter, filter)
		return selectKunjunganRujukan + where, args

	default:
		where1, args1 := buildBranchConditions("r.kd_dokter", kodeDokter, filter)
		where2, args2 := buildBranchConditions("rip.kd_dokter", kodeDokter, filter)

		query1 := selectKunjunganBukanRujukan + where1
		query2 := selectKunjunganRujukan + where2

		combinedQuery := fmt.Sprintf("%s UNION ALL %s", query1, query2)
		combinedArgs := append(args1, args2...)
		return combinedQuery, combinedArgs
	}
}

func (r *repository) DaftarAntreanDokter(ctx context.Context, kodeDokter string, filter FilterAntreanDokter) ([]KunjunganRawatJalan, int, error) {
	baseQuery, baseArgs := buildBaseQuery(kodeDokter, filter)

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS t", baseQuery)

	var totalData int
	if err := r.db.QueryRowContext(ctx, countQuery, baseArgs...).Scan(&totalData); err != nil || totalData == 0 {
		return []KunjunganRawatJalan{}, 0, nil
	}

	builder, exists := orderByMapping[filter.OrderBy]
	if !exists {
		builder = orderByMapping["waktu_registrasi"]
	}
	orderClause := " ORDER BY " + builder(filter.SortOrder)

	dataQuery := fmt.Sprintf("SELECT * FROM (%s) AS t %s LIMIT ? OFFSET ?", baseQuery, orderClause)

	dataArgs := append(baseArgs, filter.Batas, filter.Offset())

	rows, err := r.db.QueryContext(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal query data antrean: %w", err)
	}
	defer rows.Close()

	var daftarAntrean []KunjunganRawatJalan
	for rows.Next() {
		kunjungan, err := scanKunjungan(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("gagal scan data antrean: %w", err)
		}
		daftarAntrean = append(daftarAntrean, *kunjungan)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error saat iterasi data antrean: %w", err)
	}

	return daftarAntrean, totalData, nil
}

func (r *repository) DetailKunjungan(ctx context.Context, noRawat string, kodeDokter string) (*KunjunganRawatJalan, error) {
	baseQuery, baseArgs := buildBaseQuery(kodeDokter, FilterAntreanDokter{})

	query := fmt.Sprintf(`
		SELECT * 
		FROM (%s) AS t 
		WHERE t.status_pemeriksaan <> 'Batal' AND t.no_rawat = ? 
		LIMIT 1
	`, baseQuery)

	args := append(baseArgs, noRawat)
	row := r.db.QueryRowContext(ctx, query, args...)
	kunjungan, err := scanKunjungan(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("gagal query detail kunjungan: %w", err)
	}
	return kunjungan, nil
}

func (r *repository) RiwayatKunjunganPasien(ctx context.Context, noRekamMedis string) ([]KunjunganRawatJalan, error) {
	baseQuery, baseArgs := buildBaseQuery("", FilterAntreanDokter{})

	query := fmt.Sprintf(`
		SELECT * 
		FROM (%s) AS t 
		WHERE t.status_pemeriksaan <> 'Batal' AND t.no_rekam_medis = ? 
		ORDER BY t.tanggal_registrasi DESC, t.jam_registrasi DESC
	`, baseQuery)

	args := append(baseArgs, noRekamMedis)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("gagal query riwayat kunjungan pasien: %w", err)
	}
	defer rows.Close()

	var listKunjungan []KunjunganRawatJalan
	for rows.Next() {
		kunjungan, err := scanKunjungan(rows)
		if err != nil {
			return nil, fmt.Errorf("gagal scan riwayat kunjungan: %w", err)
		}
		listKunjungan = append(listKunjungan, *kunjungan)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterasi riwayat kunjungan: %w", err)
	}

	return listKunjungan, nil
}

const selectWaktuRegistrasi = `
	SELECT 
		DATE_FORMAT(tgl_registrasi, '%Y-%m-%d') AS tgl_registrasi,
		jam_reg
	FROM reg_periksa
	WHERE no_rawat = ?
`

func (r *repository) GetWaktuRegistrasi(ctx context.Context, noRawat string) (string, string, bool, error) {
	var tglReg, jamReg string
	err := r.db.QueryRowContext(ctx, selectWaktuRegistrasi, noRawat).Scan(&tglReg, &jamReg)
	if errors.Is(err, sql.ErrNoRows) {
		return "", "", false, nil
	}
	if err != nil {
		return "", "", false, fmt.Errorf("gagal query waktu registrasi: %w", err)
	}
	return tglReg, jamReg, true, nil
}

