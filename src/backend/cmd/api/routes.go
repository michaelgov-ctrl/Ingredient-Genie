package main

import (
	"net/http"
)

// # out of place testing note
//
//	$body = @{
//		"ingredients" = @(
//			"Garlic",
//			"Red Onions",
//			"Vegetable Oil",
//			"Lime"
//		)
//		"filters" = @{
//			"page" = 1
//			"pageSize" = 10
//			"sort" = "-ratio"
//		}
//	} | ConvertTo-Json
//
// $resp = irm http://localhost:4269/v1/meals/search -Method POST -Body $body -ContentType "application/json"
// # use iwr to check metadata
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)
	mux.HandleFunc("GET /v1/meals/search/sort", app.searchMealsSortTypesHandler)
	mux.HandleFunc("POST /v1/meals/search", app.searchMealsByIngredientsHandler)

	return app.recoverPanic(app.enableCORS(app.logRequest(app.rateLimit(mux))))
}
