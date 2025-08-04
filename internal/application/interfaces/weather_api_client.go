package interfaces

import (
	"context"
	"time"
)

type WeatherAPIResponse struct {
	Temperature float64
	Description string
	Humidity    int
	WindSpeed   float64
	FetchedAt   time.Time
}

type WeatherAPIClient interface {
	FetchWeatherData(ctx context.Context, city string, country string) (*WeatherAPIResponse, error)
}
