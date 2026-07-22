package data

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/michaelgov-ctrl/Ingredient-Genie-backend/internal/validator"
)

const (
	MealsTable          = "Meal"
	IngredientTable     = "Ingredient"
	MealIngredientTable = "MealIngredient"
)

type MealModel struct {
	DB *sql.DB
}

type Meal struct {
	ID            int64            `json:"id"`
	Name          string           `json:"name"`
	AlternateName string           `json:"alternateName"`
	Category      string           `json:"category"`
	Area          string           `json:"area"`
	Country       string           `json:"country"`
	Instructions  string           `json:"instructions"`
	ThumbnailURL  string           `json:"thumbnailUrl"`
	YoutubeURL    string           `json:"youtubeUrl"`
	SourceURL     string           `json:"sourceUrl"`
	Ingredients   []MealIngredient `json:"ingredients"`
}

func (m Meal) ToSqlSafeMeal() SqlSafeMeal {
	var ssm SqlSafeMeal

	ssm.ID = m.ID
	ssm.Name = StringToSqlNullString(m.Name)
	ssm.AlternateName = StringToSqlNullString(m.AlternateName)
	ssm.Category = StringToSqlNullString(m.Category)
	ssm.Area = StringToSqlNullString(m.Area)
	ssm.Country = StringToSqlNullString(m.Country)
	ssm.Instructions = StringToSqlNullString(m.Instructions)
	ssm.ThumbnailURL = StringToSqlNullString(m.ThumbnailURL)
	ssm.YoutubeURL = StringToSqlNullString(m.YoutubeURL)
	ssm.SourceURL = StringToSqlNullString(m.SourceURL)

	return ssm
}

type SqlSafeMeal struct {
	ID            int64          `db:"MealId"`
	Name          sql.NullString `db:"Name"`
	AlternateName sql.NullString `db:"AlternateName"`
	Category      sql.NullString `db:"Category"`
	Area          sql.NullString `db:"Area"`
	Country       sql.NullString `db:"Country"`
	Instructions  sql.NullString `db:"Instructions"`
	ThumbnailURL  sql.NullString `db:"ThumbnailUrl"`
	YoutubeURL    sql.NullString `db:"YoutubeUrl"`
	SourceURL     sql.NullString `db:"SourceUrl"`
}

func (ssm SqlSafeMeal) ToMeal() Meal {
	return Meal{
		ID:            ssm.ID,
		Name:          ssm.Name.String,
		AlternateName: ssm.AlternateName.String,
		Category:      ssm.Category.String,
		Area:          ssm.Area.String,
		Country:       ssm.Country.String,
		Instructions:  ssm.Instructions.String,
		ThumbnailURL:  ssm.ThumbnailURL.String,
		YoutubeURL:    ssm.YoutubeURL.String,
		SourceURL:     ssm.SourceURL.String,
		Ingredients:   make([]MealIngredient, 0),
	}
}

type Ingredient struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	NormalizedName string `json:"normalizedName"`
}

func (i Ingredient) ToSqlSafeIngredient() SqlSafeIngredient {
	var ssi SqlSafeIngredient

	ssi.ID = i.ID
	ssi.Name = StringToSqlNullString(i.Name)
	ssi.NormalizedName = StringToSqlNullString(i.NormalizedName)

	return ssi
}

type SqlSafeIngredient struct {
	ID             int64          `db:"IngredientId"`
	Name           sql.NullString `db:"Name"`
	NormalizedName sql.NullString `db:"NormalizedName"`
}

func (ssi SqlSafeIngredient) ToIngredient() Ingredient {
	return Ingredient{
		ID:             ssi.ID,
		Name:           ssi.Name.String,
		NormalizedName: ssi.NormalizedName.String,
	}
}

type MealIngredient struct {
	IngredientID   int64  `json:"ingredientId"`
	Name           string `json:"name"`
	Position       int64  `json:"position"`
	MeasureText    string `json:"measureText"`
	NormalizedName string `json:"-"`
}

type MealIngredientRecord struct {
	MealID       int64
	IngredientID int64
	Position     int64
	MeasureText  string
}

func (mi MealIngredientRecord) ToSqlSafeMealIngredientRecord() SqlSafeMealIngredientRecord {
	return SqlSafeMealIngredientRecord{
		MealID:       mi.MealID,
		IngredientID: mi.IngredientID,
		Position:     mi.Position,
		MeasureText:  StringToSqlNullString(mi.MeasureText),
	}
}

type SqlSafeMealIngredientRecord struct {
	MealID       int64          `db:"MealId"`
	IngredientID int64          `db:"IngredientId"`
	Position     int64          `db:"Position"`
	MeasureText  sql.NullString `db:"MeasureText"`
}

