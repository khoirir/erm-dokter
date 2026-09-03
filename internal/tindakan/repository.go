package tindakan

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"erm-dokter/internal/shared"
)

type Repository interface {
	DaftarTindakanLab(ctx context.Context, kategori shared.KategoriLab, filter FilterDaftarTindakanLab) ([]TindakanLab, int, error)
	GetDetailTindakanLab(ctx context.Context, kategori shared.KategoriLab, kodeTindakan string) (*TindakanLab, []TemplateLabDB, error)
	CekKeberadaanTindakanLab(ctx context.Context, kategori shared.KategoriLab, listKodeTindakan []string) (map[string]bool, error)
	CekKeberadaanTemplateLab(ctx context.Context, listKodeTindakan []string, templateMap map[string][]int) (map[string]map[int]bool, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) DaftarTindakanLab(ctx context.Context, kategori shared.KategoriLab, filter FilterDaftarTindakanLab) ([]TindakanLab, int, error) {
	var conditions []string
	var args []any

	conditions = append(conditions, "status = '1'")
	conditions = append(conditions, "kategori = ?")
	args = append(args, string(kategori))

	if filter.Keyword != "" {
		conditions = append(conditions, "(kd_jenis_prw LIKE ? OR nm_perawatan LIKE ?)")
		args = append(args, "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
	}

	whereSQL := "WHERE " + strings.Join(conditions, " AND ")

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM jns_perawatan_lab
		%s
	`, whereSQL)

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []TindakanLab{}, 0, nil
	}

	filter.Sanitize()
	dataQuery := fmt.Sprintf(`
		SELECT 
			kd_jenis_prw,
			nm_perawatan,
			total_byr
		FROM jns_perawatan_lab
		%s
		ORDER BY nm_perawatan ASC
		LIMIT ? OFFSET ?
	`, whereSQL)

	argsWithPaging := append(args, filter.Limit, filter.Offset())
	rows, err := r.db.QueryContext(ctx, dataQuery, argsWithPaging...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []TindakanLab
	for rows.Next() {
		var item TindakanLab
		err = rows.Scan(&item.KodeTindakan, &item.NamaTindakan, &item.Biaya)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, item)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *repository) GetDetailTindakanLab(ctx context.Context, kategori shared.KategoriLab, kodeTindakan string) (*TindakanLab, []TemplateLabDB, error) {
	tindakanQuery := `
		SELECT 
			kd_jenis_prw,
			nm_perawatan,
			total_byr
		FROM jns_perawatan_lab
		WHERE status = '1' AND kategori = ? AND kd_jenis_prw = ?
	`
	var t TindakanLab
	err := r.db.QueryRowContext(ctx, tindakanQuery, string(kategori), kodeTindakan).Scan(
		&t.KodeTindakan,
		&t.NamaTindakan,
		&t.Biaya,
	)
	if err != nil {
		return nil, nil, err
	}

	templatesQuery := `
		SELECT 
			id_template,
			kd_jenis_prw,
			Pemeriksaan AS nama_pemeriksaan,
			satuan,
			nilai_rujukan_ld,
			nilai_rujukan_la,
			nilai_rujukan_pd,
			nilai_rujukan_pa
		FROM template_laboratorium
		WHERE kd_jenis_prw = ?
		ORDER BY urut ASC, id_template ASC
	`
	rows, err := r.db.QueryContext(ctx, templatesQuery, kodeTindakan)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var templates []TemplateLabDB
	for rows.Next() {
		var tmpl TemplateLabDB
		err = rows.Scan(
			&tmpl.IdTemplate,
			&tmpl.KodeTindakan,
			&tmpl.NamaPemeriksaan,
			&tmpl.Satuan,
			&tmpl.NilaiRujukanLD,
			&tmpl.NilaiRujukanLA,
			&tmpl.NilaiRujukanPD,
			&tmpl.NilaiRujukanPA,
		)
		if err != nil {
			return nil, nil, err
		}
		templates = append(templates, tmpl)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, err
	}

	if templates == nil {
		templates = []TemplateLabDB{}
	}

	return &t, templates, nil
}

func (r *repository) CekKeberadaanTindakanLab(ctx context.Context, kategori shared.KategoriLab, listKodeTindakan []string) (map[string]bool, error) {
	result := make(map[string]bool)
	if len(listKodeTindakan) == 0 {
		return result, nil
	}

	placeholders := make([]string, len(listKodeTindakan))
	args := make([]any, 0, len(listKodeTindakan)+1)
	args = append(args, string(kategori))
	for i, kode := range listKodeTindakan {
		placeholders[i] = "?"
		args = append(args, kode)
	}

	query := fmt.Sprintf(`
		SELECT kd_jenis_prw
		FROM jns_perawatan_lab
		WHERE status = '1' AND kategori = ? AND kd_jenis_prw IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var kode string
		if err := rows.Scan(&kode); err != nil {
			return nil, err
		}
		result[kode] = true
	}

	return result, rows.Err()
}

func (r *repository) CekKeberadaanTemplateLab(ctx context.Context, listKodeTindakan []string, templateMap map[string][]int) (map[string]map[int]bool, error) {
	result := make(map[string]map[int]bool)
	for _, kode := range listKodeTindakan {
		result[kode] = make(map[int]bool)
	}

	var conditions []string
	var args []any

	for kodeTindakan, idList := range templateMap {
		if len(idList) == 0 {
			continue
		}
		placeholders := make([]string, len(idList))
		args = append(args, kodeTindakan)
		for i, id := range idList {
			placeholders[i] = "?"
			args = append(args, id)
		}
		conditions = append(conditions, fmt.Sprintf("(kd_jenis_prw = ? AND id_template IN (%s))", strings.Join(placeholders, ",")))
	}

	if len(conditions) == 0 {
		return result, nil
	}

	query := fmt.Sprintf(`
		SELECT kd_jenis_prw, id_template
		FROM template_laboratorium
		WHERE %s
	`, strings.Join(conditions, " OR "))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var kode string
		var id int
		if err := rows.Scan(&kode, &id); err != nil {
			return nil, err
		}
		if result[kode] == nil {
			result[kode] = make(map[int]bool)
		}
		result[kode][id] = true
	}

	return result, rows.Err()
}

