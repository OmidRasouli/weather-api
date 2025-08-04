package main

import (
	"strconv"

	"github.com/OmidRasouli/weather-api/infrastructure/database/postgres"
	configs "github.com/OmidRasouli/weather-api/internal/configs"
	database "github.com/OmidRasouli/weather-api/internal/database/migrations"
	router "github.com/OmidRasouli/weather-api/internal/interfaces/http/routers"
	"github.com/OmidRasouli/weather-api/pkg/logger"
)

func main() {
	logger.Info("The application is starting...")
	cfg := configs.MustLoad("config/config.yaml")
	RunServer(cfg)
}

func RunServer(cfg *configs.Config) {
	// Set up the HTTP router and register all routes.
	router := router.Setup()

	// Build the server address string using the configured port.
	port := cfg.GetServerConfig().Port
	addr := ":" + strconv.Itoa(port)

	logger.Infof("Server is starting on port %d", port)
	// Start the HTTP server; log an error if the server fails to start.
	if err := router.Run(addr); err != nil {
		logger.Errorf("failed to start server: %v", err)
	}
}

func RunDatabase(cfg *configs.Config) {
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
	migrationInstance, err := database.NewMigrateInstance(db, "db/migrations", cfg.GetDatabaseConfig().DBName)
	if err != nil {
		logger.Errorf("failed to create migration instance: %v", err)
	}

	// Initialize the MigrationManager with the database connection and migration instance.
	migrationManager := database.NewMigrationManager(db, migrationInstance)

	// Run all pending database migrations and log any errors.
	if err := migrationManager.RunMigrations(); err != nil {
		logger.Errorf("failed to run migrations: %v", err)
	}
}
