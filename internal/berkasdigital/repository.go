package berkasdigital

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Repository interface {
	FetchMasterBerkas(ctx context.Context) ([]MasterBerkasDigital, error)
	FetchBerkasByNoRawat(ctx context.Context, noRawat string, kodeList []string) ([]BerkasDigitalDB, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FetchMasterBerkas(ctx context.Context) ([]MasterBerkasDigital, error) {
	query := `
		SELECT kode, nama 
		FROM master_berkas_digital 
		ORDER BY kode ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []MasterBerkasDigital
	for rows.Next() {
		var item MasterBerkasDigital
		if err := rows.Scan(&item.Kode, &item.Nama); err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if list == nil {
		list = []MasterBerkasDigital{}
	}

	return list, nil
}

func (r *repository) FetchBerkasByNoRawat(ctx context.Context, noRawat string, kodeList []string) ([]BerkasDigitalDB, error) {
	if len(kodeList) == 0 {
		return []BerkasDigitalDB{}, nil
	}

	placeholders := make([]string, len(kodeList))
	args := make([]any, 0, len(kodeList)+1)
	args = append(args, noRawat)

	for i, k := range kodeList {
		placeholders[i] = "?"
		args = append(args, k)
	}

	query := fmt.Sprintf(`
		SELECT 
			bdp.no_rawat,
			bdp.kode,
			COALESCE(mbd.nama, bdp.kode) AS nama_berkas,
			COALESCE(bdp.lokasi_file, '') AS lokasi_file
		FROM berkas_digital_perawatan bdp
		LEFT JOIN master_berkas_digital mbd ON mbd.kode = bdp.kode
		WHERE bdp.no_rawat = ? AND bdp.kode IN (%s)
		ORDER BY bdp.kode ASC
	`, strings.Join(placeholders, ","))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []BerkasDigitalDB
	for rows.Next() {
		var item BerkasDigitalDB
		err := rows.Scan(&item.NoRawat, &item.Kode, &item.NamaBerkas, &item.LokasiFile)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if list == nil {
		list = []BerkasDigitalDB{}
	}

	return list, nil
}
