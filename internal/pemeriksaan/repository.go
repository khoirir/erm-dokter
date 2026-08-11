package pemeriksaan

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	DaftarPemeriksaan(ctx context.Context, noRawat string, statusLanjut string) ([]Pemeriksaan, error)
	DetailPemeriksaan(ctx context.Context, noRawat string, tanggalPemriksaan string, jamPemeriksaan string) (*Pemeriksaan, error)
	// UpdatePemeriksaan(ctx context.Context, noRawat string, pemeriksaan *Pemeriksaan) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

const selectPemeriksaanRalan = `
	SELECT 
		pr.no_rawat,
		DATE_FORMAT(pr.tgl_perawatan, '%Y-%m-%d') AS tanggal_pemeriksaan,
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
	WHERE pr.no_rawat = ?
`

const selectPemeriksaanRanap = `
	SELECT 
		pr.no_rawat,
		DATE_FORMAT(pr.tgl_perawatan, '%Y-%m-%d') AS tanggal_pemeriksaan,
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
	WHERE pr.no_rawat = ?
`

func (r *repository) DaftarPemeriksaan(ctx context.Context, noRawat string, statusLanjut string) ([]Pemeriksaan, error) {
	var query string
	var args []interface{}

	switch statusLanjut {
	case "rawat_jalan":
		query = selectPemeriksaanRalan + " ORDER BY pr.tgl_perawatan DESC, pr.jam_rawat DESC"
		args = append(args, noRawat)
	case "rawat_inap":
		query = selectPemeriksaanRanap + " ORDER BY pr.tgl_perawatan DESC, pr.jam_rawat DESC"
		args = append(args, noRawat)
	default:
		query = fmt.Sprintf(`
			SELECT * FROM (
				%s 
				UNION ALL 
				%s
			) AS t 
			ORDER BY t.tanggal_pemeriksaan DESC, t.jam_pemeriksaan DESC
		`, selectPemeriksaanRalan, selectPemeriksaanRanap)
		args = append(args, noRawat, noRawat)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("gagal query daftar pemeriksaan: %w", err)
	}
	defer rows.Close()
	var daftarPemeriksaan []Pemeriksaan
	for rows.Next() {
		var p Pemeriksaan
		err := rows.Scan(
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
			return nil, fmt.Errorf("gagal scan data pemeriksaan: %w", err)
		}
		daftarPemeriksaan = append(daftarPemeriksaan, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterasi data pemeriksaan: %w", err)
	}
	return daftarPemeriksaan, nil
}

func (r *repository) DetailPemeriksaan(ctx context.Context, noRawat string, tanggalPemeriksaan string, jamPemeriksaan string) (*Pemeriksaan, error) {
	return nil, nil
}
