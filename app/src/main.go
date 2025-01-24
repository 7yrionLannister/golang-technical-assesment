package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/7yrionLannister/golang-technical-assesment/config"
	"github.com/7yrionLannister/golang-technical-assesment/config/logger"
	"github.com/7yrionLannister/golang-technical-assesment/db"
	"github.com/jackc/pgx/v5"
)

func main() {
	config.Setup()
	err := db.ImportTestData()
	if err != nil {
		logger.Debug("Error importing test data", slog.Any("error", err))
	}
	fmt.Println("Connecting to database: ", config.Env.DataBaseUrl)
	connection, err := pgx.Connect(context.Background(), config.Env.DataBaseUrl)
	if err != nil {
		fmt.Println("Error connecting to database: ", err)
		return
	}
	fmt.Println("Connected to database, ping: ", connection.Ping(context.Background()))
}
