// Package app config provides utilities for initializing and starting application
package app

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/TapokGo/TapokURLShortener/internal/config"
	v1 "github.com/TapokGo/TapokURLShortener/internal/handler/v1"
	"github.com/TapokGo/TapokURLShortener/internal/logger"
	"github.com/TapokGo/TapokURLShortener/internal/logger/slog"
	"github.com/TapokGo/TapokURLShortener/internal/repo/sqlite"
	"github.com/TapokGo/TapokURLShortener/internal/service/url"
	"github.com/go-chi/chi/v5"
)

// App is a model of application dependencies
type App struct {
	cfg        config.Config
	Logger     logger.Logger
	logCloser  io.Closer
	repoCloser io.Closer
	router     http.Handler
}

// New init all dependencies
func New(cfg config.Config) (*App, error) {
	// Logger
	logger, err := slog.New(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init logger: %w", err)
	}

	// Storage
	repo, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	// Service
	service := url.New(repo)

	// Router + handler
	baseURL := cfg.HTTPServer.Address + ":" + strconv.Itoa(cfg.HTTPServer.Port)
	handler := v1.New(service, logger, baseURL)
	r := chi.NewRouter()
	handler.Register(r)

	return &App{
		cfg:        cfg,
		Logger:     logger,
		logCloser:  logger,
		repoCloser: repo,
		router:     r,
	}, nil
}

// Run start application
func (a *App) Run() error {
	addr := a.cfg.HTTPServer.Address + ":" + strconv.Itoa(a.cfg.HTTPServer.Port)
	a.Logger.Info("Application started", "address", addr)

	return http.ListenAndServe(addr, a.router)
}

// Close close all dependencies
func (a *App) Close() error {
	closeErrors := make([]error, 0, 2)
	// Close logger
	if a.logCloser != nil {
		err := a.logCloser.Close()
		a.logCloser = nil
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
