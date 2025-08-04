package logger

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// Custom color codes based on hex values
const (
	// Error: cc0000 (Red)
	colorError = "\033[38;2;204;0;0m"

	// Info: 0099cc (Blue)
	colorInfo = "\033[38;2;0;153;204m"

	// Success: 007e33 (Green)
	colorSuccess = "\033[38;2;0;126;51m"

	// Warning: ff8800 (Orange)
	colorWarn = "\033[38;2;255;136;0m"

	// Reset color
	colorReset = "\033[0m"
)

// CustomFormatter implements logrus.Formatter interface
type CustomFormatter struct {
	*logrus.TextFormatter
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Format message with colors based on level
	var levelColor string
	switch entry.Level {
	case logrus.ErrorLevel:
		levelColor = colorError
	case logrus.WarnLevel:
		levelColor = colorWarn
	case logrus.InfoLevel:
		levelColor = colorInfo
	case logrus.DebugLevel:
		levelColor = colorSuccess
	case logrus.TraceLevel:
		levelColor = colorSuccess
	default:
		levelColor = colorReset
	}

	// Get the caller info and message from entry data
	var message string
	if msg, ok := entry.Data["message"].(string); ok {
		message = msg
	} else {
		message = entry.Message // fallback to entry message if not in data
	}

	var fileLink string
	if caller, ok := entry.Data["caller"].(string); ok {
		file, line := parseCallerInfo(caller)
		if file != "unknown" {
			fileLink = fmt.Sprintf("file://%s:%s", file, line)
		}
	}

	// Create the formatted message
	var msg string
	if fileLink != "" {
		msg = fmt.Sprintf("%s[%s] %s (%s)%s\n",
			levelColor,
			strings.ToUpper(entry.Level.String()),
			message,
			fileLink,
			colorReset,
		)
	} else {
		msg = fmt.Sprintf("%s[%s] %s%s\n",
			levelColor,
			strings.ToUpper(entry.Level.String()),
			message,
			colorReset,
		)
	}

	return []byte(msg), nil
}

func parseCallerInfo(caller string) (string, string) {
	parts := strings.Split(caller, ":")
	if len(parts) != 2 {
		return "unknown", "0"
	}
	return parts[0], parts[1]
}

func NewCustomFormatter() *CustomFormatter {
	return &CustomFormatter{
		TextFormatter: &logrus.TextFormatter{
			DisableTimestamp: false,
			FullTimestamp:    true,
		},
	}
}
