package data

import (
	"context"
	"database/sql"
	"time"
)

const (
	MealsTable          = "Meal"
	IngredientTable     = "Ingredient"
	MealIngredientTable = "MealIngredient"
)

type MealModel struct {
	DB *sql.DB
}

type Meal struct {
	ID            int64  `json:"id"`
	MealDBID      int64  `json:"mealDbId"`
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

type SqlSafeMeal struct {
	ID            sql.NullInt64  `db:"MealId"`
	MealDBID      sql.NullInt64  `db:"ExternalMealId"`
	Name          sql.NullString `db:"Name"`
	AlternateName sql.NullString `db:"AlternateName"`
	Category      sql.NullString `db:"Category"`
	Area          sql.NullString `db:"Area"`
	Country       sql.NullString `db:"Country"`
	Instructions  sql.NullString `db:"Instructions"`
	ThumbnailUrl  sql.NullString `db:"ThumbnailUrl"`
	YoutubeUrl    sql.NullString `db:"YoutubeUrl"`
	SourceUrl     sql.NullString `db:"SourceUrl"`
}

func (ssm SqlSafeMeal) ToMeal() Meal {
	return Meal{
		ID:            ssm.ID.Int64,
		MealDBID:      ssm.MealDBID.Int64,
		Name:          ssm.Name.String,
		AlternateName: ssm.AlternateName.String,
		Category:      ssm.Category.String,
		Area:          ssm.Area.String,
		Country:       ssm.Country.String,
		Instructions:  ssm.Instructions.String,
		ThumbnailUrl:  ssm.ThumbnailUrl.String,
		YoutubeUrl:    ssm.YoutubeUrl.String,
		SourceUrl:     ssm.SourceUrl.String,
	}
}

func (mm MealModel) GetAllMeals() ([]Meal, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT * FROM Meal"
	rows, err := mm.DB.QueryContext(ctx, query)
	if err != nil {
		return []Meal{}, nil
	}
	defer rows.Close()

	var Meals []Meal
	for rows.Next() {
		var ssm SqlSafeMeal
		if err := rows.Scan(
			&ssm.ID,
			&ssm.MealDBID,
			&ssm.Name,
			&ssm.AlternateName,
			&ssm.Category,
			&ssm.Area,
			&ssm.Country,
			&ssm.Instructions,
			&ssm.ThumbnailUrl,
			&ssm.YoutubeUrl,
			&ssm.SourceUrl,
		); err != nil {
			return []Meal{}, nil
		}

		Meals = append(Meals, ssm.ToMeal())
	}

	return Meals, nil
}
