package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/michaelgov-ctrl/Ingredient-Genie-frontend/internal/data"
)

type config struct {
	port        int
	logLevel    string
	mealsApiUri string
}

type application struct {
	config config
	logger *slog.Logger
	models data.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4270, "API server port")
	flag.StringVar(&cfg.logLevel, "log-level", "error", "Logging level (trace|debug|info|warning|error)")

	flag.StringVar(&cfg.mealsApiUri, "meals-api-uri", "http://localhost:4269", "URI of the meals API to target")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(logger, cfg.mealsApiUri),
	}

	if err := app.serve(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
