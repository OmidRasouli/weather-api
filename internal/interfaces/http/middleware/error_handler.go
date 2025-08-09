package middleware

import (
	"github.com/OmidRasouli/weather-api/pkg/errors"
	"github.com/OmidRasouli/weather-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

// ErrorHandler middleware catches panics and converts errors to appropriate responses
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Recover from any panics
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("Panic recovered: %v", err)
				c.AbortWithStatusJSON(500, gin.H{
					"code":    500,
					"message": "Internal server error",
				})
			}
		}()

		c.Next()

		// Check if there were any errors during the request handling
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Check if it's our custom AppError
			if appErr, ok := err.(*errors.AppError); ok {
				// Log error with details but don't expose internal error
				if appErr.Err != nil {
					logger.Errorf("%s: %v", appErr.Message, appErr.Err)
				}

				response := gin.H{
					"code":    appErr.Code,
					"message": appErr.Message,
				}

				// Include validation details if present
				if appErr.Details != nil && len(appErr.Details) > 0 {
					response["details"] = appErr.Details
				}

				c.AbortWithStatusJSON(appErr.Code, response)
				return
			}

			// Default to internal server error for unknown error types
			logger.Errorf("Unhandled error: %v", err)
			c.AbortWithStatusJSON(500, gin.H{
				"code":    500,
				"message": "Internal server error",
			})
		}
	}
}
