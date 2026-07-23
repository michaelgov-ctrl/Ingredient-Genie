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
	logger                      *slog.Logger
	addr                        string
	mealSearchEndpoint          string
	mealSearchSortTypesEndpoint string
	httpClient                  *http.Client
}

func NewMealsClient(logger *slog.Logger, addr string) MealsClient {
	client := MealsClient{
		logger:                      logger,
		addr:                        addr,
		mealSearchEndpoint:          "/v1/meals/search",
		mealSearchSortTypesEndpoint: "/v1/meals/search/sort",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	go func() {
		// TODO: on timer verify health of API use logger for result
	}()

	return client
}

func (mc MealsClient) GetMealList(s string) (MealListResponse, error) {
	// TODO: set this up in the backend
	return MealListResponse{}, nil
}

func (mc MealsClient) GetSearchSortTypes() ([]SortType, error) {
	// TODO: get filter typess from the backend, gotta add that to the backend...
	req, err := http.NewRequest(http.MethodGet, mc.addr+mc.mealSearchSortTypesEndpoint, nil)
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

	req, err := http.NewRequest(http.MethodPost, mc.addr+mc.mealSearchEndpoint, bytes.NewReader(payload))
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
