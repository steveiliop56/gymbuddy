package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/goccy/go-yaml"
)

func main() {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(handler)

	configEnv := os.Getenv("CONFIG_PATH")
	configFlag := flag.String("config-path", "config.yml", "Path to the configuration file.")

	flag.Parse()

	configPath := configEnv

	if configPath == "" {
		configPath = *configFlag
	}

	content, err := os.ReadFile(configPath)

	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	var config AppConfig

	err = yaml.Unmarshal(content, &config)

	if err != nil {
		fmt.Println("Error parsing config file:", err)
		return
	}

	logger.Info("successfully read config file")
}