func (ssmi SqlSafeMealIngredientRecord) ToMealIngredientRecord() MealIngredientRecord {
	return MealIngredientRecord{
		MealID:       ssmi.MealID,
		IngredientID: ssmi.IngredientID,
		Position:     ssmi.Position,
		MeasureText:  ssmi.MeasureText.String,
	}
}

func StringToSqlNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{
			String: "",
			Valid:  false,
		}
	}

	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

type MealMatch struct {
	Meal                   Meal     `json:"meal"`
	MissingIngredients     []string `json:"missingIngredients"`
	MatchedIngredientCount int64    `json:"matchedIngredientCount"`
	TotalIngredientCount   int64    `json:"totalIngredientCount"`
	MatchRatio             float64  `json:"matchRatio"`
}

func normalizeIngredientName(name string) string {
	normalized := strings.Join(strings.Fields(name), " ")
	return strings.ToLower(normalized)
}

func normalizeIngredientNames(ingredients []string) []string {
	seen := make(map[string]struct{}, len(ingredients))
	normalized := make([]string, 0, len(ingredients))

	for _, ingredient := range ingredients {
		name := normalizeIngredientName(ingredient)
		if name == "" {
			continue
		}

		if _, exists := seen[name]; exists {
			continue
		}

		seen[name] = struct{}{}
		normalized = append(normalized, name)
	}

	return normalized
}

func ValidateIngredientSearch(v *validator.Validator, ingredients []string) {
	v.Check(len(ingredients) > 0, "ingredients", "must contain at least one ingredient")
	v.Check(len(ingredients) <= 20, "ingredients", "must contain no more than 20 ingredients")

	for _, ingredient := range ingredients {
		ingredient = strings.TrimSpace(ingredient)
		v.Check(ingredient != "", "ingredients", "must not contain empty values")
		v.Check(len(ingredient) <= 100, "ingredients", "must not contain values longer than 100 characters")
	}
}

func (m *MealModel) FindByIngredients(ctx context.Context, ingredients []string, filters Filters) ([]MealMatch, Metadata, error) {
	normalizedIngredients := normalizeIngredientNames(ingredients)
	if len(normalizedIngredients) == 0 {
		return []MealMatch{}, calculateMetadata(0, filters.Page, filters.PageSize), nil
	}

	valueRows := make([]string, len(normalizedIngredients))
	args := make([]any, 0, len(normalizedIngredients)+2)

	for i, ingredient := range normalizedIngredients {
		valueRows[i] = "(?)"
		args = append(args, ingredient)
	}

	args = append(args, filters.limit(), filters.offset())

	query := fmt.Sprintf(`
		WITH InputIngredients (NormalizedName) AS (
			VALUES
				%s
		),
		InputIngredientIds AS (
			SELECT DISTINCT
				i.IngredientId
			FROM %s AS i
			INNER JOIN InputIngredients AS input
				ON input.NormalizedName = i.NormalizedName
		),
		MealMatches AS (
			SELECT
				m.MealId,
				m.Name,
				m.AlternateName,
				m.Category,
				m.Area,
				m.Country,
				m.Instructions,
				m.ThumbnailUrl,
				m.YoutubeUrl,
				m.SourceUrl,

				COUNT(
					DISTINCT CASE
						WHEN input.IngredientId IS NOT NULL
						THEN mi.IngredientId
					END
				) AS MatchedIngredientCount,

				COUNT(DISTINCT mi.IngredientId)
					AS TotalIngredientCount

			FROM %s AS m

			INNER JOIN %s AS mi
				ON mi.MealId = m.MealId

			LEFT JOIN InputIngredientIds AS input
				ON input.IngredientId = mi.IngredientId

			GROUP BY
				m.MealId,
				m.Name,
				m.AlternateName,
				m.Category,
				m.Area,
				m.Country,
				m.Instructions,
				m.ThumbnailUrl,
				m.YoutubeUrl,
				m.SourceUrl
		)
		SELECT
			MealId,
			Name,
			AlternateName,
			Category,
			Area,
			Country,
			Instructions,
			ThumbnailUrl,
			YoutubeUrl,
			SourceUrl,
			MatchedIngredientCount,
			TotalIngredientCount,
			CAST(MatchedIngredientCount AS REAL)
				/ NULLIF(TotalIngredientCount, 0)
				AS MatchRatio,
			COUNT(*) OVER() AS TotalRecords
		FROM MealMatches
		WHERE MatchedIngredientCount > 0
		ORDER BY %s
		LIMIT ? OFFSET ?;
		`,
		strings.Join(valueRows, ",\n\t\t\t\t"),
		IngredientTable,
		MealsTable,
		MealIngredientTable,
		filters.orderBy(),
	)

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, fmt.Errorf("find meals by ingredients: %w", err)
	}
	defer rows.Close()

	matches := make([]MealMatch, 0, filters.PageSize)
	totalRecords := 0

	for rows.Next() {
		var sqlMeal SqlSafeMeal
		var match MealMatch

		err := rows.Scan(
			&sqlMeal.ID,
			&sqlMeal.Name,
			&sqlMeal.AlternateName,
			&sqlMeal.Category,
			&sqlMeal.Area,
			&sqlMeal.Country,
			&sqlMeal.Instructions,
			&sqlMeal.ThumbnailURL,
			&sqlMeal.YoutubeURL,
			&sqlMeal.SourceURL,
			&match.MatchedIngredientCount,
			&match.TotalIngredientCount,
			&match.MatchRatio,
			&totalRecords,
		)
		if err != nil {
			return nil, Metadata{}, fmt.Errorf("scan meal match: %w", err)
		}

		match.Meal = sqlMeal.ToMeal()
		match.MissingIngredients = make([]string, 0)

		matches = append(matches, match)
	}

	if err := rows.Err(); err != nil {
		return nil, Metadata{}, fmt.Errorf("iterate meal matches: %w", err)
	}

	if err := m.loadIngredientsForMealMatches(ctx, matches); err != nil {
		return nil, Metadata{}, err
	}

	populateMissingIngredients(matches, normalizedIngredients)

	return matches, calculateMetadata(totalRecords, filters.Page, filters.PageSize), nil
}

