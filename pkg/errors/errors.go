package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application error with HTTP status code and optional details
type AppError struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
	Err     error             `json:"-"` // Not exposed in JSON response
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s (code: %d)", e.Message, e.Code)
}

// NewBadRequest returns a 400 Bad Request error
func NewBadRequest(message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
		Err:     err,
	}
}

// NewNotFound returns a 404 Not Found error
func NewNotFound(message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: message,
		Err:     err,
	}
}

// NewInternalServerError returns a 500 Internal Server Error
func NewInternalServerError(message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
		Err:     err,
	}
}

// NewExternalAPIError returns an appropriate error for external API issues
func NewExternalAPIError(message string, err error, statusCode int) *AppError {
	// Map external API errors to appropriate HTTP status codes
	// Default to 502 Bad Gateway if not specified
	if statusCode == 0 {
		statusCode = http.StatusBadGateway
	}

	return &AppError{
		Code:    statusCode,
		Message: message,
		Err:     err,
	}
}

// ValidationError returns a 400 Bad Request with field validation details
func ValidationError(message string, details map[string]string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
		Details: details,
	}
}
