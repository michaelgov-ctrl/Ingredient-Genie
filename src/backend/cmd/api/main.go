package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/michaelgov-ctrl/Ingredient-Genie-backend/internal/data"
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

	meals, err := app.models.Meals.GetAllMeals()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	fmt.Println(len(meals), "\n", meals[len(meals)-1])
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
