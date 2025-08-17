package utils

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger with additional convenience methods
type Logger struct {
	*zap.Logger
}

// NewLogger creates a new logger instance based on environment
func NewLogger() (*Logger, error) {
	var zapLogger *zap.Logger
	var err error

	if os.Getenv("GO_ENV") == "development" {
		// Development logger with human-readable output
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		zapLogger, err = config.Build()
	} else {
		// Production logger with JSON output
		config := zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		zapLogger, err = config.Build()
	}

	if err != nil {
		return nil, err
	}

	return &Logger{Logger: zapLogger}, nil
}

// WithRequestID adds request ID to logger context
func (l *Logger) WithRequestID(requestID string) *Logger {
	return &Logger{Logger: l.With(zap.String("request_id", requestID))}
}

// WithUserContext adds user context information to logger
func (l *Logger) WithUserContext(editToken, publicToken string) *Logger {
	fields := []zap.Field{}

	if editToken != "" {
		fields = append(fields, zap.String("edit_token", editToken))
	}

	if publicToken != "" {
		fields = append(fields, zap.String("public_token", publicToken))
	}

	return &Logger{Logger: l.With(fields...)}
}

// LogError logs an error with additional context
func (l *Logger) LogError(err error, message string, fields ...zap.Field) {
	allFields := append(fields, zap.Error(err))
	l.Error(message, allFields...)
}

// LogRequest logs an HTTP request with standard fields
func (l *Logger) LogRequest(method, path, userAgent, ip string, statusCode int, duration int64) {
	l.Info("HTTP Request",
		zap.String("method", method),
		zap.String("path", path),
		zap.String("user_agent", userAgent),
		zap.String("ip", ip),
		zap.Int("status_code", statusCode),
		zap.Int64("duration_ms", duration),
	)
}

// LogDatabaseOperation logs database operations
func (l *Logger) LogDatabaseOperation(operation, table string, duration int64, err error) {
	fields := []zap.Field{
		zap.String("operation", operation),
		zap.String("table", table),
		zap.Int64("duration_ms", duration),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		l.Error("Database operation failed", fields...)
	} else {
		l.Debug("Database operation completed", fields...)
	}
}

// Sugar logger methods for convenience
func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.Sugar().Infow(msg, keysAndValues...)
}

func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.Sugar().Warnw(msg, keysAndValues...)
}

func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.Sugar().Errorw(msg, keysAndValues...)
}

func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.Sugar().Debugw(msg, keysAndValues...)
}
