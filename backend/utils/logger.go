package utils

import (
	"log"
)

// Logger defines the logging interface.
type Logger interface {
	Info(format string, args ...interface{})
	Error(format string, args ...interface{})
	Debug(format string, args ...interface{})
}

// DefaultLogger implements Logger using the standard log package.
type DefaultLogger struct {
	debug bool
}

// NewDefaultLogger creates a new DefaultLogger.
func NewDefaultLogger(debug bool) *DefaultLogger {
	return &DefaultLogger{debug: debug}
}

// Info logs informational messages.
func (l *DefaultLogger) Info(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}

// Error logs error messages.
func (l *DefaultLogger) Error(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}

// Debug logs debug messages (only if debug mode is enabled).
func (l *DefaultLogger) Debug(format string, args ...interface{}) {
	if l.debug {
		log.Printf("[DEBUG] "+format, args...)
	}
}

// NopLogger is a no-op logger that discards all messages.
type NopLogger struct{}

// NewNopLogger creates a new NopLogger.
func NewNopLogger() *NopLogger {
	return &NopLogger{}
}

// Info does nothing.
func (l *NopLogger) Info(format string, args ...interface{}) {}

// Error does nothing.
func (l *NopLogger) Error(format string, args ...interface{}) {}

// Debug does nothing.
func (l *NopLogger) Debug(format string, args ...interface{}) {}
