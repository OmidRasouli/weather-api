package services_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/OmidRasouli/weather-api/internal/application/interfaces"
	"github.com/OmidRasouli/weather-api/internal/application/services"
	"github.com/OmidRasouli/weather-api/internal/application/services/mocks"
	"github.com/OmidRasouli/weather-api/internal/domain/weather"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFetchAndStoreWeather_CacheHit(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockWeatherRepository)
	mockAPI := new(mocks.MockAPIClient)
	mockRedis := new(mocks.MockRedisClient)
	service := services.NewWeatherService(mockRepo, mockAPI, mockRedis)

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

	// Setup Redis to return the cached entry
	mockRedis.On("Get", ctx, cacheKey, mock.Anything).Run(func(args mock.Arguments) {
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

	mockRedis.AssertExpectations(t)
	mockAPI.AssertNotCalled(t, "FetchWeatherData")
	mockRepo.AssertNotCalled(t, "Save")
}

func TestFetchAndStoreWeather_Success(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockWeatherRepository)
	mockAPI := new(mocks.MockAPIClient)
	mockRedis := new(mocks.MockRedisClient)
	service := services.NewWeatherService(mockRepo, mockAPI, mockRedis)

	ctx := context.TODO()
	apiResp := &interfaces.WeatherAPIResponse{
		Temperature: 30.5,
		Description: "sunny",
		Humidity:    40,
		WindSpeed:   5.5,
		FetchedAt:   time.Now(),
	}

	cacheKey := mocks.CreateCacheKey("tehran", "IR")
	mockRedis.On("Get", ctx, cacheKey, mock.Anything).Return(fmt.Errorf("cache miss"))
	mockAPI.On("FetchWeatherData", ctx, "tehran", "IR").Return(apiResp, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*weather.Weather")).Return(nil)
	mockRedis.On("Set", ctx, cacheKey, mock.Anything).Return(nil)

	// Execute
	result, err := service.FetchAndStoreWeather(ctx, "tehran", "IR")

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, "tehran", result.City)
	assert.Equal(t, "IR", result.Country)

	mockAPI.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockRedis.AssertExpectations(t)
}
