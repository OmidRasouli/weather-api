package service

import (
	"context"
	"fmt"
	"time"

	"github.com/OmidRasouli/weather-api/internal/application/interfaces"
	"github.com/OmidRasouli/weather-api/internal/domain/weather"
	"github.com/OmidRasouli/weather-api/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WeatherService struct {
	repo       interfaces.WeatherRepository
	apiClient  interfaces.WeatherAPIClient
	cache      interfaces.Cache
	timeSource func() time.Time // testable clock
}

func (s *WeatherService) GetWeather(ctx *gin.Context, param any) (any, any) {
	panic("unimplemented")
}

func NewWeatherService(repo interfaces.WeatherRepository, api interfaces.WeatherAPIClient, cache interfaces.Cache) *WeatherService {
	return &WeatherService{
		repo:       repo,
		apiClient:  api,
		cache:      cache,
		timeSource: time.Now,
	}
}

// FetchAndStoreWeather fetches weather data from the API or cache and stores it
func (s *WeatherService) FetchAndStoreWeather(ctx context.Context, city string, country string) (*weather.Weather, error) {
	// Create a cache key based on city and country
	cacheKey := fmt.Sprintf("weather:%s:%s", city, country)

	// Try to get from cache first
	var weatherData *weather.Weather
	err := s.cache.Get(ctx, cacheKey, &weatherData)
	if err == nil {
		// Cache hit!
		logger.Infof("Retrieved weather data from cache for %s, %s", city, country)
		return weatherData, nil
	}

	// Cache miss, fetch from API
	logger.Infof("Cache miss for %s, %s. Fetching from API", city, country)
	apiData, err := s.apiClient.FetchWeatherData(ctx, city, country)
	if err != nil {
		return nil, err
	}

	weatherData = &weather.Weather{
		ID:          uuid.New(),
		City:        city,
		Country:     country,
		Temperature: apiData.Temperature,
		Description: apiData.Description,
		Humidity:    apiData.Humidity,
		WindSpeed:   apiData.WindSpeed,
		FetchedAt:   apiData.FetchedAt,
		CreatedAt:   s.timeSource(),
		UpdatedAt:   s.timeSource(),
	}

	// Store in the database
	if err := s.repo.Save(ctx, weatherData); err != nil {
		return nil, err
	}

	// Store in cache for future requests
	if err := s.cache.Set(ctx, cacheKey, weatherData); err != nil {
		// Log the error but don't fail the request
		logger.Errorf("Failed to cache weather data: %v", err)
	}

	return weatherData, nil
}

func (s *WeatherService) GetLatestWeatherByCity(ctx context.Context, city string) (*weather.Weather, error) {
	return s.repo.FindLatestByCity(ctx, city)
}

func (s *WeatherService) GetAllWeather(ctx context.Context) ([]*weather.Weather, error) {
	return s.repo.FindAll(ctx)
}

func (s *WeatherService) GetWeatherByID(ctx context.Context, id string) (*weather.Weather, error) {
	var cachedWeather *weather.Weather
	err := s.cache.Get(ctx, id, &cachedWeather)
	if err == nil && cachedWeather != nil {
		return cachedWeather, nil
	}

	w, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if w != nil {
		_ = s.cache.Set(ctx, id, w)
	}

	return w, nil
}

func (s *WeatherService) UpdateWeather(ctx context.Context, id string, update *weather.Weather) (*weather.Weather, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if update.City != "" {
		existing.City = update.City
	}
	if update.Country != "" {
		existing.Country = update.Country
	}
	if update.Temperature != 0 {
		existing.Temperature = update.Temperature
	}
	if update.Description != "" {
		existing.Description = update.Description
	}
	if update.Humidity != 0 {
		existing.Humidity = update.Humidity
	}
	if update.WindSpeed != 0 {
		existing.WindSpeed = update.WindSpeed
	}
	existing.UpdatedAt = s.timeSource()

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *WeatherService) DeleteWeather(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
