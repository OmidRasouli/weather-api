package services_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/OmidRasouli/weather-api/internal/application/interfaces"
	"github.com/OmidRasouli/weather-api/internal/application/services"
	"github.com/OmidRasouli/weather-api/internal/application/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFetchAndStoreWeather_RedisCacheError(t *testing.T) {
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
	mockRedis.On("Get", ctx, cacheKey, mock.Anything).Return(fmt.Errorf("cache error"))
	mockAPI.On("FetchWeatherData", ctx, "tehran", "IR").Return(apiResp, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*weather.Weather")).Return(nil)
	mockRedis.On("Set", ctx, cacheKey, mock.Anything).Return(fmt.Errorf("cache write error"))

	// Execute - service should still work even if Redis fails
	result, err := service.FetchAndStoreWeather(ctx, "tehran", "IR")

	// Verify
	assert.NoError(t, err) // Service should not fail if only Redis fails
	assert.NotNil(t, result)

	mockAPI.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockRedis.AssertExpectations(t)
}
