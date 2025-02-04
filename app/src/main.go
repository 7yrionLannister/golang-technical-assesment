package main

import (
	"log/slog"
	"os"

	"github.com/7yrionLannister/golang-technical-assesment/config"
	"github.com/7yrionLannister/golang-technical-assesment/config/logger"
	"github.com/7yrionLannister/golang-technical-assesment/db"
	"github.com/7yrionLannister/golang-technical-assesment/middleware"
	"github.com/7yrionLannister/golang-technical-assesment/router"
	"github.com/gin-gonic/gin"
)

// @title           Energy Consumption API
// @version         1.0
// @description     Report the energy consumption of a set of electricity meters.

// @contact.name   Daniel Fernández
// @contact.email  daniel.fernandez3@u.icesi.edu.co

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8181
// @BasePath  /
func main() {
	// TODO use zap for logging
	// TODO in the repository, make a query that brings the report from DB instead of computing it in the service
	app := setup()
	app.Run()
}

func setup() *gin.Engine {
	config.Setup()
	dbInit()
	app := gin.Default()
	middleware.Setup(app)
	router.Setup(app)
	return app
}

func dbInit() {
	// Connect to the database
	db.DB = new(db.GormDatabase)
	err := db.DB.InitDatabaseConnection()
	if err != nil {
		logger.L.Error("Error initializing database connection", slog.Any("error", err))
		os.Exit(1)
	}
	err = db.ImportTestData()
	if err != nil {
		logger.L.Warn("Error importing test data", slog.Any("error", err))
	}
}
