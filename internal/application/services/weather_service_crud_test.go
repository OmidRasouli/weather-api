package services_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/OmidRasouli/weather-api/internal/application/services"
	"github.com/OmidRasouli/weather-api/internal/application/services/mocks"
	"github.com/OmidRasouli/weather-api/internal/domain/weather"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetWeatherByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockWeatherRepository)
	mockAPI := new(mocks.MockAPIClient)
	mockRedis := new(mocks.MockRedisClient)
	service := services.NewWeatherService(mockRepo, mockAPI, mockRedis)

	ctx := context.TODO()
	id := uuid.New()
	expected := &weather.Weather{
		ID:          id,
		City:        "tehran",
		Country:     "IR",
		Temperature: 30.5,
		Description: "sunny",
		Humidity:    40,
		WindSpeed:   5.5,
		FetchedAt:   time.Now(),
	}

	mockRepo.On("FindByID", ctx, id).Return(expected, nil)

	result, err := service.GetWeatherByID(ctx, id.String())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expected.ID, result.ID)
	assert.Equal(t, expected.City, result.City)
	mockRepo.AssertExpectations(t)
}

func TestGetWeatherByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockWeatherRepository)
	mockAPI := new(mocks.MockAPIClient)
	mockRedis := new(mocks.MockRedisClient)
	service := services.NewWeatherService(mockRepo, mockAPI, mockRedis)

	ctx := context.TODO()
	id := uuid.New()
	mockRepo.On("FindByID", ctx, id).Return((*weather.Weather)(nil), fmt.Errorf("not found"))

	result, err := service.GetWeatherByID(ctx, id.String())

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
	mockRepo.AssertExpectations(t)
}
