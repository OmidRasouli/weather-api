package main

import (
	"strconv"

	"github.com/OmidRasouli/weather-api/config"
	_ "github.com/OmidRasouli/weather-api/docs"
	"github.com/OmidRasouli/weather-api/infrastructure/database/cache"
	postgres "github.com/OmidRasouli/weather-api/infrastructure/database/database"
	authUseCase "github.com/OmidRasouli/weather-api/internal/application/auth"
	"github.com/OmidRasouli/weather-api/internal/application/interfaces"
	"github.com/OmidRasouli/weather-api/internal/application/service"
	migration "github.com/OmidRasouli/weather-api/internal/database/migrations"
	authDomain "github.com/OmidRasouli/weather-api/internal/domain/services"
	"github.com/OmidRasouli/weather-api/internal/infrastructure/database/postgres/weather"
	"github.com/OmidRasouli/weather-api/internal/infrastructure/openweather"
	"github.com/OmidRasouli/weather-api/internal/interfaces/http/controller"
	router "github.com/OmidRasouli/weather-api/internal/interfaces/http/routers"
	"github.com/OmidRasouli/weather-api/pkg/logger"
	"github.com/OmidRasouli/weather-api/pkg/validator"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Weather APIServerPort
// @version         1.0
// @description     A RESTful API for weather data management
// @host            localhost:8080
// @BasePath        /
// @schemes         http
func main() {
	logger.InitLogger()
	logger.Info("The application is starting...")
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("failed to load config: %v", err)
	}

	db, redisClient := RunDatabase(cfg)

	// Handle cleanup when application shuts down
	defer func() {
		if db != nil {
			db.Close()
		}
		if redisClient != nil {
			redisClient.Close()
		}
	}()

	RunServer(cfg, db, redisClient)

	// Initialize validator with custom validations
	validator.Initialize()
}

func RunServer(cfg *config.Config, db interfaces.Database, rd interfaces.Cache) {
	weatherRepo := weather.NewWeatherPostgresRepository(db)
	apiClient := openweather.NewClient(cfg.OpenWeather.APIKey)

	// Pass Redis client to the weather service
	weatherService := service.NewWeatherService(weatherRepo, apiClient, rd)
	weatherController := controller.NewWeatherController(weatherService)
	authService := authDomain.NewAuthService()
	authUC := authUseCase.NewUseCase(authService)
	authController := controller.NewAuthController(authUC)
	r := router.Setup(weatherController, authController, authUC, db, rd)
	port := cfg.Server.Port
	addr := ":" + strconv.Itoa(port)
	logger.Infof("Server is starting on port %d", port)

	// Add Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(addr); err != nil {
		logger.Errorf("failed to start server: %v", err)
	}
}

func RunDatabase(cfg *config.Config) (interfaces.Database, interfaces.Cache) {
	// Create a new database connection using the configuration values.
	dbConfig := postgres.PostgresConfig{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		DBName:          cfg.Database.DBName,
		SSLMode:         cfg.Database.SSLMode,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	}
	db, err := postgres.NewPostgresConnection(dbConfig)
	if err != nil {
		logger.Errorf("failed to connect to postgres: %v", err)
	}

	// Initialize Redis client
	redisClient, err := cache.NewRedisConnection(cfg.Redis)
	if err != nil {
		logger.Warnf("Failed to connect to Redis: %v. Continuing without caching.", err)
	}

	// Create a new migration instance for managing database migrations.
	migrationInstance, err := migration.NewMigrateInstance(db, "/internal/database/migrations", cfg.Database.DBName)
	if err != nil {
		logger.Errorf("failed to create migration instance: %v", err)
		return db, redisClient // Prevent further migration logic if migration instance creation fails
	}

	// Initialize the MigrationManager with the database connection and migration instance.
	migrationManager := migration.NewMigrationManager(db, migrationInstance)

	// Run all pending database migrations and log any errors.
	if err := migrationManager.RunMigrations(); err != nil {
		logger.Errorf("failed to run migrations: %v", err)
	}

	return db, redisClient
}
