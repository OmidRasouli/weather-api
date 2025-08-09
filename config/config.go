package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/OmidRasouli/weather-api/pkg/logger"
)

type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	OpenWeather OpenWeatherConfig
	Redis       RedisConfig
}

type ServerConfig struct {
	Port int `envconfig:"SERVER_PORT"`
}

type DatabaseConfig struct {
	Host            string        `envconfig:"DB_HOST"`
	Port            string        `envconfig:"DB_PORT"`
	User            string        `envconfig:"DB_USER"`
	Password        string        `envconfig:"DB_PASSWORD"`
	DBName          string        `envconfig:"DB_NAME"`
	SSLMode         string        `envconfig:"DB_SSLMODE"`
	MaxIdleConns    int           `envconfig:"DB_MAX_IDLE_CONNS"`
	MaxOpenConns    int           `envconfig:"DB_MAX_OPEN_CONNS"`
	ConnMaxLifetime time.Duration `envconfig:"DB_CONN_MAX_LIFETIME"`
}

type RedisConfig struct {
	Host     string `envconfig:"REDIS_HOST"`
	Port     int    `envconfig:"REDIS_PORT"`
	Password string `envconfig:"REDIS_PASSWORD"`
	DB       int    `envconfig:"REDIS_DB"`
	TTL      int    `envconfig:"REDIS_TTL"`
}

type OpenWeatherConfig struct {
	APIKey string `envconfig:"OPENWEATHER_API_KEY"`
}

func Load() (*Config, error) {
	// Try to load .env
	_ = godotenv.Load(
		".env",
		"/app/.env",
		"/.env",
		"internal/configs/.env",
		"internal/configs/config.env",
	)

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		logger.Fatalf("failed to process env vars: %v", err)
		return nil, err
	}

	// Debug loaded config
	logger.DebugObject("cfg", cfg)
	return &cfg, nil
}
