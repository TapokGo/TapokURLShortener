// Package logger provides an interface for work with the logger
package logger

// Logger defines the contract for logger
type Logger interface {
	// Info logs a message at Info level.
	// Args must be provided as key-value pairs
	Info(msg string, args ...any)

	// Warn logs a message at Warn level.
	// Args must be provided as key-value pairs
	Warn(msg string, args ...any)

	// Error logs a message at Error level.
	// Args must be provided as key-value pairs
	Error(msg string, args ...any)

	// Debug logs a message at Debug level.
	// Args must be provided as key-value pairs
	Debug(msg string, args ...any)

	// Return new logger with context(static key-value pair).
	// Args must be a provided ad key-value pairs
	With(args ...any) Logger
}
