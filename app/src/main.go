package main

import (
	"context"
	"fmt"
	"os"

	"github.com/7yrionLannister/golang-technical-assesment/config"
	"github.com/jackc/pgx/v5"
)

func main() {
	err := config.LoadEnv()
	if err != nil {
		fmt.Println("Error loading environment variables: ", err)
		os.Exit(1)
	}
	fmt.Println("Connecting to database: ", config.Env.DataBaseUrl)
	connection, err := pgx.Connect(context.Background(), config.Env.DataBaseUrl)
	if err != nil {
		fmt.Println("Error connecting to database: ", err)
		return
	}
	fmt.Println("Connected to database, ping: ", connection.Ping(context.Background()))
}
