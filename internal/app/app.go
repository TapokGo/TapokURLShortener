// Package app config provides utilities for initializing and starting application
package app

import (
	"errors"
	"fmt"
	"os"

	"github.com/Tapok-Go/TestURLShortener/internal/config"
	"github.com/Tapok-Go/TestURLShortener/internal/logger"
	"github.com/Tapok-Go/TestURLShortener/internal/logger/slog"
	"github.com/Tapok-Go/TestURLShortener/internal/repo"
	"github.com/Tapok-Go/TestURLShortener/internal/repo/sqlite"
	"github.com/Tapok-Go/TestURLShortener/internal/service"
)

// App is a model of application dependencies
type App struct {
	cfg        config.Config
	Logger     logger.Logger
	logFile    *os.File
	urlService *service.URLService
	repo       repo.URLStorage
}

// New allows init all dependencies
func New(cfg config.Config) (*App, error) {
	logger, logFile, err := slog.NewSlogLogger(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init logger: %w", err)
	}

	// Storage
	repo, err := sqlite.New(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	// Service
	urlService := service.NewURLService(repo)

	//TODO: init router - chi, chi-render

	return &App{
		cfg:        cfg,
		Logger:     logger,
		logFile:    logFile,
		urlService: urlService,
		repo:       repo,
	}, nil
}

// Run allows start all program
func (a *App) Run() error {
	a.Logger.Info("Application started", "env", a.cfg.Env)

	fmt.Println(a.cfg)

	return nil

	//TODO: start server
}

// Close allows close all dependencies.
func (a *App) Close() error {
	closeErrros := make([]error, 0, 2)
	// Close logger
	if a.logFile != nil {
		err := a.logFile.Close()
		a.logFile = nil
		closeErrros = append(closeErrros, err)
	}

	// Close repo
	if a.repo != nil {
		err := a.repo.Close()
		a.repo = nil
		closeErrros = append(closeErrros, err)
	}
	return errors.Join(closeErrros...)
}
