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

// New function allows init all dependencies.
// Get the config.Config struct, return App struct and error
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

// Run function allows start all program.
// Return error
func (a *App) Run() error {
	a.Logger.Info("Application started", "env", a.cfg.Env)

	fmt.Println(a.cfg)

	return nil

	//TODO: start server
}

// Close function allows close all dependencies.
// Return error
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
