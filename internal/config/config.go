// Package config provides utilities for getting application settings
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

// Config is a model of application config
type Config struct {
	Env         string `yaml:"env"`
	StoragePath string `yaml:"storage_path"`
	LogPath     string `yaml:"log_path"`
	HTTPServer  struct {
		Address     string        `yaml:"address"`
		Port        int           `yaml:"port"`
		Timeout     time.Duration `yaml:"timeout"`
		IdleTimeout time.Duration `yaml:"idle_timeout"`
	} `yaml:"http_server"`
}

// LoadConfig loads config from YAML with overrides from env.
// Gets the path to the YAMl config, return Config{} and error
func LoadConfig(path string) (*Config, error) {
	var cfg Config

	setDefaults(&cfg)

	if path != "" {
		if err := loadFromYAML(path, &cfg); err != nil {
			return nil, err
		}
	}

	applyEnvOverrides(&cfg)

	if err := validateSettings(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Set default settings
func setDefaults(cfg *Config) {
	cfg.Env = "dev"
	cfg.StoragePath = "./storage/storage.db"
	cfg.Env = "./app.log"

	cfg.HTTPServer.Address = "localhost"
	cfg.HTTPServer.Port = 8080
	cfg.HTTPServer.Timeout, _ = time.ParseDuration("4s")
	cfg.HTTPServer.IdleTimeout, _ = time.ParseDuration("60s")
}

// Validate settings
func validateSettings(cfg *Config) error {
	if cfg.StoragePath == "" {
		return fmt.Errorf("storage_path cannot be empty")
	}

	if cfg.HTTPServer.Port < 1 || cfg.HTTPServer.Port > 65535 {
		return fmt.Errorf("invalid port: %d, must be 1-65535", cfg.HTTPServer.Port)
	}

	if cfg.Env != "prod" && cfg.Env != "dev" {
		return fmt.Errorf("invalid env: %s, must be dev or prod", cfg.Env)
	}

	return nil
}

// Set env config if exists
func applyEnvOverrides(cfg *Config) {
	// Add logging after parsing port and timeout if values incorrect
	if env := os.Getenv("URL_SHORTENER_ENV"); env != "" {
		cfg.Env = env
	}

	if storagePath := os.Getenv("URL_SHORTENER_STORAGE_PATH"); storagePath != "" {
		cfg.StoragePath = storagePath
	}

	if logPath := os.Getenv("URL_SHORTENER_LOG_PATH"); logPath != "" {
		cfg.LogPath = logPath
	}

	if address := os.Getenv("URL_SHORTENER_ADDRESS"); address != "" {
		cfg.HTTPServer.Address = address
	}

	if portStr := os.Getenv("URL_SHORTENER_PORT"); portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err == nil {
			cfg.HTTPServer.Port = port
		}
	}

	if timeoutStr := os.Getenv("URL_SHORTENER_TIMEOUT"); timeoutStr != "" {
		timeout, err := time.ParseDuration(timeoutStr)
		if err == nil {
			cfg.HTTPServer.Timeout = timeout
		}
	}

	if idleTimeoutStr := os.Getenv("URL_SHORTENER_IDLE_TIMEOUT"); idleTimeoutStr != "" {
		timeout, err := time.ParseDuration(idleTimeoutStr)
		if err == nil {
			cfg.HTTPServer.IdleTimeout = timeout
		}
	}
}

func loadFromYAML(path string, cfg *Config) error {
	filePath := filepath.Clean(path)
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read config file %q: %w", path, err)
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	return nil
}
