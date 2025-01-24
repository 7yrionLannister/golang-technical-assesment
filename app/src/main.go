package main

import (
	"context"
	"fmt"

	"github.com/7yrionLannister/golang-technical-assesment/config"
	"github.com/jackc/pgx/v5"
)

func main() {
	config.Setup()
	fmt.Println("Connecting to database: ", config.Env.DataBaseUrl)
	connection, err := pgx.Connect(context.Background(), config.Env.DataBaseUrl)
	if err != nil {
		fmt.Println("Error connecting to database: ", err)
		return
	}
	fmt.Println("Connected to database, ping: ", connection.Ping(context.Background()))
}
