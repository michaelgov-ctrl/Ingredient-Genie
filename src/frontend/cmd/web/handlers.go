package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/michaelgov-ctrl/Ingredient-Genie-frontend/internal/data"
	"github.com/michaelgov-ctrl/Ingredient-Genie-frontend/internal/validator"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	resp, err := app.models.Meals.GetMealList("")
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Memes = resp.Meals
	data.Metadata = resp.Metadata

	app.render(w, r, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) mealView(w http.ResponseWriter, r *http.Request) {
	// TODO: do
}

func (app *application) mealCreate(w http.ResponseWriter, r *http.Request) {
	// TODO: do
}

func (app *application) mealCreatePost(w http.ResponseWriter, r *http.Request) {
	// TODO: do
}

type mealSearchForm struct {
	Name                string        `form:"ingredients"`
	Sort                data.SortType `form:"sort"`
	validator.Validator `form:"-"`
}

func (app *application) searchMealsByIngredients(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = mealSearchForm{}
	app.render(w, r, http.StatusOK, "ingredientsearch.tmpl.html", data)
}

func (app *application) searchMealsByIngredientsPost(w http.ResponseWriter, r *http.Request) {
	// gather sort types to populate form..
	sortTypes, err := app.models.Meals.GetSearchSortTypes()
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	fmt.Println(sortTypes)

	// TODO: provide a pretty form to actually process user input for a request
	req := data.IngredientMealSearchRequest{
		Ingredients: []string{"Garlic"},
		Filters: data.Filters{
			Page:     1,
			PageSize: 10,
			Sort:     "-ratio",
		},
	}

	resp, err := app.models.Meals.SearchByIngredients(req)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	bytes, err := json.Marshal(resp.Meals)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Fprint(w, string(bytes))
}
