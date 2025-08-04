package main

import (
	"strconv"

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
)

func main() {
	logger.Info("The application is starting...")
	cfg := configs.MustLoad("config/config.yaml")
	db := RunDatabase(cfg)
	RunServer(cfg, db)
}

func RunServer(cfg *configs.Config, db database.Database) {
	weatherRepo := postgresRepo.NewWeatherPostgresRepository(db)
	apiClient := openweather.NewClient(cfg.GetOpenWeatherConfig().APIKey)
	weatherService := services.NewWeatherService(weatherRepo, apiClient)
	weatherController := controller.NewWeatherController(weatherService)
	r := router.Setup(weatherController)
	port := cfg.GetServerConfig().Port
	addr := ":" + strconv.Itoa(port)
	logger.Infof("Server is starting on port %d", port)

	if err := r.Run(addr); err != nil {
		logger.Errorf("failed to start server: %v", err)
	}
}

func RunDatabase(cfg *configs.Config) database.Database {
	// Create a new database connection using the configuration values.
	db, err := postgres.NewPostgresConnection(postgres.PostgresConfig{
		Host:     cfg.GetDatabaseConfig().Host,
		Port:     strconv.Itoa(cfg.GetDatabaseConfig().Port),
		User:     cfg.GetDatabaseConfig().User,
		Password: cfg.GetDatabaseConfig().Password,
		DBName:   cfg.GetDatabaseConfig().DBName,
		SSLMode:  cfg.GetDatabaseConfig().SSLMode,
	})
	if err != nil {
		logger.Errorf("failed to connect to postgres: %v", err)
	}

	// Create a new migration instance for managing database migrations.
	migrationInstance, err := databaseMigration.NewMigrateInstance(db, "db/migrations", cfg.GetDatabaseConfig().DBName)
	if err != nil {
		logger.Errorf("failed to create migration instance: %v", err)
	}

	// Initialize the MigrationManager with the database connection and migration instance.
	migrationManager := databaseMigration.NewMigrationManager(db, migrationInstance)

	// Run all pending database migrations and log any errors.
	if err := migrationManager.RunMigrations(); err != nil {
		logger.Errorf("failed to run migrations: %v", err)
	}

	return db
}
