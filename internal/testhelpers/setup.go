package testhelpers

import (
	"sync"

	"github.com/OmidRasouli/weather-api/pkg/logger"
)

var (
	loggerInitOnce sync.Once
)

// InitTestLogger initializes the logger for tests
// It's safe to call multiple times - will only initialize once
func InitTestLogger() {
	loggerInitOnce.Do(func() {
		logger.InitLogger()
	})
}
