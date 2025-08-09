package service

import (
	"os"
	"testing"

	"github.com/OmidRasouli/weather-api/internal/testhelpers"
)

func TestMain(m *testing.M) {
	// Initialize logger for all tests in this package
	testhelpers.InitTestLogger()

	// Run tests
	code := m.Run()
	os.Exit(code)
}
