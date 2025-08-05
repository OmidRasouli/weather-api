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
