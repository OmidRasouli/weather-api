package weather

import (
	"time"

	"github.com/google/uuid"
)

type weatherModel struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	CityName    string
	Country     string
	Temperature float64
	Description string
	Humidity    int
	WindSpeed   float64
	FetchedAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
