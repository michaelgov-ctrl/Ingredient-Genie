package main

import (
	"net/http"

	"github.com/michaelgov-ctrl/Ingredient-Genie-frontend/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	// TODO: create meal endpoint
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /meal/view/{id}", app.mealView)
	mux.HandleFunc("GET /meal/search", app.searchMealsByIngredients)
	mux.HandleFunc("POST /meal/search", app.searchMealsByIngredientsPost)

	return app.recoverPanic(app.logRequest(commonHeaders(noSurf(mux))))
}
