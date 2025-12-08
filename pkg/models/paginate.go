package models

type Pagination struct {
	Page    int   `query:"page" json:"page"`
	PerPage int   `query:"per_page" json:"per_page"`
	Total   int64 `json:"total"`
}
