package mocks

import (
	"context"
	"fmt"
	"time"

	"github.com/OmidRasouli/weather-api/internal/application/interfaces"
	"github.com/OmidRasouli/weather-api/internal/domain/weather"
	"github.com/stretchr/testify/mock"
)

// MockWeatherRepository mocks the weather repository
type MockWeatherRepository struct {
	mock.Mock
}

// MockAPIClient mocks the API client
type MockAPIClient struct {
	mock.Mock
}

// MockCache mocks the Redis client
type MockCache struct {
	mock.Mock
}

// MockWeatherRepository methods
func (m *MockWeatherRepository) Save(ctx context.Context, w *weather.Weather) error {
	args := m.Called(ctx, w)
	return args.Error(0)
}

func (m *MockWeatherRepository) FindByID(ctx context.Context, id string) (*weather.Weather, error) {
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

func (m *MockWeatherRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockAPIClient methods
func (m *MockAPIClient) FetchWeatherData(ctx context.Context, city string, country string) (*interfaces.WeatherAPIResponse, error) {
	args := m.Called(ctx, city, country)
	return args.Get(0).(*interfaces.WeatherAPIResponse), args.Error(1)
}

// MockRedisClient methods
func (m *MockCache) Get(ctx context.Context, key string, dest interface{}) error {
	args := m.Called(ctx, key, dest)
	return args.Error(0)
}

func (m *MockCache) Set(ctx context.Context, key string, value interface{}) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *MockCache) SetWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *MockCache) Delete(ctx context.Context, keys ...string) error {
	args := m.Called(ctx, keys)
	return args.Error(0)
}

func (m *MockCache) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *MockCache) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(time.Duration), args.Error(1)
}

func (m *MockCache) Flush(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCache) GetKeys(ctx context.Context, pattern string) ([]string, error) {
	args := m.Called(ctx, pattern)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	args := m.Called(ctx, key, ttl)
	return args.Error(0)
}

func (m *MockCache) Increment(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCache) HealthCheck(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCache) Close() error {
	args := m.Called()
	return args.Error(0)
}

// Helper function to create cache key
func CreateCacheKey(city, country string) string {
	return fmt.Sprintf("weather:%s:%s", city, country)
}
