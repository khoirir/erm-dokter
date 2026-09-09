package pasien

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/formatter"
)

type Repository interface {
	CariByNoRM(ctx context.Context, noRM string) (*Pasien, error)
	CariByNIK(ctx context.Context, nik string) (*Pasien, error)
	CariPasien(ctx context.Context, kataKunci string, limit int) ([]Pasien, error)
	GetNoRMByNoRawat(ctx context.Context, noRawat string) (string, error)
	RiwayatKunjunganPasien(ctx context.Context, noRM string, filter FilterRiwayatKunjungan) ([]RiwayatKunjungan, int, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

const selectPasien = `
	SELECT 
		no_rkm_medis,
		nm_pasien,
		no_ktp,
		no_peserta,
		jk,
		DATE_FORMAT(tgl_lahir, '%Y-%m-%d') AS tgl_lahir,
		tmp_lahir,
		alamat,
		no_tlp
	FROM pasien
`

func (r *repository) CariByNoRM(ctx context.Context, noRM string) (*Pasien, error) {
	query := selectPasien + " WHERE no_rkm_medis = ? LIMIT 1"
	row := r.db.QueryRowContext(ctx, query, noRM)

	var p Pasien
	err := row.Scan(
		&p.NoRM,
		&p.Nama,
		&p.NIK,
		&p.NoBPJS,
		&p.JenisKelamin,
		&p.TanggalLahir,
		&p.TempatLahir,
		&p.Alamat,
		&p.NoTelp,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("gagal query pasien by no_rkm_medis: %w", err)
	}

	return &p, nil
}

func (r *repository) CariByNIK(ctx context.Context, nik string) (*Pasien, error) {
	query := selectPasien + " WHERE no_ktp = ? LIMIT 1"
	row := r.db.QueryRowContext(ctx, query, nik)

	var p Pasien
	err := row.Scan(
		&p.NoRM,
		&p.Nama,
		&p.NIK,
		&p.NoBPJS,
		&p.JenisKelamin,
		&p.TanggalLahir,
		&p.TempatLahir,
		&p.Alamat,
		&p.NoTelp,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("gagal query pasien by NIK: %w", err)
	}

	return &p, nil
}

func (r *repository) CariPasien(ctx context.Context, kataKunci string, limit int) ([]Pasien, error) {
	if limit <= 0 {
		limit = 20
	}
	pattern := "%" + kataKunci + "%"
	query := selectPasien + " WHERE nm_pasien LIKE ? OR no_rkm_medis LIKE ? OR no_ktp LIKE ? ORDER BY nm_pasien ASC LIMIT ?"

	rows, err := r.db.QueryContext(ctx, query, pattern, pattern, pattern, limit)
	if err != nil {
		return nil, fmt.Errorf("gagal mencari pasien: %w", err)
	}
	defer rows.Close()

	list := make([]Pasien, 0)
	for rows.Next() {
		var p Pasien
		err := rows.Scan(
			&p.NoRM,
			&p.Nama,
			&p.NIK,
			&p.NoBPJS,
			&p.JenisKelamin,
			&p.TanggalLahir,
			&p.TempatLahir,
			&p.Alamat,
			&p.NoTelp,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal scan data pencarian pasien: %w", err)
		}
		list = append(list, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterasi pencarian pasien: %w", err)
	}

	return list, nil
}

func (r *repository) GetNoRMByNoRawat(ctx context.Context, noRawat string) (string, error) {
	query := "SELECT no_rkm_medis FROM reg_periksa WHERE no_rawat = ? LIMIT 1"
	var noRM string
	err := r.db.QueryRowContext(ctx, query, noRawat).Scan(&noRM)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("gagal query no_rkm_medis dari no_rawat: %w", err)
	}
	return noRM, nil
}

const selectRiwayatKunjungan = `
	SELECT 
		rp.no_rawat,
		rp.no_reg AS no_registrasi,
		DATE_FORMAT(rp.tgl_registrasi, '%Y-%m-%d') AS tanggal_registrasi,
		rp.jam_reg AS jam_registrasi,
		rp.no_rkm_medis,
		COALESCE(p.nm_pasien, '') AS nama_pasien,
		rp.status_lanjut,
		rp.status_bayar,
		rp.kd_pj AS kode_penjamin,
		COALESCE(pj.png_jawab, '') AS nama_penjamin,
		rp.kd_dokter AS kode_dokter,
		COALESCE(d.nm_dokter, '') AS nama_dokter,
		COALESCE(pol.kd_poli, '') AS kode_poli,
		COALESCE(pol.nm_poli, '') AS nama_poli,
		rp.stts AS status_periksa,
		COALESCE((
			SELECT GROUP_CONCAT(d_dpjp.nm_dokter ORDER BY d_dpjp.nm_dokter SEPARATOR ';;;')
			FROM dpjp_ranap dr
			INNER JOIN dokter d_dpjp ON dr.kd_dokter = d_dpjp.kd_dokter
			WHERE dr.no_rawat = rp.no_rawat
		), '') AS dpjp
`

const fromAndJoinRiwayat = `
	FROM reg_periksa rp
	INNER JOIN pasien p ON rp.no_rkm_medis = p.no_rkm_medis
	LEFT JOIN penjab pj ON rp.kd_pj = pj.kd_pj
	LEFT JOIN dokter d ON rp.kd_dokter = d.kd_dokter
	LEFT JOIN poliklinik pol ON rp.kd_poli = pol.kd_poli
`

func (r *repository) buildRiwayatConditions(noRM string, filter FilterRiwayatKunjungan) ([]string, []any) {
	conditions := []string{
		"rp.no_rkm_medis = ?",
		"rp.stts <> 'Batal'",
	}
	args := []any{noRM}

	if filter.Tanggal != "" {
		tglAwal, tglAkhir := formatter.ParseRentangTanggal(filter.Tanggal)
		if tglAwal != "" && tglAkhir != "" {
			conditions = append(conditions, "rp.tgl_registrasi BETWEEN ? AND ?")
			args = append(args, tglAwal, tglAkhir)
		}
	}

	return conditions, args
}

func (r *repository) countRiwayat(ctx context.Context, whereClause string, args []any) (int, error) {
	countQuery := fmt.Sprintf(`
		SELECT COUNT(rp.no_rawat)
		FROM reg_periksa rp
		%s
	`, whereClause)

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("gagal menghitung total riwayat kunjungan: %w", err)
	}

	return total, nil
}

func (r *repository) RiwayatKunjunganPasien(ctx context.Context, noRM string, filter FilterRiwayatKunjungan) ([]RiwayatKunjungan, int, error) {
	conditions, whereArgs := r.buildRiwayatConditions(noRM, filter)
	whereClause := "WHERE " + strings.Join(conditions, " AND ")

	total, err := r.countRiwayat(ctx, whereClause, whereArgs)
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return make([]RiwayatKunjungan, 0), 0, nil
	}

	dataQuery := fmt.Sprintf(`
		%s
		%s
		%s
		ORDER BY rp.tgl_registrasi DESC, rp.jam_reg DESC
		LIMIT ? OFFSET ?
	`, selectRiwayatKunjungan, fromAndJoinRiwayat, whereClause)

	finalArgs := make([]any, 0, len(whereArgs)+2)
	finalArgs = append(finalArgs, whereArgs...)
	finalArgs = append(finalArgs, filter.Limit, filter.Offset())

	rows, err := r.db.QueryContext(ctx, dataQuery, finalArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal query riwayat kunjungan: %w", err)
	}
	defer rows.Close()

	list := make([]RiwayatKunjungan, 0, filter.Limit)
	ranapNoRawats := make([]string, 0)
	ranapIndexMap := make(map[string]int)
	ralanNoRawats := make([]string, 0)
	ralanIndexMap := make(map[string]int)
	dpjpMap := make(map[string][]string)

	for rows.Next() {
		var (
			item    RiwayatKunjungan
			dpjpStr string
		)
		err := rows.Scan(
			&item.NoRawat,
			&item.NoRegistrasi,
			&item.TanggalRegistrasi,
			&item.JamRegistrasi,
			&item.NoRekamMedis,
			&item.NamaPasien,
			&item.StatusLanjut,
			&item.StatusBayar,
			&item.KodePenjamin,
			&item.NamaPenjamin,
			&item.KodeDokter,
			&item.NamaDokter,
			&item.KodePoli,
			&item.NamaPoli,
			&item.StatusPeriksa,
			&dpjpStr,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("gagal scan riwayat kunjungan: %w", err)
		}

		if dpjpStr != "" {
			parts := strings.Split(dpjpStr, ";;;")
			cleanDPJP := make([]string, 0, len(parts))
			for _, p := range parts {
				clean := strings.TrimSpace(p)
				if clean != "" {
					cleanDPJP = append(cleanDPJP, clean)
				}
			}
			if len(cleanDPJP) > 0 {
				dpjpMap[item.NoRawat] = cleanDPJP
			}
		}

		switch item.StatusLanjut {
		case shared.StatusLanjutRawatJalan:
			ralanNoRawats = append(ralanNoRawats, item.NoRawat)
			ralanIndexMap[item.NoRawat] = len(list)
		case shared.StatusLanjutRawatInap:
			ranapNoRawats = append(ranapNoRawats, item.NoRawat)
			ranapIndexMap[item.NoRawat] = len(list)
			item.KamarInap = &RiwayatRawatInap{
				DPJP:  make([]string, 0),
				Kamar: make([]RiwayatKamarInap, 0),
			}
		}

		list = append(list, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterasi riwayat kunjungan: %w", err)
	}

	if len(ralanNoRawats) > 0 {
		rujukanMap, err := r.getRiwayatRujukanInternalByNoRawats(ctx, ralanNoRawats)
		if err != nil {
			return nil, 0, err
		}
		for noRawat, rujukans := range rujukanMap {
			if idx, ok := ralanIndexMap[noRawat]; ok {
				list[idx].RujukanInternal = rujukans
			}
		}
	}

	if len(ranapNoRawats) > 0 {
		kamarMap, err := r.getRiwayatKamarInapByNoRawats(ctx, ranapNoRawats)
		if err != nil {
			return nil, 0, err
		}
		for _, noRawat := range ranapNoRawats {
			if idx, ok := ranapIndexMap[noRawat]; ok {
				dpjps := dpjpMap[noRawat]
				if dpjps == nil {
					dpjps = make([]string, 0)
				}
				kamars := kamarMap[noRawat]
				if kamars == nil {
					kamars = make([]RiwayatKamarInap, 0)
				}
				list[idx].KamarInap = &RiwayatRawatInap{
					DPJP:  dpjps,
					Kamar: kamars,
				}
			}
		}
	}

	return list, total, nil
}

func (r *repository) getRiwayatRujukanInternalByNoRawats(ctx context.Context, noRawats []string) (map[string][]RiwayatRujukanInternal, error) {
	if len(noRawats) == 0 {
		return make(map[string][]RiwayatRujukanInternal), nil
	}

	placeholders := make([]string, len(noRawats))
	args := make([]any, len(noRawats))
	for i, nr := range noRawats {
		placeholders[i] = "?"
		args[i] = nr
	}

	query := fmt.Sprintf(`
		SELECT 
			rip.no_rawat,
			rip.kd_poli,
			COALESCE(p.nm_poli, '') AS nama_poli,
			rip.kd_dokter,
			COALESCE(d.nm_dokter, '') AS nama_dokter
		FROM rujukan_internal_poli rip
		INNER JOIN poliklinik p ON rip.kd_poli = p.kd_poli
		INNER JOIN dokter d ON rip.kd_dokter = d.kd_dokter
		WHERE rip.no_rawat IN (%s)
		ORDER BY p.nm_poli ASC, d.nm_dokter ASC
	`, strings.Join(placeholders, ", "))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("gagal query riwayat rujukan internal: %w", err)
	}
	defer rows.Close()

	result := make(map[string][]RiwayatRujukanInternal)
	for rows.Next() {
		var (
			noRawat string
			rujukan RiwayatRujukanInternal
		)
		err := rows.Scan(
			&noRawat,
			&rujukan.KodePoli,
			&rujukan.NamaPoli,
			&rujukan.KodeDokter,
			&rujukan.NamaDokter,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal scan riwayat rujukan internal: %w", err)
		}

		result[noRawat] = append(result[noRawat], rujukan)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterasi riwayat rujukan internal: %w", err)
	}

	return result, nil
}

