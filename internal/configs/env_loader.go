// env_loader.go
// Handles loading configuration from environment variables

package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/OmidRasouli/weather-api/pkg/logger"
)

// loadEnvFile loads environment variables from .env file
func loadEnvFile() error {
	// Look for .env file in the project root or configs directory
	locations := []string{
		".env",
		"internal/configs/.env",
		"internal/configs/config.env",
	}

	for _, loc := range locations {
		err := godotenv.Load(loc)
		if err == nil {
			logger.Infof("Loaded environment from %s", loc)
			return nil
		}
	}

	// If no .env file is found, log a warning but continue
	logger.Warn("No .env file found, using existing environment variables")
	return nil
}

// loadConfig reads configuration from environment variables.
// Returns a Config struct or an error if loading fails.
func loadConfig() (*Config, error) {
	// First load variables from .env file
	if err := loadEnvFile(); err != nil {
		return nil, err
	}

	// Server config with default fallback
	serverPortStr := os.Getenv("SERVER_PORT")
	serverPort := 8080 // Default port
	if serverPortStr != "" {
		port, err := strconv.Atoi(serverPortStr)
		if err != nil {
			logger.Warnf("Invalid SERVER_PORT value: %s, using default: %d", serverPortStr, serverPort)
		} else {
			serverPort = port
		}
	} else {
		logger.Warn("SERVER_PORT not set, using default: 8080")
	}

	// Database config with better error handling
	dbPortStr := os.Getenv("DB_PORT")
	dbPort := 5432 // Default PostgreSQL port
	if dbPortStr != "" {
		port, err := strconv.Atoi(dbPortStr)
		if err != nil {
			logger.Warnf("Invalid DB_PORT value: %s, using default: %d", dbPortStr, dbPort)
		} else {
			dbPort = port
		}
	} else {
		logger.Warn("DB_PORT not set, using default: 5432")
	}

	// Get other DB settings with defaults
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbUser := getEnvOrDefault("DB_USER", "postgres")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "")
	dbName := getEnvOrDefault("DB_NAME", "weather")
	dbSSLMode := getEnvOrDefault("DB_SSLMODE", "disable")

	// OpenWeather config
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		logger.Warn("OPENWEATHER_API_KEY not set, API calls will likely fail")
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
	}
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
		logger.Warnf("%s not set, using default: %s", key, defaultValue)
		return defaultValue
	}
	return value
}
