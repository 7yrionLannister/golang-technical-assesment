package db

import (
	"github.com/7yrionLannister/golang-technical-assesment/config"
	"github.com/7yrionLannister/golang-technical-assesment/config/logger"
	"github.com/7yrionLannister/golang-technical-assesment/util"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var gormConfig = &gorm.Config{
	NamingStrategy: schema.NamingStrategy{
		SingularTable: true,
	},
	Logger: gormLogger.Default.LogMode(gormLogger.Silent),
}

// Database abstraction
type Database interface {
	InitDatabaseConnection() error                  // Setup global database connection
	Find(out any, query string, args ...any) error  // Find records that match the query
	CreateInBatches(value any, batchSize int) error // Create records in batches
}

// Specific implementation of the Database interface using gorm
type GormDatabase struct {
	DB *gorm.DB
}

func (g *GormDatabase) Find(out any, query string, args ...any) error {
	return g.DB.Where(query, args...).Find(out).Error
}

func (g *GormDatabase) CreateInBatches(value any, batchSize int) error {
	return g.DB.CreateInBatches(value, batchSize).Error
}

// Global gorm database connection.
// Call [InitDatabaseConnection] once at the beggining of the application to initialize the connection.
var DB Database

func (g *GormDatabase) InitDatabaseConnection() error {
	// Connect gorm to database
	logger.L.Debug("Connecting to database at " + config.Env.DataBaseUrl)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: "pgx",
		DSN:        config.Env.DataBaseUrl,
	}), gormConfig)
	if err != nil {
		return util.HandleError(err, "Failed to connect to database")
	}
	g.DB = db
	logger.L.Debug("Connected to database")
	return nil
}
