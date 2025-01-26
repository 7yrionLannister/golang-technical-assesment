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

const (
	dataFile = "../data/test.csv"
)

var gormConfig = &gorm.Config{
	NamingStrategy: schema.NamingStrategy{
		SingularTable: true,
	},
	Logger: gormLogger.Default.LogMode(gormLogger.Silent),
}

// Database abstraction
type Database interface {
	Find(out any, query string, args ...any) error  // Find records that match the query
	CreateInBatches(value any, batchSize int) error // Create records in batches
}

// Specific implementation of the Database interface using gorm
type gormDatabase struct {
	DB *gorm.DB
}

func (g *gormDatabase) Find(out any, query string, args ...any) error {
	return g.DB.Where(query, args...).Find(out).Error
}

func (g *gormDatabase) CreateInBatches(value any, batchSize int) error {
	return g.DB.CreateInBatches(value, batchSize).Error
}

// Global gorm database connection.
// Call [InitDatabaseConnection] once at the beggining of the application to initialize the connection.
var DB *gormDatabase

func InitDatabaseConnection() error {
	// Connect gorm to database
	logger.Debug("Connecting to database at " + config.Env.DataBaseUrl)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: "pgx",
		DSN:        config.Env.DataBaseUrl,
	}), gormConfig)
	if err != nil {
		return util.HandleError(err, "Failed to connect to database")
	}
	DB = &gormDatabase{db}
	logger.Debug("Connected to database")
	return nil
}
