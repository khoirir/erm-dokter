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

