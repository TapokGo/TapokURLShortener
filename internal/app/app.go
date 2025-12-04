package app

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Tapok-Go/TestURLShortener/internal/config"
)

type App struct {
	Cfg     config.Config
	Logger  *slog.Logger
	logFile *os.File
}

// Init app
func New(cfg config.Config) (*App, error) {
	logger, logFile, err := initLogger(&cfg)
	if err != nil {
		return nil, fmt.Errorf("faield to init logger: %w", err)
	}

	//TODO: init storage - sqlite cuz pet-project

	//TODO: init router - chi, chi-render

	return &App{
		Cfg:     cfg,
		Logger:  logger,
		logFile: logFile,
	}, nil
}

// Start app
func (a *App) Run() error {
	a.Logger.Info("Application started", "env", a.Cfg.Env)

	fmt.Println(a.Cfg)

	return nil

	//TODO: start server
}

// Close logs file
func (a *App) Close() error {
	if a.logFile != nil {
		return a.logFile.Close()
	}
	return nil
}

// Init logger
func initLogger(cfg *config.Config) (*slog.Logger, *os.File, error) {
	var handler slog.Handler
	var log_file *os.File

	if cfg.Env == "dev" {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	} else {
		file, err := os.OpenFile(cfg.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to open log file: %w", err)
		}
		log_file = file

		handler = slog.NewJSONHandler(log_file, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	logger := slog.New(handler)
	return logger, log_file, nil
}
