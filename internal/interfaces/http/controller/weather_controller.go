package controller

import (
	"net/http"

	"github.com/OmidRasouli/weather-api/internal/application/services"
	"github.com/gin-gonic/gin"
)

type WeatherController struct {
	service *services.WeatherService
}

func NewWeatherController(service *services.WeatherService) *WeatherController {
	return &WeatherController{service: service}
}

type FetchWeatherRequest struct {
	City    string `json:"city" binding:"required"`
	Country string `json:"country" binding:"required"`
}

func (wc *WeatherController) FetchAndStore(c *gin.Context) {
	var req FetchWeatherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	result, err := wc.service.FetchAndStoreWeather(c, req.City, req.Country)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (wc *WeatherController) GetByCity(c *gin.Context) {
	city := c.Param("city")
	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "city is required"})
		return
	}

	result, err := wc.service.GetLatestWeatherByCity(c, city)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "weather data not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (wc *WeatherController) GetAll(c *gin.Context) {
	result, err := wc.service.GetAllWeather(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch weather records"})
		return
	}

	c.JSON(http.StatusOK, result)
}
