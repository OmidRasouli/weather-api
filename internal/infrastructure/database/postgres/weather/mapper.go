package weather

import "github.com/OmidRasouli/weather-api/internal/domain/weather"

// map from domain to db model
func toDBModel(w *weather.Weather) *weatherModel {
	return &weatherModel{
		ID:          w.ID,
		CityName:    w.CityName,
		Country:     w.Country,
		Temperature: w.Temperature,
		Description: w.Description,
		Humidity:    w.Humidity,
		WindSpeed:   w.WindSpeed,
		FetchedAt:   w.FetchedAt,
		CreatedAt:   w.CreatedAt,
		UpdatedAt:   w.UpdatedAt,
	}
}

// map from db model to domain
func toDomainModel(m *weatherModel) *weather.Weather {
	return &weather.Weather{
		ID:          m.ID,
		CityName:    m.CityName,
		Country:     m.Country,
		Temperature: m.Temperature,
		Description: m.Description,
		Humidity:    m.Humidity,
		WindSpeed:   m.WindSpeed,
		FetchedAt:   m.FetchedAt,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
