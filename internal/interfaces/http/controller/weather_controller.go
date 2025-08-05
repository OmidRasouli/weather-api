package controller

import (
	"context"
	"net/http"

	"github.com/OmidRasouli/weather-api/internal/domain/weather"
	"github.com/OmidRasouli/weather-api/pkg/errors"
	"github.com/gin-gonic/gin"
)

// WeatherServicer defines the interface for weather service operations.
// This interface decouples the controller from the concrete service implementation.
type WeatherServicer interface {
	FetchAndStoreWeather(ctx context.Context, city, country string) (*weather.Weather, error)
	GetLatestWeatherByCity(ctx context.Context, city string) (*weather.Weather, error)
	GetAllWeather(ctx context.Context) ([]*weather.Weather, error)
	GetWeatherByID(ctx context.Context, id string) (*weather.Weather, error)
	UpdateWeather(ctx context.Context, id string, update *weather.Weather) (*weather.Weather, error)
	DeleteWeather(ctx context.Context, id string) error
}

type WeatherController struct {
	service WeatherServicer
}

// NewWeatherController creates a new weather controller with the provided service.
// Using the interface instead of the concrete type enables better testability and flexibility.
func NewWeatherController(service WeatherServicer) *WeatherController {
	return &WeatherController{service: service}
}

type FetchWeatherRequest struct {
	City    string `json:"city" binding:"required"`
	Country string `json:"country" binding:"required"`
}

func (wc *WeatherController) FetchAndStore(c *gin.Context) {
	var req FetchWeatherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errors.NewBadRequest("Invalid request body", err))
		return
	}

	result, err := wc.service.FetchAndStoreWeather(c, req.City, req.Country)
	if err != nil {
		_ = c.Error(errors.NewExternalAPIError("Failed to fetch weather data", err, http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, result)
}

func (wc *WeatherController) GetByCity(c *gin.Context) {
	city := c.Param("city")
	if city == "" {
		_ = c.Error(errors.NewBadRequest("City name is required", nil))
		return
	}

	result, err := wc.service.GetLatestWeatherByCity(c, city)
	if err != nil {
		_ = c.Error(errors.NewNotFound("Weather data not found for the city", err))
		return
	}

	c.JSON(http.StatusOK, result)
}

func (wc *WeatherController) GetAll(c *gin.Context) {
	result, err := wc.service.GetAllWeather(c)
	if err != nil {
		_ = c.Error(errors.NewInternalServerError("Failed to fetch weather records", err))
		return
	}

	c.JSON(http.StatusOK, result)
}

func (wc *WeatherController) GetByID(c *gin.Context) {
	id := c.Param("id")
	result, err := wc.service.GetWeatherByID(c, id)
	if err != nil {
		_ = c.Error(errors.NewNotFound("Weather data not found", err))
		return
	}
	c.JSON(http.StatusOK, result)
}

func (wc *WeatherController) Update(c *gin.Context) {
	id := c.Param("id")
	var req weather.Weather
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	result, err := wc.service.UpdateWeather(c, id, &req)
	if err != nil {
		_ = c.Error(errors.NewInternalServerError("Failed to update weather data", err))
		return
	}
	c.JSON(http.StatusOK, result)
}

func (wc *WeatherController) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := wc.service.DeleteWeather(c, id); err != nil {
		_ = c.Error(errors.NewInternalServerError("Failed to delete weather record", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Weather record deleted"})
}

func (wc *WeatherController) GetLatestByCity(c *gin.Context) {
	city := c.Param("city")
	result, err := wc.service.GetLatestWeatherByCity(c, city)
	if err != nil {
		_ = c.Error(errors.NewNotFound("Weather data not found for the city", err))
		return
	}
	c.JSON(http.StatusOK, result)
}
