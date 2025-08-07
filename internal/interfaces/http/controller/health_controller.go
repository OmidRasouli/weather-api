package controller

import (
	"net/http"

	"github.com/OmidRasouli/weather-api/infrastructure/database"
	"github.com/OmidRasouli/weather-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

// HealthController handles health check endpoints
type HealthController struct {
	db    database.Database
	redis database.RedisClient
}

// NewHealthController creates a new health controller
func NewHealthController(db database.Database, redis database.RedisClient) *HealthController {
	return &HealthController{
		db:    db,
		redis: redis,
	}
}

// HealthResponse represents the health check response structure
type HealthResponse struct {
	Status     string            `json:"status"`
	Components map[string]string `json:"components"`
	Version    string            `json:"version"`
}

// BasicHealth godoc
// @Summary      Basic health check
// @Description  Returns 200 OK if the service is running
// @Tags         health
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Router       /health [get]
func (hc *HealthController) BasicHealth(c *gin.Context) {
	response := HealthResponse{
		Status:     "UP",
		Components: map[string]string{"api": "UP"},
		Version:    "1.0.0",
	}

	c.JSON(http.StatusOK, response)
}

// ReadinessCheck godoc
// @Summary      Readiness check
// @Description  Verifies connections to PostgreSQL and Redis
// @Tags         health
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Failure      503  {object}  HealthResponse
// @Router       /health/ready [get]
func (hc *HealthController) ReadinessCheck(c *gin.Context) {
	status := "UP"
	statusCode := http.StatusOK
	components := map[string]string{
		"api":      "UP",
		"database": "UP",
		"redis":    "UP",
	}

	ctx := c.Request.Context()

	// Check database connection
	if err := hc.db.Ping(ctx); err != nil {
		logger.Errorf("Database health check failed: %v", err)
		components["database"] = "DOWN"
		status = "DOWN"
		statusCode = http.StatusServiceUnavailable
	}

	// Check Redis connection
	if hc.redis == nil {
		components["redis"] = "DOWN"
		status = "DOWN"
		statusCode = http.StatusServiceUnavailable
	} else if err := hc.redis.HealthCheck(ctx); err != nil {
		logger.Errorf("Redis health check failed: %v", err)
		components["redis"] = "DOWN"
		status = "DOWN"
		statusCode = http.StatusServiceUnavailable
	}

	response := HealthResponse{
		Status:     status,
		Components: components,
		Version:    "1.0.0",
	}

	c.JSON(statusCode, response)
}

// LivenessCheck godoc
// @Summary      Liveness check
// @Description  Simple check for container orchestration
// @Tags         health
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Router       /health/live [get]
func (hc *HealthController) LivenessCheck(c *gin.Context) {
	response := HealthResponse{
		Status:     "UP",
		Components: map[string]string{"api": "UP"},
		Version:    "1.0.0",
	}

	c.JSON(http.StatusOK, response)
}
