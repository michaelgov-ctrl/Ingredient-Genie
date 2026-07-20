package data

import (
	"github.com/michaelgov-ctrl/Ingredient-Genie-backend/internal/validator"
)

type Filters struct {
	Page     int
	PageSize int
	Sort     string
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than 0")
	v.Check(f.Page <= 1_000, "page", "must be less than or equal to 1,000") // this may cause an unexpected issue if we surpass 1000 pages, unlikely..
	v.Check(f.PageSize > 0, "page_size", "must be greater than 0")
	v.Check(f.PageSize <= 10, "page_size", "must be less than or equal to 10")

	_, validSort := mealSortStmts[f.Sort]
	v.Check(validSort, "sort", "invalid sort value")
}

var mealSortStmts = map[string]string{
	"id":     "MealId ASC",
	"-id":    "MealId DESC",
	"ratio":  "MatchRatio ASC, MealId ASC",
	"-ratio": "MatchRatio DESC, MealId ASC",
	"name":   "Name COLLATE NOCASE ASC, MealId ASC",
	"-name":  "Name COLLATE NOCASE DESC, MealId ASC",
}

func (f Filters) orderBy() string {
	expression, ok := mealSortStmts[f.Sort]
	if !ok {
		panic("unsafe sort parameter: " + f.Sort)
	}

	return expression
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

type Metadata struct {
	CurrentPage  int `json:"current_page"`
	PageSize     int `json:"page_size"`
	FirstPage    int `json:"first_page"`
	LastPage     int `json:"last_page"`
	TotalRecords int `json:"total_records"`
}

func calculateMetadata(totalMatchingRecords, page, pageSize int) Metadata {
	metadata := Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		TotalRecords: totalMatchingRecords,
	}

	if totalMatchingRecords == 0 {
		return metadata
	}

	metadata.FirstPage = 1
	metadata.LastPage = (totalMatchingRecords + pageSize - 1) / pageSize

	return metadata
}
