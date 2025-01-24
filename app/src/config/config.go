package config

import (
	"log"
	"log/slog"
	"os"

	"github.com/7yrionLannister/golang-technical-assesment/config/logger"
)

func Setup() {
	// Load environment variables
	err := LoadEnv()
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err) // use log as logger is not initialized yet
	}

	// Set the structured logger
	logger.InitLogger(Env.LogLevel)

	// Migrate the database
	err = MigrateUp(Env.DataBaseUrl)
	if err != nil {
		logger.Error("Error migrating database", slog.Any("error", err))
		os.Exit(1)
	}
}
