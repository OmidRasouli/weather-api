package main

import (
	"strconv"

	_ "github.com/OmidRasouli/weather-api/docs" 
	"github.com/OmidRasouli/weather-api/infrastructure/database"
	"github.com/OmidRasouli/weather-api/infrastructure/database/postgres"
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
	db := RunDatabase(cfg)
	RunServer(cfg, db)

	// Initialize validator with custom validations
	validator.Initialize()
}

func RunServer(cfg *configs.Config, db database.Database) {
	weatherRepo := postgresRepo.NewWeatherPostgresRepository(db)
	apiClient := openweather.NewClient(cfg.GetOpenWeather().APIKey)

	weatherController := controller.NewWeatherController(weatherService)
	r := router.Setup(weatherController)
	port := cfg.Server.Port
	addr := ":" + strconv.Itoa(port)
	logger.Infof("Server is starting on port %d", port)

	// Add Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(addr); err != nil {
		logger.Errorf("failed to start server: %v", err)
	}
}

func RunDatabase(cfg *configs.Config) database.Database {
	// Create a new database connection using the configuration values.
	db, err := postgres.NewPostgresConnection(postgres.PostgresConfig{
		Host:     cfg.Database.Host,
		Port:     strconv.Itoa(cfg.Database.Port),
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	})
	if err != nil {
		logger.Errorf("failed to connect to postgres: %v", err)
	}

	// Create a new migration instance for managing database migrations.
	migrationInstance, err := databaseMigration.NewMigrateInstance(db, "/internal/database/migrations", cfg.Database.DBName)
	if err != nil {
		logger.Errorf("failed to create migration instance: %v", err)
		return db // Prevent further migration logic if migration instance creation fails
	}

	// Initialize the MigrationManager with the database connection and migration instance.
	migrationManager := databaseMigration.NewMigrationManager(db, migrationInstance)

	// Run all pending database migrations and log any errors.
	if err := migrationManager.RunMigrations(); err != nil {
		logger.Errorf("failed to run migrations: %v", err)
	}

	return db
}
