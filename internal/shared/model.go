package shared

import "math"

type PaginationMeta struct {
	TotalRecords int `json:"total_records"`
	TotalPages   int `json:"total_pages"`
	CurrentPage  int `json:"current_page"`
	PerPage      int `json:"per_page"`
}

func NewPaginationMeta(totalRecords, page, limit int) PaginationMeta {
	if limit <= 0 {
		limit = 20
	}
	totalPages := int(math.Ceil(float64(totalRecords) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}
	return PaginationMeta{
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		CurrentPage:  page,
		PerPage:      limit,
	}
}
