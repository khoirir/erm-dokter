package shared

type PaginationMeta struct {
	TotalRecords int `json:"total_records"`
	TotalPages   int `json:"total_pages"`
	CurrentPage  int `json:"current_page"`
	PerPage      int `json:"per_page"`
}