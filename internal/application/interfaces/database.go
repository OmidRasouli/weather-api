package interfaces

import (
	"context"

	"github.com/OmidRasouli/weather-api/internal/domain/weather"
)

type WeatherRepository interface {
	Save(ctx context.Context, w *weather.Weather) error
	FindByID(ctx context.Context, id string) (*weather.Weather, error)
	FindAll(ctx context.Context) ([]*weather.Weather, error)
	FindLatestByCity(ctx context.Context, city string) (*weather.Weather, error)
	Update(ctx context.Context, w *weather.Weather) error
	Delete(ctx context.Context, id string) error
}
