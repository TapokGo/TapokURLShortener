package main

import (
	"flag"
	"log"

	"github.com/Tapok-Go/TestURLShortener/internal/app"
	"github.com/Tapok-Go/TestURLShortener/internal/config"
)

func main() {
	// Get path to config file
	var config_path string
	flag.StringVar(&config_path, "config", "./dev.yaml", "Path to the config file")
	flag.Parse()

	// Loads the application config
	cfg, err := config.LoadConfig(config_path)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Init app
	app, err := app.New(cfg)
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}
	defer app.Close()

	// Start app
	if err := app.Run(); err != nil {
		app.Logger.Error("application failed", "error", err)
	}
}
