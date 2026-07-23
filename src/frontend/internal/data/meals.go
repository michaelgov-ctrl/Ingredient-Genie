package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// TODO: some kind of validation on addr
type MealsClient struct {
	logger                 *slog.Logger
	addr                   string
	healthcheckEndpoint    string
	mealsGetEndpoint       string
	mealsListEndpoint      string
	mealsSearchEndpoint    string
	mealsSortTypesEndpoint string
	httpClient             *http.Client
}

func NewMealsClient(logger *slog.Logger, addr string) MealsClient {
	version := "/v1"

	client := MealsClient{
		logger:                 logger,
		addr:                   addr,
		healthcheckEndpoint:    version + "/healthcheck",
		mealsGetEndpoint:       version + "/meals/get",
		mealsListEndpoint:      version + "/meals/list",
		mealsSearchEndpoint:    version + "/meals/search",
		mealsSortTypesEndpoint: version + "/meals/sort",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	go func() {
		// TODO: on timer verify health of API use logger for result
	}()

	return client
}

func (mc MealsClient) GetMeal(id int) (Meal, error) {
	input := struct {
		ID int `json:"id"`
	}{
		ID: id,
	}

	body, err := json.Marshal(input)
	if err != nil {
		return Meal{}, err
	}

	req, err := http.NewRequest(http.MethodPost, mc.addr+mc.mealsGetEndpoint, bytes.NewReader(body))
	if err != nil {
		return Meal{}, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := mc.httpClient.Do(req)
	if err != nil {
		return Meal{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return Meal{}, ErrNoMeal
	}

	if resp.StatusCode != http.StatusOK {
		return Meal{}, fmt.Errorf("failed to get meal with http status code: %d", resp.StatusCode)
	}

	var response struct {
		Meal Meal `json:"meal"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return Meal{}, err
	}

	return response.Meal, nil
}

func (mc MealsClient) GetMealList(filters Filters) (MealListResponse, error) {
	input := struct {
		Filters Filters `json:"filters"`
	}{
		Filters: filters,
	}

	body, err := json.Marshal(input)
	if err != nil {
		return MealListResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, mc.addr+mc.mealsListEndpoint, bytes.NewReader(body))
	if err != nil {
		return MealListResponse{}, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := mc.httpClient.Do(req)
	if err != nil {
		return MealListResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return MealListResponse{}, fmt.Errorf("failed to get meal list with http status code: %d", resp.StatusCode)
	}

	var response MealListResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return MealListResponse{}, err
	}

	return response, nil
}

func (mc MealsClient) GetSortTypes() ([]SortType, error) {
	req, err := http.NewRequest(http.MethodGet, mc.addr+mc.mealsSortTypesEndpoint, nil)
	if err != nil {
		return []SortType{}, err
	}

	req.Header.Add("Accept", "application/json")

	resp, err := mc.httpClient.Do(req)
	if err != nil {
		return []SortType{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []SortType{}, fmt.Errorf("failed to enumerate sort types with http status code: %d", resp.StatusCode)
	}

	var Response struct {
		SortTypes []SortType `json:"sortTypes"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&Response); err != nil {
		return []SortType{}, err
	}

	return Response.SortTypes, nil
}

func (mc MealsClient) SearchByIngredients(body IngredientMealSearchRequest) (MealSearchResponse, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return MealSearchResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, mc.addr+mc.mealsSearchEndpoint, bytes.NewReader(payload))
	if err != nil {
		return MealSearchResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := mc.httpClient.Do(req)
	if err != nil {
		return MealSearchResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return MealSearchResponse{}, fmt.Errorf("failed to search ingredients with http status code: %d", resp.StatusCode)
	}

	var msr MealSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&msr); err != nil {
		return MealSearchResponse{}, err

	}

	return msr, nil
}
