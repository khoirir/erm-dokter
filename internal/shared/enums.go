package shared

type StatusLanjut string

const (
	StatusLanjutRawatJalan StatusLanjut = "Ralan"
	StatusLanjutRawatInap  StatusLanjut = "Ranap"
)

func (s StatusLanjut) IsValid() bool {
	switch s {
	case StatusLanjutRawatJalan, StatusLanjutRawatInap:
		return true
	default:
		return false
	}
}

type SortOrder string

const (
	SortASC  SortOrder = "ASC"
	SortDESC SortOrder = "DESC"
)

func (s SortOrder) IsValid() bool {
	switch s {
	case SortASC, SortDESC:
		return true
	default:
		return false
	}
}

type KategoriLab string

const (
	KategoriLabPK KategoriLab = "PK"
	KategoriLabPA KategoriLab = "PA"
	KategoriLabMB KategoriLab = "MB"
)

func (k KategoriLab) IsValid() bool {
	switch k {
	case KategoriLabPK, KategoriLabPA, KategoriLabMB:
		return true
	default:
		return false
	}
}
