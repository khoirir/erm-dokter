package resumepasien

import (
	"context"
	"database/sql"
	"errors"
)

type Repository interface {
	DetailResumePasien(ctx context.Context, noRawat string) (*ResumePasien, error)
	RiwayatResumePasienByNoRM(ctx context.Context, noRM string) ([]ResumePasien, error)
	CekResumePasienAda(ctx context.Context, noRawat string) (bool, error)
	SimpanResumePasien(ctx context.Context, noRawat, kodeDokter string, req SimpanResumePasienRequest) (*ResumePasien, error)
	UpdateResumePasien(ctx context.Context, noRawat string, req UpdateResumePasienRequest) (*ResumePasien, error)
	HapusResumePasien(ctx context.Context, noRawat string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

const selectResumePasien = `
	SELECT 
		rp.no_rawat,
		rp.kd_dokter,
		COALESCE(d.nm_dokter, '') AS nama_dokter,
		COALESCE(rp.keluhan_utama, '') AS keluhan_utama,
		COALESCE(rp.jalannya_penyakit, '') AS jalannya_penyakit,
		COALESCE(rp.pemeriksaan_penunjang, '') AS pemeriksaan_penunjang,
		COALESCE(rp.hasil_laborat, '') AS hasil_laborat,
		COALESCE(rp.diagnosa_utama, '') AS diagnosa_utama,
		COALESCE(rp.kd_diagnosa_utama, '') AS kd_diagnosa_utama,
		COALESCE(rp.diagnosa_sekunder, '') AS diagnosa_sekunder,
		COALESCE(rp.kd_diagnosa_sekunder, '') AS kd_diagnosa_sekunder,
		COALESCE(rp.diagnosa_sekunder2, '') AS diagnosa_sekunder2,
		COALESCE(rp.kd_diagnosa_sekunder2, '') AS kd_diagnosa_sekunder2,
		COALESCE(rp.diagnosa_sekunder3, '') AS diagnosa_sekunder3,
		COALESCE(rp.kd_diagnosa_sekunder3, '') AS kd_diagnosa_sekunder3,
		COALESCE(rp.diagnosa_sekunder4, '') AS diagnosa_sekunder4,
		COALESCE(rp.kd_diagnosa_sekunder4, '') AS kd_diagnosa_sekunder4,
		COALESCE(rp.prosedur_utama, '') AS prosedur_utama,
		COALESCE(rp.kd_prosedur_utama, '') AS kd_prosedur_utama,
		COALESCE(rp.prosedur_sekunder, '') AS prosedur_sekunder,
		COALESCE(rp.kd_prosedur_sekunder, '') AS kd_prosedur_sekunder,
		COALESCE(rp.prosedur_sekunder2, '') AS prosedur_sekunder2,
		COALESCE(rp.kd_prosedur_sekunder2, '') AS kd_prosedur_sekunder2,
		COALESCE(rp.prosedur_sekunder3, '') AS prosedur_sekunder3,
		COALESCE(rp.kd_prosedur_sekunder3, '') AS kd_prosedur_sekunder3,
		rp.kondisi_pulang,
		COALESCE(rp.obat_pulang, '') AS obat_pulang
	FROM resume_pasien rp
	INNER JOIN dokter d ON d.kd_dokter = rp.kd_dokter
`

func scanResumePasien(scanner interface{ Scan(dest ...any) error }) (*ResumePasien, error) {
	var item ResumePasien
	err := scanner.Scan(
		&item.NoRawat,
		&item.KodeDokter,
		&item.NamaDokter,
		&item.KeluhanUtama,
		&item.JalannyaPenyakit,
		&item.PemeriksaanPenunjang,
		&item.HasilLaborat,
		&item.DiagnosaUtama,
		&item.KodeDiagnosaUtama,
		&item.DiagnosaSekunder,
		&item.KodeDiagnosaSekunder,
		&item.DiagnosaSekunder2,
		&item.KodeDiagnosaSekunder2,
		&item.DiagnosaSekunder3,
		&item.KodeDiagnosaSekunder3,
		&item.DiagnosaSekunder4,
		&item.KodeDiagnosaSekunder4,
		&item.ProsedurUtama,
		&item.KodeProsedurUtama,
		&item.ProsedurSekunder,
		&item.KodeProsedurSekunder,
		&item.ProsedurSekunder2,
		&item.KodeProsedurSekunder2,
		&item.ProsedurSekunder3,
		&item.KodeProsedurSekunder3,
		&item.KondisiPulang,
		&item.ObatPulang,
	)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *repository) DetailResumePasien(ctx context.Context, noRawat string) (*ResumePasien, error) {
	query := selectResumePasien + ` WHERE rp.no_rawat = ? LIMIT 1`
	row := r.db.QueryRowContext(ctx, query, noRawat)

	item, err := scanResumePasien(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return item, nil
}

func (r *repository) RiwayatResumePasienByNoRM(ctx context.Context, noRM string) ([]ResumePasien, error) {
	query := selectResumePasien + ` 
		INNER JOIN reg_periksa reg ON reg.no_rawat = rp.no_rawat 
		WHERE reg.no_rkm_medis = ? 
		ORDER BY reg.tgl_registrasi DESC, reg.jam_reg DESC
	`
	rows, err := r.db.QueryContext(ctx, query, noRM)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]ResumePasien, 0)
	for rows.Next() {
		item, scanErr := scanResumePasien(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		list = append(list, *item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *repository) CekResumePasienAda(ctx context.Context, noRawat string) (bool, error) {
	query := `SELECT COUNT(1) FROM resume_pasien WHERE no_rawat = ?`
	var count int
	err := r.db.QueryRowContext(ctx, query, noRawat).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *repository) SimpanResumePasien(ctx context.Context, noRawat, kodeDokter string, req SimpanResumePasienRequest) (*ResumePasien, error) {
	query := `
		INSERT INTO resume_pasien (
			no_rawat,
			kd_dokter,
			keluhan_utama,
			jalannya_penyakit,
			pemeriksaan_penunjang,
			hasil_laborat,
			diagnosa_utama,
			kd_diagnosa_utama,
			diagnosa_sekunder,
			kd_diagnosa_sekunder,
			diagnosa_sekunder2,
			kd_diagnosa_sekunder2,
			diagnosa_sekunder3,
			kd_diagnosa_sekunder3,
			diagnosa_sekunder4,
			kd_diagnosa_sekunder4,
			prosedur_utama,
			kd_prosedur_utama,
			prosedur_sekunder,
			kd_prosedur_sekunder,
			prosedur_sekunder2,
			kd_prosedur_sekunder2,
			prosedur_sekunder3,
			kd_prosedur_sekunder3,
			kondisi_pulang,
			obat_pulang
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		noRawat,
		kodeDokter,
		req.KeluhanUtama,
		req.JalannyaPenyakit,
		req.PemeriksaanPenunjang,
		req.HasilLaborat,
		req.DiagnosaUtama,
		req.KodeDiagnosaUtama,
		req.DiagnosaSekunder,
		req.KodeDiagnosaSekunder,
		req.DiagnosaSekunder2,
		req.KodeDiagnosaSekunder2,
		req.DiagnosaSekunder3,
		req.KodeDiagnosaSekunder3,
		req.DiagnosaSekunder4,
		req.KodeDiagnosaSekunder4,
		req.ProsedurUtama,
		req.KodeProsedurUtama,
		req.ProsedurSekunder,
		req.KodeProsedurSekunder,
		req.ProsedurSekunder2,
		req.KodeProsedurSekunder2,
		req.ProsedurSekunder3,
		req.KodeProsedurSekunder3,
		req.KondisiPulang,
		req.ObatPulang,
	)
	if err != nil {
		return nil, err
	}

	return r.DetailResumePasien(ctx, noRawat)
}

func (r *repository) UpdateResumePasien(ctx context.Context, noRawat string, req UpdateResumePasienRequest) (*ResumePasien, error) {
	query := `
		UPDATE resume_pasien SET
			keluhan_utama = ?,
			jalannya_penyakit = ?,
			pemeriksaan_penunjang = ?,
			hasil_laborat = ?,
			diagnosa_utama = ?,
			kd_diagnosa_utama = ?,
			diagnosa_sekunder = ?,
			kd_diagnosa_sekunder = ?,
			diagnosa_sekunder2 = ?,
			kd_diagnosa_sekunder2 = ?,
			diagnosa_sekunder3 = ?,
			kd_diagnosa_sekunder3 = ?,
			diagnosa_sekunder4 = ?,
			kd_diagnosa_sekunder4 = ?,
			prosedur_utama = ?,
			kd_prosedur_utama = ?,
			prosedur_sekunder = ?,
			kd_prosedur_sekunder = ?,
			prosedur_sekunder2 = ?,
			kd_prosedur_sekunder2 = ?,
			prosedur_sekunder3 = ?,
			kd_prosedur_sekunder3 = ?,
			kondisi_pulang = ?,
			obat_pulang = ?
		WHERE no_rawat = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		req.KeluhanUtama,
		req.JalannyaPenyakit,
		req.PemeriksaanPenunjang,
		req.HasilLaborat,
		req.DiagnosaUtama,
		req.KodeDiagnosaUtama,
		req.DiagnosaSekunder,
		req.KodeDiagnosaSekunder,
		req.DiagnosaSekunder2,
		req.KodeDiagnosaSekunder2,
		req.DiagnosaSekunder3,
		req.KodeDiagnosaSekunder3,
		req.DiagnosaSekunder4,
		req.KodeDiagnosaSekunder4,
		req.ProsedurUtama,
		req.KodeProsedurUtama,
		req.ProsedurSekunder,
		req.KodeProsedurSekunder,
		req.ProsedurSekunder2,
		req.KodeProsedurSekunder2,
		req.ProsedurSekunder3,
		req.KodeProsedurSekunder3,
		req.KondisiPulang,
		req.ObatPulang,
		noRawat,
	)
	if err != nil {
		return nil, err
	}

	return r.DetailResumePasien(ctx, noRawat)
}

func (r *repository) HapusResumePasien(ctx context.Context, noRawat string) error {
	query := `DELETE FROM resume_pasien WHERE no_rawat = ?`
	_, err := r.db.ExecContext(ctx, query, noRawat)
	return err
}
