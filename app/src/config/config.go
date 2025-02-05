package config

import (
	"log"
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
	logger.L.InitLogger(Env.LogLevel)

	// Migrate the database
	err = MigrateUp()
	if err != nil {
		logger.L.Error("Error migrating database", "error", err)
		os.Exit(1)
	}
	logger.L.Debug("Initialized application")
}
