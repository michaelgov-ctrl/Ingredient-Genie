package data

import "database/sql"

type Models struct {
	Meals MealModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Meals: MealModel{DB: db},
	}
}
