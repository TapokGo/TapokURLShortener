package app

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/Tapok-Go/TestURLShortener/internal/config"
)

func Run() {
	// Get path to config file
	var config_path string
	flag.StringVar(&config_path, "config", "./dev.yaml", "Path to the config file")
	flag.Parse()

	cfg, err := config.LoadConfig(config_path)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	fmt.Println(cfg)

	//TODO: init logger - slog

	//TODO: init storage - sqlite cuz pet-project

	//TODO: init router - chi, chi-render

	//TODO: start server
}
