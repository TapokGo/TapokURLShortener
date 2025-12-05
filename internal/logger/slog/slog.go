package slog

import (
	"log/slog"

	"github.com/Tapok-Go/TestURLShortener/internal/config"
	"github.com/Tapok-Go/TestURLShortener/internal/logger"
)

type slogLogger struct {
	cfg    *config.Config
	logger *slog.Logger
}

func NewSlogLogger(cfg *config.Config, logger *slog.Logger) *slogLogger {
	return &slogLogger{
		cfg:    cfg,
		logger: logger,
	}
}

func (s *slogLogger) Info(msg string, args ...any) {
	s.logger.Info(msg, args...)
}

func (s *slogLogger) Error(msg string, args ...any) {
	s.logger.Error(msg, args...)
}

func (s *slogLogger) Warn(msg string, args ...any) {
	s.logger.Warn(msg, args...)
}

func (s *slogLogger) Debug(msg string, args ...any) {
	s.logger.Debug(msg, args...)
}

func (s *slogLogger) With(args ...any) logger.Logger {
	return &slogLogger{
		logger: s.logger.With(args...),
		cfg:    s.cfg,
	}
}
