package main

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/michaelgov-ctrl/Ingredient-Genie-frontend/internal/data"
	"github.com/michaelgov-ctrl/Ingredient-Genie-frontend/internal/validator"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	resp, err := app.models.Meals.GetMealList(data.Filters{
		Page:     1,
		PageSize: 20,
		Sort:     "name",
	})
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Meals = resp.Meals
	data.Metadata = resp.Metadata

	app.render(w, r, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) mealView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	meal, err := app.models.Meals.GetMeal(id)
	if err != nil {
		if errors.Is(err, data.ErrNoMeal) {
			http.NotFound(w, r)
			return
		}

		app.serverError(w, r, err)
		return
	}

	templateData := app.newTemplateData(r)
	templateData.Meal = meal

	app.render(w, r, http.StatusOK, "view.tmpl.html", templateData)
}

type mealSearchForm struct {
	Ingredients         []string `form:"ingredients"`
	Page                int      `form:"page"`
	PageSize            int      `form:"pagesize"`
	Sort                string   `form:"sort"`
	validator.Validator `form:"-"`
}

func (app *application) searchMealsByIngredients(w http.ResponseWriter, r *http.Request) {
	sortTypes, err := app.models.Meals.GetSortTypes()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// defaults
	form := mealSearchForm{
		Page:     1,
		PageSize: 10,
		Sort:     "-ratio",
	}

	query := r.URL.Query()

	// first visit
	if len(query) == 0 {
		templateData := app.newTemplateData(r)
		templateData.Form = form
		templateData.SortTypes = sortTypes

		app.render(w, r, http.StatusOK, "ingredientsearch.tmpl.html", templateData)

		return
	}

	form.Ingredients = normalizeIngredients(query["ingredients"])

	if value := query.Get("page"); value != "" {
		form.Page, err = strconv.Atoi(value)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
	}

	if value := query.Get("pagesize"); value != "" {
		form.PageSize, err = strconv.Atoi(value)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
	}

	if value := query.Get("sort"); value != "" {
		form.Sort = value
	}

	form.CheckField(validator.NotEmpty(form.Ingredients), "ingredients", "At least one ingredient must be provided")
	form.CheckField(form.Page >= 1, "page", "Page must be greater than zero")
	form.CheckField(form.PageSize >= 1 && form.PageSize <= 100, "pagesize", "Page size must be between 1 and 100")

	if !form.Valid() {
		templateData := app.newTemplateData(r)
		templateData.Form = form
		templateData.SortTypes = sortTypes

		app.render(w, r, http.StatusUnprocessableEntity, "ingredientsearch.tmpl.html", templateData)

		return
	}

	req := data.IngredientMealSearchRequest{
		Ingredients: form.Ingredients,
		Filters: data.Filters{
			Page:     form.Page,
			PageSize: form.PageSize,
			Sort:     data.SortType(form.Sort),
		},
	}

	resp, err := app.models.Meals.SearchByIngredients(req)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	templateData := app.newTemplateData(r)
	templateData.Form = form
	templateData.SortTypes = sortTypes
	templateData.MealMatches = resp.MealMatches
	templateData.Metadata = resp.Metadata

	app.render(w, r, http.StatusOK, "ingredientsearch.tmpl.html", templateData)
}

func (app *application) searchMealsByIngredientsPost(w http.ResponseWriter, r *http.Request) {
	var form mealSearchForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.Ingredients = normalizeIngredients(form.Ingredients)

	if form.PageSize == 0 {
		form.PageSize = 10
	}

	if form.Sort == "" {
		form.Sort = "-ratio"
	}

	// A new search always starts on page 1.
	form.Page = 1

	form.CheckField(validator.NotEmpty(form.Ingredients), "ingredients", "At least one ingredient must be provided")
	form.CheckField(form.PageSize >= 1 && form.PageSize <= 100, "pagesize", "Page size must be between 1 and 100")

	if !form.Valid() {
		sortTypes, err := app.models.Meals.GetSortTypes()
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		templateData := app.newTemplateData(r)
		templateData.Form = form
		templateData.SortTypes = sortTypes

		app.render(w, r, http.StatusUnprocessableEntity, "ingredientsearch.tmpl.html", templateData)

		return
	}

	query := url.Values{}

	for _, ingredient := range form.Ingredients {
		query.Add("ingredients", ingredient)
	}

	query.Set("page", strconv.Itoa(form.Page))
	query.Set("pagesize", strconv.Itoa(form.PageSize))
	query.Set("sort", form.Sort)

	http.Redirect(w, r, "/meal/search?"+query.Encode(), http.StatusSeeOther)
}

// This could just be done by the backend...
// But frontend validation is kind, even if redundant maybe
func normalizeIngredients(ingredients []string) []string {
	normalized := make([]string, 0, len(ingredients))
	seen := make(map[string]struct{}, len(ingredients))

	for _, ingredient := range ingredients {
		ingredient = strings.TrimSpace(ingredient)

		if ingredient == "" {
			continue
		}

		key := strings.ToLower(ingredient)

		if _, exists := seen[key]; exists {
			continue
		}

		seen[key] = struct{}{}
		normalized = append(normalized, ingredient)
	}

	return normalized
}
