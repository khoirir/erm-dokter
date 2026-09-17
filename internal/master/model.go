package master

import (
	"strings"
	"time"

	"erm-dokter/internal/shared/apperror"
)

type ItemMaster struct {
	Kode string `json:"kode"`
	Nama string `json:"nama"`
}

type Penjamin struct {
	ItemMaster
}

type Depo struct {
	ItemMaster
}

type Poliklinik struct {
	ItemMaster
}

type Bangsal struct {
	ItemMaster
}

type KelasKamar struct {
	ItemMaster
}

type ICD10 struct {
	ItemMaster
}

type ICD9 struct {
	ItemMaster
}

type FilterMasterICD struct {
	Keyword string `json:"keyword,omitempty"`
	Page    int    `json:"page,omitempty"`
	Limit   int    `json:"limit,omitempty"`
}

func (f *FilterMasterICD) Sanitize() {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Limit <= 0 {
		f.Limit = 20
	} else if f.Limit > 100 {
		f.Limit = 100
	}
	f.Keyword = strings.TrimSpace(f.Keyword)
}

func (f FilterMasterICD) Offset() int {
	return (f.Page - 1) * f.Limit
}

func (f FilterMasterICD) Validate() apperror.ValidationError {
	errs := make(apperror.ValidationError)

	if f.Keyword != "" && len(f.Keyword) < 3 {
		errs["keyword"] = "Kata kunci pencarian minimal 3 karakter"
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

type SyncICDResult struct {
	TotalICD10 int       `json:"total_icd10"`
	TotalICD9  int       `json:"total_icd9"`
	SyncedAt   time.Time `json:"synced_at"`
}


