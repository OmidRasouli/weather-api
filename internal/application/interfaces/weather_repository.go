package interfaces

import (
	"context"

	"github.com/OmidRasouli/weather-api/internal/domain/weather"
	"github.com/google/uuid"
)

type WeatherRepository interface {
	Save(ctx context.Context, w *weather.Weather) error
	FindByID(ctx context.Context, id uuid.UUID) (*weather.Weather, error)
	FindAll(ctx context.Context) ([]*weather.Weather, error)
	FindLatestByCity(ctx context.Context, city string) (*weather.Weather, error)
	Update(ctx context.Context, w *weather.Weather) error
	Delete(ctx context.Context, id uuid.UUID) error
}
