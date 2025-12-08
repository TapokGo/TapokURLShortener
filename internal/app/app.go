// Package app config provides utilities for initializing and starting application
package app

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Tapok-Go/TestURLShortener/internal/config"
	"github.com/Tapok-Go/TestURLShortener/internal/logger"
	"github.com/Tapok-Go/TestURLShortener/internal/logger/slog"
	"github.com/Tapok-Go/TestURLShortener/internal/repo/sqlite"
	"github.com/Tapok-Go/TestURLShortener/internal/service/url_service"
)

// App is a model of application dependencies
type App struct {
	cfg        config.Config
	Logger     logger.Logger
	logFile    *os.File
	repoCloser io.Closer
}

// New init all dependencies
func New(cfg config.Config) (*App, error) {
	// Logger
	logger, logFile, err := slog.New(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init logger: %w", err)
	}

	// Storage
	repo, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	// Service
	_ = url_service.New(repo)

	//TODO: init router - chi, chi-render

	return &App{
		cfg:        cfg,
		Logger:     logger,
		logFile:    logFile,
		repoCloser: repo,
	}, nil
}

// Run start application
func (a *App) Run() error {
	a.Logger.Info("Application started", "env", a.cfg.Env)

	fmt.Println(a.cfg)

	return nil
}

// Close close all dependencies
func (a *App) Close() error {
	closeErrors := make([]error, 0, 2)
	// Close logger
	if a.logFile != nil {
		err := a.logFile.Close()
		a.logFile = nil
		closeErrors = append(closeErrors, err)
	}

	// Close repo
	if a.repoCloser != nil {
		err := a.repoCloser.Close()
		a.repoCloser = nil
		closeErrors = append(closeErrors, err)
	}

	return errors.Join(closeErrors...)
}
