package data

import (
	"log/slog"
)

type MealsApi interface {
	GetMeal(int) (Meal, error)
	GetMealList(Filters) (MealListResponse, error)
	GetSortTypes() ([]SortType, error)
	SearchByIngredients(IngredientMealSearchRequest) (MealSearchResponse, error)
}

type Models struct {
	Meals MealsApi
}

func NewModels(logger *slog.Logger, mealApiAddr string) Models {
	return Models{
		Meals: NewMealsClient(logger, mealApiAddr),
	}
}

type MealMatch struct {
	Meal                   Meal     `json:"meal"`
	MissingIngredients     []string `json:"missingIngredients"`
	MatchedIngredientCount int64    `json:"matchedIngredientCount"`
	TotalIngredientCount   int64    `json:"totalIngredientCount"`
	MatchRatio             float64  `json:"matchRatio"`
}

type Meal struct {
	ID            int64            `json:"id"`
	Name          string           `json:"name"`
	AlternateName string           `json:"alternateName"`
	Category      string           `json:"category"`
	Area          string           `json:"area"`
	Country       string           `json:"country"`
	Instructions  string           `json:"instructions"`
	ThumbnailURL  string           `json:"thumbnailUrl"`
	YoutubeURL    string           `json:"youtubeUrl"`
	SourceURL     string           `json:"sourceUrl"`
	Ingredients   []MealIngredient `json:"ingredients"`
}

type MealIngredient struct {
	IngredientID int64  `json:"ingredientId"`
	Name         string `json:"name"`
	Position     int64  `json:"position"`
	MeasureText  string `json:"measureText"`
}

type Metadata struct {
	CurrentPage  int `json:"current_page"`
	PageSize     int `json:"page_size"`
	FirstPage    int `json:"first_page"`
	LastPage     int `json:"last_page"`
	TotalRecords int `json:"total_records"`
}

type IngredientMealSearchRequest struct {
	Ingredients []string `json:"ingredients"`
	Filters     Filters  `json:"filters"`
}

type MealListResponse struct {
	Meals    []Meal   `json:"meals"`
	Metadata Metadata `json:"metadata"`
}

type MealSearchResponse struct {
	MealMatches []MealMatch `json:"mealMatches"`
	Metadata    Metadata    `json:"metadata"`
}
