package weather

import (
	"context"

	"github.com/OmidRasouli/weather-api/infrastructure/database"
	"github.com/OmidRasouli/weather-api/internal/application/interfaces"
	"github.com/OmidRasouli/weather-api/internal/domain/weather"
)

type WeatherPostgresRepository struct {
	db database.Database
}

func NewWeatherPostgresRepository(db database.Database) interfaces.WeatherRepository {
	return &WeatherPostgresRepository{db: db}
}
func (r *WeatherPostgresRepository) Save(ctx context.Context, w *weather.Weather) error {
	return r.db.WithContext(ctx).Create(toDBModel(w)).Error
}

func (r *WeatherPostgresRepository) FindByID(ctx context.Context, id string) (*weather.Weather, error) {
	var model weatherModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return toDomainModel(&model), nil
}

func (r *WeatherPostgresRepository) FindAll(ctx context.Context) ([]*weather.Weather, error) {
	var models []weatherModel
	err := r.db.WithContext(ctx).Find(&models).Error
	if err != nil {
		return nil, err
	}

	var result []*weather.Weather
	for _, m := range models {
		result = append(result, toDomainModel(&m))
	}
	return result, nil
}

func (r *WeatherPostgresRepository) FindLatestByCity(ctx context.Context, city string) (*weather.Weather, error) {
	var m weatherModel
	err := r.db.WithContext(ctx).
		Where("city_name = ?", city).
		Order("fetched_at DESC").
		First(&m).Error
	if err != nil {
		return nil, err
	}
	return toDomainModel(&m), nil
}

func (r *WeatherPostgresRepository) Update(ctx context.Context, w *weather.Weather) error {
	return r.db.WithContext(ctx).Save(toDBModel(w)).Error
}

func (r *WeatherPostgresRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&weatherModel{}, "id = ?", id).Error
}
