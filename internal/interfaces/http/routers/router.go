package router

import (
	authUseCase "github.com/OmidRasouli/weather-api/internal/application/auth"
	"github.com/OmidRasouli/weather-api/internal/application/interfaces"
	"github.com/OmidRasouli/weather-api/internal/interfaces/http/controller"
	"github.com/OmidRasouli/weather-api/internal/interfaces/http/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(
	weatherController *controller.WeatherController,
	authController *controller.AuthController,
	authUC *authUseCase.UseCase,
	db interfaces.Database,
	redisClient interfaces.Cache) *gin.Engine {
	router := gin.Default()

	// Add CORS middleware to allow cross-origin requests (useful for frontend integration).
	router.Use(cors.Default())

	// Auth routes (public)
	router.POST("/login", authController.Login)

	// Public weather routes (read-only)
	weatherPublic := router.Group("/weather")
	{
		weatherPublic.GET("", weatherController.GetAll)
		// Register static path before parameterized to avoid shadowing
		weatherPublic.GET("/latest/:city", weatherController.GetLatestByCity)
		weatherPublic.GET("/:id", weatherController.GetByID)
	}

	// Protected weather routes (mutating operations require JWT)
	weatherProtected := router.Group("/weather", middleware.JWTAuth(authUC))
	{
		weatherProtected.POST("", weatherController.FetchAndStore)
		weatherProtected.PUT("/:id", weatherController.Update)
		weatherProtected.DELETE("/:id", weatherController.Delete)
	}

	// Add health check routes
	healthController := controller.NewHealthController(db, redisClient)
	router.GET("/health", healthController.BasicHealth)
	router.GET("/health/ready", healthController.ReadinessCheck)
	router.GET("/health/live", healthController.LivenessCheck)

	return router
}
