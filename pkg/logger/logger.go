package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	once   sync.Once
	logger *logrus.Logger
)

func getCallerInfo() string {
	for i := 1; i < 10; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		funcName := runtime.FuncForPC(pc).Name()

		if !strings.Contains(funcName, "logger") &&
			!strings.Contains(funcName, "errors") &&
			!strings.Contains(file, "logger") &&
			!strings.Contains(file, "errors") {

			absPath, err := filepath.Abs(file)
			if err == nil {
				file = absPath
			}

			return fmt.Sprintf("%s:%d", file, line)
		}
	}

	return "unknown:0"
}

func InitLogger() {
	once.Do(func() {
		logger = logrus.New()
		logger.SetReportCaller(true)
		logger.SetFormatter(NewCustomFormatter())
		logger.SetLevel(logrus.DebugLevel)
	})
}

func GetLogger() *logrus.Logger {
	if logger == nil {
		log.Fatal("Logger has not been initialized. Please call InitLogger first.")
	}
	return logger
}

func createMap(message, caller string, others ...map[string]interface{}) map[string]interface{} {
	log := make(map[string]interface{})
	for _, other := range others {
		for k, v := range other {
			log[k] = v
		}
	}
	log["caller"] = caller
	log["message"] = message
	return log
}

func Log(message string, others ...map[string]interface{}) {
	caller := getCallerInfo()
	logMap := createMap(message, caller, others...)
	GetLogger().WithFields(logMap).Log(logrus.DebugLevel, message)
}

func Logf(format string, args ...interface{}) {
	caller := getCallerInfo()
	message := fmt.Sprintf(format, args...)
	logMap := createMap(message, caller)
	logger := GetLogger()
	logger.WithFields(logMap).Log(logrus.DebugLevel, message)
}

func Info(message string, others ...map[string]interface{}) {
	caller := getCallerInfo()
	logMap := createMap(message, caller, others...)
	GetLogger().WithFields(logMap).Info(message)
}

func Infof(format string, args ...interface{}) {
	caller := getCallerInfo()
	message := fmt.Sprintf(format, args...)
	logMap := createMap(message, caller)
	logger := GetLogger()
	logger.WithFields(logMap).Info(message)
}

func Debug(message string, others ...map[string]interface{}) {
	caller := getCallerInfo()
	logMap := createMap(message, caller, others...)
	GetLogger().WithFields(logMap).Debug(message)
}

func Debugf(format string, args ...interface{}) {
	caller := getCallerInfo()
	message := fmt.Sprintf(format, args...)
	logMap := createMap(message, caller)
	logger := GetLogger()
	logger.WithFields(logMap).Debug(message)
}

func Warn(message string, others ...map[string]interface{}) {
	caller := getCallerInfo()
	logMap := createMap(message, caller, others...)
	GetLogger().WithFields(logMap).Warn(message)
}

func Warnf(format string, args ...interface{}) {
	caller := getCallerInfo()
	message := fmt.Sprintf(format, args...)
	logMap := createMap(message, caller)
	logger := GetLogger()
	logger.WithFields(logMap).Warn(message)
}

func Error(message string, others ...map[string]interface{}) {
	caller := getCallerInfo()
	logMap := createMap(message, caller, others...)
	GetLogger().WithFields(logMap).Error(message)
}

func Errorf(format string, args ...interface{}) {
	caller := getCallerInfo()
	message := fmt.Sprintf(format, args...)
	logMap := createMap(message, caller)
	logger := GetLogger()
	logger.WithFields(logMap).Error(message)
}

// LogObject prints any object in a readable format with the given message
func LogObject(message string, obj interface{}, others ...map[string]interface{}) {
	caller := getCallerInfo()
	objStr := formatObject(obj)
	fullMessage := fmt.Sprintf("%s\n%s", message, objStr)
	logMap := createMap(fullMessage, caller, others...)
	GetLogger().WithFields(logMap).Log(logrus.DebugLevel, fullMessage)
}

// formatObject converts any object to a pretty-printed string
func formatObject(obj interface{}) string {
	bytes, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error formatting object: %v", err)
	}
	return string(bytes)
}

// LogObjectf prints any object in a readable format with a formatted message
func LogObjectf(format string, obj interface{}, args ...interface{}) {
	caller := getCallerInfo()
	message := fmt.Sprintf(format, args...)
	objStr := formatObject(obj)
	fullMessage := fmt.Sprintf("%s\n%s", message, objStr)
	logMap := createMap(fullMessage, caller)
	logger := GetLogger()
	logger.WithFields(logMap).Log(logrus.DebugLevel, fullMessage)
}

// InfoObject prints any object in a readable format with INFO level
func InfoObject(message string, obj interface{}, others ...map[string]interface{}) {
	caller := getCallerInfo()
	objStr := formatObject(obj)
	fullMessage := fmt.Sprintf("%s\n%s", message, objStr)
	logMap := createMap(fullMessage, caller, others...)
	GetLogger().WithFields(logMap).Info(fullMessage)
}

// DebugObject prints any object in a readable format with DEBUG level
func DebugObject(message string, obj interface{}, others ...map[string]interface{}) {
	caller := getCallerInfo()
	objStr := formatObject(obj)
	fullMessage := fmt.Sprintf("%s\n%s", message, objStr)
	logMap := createMap(fullMessage, caller, others...)
	GetLogger().WithFields(logMap).Debug(fullMessage)
}

// WarnObject prints any object in a readable format with WARN level
func WarnObject(message string, obj interface{}, others ...map[string]interface{}) {
	caller := getCallerInfo()
	objStr := formatObject(obj)
	fullMessage := fmt.Sprintf("%s\n%s", message, objStr)
	logMap := createMap(fullMessage, caller, others...)
	GetLogger().WithFields(logMap).Warn(fullMessage)
}

// ErrorObject prints any object in a readable format with ERROR level
func ErrorObject(message string, obj interface{}, others ...map[string]interface{}) {
	caller := getCallerInfo()
	objStr := formatObject(obj)
	fullMessage := fmt.Sprintf("%s\n%s", message, objStr)
	logMap := createMap(fullMessage, caller, others...)
	GetLogger().WithFields(logMap).Error(fullMessage)
}

// Trace level logging functions
func Trace(message string, others ...map[string]interface{}) {
	caller := getCallerInfo()
	logMap := createMap(message, caller, others...)
	GetLogger().WithFields(logMap).Trace(message)
}

func Tracef(format string, args ...interface{}) {
	caller := getCallerInfo()
	message := fmt.Sprintf(format, args...)
	logMap := createMap(message, caller)
	logger := GetLogger()
	logger.WithFields(logMap).Trace(message)
}

func TraceObject(message string, obj interface{}, others ...map[string]interface{}) {
	caller := getCallerInfo()
	objStr := formatObject(obj)
	fullMessage := fmt.Sprintf("%s\n%s", message, objStr)
	logMap := createMap(fullMessage, caller, others...)
	GetLogger().WithFields(logMap).Trace(fullMessage)
}

func Fatal(message string, others ...map[string]interface{}) {
	caller := getCallerInfo()
	logMap := createMap(message, caller, others...)
	GetLogger().WithFields(logMap).Fatal(message)
	os.Exit(1)
}

func Fatalf(format string, args ...interface{}) {
	caller := getCallerInfo()
	message := fmt.Sprintf(format, args...)
	logMap := createMap(message, caller)
	logger := GetLogger()
	logger.WithFields(logMap).Fatal(message)
	os.Exit(1)
}
