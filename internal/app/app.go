// Package app config provides utilities for initializing and starting application
package app

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Tapok-Go/TestURLShortener/internal/config"
)

// App is a model of application dependencies
type App struct {
	Cfg     config.Config
	Logger  *slog.Logger
	logFile *os.File
}

// New function allows init all dependencies.
// Get the config.Config struct, return App struct and error
func New(cfg config.Config) (*App, error) {
	logger, logFile, err := initLogger(&cfg)
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

// TODO:extract into a separate package in yhe future perhaps
func initLogger(cfg *config.Config) (*slog.Logger, *os.File, error) {
	var handler slog.Handler
	var logFile *os.File

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

	logger := slog.New(handler)
	return logger, logFile, nil
}
