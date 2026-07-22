package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/michaelgov-ctrl/Ingredient-Genie-frontend/internal/data"
)

func (app *application) searchMealsByIngredientsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: provide a pretty form to actually process user input for a request
	req := data.IngredientMealSearchRequest{
		Ingredients: []string{"Garlic"},
		Filters: data.Filters{
			Page:     1,
			PageSize: 10,
			Sort:     "-ratio",
		},
	}

	meals, metadata, err := app.models.Meals.SearchByIngredients(req)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	mealsJson, err := json.Marshal(meals)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	metadataJson, err := json.Marshal(metadata)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Fprint(w, string(mealsJson), string(metadataJson))
}
