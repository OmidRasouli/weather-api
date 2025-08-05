package router

import (
	"github.com/OmidRasouli/weather-api/internal/interfaces/http/controller"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(weatherController *controller.WeatherController) *gin.Engine {
	router := gin.Default()

	// Add CORS middleware to allow cross-origin requests (useful for frontend integration).
	router.Use(cors.Default())

	// Health check
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	router.POST("/weather", weatherController.FetchAndStore)
	router.GET("/weather", weatherController.GetAll)
	router.GET("/weather/:id", weatherController.GetByID)                      // NEW
	router.PUT("/weather/:id", weatherController.Update)                       // NEW
	router.DELETE("/weather/:id", weatherController.Delete)                    // NEW
	router.GET("/weather/latest/:cityName", weatherController.GetLatestByCity) // NEW

	return router
}
