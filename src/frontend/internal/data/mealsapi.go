package data

import (
	"log/slog"
	"net"
	"net/http"
	"time"
)

type MealsApiClient struct {
	logger     *slog.Logger
	addr       net.Addr
	httpClient *http.Client
}

func NewMealsApiClient(logger *slog.Logger, addr net.Addr) MealsApiClient {
	client := MealsApiClient{
		logger: logger,
		addr:   addr,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	go func() {
		// TODO: on timer verify health of API use logger for result
	}()

	return client
}

func (mac MealsApiClient) SearchByIngredients() (MealResponse, Metadata, error) {
	// TODO: make request to api
	return MealResponse{}, Metadata{}, nil
}
