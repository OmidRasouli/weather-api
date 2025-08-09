package test

import (
	"os"
	"testing"

	"github.com/OmidRasouli/weather-api/config"
	"github.com/OmidRasouli/weather-api/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestConfigWithEnvVars(t *testing.T) {
	// Initialize logger first
	logger.InitLogger()

	// Store original values to restore later
	originalServerPort := os.Getenv("SERVER_PORT")
	originalDBHost := os.Getenv("DB_HOST")
	originalRedisHost := os.Getenv("REDIS_HOST")

	// Set test environment variables
	os.Setenv("SERVER_PORT", "9000")
	os.Setenv("DB_HOST", "testdb")
	os.Setenv("REDIS_HOST", "testredis")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	assert.NotNil(t, cfg)
	assert.Equal(t, 9000, cfg.Server.Port)
	assert.Equal(t, "testdb", cfg.Database.Host)
	assert.Equal(t, "testredis", cfg.Redis.Host)

	// Cleanup - restore original values
	if originalServerPort == "" {
		os.Unsetenv("SERVER_PORT")
	} else {
		os.Setenv("SERVER_PORT", originalServerPort)
	}
	if originalDBHost == "" {
		os.Unsetenv("DB_HOST")
	} else {
		os.Setenv("DB_HOST", originalDBHost)
	}
	if originalRedisHost == "" {
		os.Unsetenv("REDIS_HOST")
	} else {
		os.Setenv("REDIS_HOST", originalRedisHost)
	}
}

func TestConfigDefaults(t *testing.T) {
	// Initialize logger first
	logger.InitLogger()

	// Store original values to restore later
	originalServerPort := os.Getenv("SERVER_PORT")
	originalDBHost := os.Getenv("DB_HOST")
	originalRedisHost := os.Getenv("REDIS_HOST")

	// Clear any existing env vars that might interfere
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("REDIS_HOST")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// Test that defaults are set properly
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, 5432, cfg.Database.Port)
	assert.Equal(t, "localhost", cfg.Redis.Host)
	assert.Equal(t, 6379, cfg.Redis.Port)

	// Cleanup - restore original values
	if originalServerPort != "" {
		os.Setenv("SERVER_PORT", originalServerPort)
	}
	if originalDBHost != "" {
		os.Setenv("DB_HOST", originalDBHost)
	}
	if originalRedisHost != "" {
		os.Setenv("REDIS_HOST", originalRedisHost)
	}
}
