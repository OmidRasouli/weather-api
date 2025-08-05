package router

import (
	"github.com/OmidRasouli/weather-api/internal/interfaces/http/controller"
	"github.com/OmidRasouli/weather-api/internal/interfaces/http/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(weatherController *controller.WeatherController) *gin.Engine {
	router := gin.Default()

	// Add CORS middleware to allow cross-origin requests (useful for frontend integration).
	router.Use(cors.Default())

	router.Use(middleware.ErrorHandler())

	// Health check
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	router.POST("/weather", weatherController.FetchAndStore)
	router.GET("/weather", weatherController.GetAll)
	router.GET("/weather/:id", weatherController.GetByID)
	router.PUT("/weather/:id", weatherController.Update)
	router.DELETE("/weather/:id", weatherController.Delete)
	router.GET("/weather/latest/:city", weatherController.GetLatestByCity)

	return router
}
