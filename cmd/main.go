package main

import (
	"strconv"

	_ "github.com/OmidRasouli/weather-api/docs"
	"github.com/OmidRasouli/weather-api/infrastructure/database"
	"github.com/OmidRasouli/weather-api/infrastructure/database/postgres"
	"github.com/OmidRasouli/weather-api/infrastructure/database/redis"
	"github.com/OmidRasouli/weather-api/internal/application/services"
	configs "github.com/OmidRasouli/weather-api/internal/configs"
	databaseMigration "github.com/OmidRasouli/weather-api/internal/database/migrations"
	postgresRepo "github.com/OmidRasouli/weather-api/internal/infrastructure/database/postgres/weather"
	"github.com/OmidRasouli/weather-api/internal/infrastructure/openweather"
	"github.com/OmidRasouli/weather-api/internal/interfaces/http/controller"
	router "github.com/OmidRasouli/weather-api/internal/interfaces/http/routers"
	"github.com/OmidRasouli/weather-api/pkg/logger"
	"github.com/OmidRasouli/weather-api/pkg/validator"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Weather API
// @version         1.0
// @description     A RESTful API for weather data management
// @host            localhost:8080
// @BasePath        /
// @schemes         http
func main() {
	logger.InitLogger()
	logger.Info("The application is starting...")
	cfg := configs.MustLoad()
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

func RunServer(cfg *configs.Config, db database.Database, rd database.RedisClient) {
	weatherRepo := postgresRepo.NewWeatherPostgresRepository(db)
	apiClient := openweather.NewClient(cfg.GetOpenWeather().APIKey)

	// Pass Redis client to the weather service
	weatherService := services.NewWeatherService(weatherRepo, apiClient, rd)
	weatherController := controller.NewWeatherController(weatherService)
	r := router.Setup(weatherController, db, rd)
	port := cfg.Server.Port
	addr := ":" + strconv.Itoa(port)
	logger.Infof("Server is starting on port %d", port)

	// Add Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(addr); err != nil {
		logger.Errorf("failed to start server: %v", err)
	}
}

func RunDatabase(cfg *configs.Config) (database.Database, database.RedisClient) {
	// Create a new database connection using the configuration values.
	db, err := postgres.NewPostgresConnection(postgres.PostgresConfig{
		Host:     cfg.Database.Host,
		Port:     "5432", //strconv.Itoa(cfg.Database.Port),
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	})
	if err != nil {
		logger.Errorf("failed to connect to postgres: %v", err)
	}

	// Initialize Redis client
	redisClient, err := redis.NewRedisConnection(cfg.GetRedis())
	if err != nil {
		logger.Warnf("Failed to connect to Redis: %v. Continuing without caching.", err)
	}

	// Create a new migration instance for managing database migrations.
	migrationInstance, err := databaseMigration.NewMigrateInstance(db, "/internal/database/migrations", cfg.Database.DBName)
	if err != nil {
		logger.Errorf("failed to create migration instance: %v", err)
		return db, redisClient // Prevent further migration logic if migration instance creation fails
	}

	// Initialize the MigrationManager with the database connection and migration instance.
	migrationManager := databaseMigration.NewMigrationManager(db, migrationInstance)

	// Run all pending database migrations and log any errors.
	if err := migrationManager.RunMigrations(); err != nil {
		logger.Errorf("failed to run migrations: %v", err)
	}

	return db, redisClient
}
