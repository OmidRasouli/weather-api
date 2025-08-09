package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/OmidRasouli/weather-api/internal/application/interfaces"
	"github.com/OmidRasouli/weather-api/internal/application/service"
	"github.com/OmidRasouli/weather-api/internal/application/service/mocks"
	"github.com/OmidRasouli/weather-api/internal/domain/weather"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFetchAndStoreWeather_CacheHit(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockWeatherRepository)
	mockAPI := new(mocks.MockAPIClient)
	mockCache := new(mocks.MockCache)
	service := service.NewWeatherService(mockRepo, mockAPI, mockCache)

	ctx := context.TODO()
	cacheKey := mocks.CreateCacheKey("tehran", "IR")

	// Create a cached weather entry
	cachedWeather := &weather.Weather{
		City:        "tehran",
		Country:     "IR",
		Temperature: 28.5,
		Description: "clear sky",
		Humidity:    35,
		WindSpeed:   4.5,
		FetchedAt:   time.Now(),
	}

	// Setup Cache to return the cached entry
	mockCache.On("Get", ctx, cacheKey, mock.Anything).Run(func(args mock.Arguments) {
		dest := args.Get(2).(**weather.Weather)
		*dest = cachedWeather
	}).Return(nil)

	// Execute
	result, err := service.FetchAndStoreWeather(ctx, "tehran", "IR")

	// Verify
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, cachedWeather.City, result.City)
	assert.Equal(t, cachedWeather.Temperature, result.Temperature)

	mockCache.AssertExpectations(t)
	mockAPI.AssertNotCalled(t, "FetchWeatherData")
	mockRepo.AssertNotCalled(t, "Save")
}

func TestFetchAndStoreWeather_Success(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockWeatherRepository)
	mockAPI := new(mocks.MockAPIClient)
	mockCache := new(mocks.MockCache)
	service := service.NewWeatherService(mockRepo, mockAPI, mockCache)

	ctx := context.TODO()
	apiResp := &interfaces.WeatherAPIResponse{
		Temperature: 30.5,
		Description: "sunny",
		Humidity:    40,
		WindSpeed:   5.5,
		FetchedAt:   time.Now(),
	}

	cacheKey := mocks.CreateCacheKey("tehran", "IR")
	mockCache.On("Get", ctx, cacheKey, mock.Anything).Return(fmt.Errorf("cache miss"))
	mockAPI.On("FetchWeatherData", ctx, "tehran", "IR").Return(apiResp, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*weather.Weather")).Return(nil)
	mockCache.On("Set", ctx, cacheKey, mock.Anything).Return(nil)

	// Execute
	result, err := service.FetchAndStoreWeather(ctx, "tehran", "IR")

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, "tehran", result.City)
	assert.Equal(t, "IR", result.Country)

	mockAPI.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}
