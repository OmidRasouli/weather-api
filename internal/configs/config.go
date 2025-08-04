// config.go
// Provides configuration loading and access for the VOD Streamer educational project.
// This file defines the structure of the configuration, loads it from a YAML file,
// and exposes helper methods to access different configuration sections.

package configs

import (
	"os"

	"github.com/OmidRasouli/weather-api/pkg/logger"
	"gopkg.in/yaml.v3"
)

var GlobalConfig *Config

// Config is the main configuration struct that holds all app settings.
// It is populated from the YAML config file.
type Config struct {
	Server      ServerConfig      `yaml:"server"`      // Server-related settings (e.g., port)
	Database    DatabaseConfig    `yaml:"database"`    // Database connection settings
	OpenWeather OpenWeatherConfig `yaml:"openweather"` // OpenWeather API settings
	// Add other configuration sections here as needed, e.g.:
	// Database DatabaseConfig `yaml:"database"`
	// OpenWeather OpenWeatherConfig `yaml:"openweather"`
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Port int `yaml:"port"` // Port on which the HTTP server will listen
}

// DatabaseConfig holds the configuration for connecting to a database.
// It includes host, port, user credentials, database name, and SSL mode.
type DatabaseConfig struct {
	Host     string `yaml:"host"`     // Database host address
	Port     int    `yaml:"port"`     // Database port
	User     string `yaml:"user"`     // Username for database authentication
	Password string `yaml:"password"` // Password for database authentication
	DBName   string `yaml:"dbname"`   // Name of the database to connect to
	SSLMode  string `yaml:"sslmode"`  // SSL mode (e.g., disable, require)
}

// OpenWeatherConfig holds the configuration for accessing the OpenWeather API.
// It includes the API key required for authentication.
type OpenWeatherConfig struct {
	APIKey string `yaml:"apiKey"` // API key for OpenWeather service
}

// loadConfig reads and parses the YAML configuration file at the given path.
// Returns a Config struct or an error if loading fails.
func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// MustLoad loads the configuration and panics if there is any error.
// Useful for ensuring the app never starts with invalid or missing config.
// It caches the loaded configuration in GlobalConfig for reuse.
func MustLoad(path string) *Config {
	if GlobalConfig != nil {
		return GlobalConfig
	}
	cfg, err := loadConfig(path)
	if err != nil {
		logger.Fatalf("failed to load config from %s: %v", path, err)
	}
	GlobalConfig = cfg
	return cfg
}

// GetServerConfig returns the server section of the config.
// Use this to access server-related settings such as the HTTP port.
func (c *Config) GetServerConfig() ServerConfig {
	return c.Server
}

func (c *Config) GetDatabaseConfig() DatabaseConfig {
	return c.Database
}

func (c *Config) GetOpenWeatherConfig() OpenWeatherConfig {
	return c.OpenWeather
}
