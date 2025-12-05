// Package app config provides utilities for initializing and starting application
package app

import (
	"fmt"
	"os"

	"github.com/Tapok-Go/TestURLShortener/internal/config"
	"github.com/Tapok-Go/TestURLShortener/internal/logger"
	"github.com/Tapok-Go/TestURLShortener/internal/logger/slog"
	"github.com/Tapok-Go/TestURLShortener/internal/repo"
	"github.com/Tapok-Go/TestURLShortener/internal/repo/sqlite"
	"github.com/Tapok-Go/TestURLShortener/internal/service"
)

// TODO: add documentation in all packages

// App is a model of application dependencies
type App struct {
	Cfg        config.Config
	Logger     logger.Logger
	logFile    *os.File
	URLService *service.URLService
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
		Cfg:        cfg,
		Logger:     logger,
		logFile:    logFile,
		URLService: urlService,
		repo:       repo,
	}, nil
}

// Run function allows start all program.
// Return error
func (a *App) Run() error {
	a.Logger.Info("Application started", "env", a.Cfg.Env)

	fmt.Println(a.Cfg)

	return nil

	//TODO: start server
}

// TODO: refactor to correct closing all dependecies
// Close function allows close all dependencies.
// Return error
func (a *App) Close() error {
	// Close logger
	if a.logFile != nil {
		err := a.logFile.Close()
		a.logFile = nil
		return err
	}

	// Close repo
	if a.repo != nil {
		err := a.repo.Close()
		a.repo = nil
		return err
	}
	return nil
}
