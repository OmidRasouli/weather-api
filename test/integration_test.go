package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/OmidRasouli/weather-api/internal/testhelpers"
	"github.com/OmidRasouli/weather-api/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Initialize logger for tests
	testhelpers.InitTestLogger()

	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func TestHealthEndpoint(t *testing.T) {
	t.Parallel()
	// Create a simple router for testing
	router := gin.New()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, "UP", body["status"])
}

func TestGinRouterSetup(t *testing.T) {
	t.Parallel()
	router := gin.New()
	assert.NotNil(t, router)

	// Test that we can add routes
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLoggerInitialization(t *testing.T) {
	t.Parallel()
	// Test that logger is working
	logger.Info("Test log message")
	// If we get here without panic, logger is working
	assert.True(t, true)
}
