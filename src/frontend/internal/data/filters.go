package data

type Filters struct {
	Page     int        `json:"page"`
	PageSize int        `json:"pageSize"`
	Sort     FilterType `json:"sort"`
}

type FilterType string
