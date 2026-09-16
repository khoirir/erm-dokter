package rujukaninternal

import (
	"errors"
	"fmt"
	"strings"

	"erm-dokter/internal/shared/apperror"
)

type IdOpsiPoliDokter struct {
	KodePoli   string
	KodeDokter string
}

func (id IdOpsiPoliDokter) CompositeKey() string {
	return fmt.Sprintf("%s~%s", id.KodePoli, id.KodeDokter)
}

func ParseIdOpsiPoliDokter(decryptedKey string) (IdOpsiPoliDokter, error) {
	parts := strings.Split(decryptedKey, "~")
	if len(parts) != 2 {
		return IdOpsiPoliDokter{}, errors.New("ID poli rujukan tidak valid")
	}
	return IdOpsiPoliDokter{
		KodePoli:   strings.TrimSpace(parts[0]),
		KodeDokter: strings.TrimSpace(parts[1]),
	}, nil
}

type InfoPoliDokter struct {
	KodePoli   string `json:"kode_poli"`
	NamaPoli   string `json:"nama_poli"`
	KodeDokter string `json:"kode_dokter"`
	NamaDokter string `json:"nama_dokter"`
}

type OpsiPoliDokter struct {
	Id string `json:"id"`
	InfoPoliDokter
}

func (o *OpsiPoliDokter) CompositeKey() string {
	return fmt.Sprintf("%s~%s", o.KodePoli, o.KodeDokter)
}

type IdRujukanInternal struct {
	NoRawat    string
	KodePoli   string
	KodeDokter string
}

func (id IdRujukanInternal) CompositeKey() string {
	return fmt.Sprintf("%s~%s~%s", id.NoRawat, id.KodePoli, id.KodeDokter)
}

func ParseIdRujukanInternal(decryptedKey string) (IdRujukanInternal, error) {
	parts := strings.Split(decryptedKey, "~")
	if len(parts) != 3 {
		return IdRujukanInternal{}, errors.New("ID rujukan tidak valid")
	}
	return IdRujukanInternal{
		NoRawat:    strings.TrimSpace(parts[0]),
		KodePoli:   strings.TrimSpace(parts[1]),
		KodeDokter: strings.TrimSpace(parts[2]),
	}, nil
}

type RujukanInternal struct {
	Id          string `json:"id"`
	IdKunjungan string `json:"id_kunjungan"`
	NoRawat     string `json:"no_rawat"`
	InfoPoliDokter
}

func (r *RujukanInternal) CompositeKey() string {
	return fmt.Sprintf("%s~%s~%s", r.NoRawat, r.KodePoli, r.KodeDokter)
}

type SimpanRujukanRequest struct {
	IdTujuan string `json:"id_tujuan"`
}

func (req *SimpanRujukanRequest) Sanitize() {
	req.IdTujuan = strings.TrimSpace(req.IdTujuan)
}

func (req *SimpanRujukanRequest) Validate() apperror.ValidationError {
	errs := make(apperror.ValidationError)

	if req.IdTujuan == "" {
		errs["id_tujuan"] = "Tujuan rujukan dokter dan poliklinik wajib dipilih"
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}
