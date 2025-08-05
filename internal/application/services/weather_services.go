package services

import (
	"context"
	"time"

	"github.com/OmidRasouli/weather-api/internal/application/interfaces"
	"github.com/OmidRasouli/weather-api/internal/domain/weather"
	"github.com/google/uuid"
)

type WeatherService struct {
	repo       interfaces.WeatherRepository
	apiClient  interfaces.WeatherAPIClient
	timeSource func() time.Time // testable clock
}

func NewWeatherService(repo interfaces.WeatherRepository, api interfaces.WeatherAPIClient) *WeatherService {
	return &WeatherService{
		repo:       repo,
		apiClient:  api,
		timeSource: time.Now,
	}
}

func (s *WeatherService) FetchAndStoreWeather(ctx context.Context, city string, country string) (*weather.Weather, error) {
	apiData, err := s.apiClient.FetchWeatherData(ctx, city, country)
	if err != nil {
		return nil, err
	}

	entity := &weather.Weather{
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

	if err := s.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *WeatherService) GetLatestWeatherByCity(ctx context.Context, city string) (*weather.Weather, error) {
	return s.repo.FindLatestByCity(ctx, city)
}

func (s *WeatherService) GetAllWeather(ctx context.Context) ([]*weather.Weather, error) {
	return s.repo.FindAll(ctx)
}

func (s *WeatherService) GetWeatherByID(ctx context.Context, id string) (*weather.Weather, error) {
	return s.repo.FindByID(ctx, id)
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
