package config

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/7yrionLannister/golang-technical-assesment/config/logger"
	"github.com/7yrionLannister/golang-technical-assesment/util"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Run the golang-migrate migrations defined in the db.migrations folder
func MigrateUp() error {
	logger.L.Debug("Applying migrations...")
	sqlDB, err := sql.Open("pgx", Env.DataBaseUrl)
	if err != nil {
		logger.L.Error("Failed to open database", slog.Any("error", err))
		os.Exit(1)
	}
	defer sqlDB.Close()

	// Run migrations with golang-migrate
	m, err := migrate.New(
		"file://./config/db.migrations",
		Env.DataBaseUrl,
	)
	if err != nil {
		return util.HandleError(err, "Failed to create migrate instance")
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return util.HandleError(err, "Failed to apply migrations")
	}

	logger.L.Debug("Migrations applied successfully!")
	return nil
}
