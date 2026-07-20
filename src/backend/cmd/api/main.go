package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/michaelgov-ctrl/Ingredient-Genie-backend/internal/data"
	"github.com/michaelgov-ctrl/Ingredient-Genie-backend/internal/validator"
	_ "modernc.org/sqlite"
)

type config struct {
	db struct {
		dsn string
	}
}

type application struct {
	config config
	logger *slog.Logger
	models data.Models
}

func main() {
	var cfg config
	cfg.db.dsn = "internal/data/meals.sqlite"

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	app.searchMealsHandler()
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("sqlite", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// starting to prep to transition to a json api
func (app *application) searchMealsHandler( /*w http.ResponseWriter, r *http.Request*/ ) {
	var input struct {
		Ingredients []string
		data.Filters
	}

	v := validator.New()

	// extract necessary info like ingredients, page, page size, and sort from `r.Body`
	input.Ingredients = []string{
		"Garlic",
		"Red Onions",
		"Vegetable Oil",
		"Lime",
	}
	input.Filters.Page = 1
	input.Filters.PageSize = 10
	input.Filters.Sort = "-ratio"

	if data.ValidateIngredientSearch(v, input.Ingredients); !v.Valid() {
		// failed validation
		return
	}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		// failed validation
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	meals, metadata, err := app.models.Meals.FindByIngredients(ctx, input.Ingredients, input.Filters)
	if err != nil {
		// return err
		return
	}

	// write the response

	fmt.Println(metadata, len(meals), "\n", meals[len(meals)-1])
}
