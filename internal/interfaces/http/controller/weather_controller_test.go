package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/OmidRasouli/weather-api/internal/domain/weather"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/mock"
)

type MockWeatherService struct {
	mock.Mock
}

func (m *MockWeatherService) FetchAndStoreWeather(ctx context.Context, city, country string) (*weather.Weather, error) {
	args := m.Called(ctx, city, country)
	return args.Get(0).(*weather.Weather), args.Error(1)
}

func (m *MockWeatherService) GetLatestWeatherByCity(ctx context.Context, city string) (*weather.Weather, error) {
	args := m.Called(ctx, city)
	return args.Get(0).(*weather.Weather), args.Error(1)
}

func (m *MockWeatherService) GetAllWeather(ctx context.Context) ([]*weather.Weather, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*weather.Weather), args.Error(1)
}

func (m *MockWeatherService) GetWeatherByID(ctx context.Context, id string) (*weather.Weather, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*weather.Weather), args.Error(1)
}

func (m *MockWeatherService) UpdateWeather(ctx context.Context, id string, update *weather.Weather) (*weather.Weather, error) {
	args := m.Called(ctx, id, update)
	return args.Get(0).(*weather.Weather), args.Error(1)
}

func (m *MockWeatherService) DeleteWeather(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestFetchAndStore_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockWeatherService)

	sut := NewWeatherController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	reqBody := `{"city":"tehran", "country":"IR"}`
	c.Request = httptest.NewRequest("POST", "/weather", strings.NewReader(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	expected := &weather.Weather{
		CityName:    "tehran",
		Country:     "IR",
		Temperature: 32.5,
		Humidity:    60,
		Description: "clear sky",
		WindSpeed:   2.1,
		FetchedAt:   time.Now(),
	}

	mockService.On("FetchAndStoreWeather", mock.Anything, "tehran", "IR").Return(expected, nil)

	sut.FetchAndStore(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockWeatherService)
	sut := NewWeatherController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "test-id"}}

	expected := &weather.Weather{
		CityName:    "tehran",
		Country:     "IR",
		Temperature: 32.5,
		Humidity:    60,
		Description: "clear sky",
		WindSpeed:   2.1,
		FetchedAt:   time.Now(),
	}

	mockService.On("GetWeatherByID", mock.Anything, "test-id").Return(expected, nil)

	sut.GetByID(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestUpdate_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockWeatherService)
	sut := NewWeatherController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "test-id"}}

	reqBody := `{"cityName":"tehran","country":"IR","temperature":33.0,"humidity":65,"description":"few clouds","windSpeed":2.5}`
	c.Request = httptest.NewRequest("PUT", "/weather/test-id", strings.NewReader(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	update := &weather.Weather{
		CityName:    "tehran",
		Country:     "IR",
		Temperature: 33.0,
		Humidity:    65,
		Description: "few clouds",
		WindSpeed:   2.5,
	}

	expected := *update

	mockService.On("UpdateWeather", mock.Anything, "test-id", mock.AnythingOfType("*weather.Weather")).Return(&expected, nil)

	sut.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestDelete_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockWeatherService)
	sut := NewWeatherController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "test-id"}}

	mockService.On("DeleteWeather", mock.Anything, "test-id").Return(nil)

	sut.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetLatestByCity_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockWeatherService)
	sut := NewWeatherController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "cityName", Value: "tehran"}}

	expected := &weather.Weather{
		CityName:    "tehran",
		Country:     "IR",
		Temperature: 32.5,
		Humidity:    60,
		Description: "clear sky",
		WindSpeed:   2.1,
		FetchedAt:   time.Now(),
	}

	mockService.On("GetLatestWeatherByCity", mock.Anything, "tehran").Return(expected, nil)

	sut.GetLatestByCity(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
