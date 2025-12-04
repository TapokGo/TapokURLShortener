// Main package of url-shortener application
package main

import (
	"flag"
	"log"

	"github.com/Tapok-Go/TestURLShortener/internal/app"
	"github.com/Tapok-Go/TestURLShortener/internal/config"
)

func main() {
	// Get path to config file
	var configPath string
	flag.StringVar(&configPath, "config", "./dev.yaml", "Path to the config file")
	flag.Parse()

	// Loads the application config
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Init app
	app, err := app.New(cfg)
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}
	defer func() {
		err := app.Close()
		if err != nil {
			app.Logger.Error("failed to close app dependencies", "error", err)
		}
	}()

	// Start app
	if err := app.Run(); err != nil {
		app.Logger.Error("application failed", "error", err)
	}
}
