// env_loader.go
// Handles loading configuration from environment variables

package configs

import (
	"os"
	"strconv"

	"github.com/OmidRasouli/weather-api/pkg/logger"
	"github.com/joho/godotenv"
)

// loadEnvFile loads environment variables from .env file
func loadEnvFile() error {
	// Look for .env file in different locations
	// Docker containers typically have .env at root level
	locations := []string{
		".env",                        // Current working directory (Docker default)
		"/app/.env",                   // If app is in /app directory
		"/.env",                       // Root of container
		"internal/configs/.env",       // Local development
		"internal/configs/config.env", // Alternative local
	}

	for _, loc := range locations {
		err := godotenv.Load(loc)
		if err == nil {
			logger.Infof("Loaded environment from %s", loc)
			return nil
		}
	}

	// If no .env file is found, log a warning but continue
	// This is normal in Docker when using docker-compose env_file
	logger.Warn("No .env file found, using existing environment variables")
	return nil
}

func loadConfig() (*Config, error) {
	if err := loadEnvFile(); err != nil {
		logger.Warnf("Error loading .env file: %v", err)
	}

	// Server config
	serverPortStr := os.Getenv("SERVER_PORT")
	serverPort := 8080
	if serverPortStr != "" {
		if port, err := strconv.Atoi(serverPortStr); err == nil {
			serverPort = port
		}
	}

	// Database config
	dbPortStr := os.Getenv("DB_PORT")
	dbPort := 5432
	if dbPortStr != "" {
		if port, err := strconv.Atoi(dbPortStr); err == nil {
			dbPort = port
		}
	}

	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbUser := getEnvOrDefault("DB_USER", "postgres")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "")
	dbName := getEnvOrDefault("DB_NAME", "weather")
	dbSSLMode := getEnvOrDefault("DB_SSLMODE", "disable")

	// OpenWeather config
	apiKey := os.Getenv("OPENWEATHER_API_KEY")

	// Redis configuration
	redisPortStr := getEnvOrDefault("REDIS_PORT", "6379")
	redisPort := 6379
	if redisPortStr != "" {
		if port, err := strconv.Atoi(redisPortStr); err == nil {
			redisPort = port
		}
	}

	redisTTLStr := getEnvOrDefault("REDIS_TTL", "600")
	redisTTL := 600
	if redisTTLStr != "" {
		if ttl, err := strconv.Atoi(redisTTLStr); err == nil {
			redisTTL = ttl
		}
	}

	redisDB := 0
	redisDBStr := getEnvOrDefault("REDIS_DB", "0")
	if redisDBStr != "" {
		if db, err := strconv.Atoi(redisDBStr); err == nil {
			redisDB = db
		}
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: serverPort,
		},
		Database: DatabaseConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			DBName:   dbName,
			SSLMode:  dbSSLMode,
		},
		OpenWeather: OpenWeatherConfig{
			APIKey: apiKey,
		},
		Redis: RedisConfig{
			Host:     getEnvOrDefault("REDIS_HOST", "localhost"),
			Port:     redisPort,
			Password: getEnvOrDefault("REDIS_PASSWORD", ""),
			DB:       redisDB,
			TTL:      redisTTL,
		},
	}

	logger.Infof("Configuration loaded - Server port: %d", cfg.Server.Port)
	return cfg, nil
}

// MustLoad loads the configuration from environment variables and panics if there is any error.
// It caches the loaded configuration in GlobalConfig for reuse.
func MustLoad() *Config {
	if Global() != nil {
		return Global()
	}
	cfg, err := loadConfig()
	if err != nil {
		logger.Fatalf("failed to load config from environment: %v", err)
	}
	SetGlobal(cfg)
	return cfg
}

// Helper function to get environment variable with fallback default
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
