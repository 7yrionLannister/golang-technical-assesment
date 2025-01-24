package config

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/7yrionLannister/golang-technical-assesment/config/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func MigrateUp(url string) error {
	sqlDB, err := sql.Open("pgx", Env.DataBaseUrl)
	if err != nil {
		logger.Error("Failed to open database", slog.Any("error", err))
		os.Exit(1)
	}
	defer sqlDB.Close()

	// Run migrations with golang-migrate
	m, err := migrate.New(
		"file://./db/migrations",
		Env.DataBaseUrl,
	)
	if err != nil {
		msg := "Failed to create migrate instance"
		e := fmt.Errorf("%s: %w", strings.ToLower(msg), err)
		logger.Error(msg, slog.Any("error", e))
		return e
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		msg := "failed to apply migrations"
		e := fmt.Errorf("%s: %w", strings.ToLower(msg), err)
		logger.Error(msg, slog.Any("error", e))
		return e
	}

	logger.Debug("Migrations applied successfully!")
	return nil
}
