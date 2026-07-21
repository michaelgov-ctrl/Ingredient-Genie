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
	ID            int64  `json:"id"`
	MealDBID      int64  `json:"mealDbId"`
	Name          string `json:"name"`
	AlternateName string `json:"alternateName"`
	Category      string `json:"category"`
	Area          string `json:"area"`
	Country       string `json:"country"`
	Instructions  string `json:"instructions"`
	ThumbnailUrl  string `json:"thumbnailUrl"`
	YoutubeUrl    string `json:"youtubeUrl"`
	SourceUrl     string `json:"sourceUrl"`
}

func (m Meal) ToSqlSafeMeal() SqlSafeMeal {
	var ssm SqlSafeMeal

	ssm.ID = m.ID
	ssm.MealDBID = m.MealDBID
	ssm.Name = StringToSqlNullString(m.Name)
	ssm.AlternateName = StringToSqlNullString(m.AlternateName)
	ssm.Category = StringToSqlNullString(m.Category)
	ssm.Area = StringToSqlNullString(m.Area)
	ssm.Country = StringToSqlNullString(m.Country)
	ssm.Instructions = StringToSqlNullString(m.Instructions)
	ssm.ThumbnailUrl = StringToSqlNullString(m.ThumbnailUrl)
	ssm.YoutubeUrl = StringToSqlNullString(m.YoutubeUrl)
	ssm.SourceUrl = StringToSqlNullString(m.SourceUrl)

	return ssm
}

type SqlSafeMeal struct {
	ID            int64          `db:"MealId"`
	MealDBID      int64          `db:"ExternalMealId"`
	Name          sql.NullString `db:"Name"`
	AlternateName sql.NullString `db:"AlternateName"`
	Category      sql.NullString `db:"Category"`
	Area          sql.NullString `db:"Area"`
	Country       sql.NullString `db:"Country"`
	Instructions  sql.NullString `db:"Instructions"`
	ThumbnailUrl  sql.NullString `db:"ThumbnailUrl"`
	YoutubeUrl    sql.NullString `db:"YoutubeUrl"`
	SourceUrl     sql.NullString `db:"SourceUrl"`
}

func (ssm SqlSafeMeal) ToMeal() Meal {
	return Meal{
		ID:            ssm.ID,
		MealDBID:      ssm.MealDBID,
		Name:          ssm.Name.String,
		AlternateName: ssm.AlternateName.String,
		Category:      ssm.Category.String,
		Area:          ssm.Area.String,
		Country:       ssm.Country.String,
		Instructions:  ssm.Instructions.String,
		ThumbnailUrl:  ssm.ThumbnailUrl.String,
		YoutubeUrl:    ssm.YoutubeUrl.String,
		SourceUrl:     ssm.SourceUrl.String,
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
	ID           int64  `json:"id"`
	IngredientID int64  `json:"ingredientId"`
	Position     int64  `json:"position"`
	MeasureText  string `json:"MeasureText"`
}

func (mi MealIngredient) ToSqlSafeMealIngredient() SqlSafeMealIngredient {
	var ssmi SqlSafeMealIngredient

	ssmi.ID = mi.ID
	ssmi.IngredientID = mi.IngredientID
	ssmi.Position = mi.Position
	ssmi.MeasureText = StringToSqlNullString(mi.MeasureText)

	return ssmi
}

type SqlSafeMealIngredient struct {
	ID           int64          `db:"MealId"`
	IngredientID int64          `db:"IngredientId"`
	Position     int64          `db:"Position"`
	MeasureText  sql.NullString `db:"MeasureText"`
}

func (ssmi SqlSafeMealIngredient) ToMealIngredient() MealIngredient {
	return MealIngredient{
		ID:           ssmi.ID,
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
	Meal                   Meal    `json:"meal"`
	MatchedIngredientCount int64   `json:"matchedIngredientCount"`
	TotalIngredientCount   int64   `json:"totalIngredientCount"`
	MissingIngredientCount int64   `json:"missingIngredientCount"`
	MatchRatio             float64 `json:"matchRatio"`
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
				CAST(m.ExternalMealId AS INTEGER) AS ExternalMealId,
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
				m.ExternalMealId,
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
			ExternalMealId,
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
			TotalIngredientCount - MatchedIngredientCount
				AS MissingIngredientCount,
			CAST(MatchedIngredientCount AS REAL)
				/ NULLIF(TotalIngredientCount, 0)
				AS MatchRatio,
			COUNT(*) OVER() AS TotalRecords
		FROM MealMatches
		WHERE MatchedIngredientCount > 0
		ORDER BY %s
		LIMIT ? OFFSET ?;
		`,
		strings.Join(valueRows, ",\n\t\t"),
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
			&sqlMeal.MealDBID,
			&sqlMeal.Name,
			&sqlMeal.AlternateName,
			&sqlMeal.Category,
			&sqlMeal.Area,
			&sqlMeal.Country,
			&sqlMeal.Instructions,
			&sqlMeal.ThumbnailUrl,
			&sqlMeal.YoutubeUrl,
			&sqlMeal.SourceUrl,
			&match.MatchedIngredientCount,
			&match.TotalIngredientCount,
			&match.MissingIngredientCount,
			&match.MatchRatio,
			&totalRecords,
		)
		if err != nil {
			return nil, Metadata{}, fmt.Errorf("scan meal match: %w", err)
		}

		match.Meal = sqlMeal.ToMeal()
		matches = append(matches, match)
	}

	if err := rows.Err(); err != nil {
		return nil, Metadata{}, fmt.Errorf("iterate meal matches: %w", err)
	}

	return matches, calculateMetadata(totalRecords, filters.Page, filters.PageSize), nil
}
