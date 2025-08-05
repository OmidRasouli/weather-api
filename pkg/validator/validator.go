package validator

import (
	"reflect"
	"strings"

	"github.com/OmidRasouli/weather-api/pkg/errors"
	"github.com/OmidRasouli/weather-api/pkg/logger"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Initialize creates and registers custom validators
func Initialize() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register custom validation tags here
		if err := v.RegisterValidation("country", validateCountry); err != nil {
			logger.Fatalf("failed to register 'country' validation: %v", err)
		}

		// Register a function to get JSON field names instead of struct field names
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

// validateCountry checks if the country code/name is valid
func validateCountry(fl validator.FieldLevel) bool {
	// This is a simplified example - in a real app, you'd check against a list of valid countries
	country := fl.Field().String()
	return len(country) >= 2 && len(country) <= 56 // Between 2-char code and longest country name
}

// ValidateRequest validates a struct and returns validation errors in a map
func ValidateRequest(obj interface{}) (map[string]string, error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.Struct(obj)
		if err != nil {
			validationErrors := make(map[string]string)

			for _, err := range err.(validator.ValidationErrors) {
				field := err.Field() // This will now be the JSON field name
				validationErrors[field] = validationErrorMessage(err)
			}

			return validationErrors, errors.ValidationError("Validation failed", validationErrors)
		}
	}

	return nil, nil
}

// validationErrorMessage returns a user-friendly error message based on the validation error
func validationErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return "Value must be greater than or equal to " + err.Param()
	case "max":
		return "Value must be less than or equal to " + err.Param()
	case "country":
		return "Invalid country name or code"
	case "oneof":
		return "Value must be one of: " + err.Param()
	case "gt":
		return "Value must be greater than " + err.Param()
	case "gte":
		return "Value must be greater than or equal to " + err.Param()
	case "lt":
		return "Value must be less than " + err.Param()
	case "lte":
		return "Value must be less than or equal to " + err.Param()
	}

	return "Invalid value"
}
