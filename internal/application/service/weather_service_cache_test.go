package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/OmidRasouli/weather-api/internal/application/interfaces"
	"github.com/OmidRasouli/weather-api/internal/application/service"
	"github.com/OmidRasouli/weather-api/internal/application/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFetchAndStoreWeather_CacheError(t *testing.T) {
	mockRepo := new(mocks.MockWeatherRepository)
	mockAPI := new(mocks.MockAPIClient)
	mockCache := new(mocks.MockCache)

	svc := service.NewWeatherService(mockRepo, mockAPI, mockCache)

	ctx := context.TODO()
	apiResp := &interfaces.WeatherAPIResponse{
		Temperature: 30.5,
		Description: "sunny",
		Humidity:    40,
		WindSpeed:   5.5,
		FetchedAt:   time.Now(),
	}

	cacheKey := mocks.CreateCacheKey("tehran", "IR")
	mockCache.On("Get", ctx, cacheKey, mock.Anything).Return(fmt.Errorf("cache error"))
	mockAPI.On("FetchWeatherData", ctx, "tehran", "IR").Return(apiResp, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*weather.Weather")).Return(nil)
	mockCache.On("Set", ctx, cacheKey, mock.Anything).Return(fmt.Errorf("cache write error"))

	result, err := svc.FetchAndStoreWeather(ctx, "tehran", "IR")

	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockAPI.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}
