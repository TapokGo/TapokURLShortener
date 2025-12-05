// Package app config provides utilities for initializing and starting application
package app

import (
	"fmt"
	"os"

	"github.com/Tapok-Go/TestURLShortener/internal/config"
	"github.com/Tapok-Go/TestURLShortener/internal/logger"
	"github.com/Tapok-Go/TestURLShortener/internal/logger/slog"

)

// App is a model of application dependencies
type App struct {
	Cfg     config.Config
	Logger  logger.Logger
	logFile *os.File
}

// New function allows init all dependencies.
// Get the config.Config struct, return App struct and error
func New(cfg config.Config) (*App, error) {
	logger, logFile, err := slog.NewSlogLogger(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init logger: %w", err)
	}

	//TODO: init storage - sqlite cuz pet-project

	//TODO: init router - chi, chi-render

	return &App{
		Cfg:     cfg,
		Logger:  logger,
		logFile: logFile,
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
	if a.logFile != nil {
		err := a.logFile.Close()
		a.logFile = nil
		return err
	}
	return nil
}