func (r *repository) getRiwayatKamarInapByNoRawats(ctx context.Context, noRawats []string) (map[string][]RiwayatKamarInap, error) {
	if len(noRawats) == 0 {
		return make(map[string][]RiwayatKamarInap), nil
	}

	placeholders := make([]string, len(noRawats))
	args := make([]any, len(noRawats))
	for i, nr := range noRawats {
		placeholders[i] = "?"
		args[i] = nr
	}

	query := fmt.Sprintf(`
		SELECT 
			ki.no_rawat,
			ki.kd_kamar,
			COALESCE(b.nm_bangsal, '') AS nama_bangsal,
			COALESCE(k.kelas, '') AS kelas,
			COALESCE(DATE_FORMAT(ki.tgl_masuk, '%%Y-%%m-%%d'), '') AS tanggal_masuk,
			COALESCE(ki.jam_masuk, '') AS jam_masuk,
			COALESCE(DATE_FORMAT(ki.tgl_keluar, '%%Y-%%m-%%d'), '') AS tanggal_keluar,
			COALESCE(ki.jam_keluar, '') AS jam_keluar,
			COALESCE(ki.stts_pulang, '') AS status_pulang,
			COALESCE(ki.lama, 0) AS lama_inap,
			COALESCE(ki.diagnosa_awal, '') AS diagnosa_awal,
			COALESCE(ki.diagnosa_akhir, '') AS diagnosa_akhir
		FROM kamar_inap ki
		LEFT JOIN kamar k ON ki.kd_kamar = k.kd_kamar
		LEFT JOIN bangsal b ON k.kd_bangsal = b.kd_bangsal
		WHERE ki.no_rawat IN (%s)
		ORDER BY ki.tgl_masuk ASC, ki.jam_masuk ASC
	`, strings.Join(placeholders, ", "))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("gagal query riwayat kamar inap: %w", err)
	}
	defer rows.Close()

	result := make(map[string][]RiwayatKamarInap)
	for rows.Next() {
		var (
			noRawat string
			kamar   RiwayatKamarInap
		)
		err := rows.Scan(
			&noRawat,
			&kamar.KodeKamar,
			&kamar.NamaBangsal,
			&kamar.Kelas,
			&kamar.TanggalMasuk,
			&kamar.JamMasuk,
			&kamar.TanggalKeluar,
			&kamar.JamKeluar,
			&kamar.StatusPulang,
			&kamar.LamaInap,
			&kamar.DiagnosaAwal,
			&kamar.DiagnosaAkhir,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal scan riwayat kamar inap: %w", err)
		}

		if kamar.TanggalKeluar == "0000-00-00" {
			kamar.TanggalKeluar = ""
		}
		if kamar.JamKeluar == "00:00:00" {
			kamar.JamKeluar = ""
		}
		if kamar.StatusPulang == "" {
			kamar.StatusPulang = "-"
		}

		result[noRawat] = append(result[noRawat], kamar)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterasi riwayat kamar inap: %w", err)
	}

	return result, nil
}
