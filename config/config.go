package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	OpenWeather OpenWeatherConfig
	Redis       RedisConfig
}

type ServerConfig struct {
	Port int `env:"SERVER_PORT"`
}

type DatabaseConfig struct {
	Host            string        `env:"DB_HOST"`
	Port            string        `env:"DB_PORT"`
	User            string        `env:"DB_USER"`
	Password        string        `env:"DB_PASS"`
	DBName          string        `env:"DB_NAME"`
	SSLMode         string        `env:"DB_SSL_MODE"`
	MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS"`
	MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS"`
	ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME"`
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST"`
	Port     int    `env:"REDIS_PORT"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB"`
	TTL      int    `env:"REDIS_TTL"`
}

type OpenWeatherConfig struct {
	APIKey string `env:"OPENWEATHER_KEY"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	cfg.Database.Port = "5432"
	return &cfg, nil
}
