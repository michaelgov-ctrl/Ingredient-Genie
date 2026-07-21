package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/michaelgov-ctrl/Ingredient-Genie-backend/internal/data"
	"github.com/michaelgov-ctrl/Ingredient-Genie-backend/internal/validator"
)

type envelope map[string]any

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	if err := app.writeJSON(w, http.StatusOK, envelope{"status": "available"}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) searchMealsByIngredientsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Ingredients  []string `json:"ingredients"`
		data.Filters `json:"filters"`
	}

	// extract necessary info like ingredients, page, page size, and sort from `r.Body`
	// input.Ingredients = []string{
	// 	"Garlic",
	// 	"Red Onions",
	// 	"Vegetable Oil",
	// 	"Lime",
	// }
	// input.Filters.Page = 1
	// input.Filters.PageSize = 10
	// input.Filters.Sort = "-ratio"
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateIngredientSearch(v, input.Ingredients); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	meals, metadata, err := app.models.Meals.FindByIngredients(ctx, input.Ingredients, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// fmt.Println(metadata, len(meals), "\n", meals[len(meals)-1])
	if err := app.writeJSON(w, http.StatusOK, envelope{"metadata": metadata, "meals": meals}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	var max_bytes int64 = 1_048_576 // 1mb
	r.Body = http.MaxBytesReader(w, r.Body, max_bytes)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		var (
			syntaxError           *json.SyntaxError
			unmarshalTypeError    *json.UnmarshalTypeError
			invalidUnmarshalError *json.InvalidUnmarshalError
			maxBytesError         *http.MaxBytesError
		)

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON at type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldname := strings.TrimPrefix(err.Error(), "json: unknown")
			return fmt.Errorf("body contains unknown key %s", fieldname)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes, that's one chonky meme", maxBytesError.Limit)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for k, v := range headers {
		w.Header()[k] = v
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
