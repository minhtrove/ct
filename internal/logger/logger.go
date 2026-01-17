// Package logger This file sets up a structured logging system using the zap logging library.
package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// InitLogger initializes the global logger
func InitLogger() error {
	env := os.Getenv("APP_ENV") // dev | prod

	var config zap.Config

	if env == "dev" {
		config = zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		config = zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	// Custom encoder config
	config.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var err error
	Logger, err = config.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
	if err != nil {
		return err
	}

	return nil
}

// ServiceLogger creates a logger with service name prefix
func ServiceLogger(serviceName string) *zap.Logger {
	if Logger == nil {
		// Fallback to basic logger if not initialized
		Logger, _ = zap.NewProduction()
	}
	return Logger.Named(serviceName)
}

// Info logs an info message with service name
func Info(serviceName, message string, fields ...zap.Field) {
	ServiceLogger(serviceName).Info(message, fields...)
}

// Error logs an error message with service name
func Error(serviceName, message string, fields ...zap.Field) {
	ServiceLogger(serviceName).Error(message, fields...)
}

// Debug logs a debug message with service name
func Debug(serviceName, message string, fields ...zap.Field) {
	ServiceLogger(serviceName).Debug(message, fields...)
}

// Warn logs a warning message with service name
func Warn(serviceName, message string, fields ...zap.Field) {
	ServiceLogger(serviceName).Warn(message, fields...)
}

// Fatal logs a fatal message with service name and exits
func Fatal(serviceName, message string, fields ...zap.Field) {
	ServiceLogger(serviceName).Fatal(message, fields...)
}
