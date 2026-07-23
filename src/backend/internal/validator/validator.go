package validator

import (
	"regexp"
	"slices"
	"strings"
)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, msg string) {
	if _, ok := v.Errors[key]; !ok {
		v.Errors[key] = msg
	}
}

func (v *Validator) Check(ok bool, key, msg string) {
	if !ok {
		v.AddError(key, msg)
	}
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func Unique[T comparable](values []T) bool {
	uniq := make(map[T]struct{})
	for _, v := range values {
		uniq[v] = struct{}{}
	}

	return len(uniq) == len(values)
}

func ValidateIngredientSearch(v *Validator, ingredients []string) {
	v.Check(len(ingredients) > 0, "ingredients", "must contain at least one ingredient")
	v.Check(len(ingredients) <= 20, "ingredients", "must contain no more than 20 ingredients")

	for _, ingredient := range ingredients {
		ingredient = strings.TrimSpace(ingredient)
		v.Check(ingredient != "", "ingredients", "must not contain empty values")
		v.Check(len(ingredient) <= 100, "ingredients", "must not contain values longer than 100 characters")
	}
}
