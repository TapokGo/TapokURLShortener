// Package slog provides slog realization of logger.Logger interface
package slog

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Tapok-Go/TestURLShortener/internal/config"
	"github.com/Tapok-Go/TestURLShortener/internal/logger"
)

type slogLogger struct {
	cfg    *config.Config
	logger *slog.Logger
}

// New creates a new logger.Logger and logs file(in env="prod")
func New(cfg *config.Config) (logger.Logger, *os.File, error) {
	var handler slog.Handler
	var logFile *os.File

	// Create logs file only env="prod"
	//Write text logs in console, JSON logs in file
	if cfg.Env == "dev" {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	} else {
		file, err := os.OpenFile(cfg.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to open log file: %w", err)
		}
		logFile = file

		handler = slog.NewJSONHandler(logFile, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	sl := &slogLogger{
		cfg:    cfg,
		logger: slog.New(handler),
	}

	return sl, logFile, nil
}

// Info logs a message at Info level
func (s *slogLogger) Info(msg string, args ...any) {
	s.logger.Info(msg, args...)
}

// Error logs a message at Error level
func (s *slogLogger) Error(msg string, args ...any) {
	s.logger.Error(msg, args...)
}

// Warn logs a message at Warn level
func (s *slogLogger) Warn(msg string, args ...any) {
	s.logger.Warn(msg, args...)
}

// Debug logs a message at Debug level
func (s *slogLogger) Debug(msg string, args ...any) {
	s.logger.Debug(msg, args...)
}

// With return new logger.Logger with context
func (s *slogLogger) With(args ...any) logger.Logger {
	return &slogLogger{
		logger: s.logger.With(args...),
		cfg:    s.cfg,
	}
}
