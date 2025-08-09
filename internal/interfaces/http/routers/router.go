package router

import (
	"github.com/OmidRasouli/weather-api/infrastructure/database"
	"github.com/OmidRasouli/weather-api/infrastructure/database/cache"
	"github.com/OmidRasouli/weather-api/internal/interfaces/http/controller"
	"github.com/OmidRasouli/weather-api/internal/interfaces/http/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(weatherController *controller.WeatherController, authController *controller.AuthController, authUC *authUseCase.UseCase, db database.Database, redisClient cache.RedisClient) *gin.Engine {
	router := gin.Default()

	// Add CORS middleware to allow cross-origin requests (useful for frontend integration).
	router.Use(cors.Default())

	router.Use(middleware.ErrorHandler())

	router.POST("/weather", weatherController.FetchAndStore)
	router.GET("/weather", weatherController.GetAll)
	router.GET("/weather/:id", weatherController.GetByID)
	router.PUT("/weather/:id", weatherController.Update)
	router.DELETE("/weather/:id", weatherController.Delete)
	router.GET("/weather/latest/:city", weatherController.GetLatestByCity)

	// Add health check routes
	healthController := controller.NewHealthController(db, redisClient)
	router.GET("/health", healthController.BasicHealth)
	router.GET("/health/ready", healthController.ReadinessCheck)
	router.GET("/health/live", healthController.LivenessCheck)

	return router
}
