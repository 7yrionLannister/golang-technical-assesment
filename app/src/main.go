package main

import (
	"log/slog"
	"os"

	"github.com/7yrionLannister/golang-technical-assesment/config"
	"github.com/7yrionLannister/golang-technical-assesment/config/logger"
	"github.com/7yrionLannister/golang-technical-assesment/db"
)

func main() {
	config.Setup()
	err := db.InitDatabaseConnection()
	if err != nil {
		logger.Error("Error initializing database connection", slog.Any("error", err))
		os.Exit(1)
	}
	err = db.ImportTestData()
	if err != nil {
		logger.Error("Error importing test data", slog.Any("error", err))
	}
}
