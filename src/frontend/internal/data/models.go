package data

import (
	"log/slog"
	"net"
)

type MealsApi interface {
	SearchByIngredients() (MealResponse, Metadata, error)
}

type Models struct {
	MealsApi MealsApi
}

func NewModels(logger *slog.Logger, mealApiAddr net.Addr) Models {
	return Models{
		MealsApi: NewMealsApiClient(logger, mealApiAddr),
	}
}

type MealResponse struct {
	Meal                   Meal     `json:"meal"`
	MissingIngredients     []string `json:"missingIngredients"`
	MatchedIngredientCount int64    `json:"matchedIngredientCount"`
	TotalIngredientCount   int64    `json:"totalIngredientCount"`
	MatchRatio             float64  `json:"matchRatio"`
}

type Meal struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	AlternateName string `json:"alternateName"`
	Category      string `json:"category"`
	Area          string `json:"area"`
	Country       string `json:"country"`
	Instructions  string `json:"instructions"`
	ThumbnailUrl  string `json:"thumbnailUrl"`
	YoutubeUrl    string `json:"youtubeUrl"`
	SourceUrl     string `json:"sourceUrl"`
}

type Metadata struct {
	CurrentPage  int `json:"current_page"`
	PageSize     int `json:"page_size"`
	FirstPage    int `json:"first_page"`
	LastPage     int `json:"last_page"`
	TotalRecords int `json:"total_records"`
}
