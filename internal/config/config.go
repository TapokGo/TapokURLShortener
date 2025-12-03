package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// TODO: make test for LoadConfig

type Config struct {
	Env         string `yaml:"env"`
	StoragePath string `yaml:"storage_path"`

	HTTPServer struct {
		Address     string        `yaml:"address"`
		Port        int           `yaml:"port"`
		Timeout     time.Duration `yaml:"timeout"`
		IdleTimeout time.Duration `yaml:"idle_timeout"`
	} `yaml:"http_server"`
}

func LoadConfig(path string) (Config, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file %q: %w", path, err)
	}

	var cfg Config
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse config: %w", err)
	}

	if err = applyDefaultAndValidate(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// Set default settings and check important fields
func applyDefaultAndValidate(cfg *Config) error {
	// Default settings
	if cfg.HTTPServer.Address == "" {
		cfg.HTTPServer.Address = "0.0.0.0"
	}
	if cfg.HTTPServer.Port == 0 {
		cfg.HTTPServer.Port = 8080
	}
	if cfg.HTTPServer.Timeout == 0 {
		cfg.HTTPServer.Timeout = 4 * time.Second
	}
	if cfg.HTTPServer.IdleTimeout == 0 {
		cfg.HTTPServer.IdleTimeout = 60 * time.Second
	}

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
