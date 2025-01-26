package main

import (
	"log/slog"
	"os"

	"github.com/7yrionLannister/golang-technical-assesment/config"
	"github.com/7yrionLannister/golang-technical-assesment/config/logger"
	"github.com/7yrionLannister/golang-technical-assesment/db"
	"github.com/7yrionLannister/golang-technical-assesment/router"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Setup()
	dbInit()
	app := gin.Default()
	// middleware.Setup(app)
	router.Setup(app)
	app.Run()
}

func dbInit() {
	// Connect to the database
	db.DB = new(db.GormDatabase)
	err := db.DB.InitDatabaseConnection()
	if err != nil {
		logger.Error("Error initializing database connection", slog.Any("error", err))
		os.Exit(1)
	}
	err = db.ImportTestData()
	if err != nil {
		logger.Error("Error importing test data", slog.Any("error", err))
	}
}
