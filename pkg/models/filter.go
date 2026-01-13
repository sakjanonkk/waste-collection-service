package models

type Search struct {
	Keyword string `query:"keyword" json:"keyword"`
	Column  string `query:"column" json:"column"`
}

func (s *Search) GetSearchString() string {
	return "keyword=" + s.Keyword + "&column=" + s.Column
}

// StaffFilter contains filter parameters for staff list
type StaffFilter struct {
	Search string `query:"search"`
	Role   string `query:"role"`
	Status string `query:"status"`
}