func (m *MealModel) loadIngredientsForMealMatches(ctx context.Context, matches []MealMatch) error {
	if len(matches) == 0 {
		return nil
	}

	placeholders := make([]string, len(matches))
	args := make([]any, len(matches))
	matchIndexByMealID := make(map[int64]int, len(matches))

	for i := range matches {
		placeholders[i] = "?"
		args[i] = matches[i].Meal.ID
		matchIndexByMealID[matches[i].Meal.ID] = i
	}

	query := fmt.Sprintf(`
		SELECT
			mi.MealId,
			i.IngredientId,
			i.Name,
			i.NormalizedName,
			mi.Position,
			mi.MeasureText
		FROM %s AS mi
		INNER JOIN %s AS i
			ON i.IngredientId = mi.IngredientId
		WHERE mi.MealId IN (%s)
		ORDER BY
			mi.MealId,
			mi.Position;
		`,
		MealIngredientTable,
		IngredientTable,
		strings.Join(placeholders, ", "),
	)

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("load ingredients for meal matches: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var mealID int64
		var ingredient MealIngredient
		var name sql.NullString
		var normalizedName sql.NullString
		var measureText sql.NullString

		err := rows.Scan(
			&mealID,
			&ingredient.IngredientID,
			&name,
			&normalizedName,
			&ingredient.Position,
			&measureText,
		)
		if err != nil {
			return fmt.Errorf("scan meal ingredient: %w", err)
		}

		ingredient.Name = name.String
		ingredient.NormalizedName = normalizedName.String
		ingredient.MeasureText = measureText.String

		matchIndex, exists := matchIndexByMealID[mealID]
		if !exists {
			return fmt.Errorf("ingredient references unexpected meal %d", mealID)
		}

		matches[matchIndex].Meal.Ingredients = append(matches[matchIndex].Meal.Ingredients, ingredient)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate meal ingredients: %w", err)
	}

	return nil
}

func populateMissingIngredients(matches []MealMatch, normalizedInputIngredients []string) {
	inputSet := make(map[string]struct{}, len(normalizedInputIngredients))

	for _, ingredient := range normalizedInputIngredients {
		inputSet[ingredient] = struct{}{}
	}

	for i := range matches {
		missingCount := matches[i].TotalIngredientCount - matches[i].MatchedIngredientCount
		missingCount = max(missingCount, 0)

		missingIngredients := make([]string, 0, int(missingCount))
		seenMissingIngredientIDs := make(map[int64]struct{}, int(missingCount))

		for _, ingredient := range matches[i].Meal.Ingredients {
			if _, matched := inputSet[ingredient.NormalizedName]; matched {
				continue
			}

			if _, alreadyAdded := seenMissingIngredientIDs[ingredient.IngredientID]; alreadyAdded {
				continue
			}

			seenMissingIngredientIDs[ingredient.IngredientID] = struct{}{}

			missingIngredients = append(missingIngredients, ingredient.Name)
		}

		matches[i].MissingIngredients = missingIngredients
	}
}
