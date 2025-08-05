package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/OmidRasouli/weather-api/internal/application/interfaces"
	"github.com/OmidRasouli/weather-api/internal/domain/weather"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWeatherRepository struct {
	mock.Mock
}

type MockAPIClient struct {
	mock.Mock
}

func (m *MockAPIClient) FetchWeatherData(ctx context.Context, city string, country string) (*interfaces.WeatherAPIResponse, error) {
	args := m.Called(ctx, city, country)
	return args.Get(0).(*interfaces.WeatherAPIResponse), args.Error(1)
}

func (m *MockWeatherRepository) Save(ctx context.Context, w *weather.Weather) error {
	args := m.Called(ctx, w)
	return args.Error(0)
}

func (m *MockWeatherRepository) FindByID(ctx context.Context, id uuid.UUID) (*weather.Weather, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*weather.Weather), args.Error(1)
}

func (m *MockWeatherRepository) FindAll(ctx context.Context) ([]*weather.Weather, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*weather.Weather), args.Error(1)
}

func (m *MockWeatherRepository) FindLatestByCity(ctx context.Context, city string) (*weather.Weather, error) {
	args := m.Called(ctx, city)
	return args.Get(0).(*weather.Weather), args.Error(1)
}

func (m *MockWeatherRepository) Update(ctx context.Context, w *weather.Weather) error {
	args := m.Called(ctx, w)
	return args.Error(0)
}

func (m *MockWeatherRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestFetchAndStoreWeather_Success(t *testing.T) {
	mockRepo := new(MockWeatherRepository)
	mockAPI := new(MockAPIClient)
	service := NewWeatherService(mockRepo, mockAPI)

	ctx := context.TODO()
	apiResp := &interfaces.WeatherAPIResponse{
		Temperature: 30.5,
		Description: "sunny",
		Humidity:    40,
		WindSpeed:   5.5,
		FetchedAt:   time.Now(),
	}

	mockAPI.On("FetchWeatherData", ctx, "tehran", "IR").Return(apiResp, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*weather.Weather")).Return(nil)

	result, err := service.FetchAndStoreWeather(ctx, "tehran", "IR")

	assert.NoError(t, err)
	assert.Equal(t, "tehran", result.CityName)
	assert.Equal(t, "IR", result.Country)

	mockAPI.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// TestFetchAndStoreWeather_APIError tests the behavior when the weather API returns an error
func TestFetchAndStoreWeather_APIError(t *testing.T) {
	mockRepo := new(MockWeatherRepository)
	mockAPI := new(MockAPIClient)
	service := NewWeatherService(mockRepo, mockAPI)

	ctx := context.TODO()
	expectedErr := fmt.Errorf("API unavailable")

	mockAPI.On("FetchWeatherData", ctx, "tehran", "IR").Return((*interfaces.WeatherAPIResponse)(nil), expectedErr)
	// Repository should not be called when API fails

	result, err := service.FetchAndStoreWeather(ctx, "tehran", "IR")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "API unavailable")

	mockAPI.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Save") // Verify repository
}
