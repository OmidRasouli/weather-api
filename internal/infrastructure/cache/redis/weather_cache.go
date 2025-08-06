package redis

import (
	"context"
	"fmt"

	"github.com/OmidRasouli/weather-api/infrastructure/database/redis"
	"github.com/OmidRasouli/weather-api/internal/domain/weather"
	"github.com/OmidRasouli/weather-api/pkg/logger"
)

const (
	weatherKeyPrefix = "weather:"
	cityPrefix       = "city:"
	latestPrefix     = "latest:"
)

// WeatherCache provides Redis caching for weather data
type WeatherCache struct {
	redis *redis.Redis
}

// NewWeatherCache creates a new weather cache
func NewWeatherCache(redis *redis.Redis) *WeatherCache {
	return &WeatherCache{redis: redis}
}

// getCityKey generates a Redis key for a city's weather
func getCityKey(city, country string) string {
	return fmt.Sprintf("%s%s:%s", weatherKeyPrefix, city, country)
}

// getLatestCityKey generates a Redis key for the latest weather of a city
func getLatestCityKey(city string) string {
	return fmt.Sprintf("%s%s%s", weatherKeyPrefix, latestPrefix, city)
}

// CacheWeather stores weather data in Redis
func (wc *WeatherCache) CacheWeather(ctx context.Context, w *weather.Weather) error {
	// Cache by city:country
	cityKey := getCityKey(w.City, w.Country)
	if err := wc.redis.Set(ctx, cityKey, w); err != nil {
		return fmt.Errorf("failed to cache weather by city/country: %w", err)
	}

	// Cache as latest for this city
	latestKey := getLatestCityKey(w.City)
	if err := wc.redis.Set(ctx, latestKey, w); err != nil {
		logger.Warnf("Failed to cache latest weather for city %s: %v", w.City, err)
		// Continue even if this fails
	}

	logger.Infof("Cached weather data for %s, %s", w.City, w.Country)
	return nil
}

// GetWeatherByCityCountry retrieves weather data for a specific city and country
func (wc *WeatherCache) GetWeatherByCityCountry(ctx context.Context, city, country string) (*weather.Weather, error) {
	key := getCityKey(city, country)

	var weatherData weather.Weather
	err := wc.redis.Get(ctx, key, &weatherData)
	if err != nil {
		return nil, fmt.Errorf("cache miss for %s: %w", key, err)
	}

	logger.Infof("Cache hit for weather data: %s, %s", city, country)
	return &weatherData, nil
}

// GetLatestWeatherByCity retrieves the latest weather data for a city
func (wc *WeatherCache) GetLatestWeatherByCity(ctx context.Context, city string) (*weather.Weather, error) {
	key := getLatestCityKey(city)

	var weatherData weather.Weather
	err := wc.redis.Get(ctx, key, &weatherData)
	if err != nil {
		return nil, fmt.Errorf("cache miss for latest weather in %s: %w", city, err)
	}

	return &weatherData, nil
}

// InvalidateWeather removes weather data from cache
func (wc *WeatherCache) InvalidateWeather(ctx context.Context, city, country string) error {
	key := getCityKey(city, country)
	latestKey := getLatestCityKey(city)

	if err := wc.redis.Delete(ctx, key, latestKey); err != nil {
		return fmt.Errorf("failed to invalidate cache: %w", err)
	}

	logger.Infof("Invalidated cache for %s, %s", city, country)
	return nil
}

// GetCachedCities returns a list of all cached cities
func (wc *WeatherCache) GetCachedCities(ctx context.Context) ([]string, error) {
	pattern := fmt.Sprintf("%s%s*", weatherKeyPrefix, cityPrefix)
	keys, err := wc.redis.GetKeys(ctx, pattern)
	if err != nil {
		return nil, err
	}

	cities := make([]string, 0, len(keys))
	for _, key := range keys {
		// Extract city name from key
		city := key[len(weatherKeyPrefix)+len(cityPrefix):]
		cities = append(cities, city)
	}

	return cities, nil
}
