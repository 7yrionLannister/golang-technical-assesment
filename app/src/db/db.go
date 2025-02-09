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
	Select(query string, args ...any) Database      // Add a select clause to the query
	Model(value any) Database                       // Specify a model to the query
	Scan(dest any) Database                         // Scan the result into the destination DTO
	Group(query string) Database                    // Add a group by clause to the query
	Where(query string, args ...any) Database       // Add a where clause to the query
	Error() error                                   // Returns the last error that occurred
	InitDatabaseConnection() error                  // Setup global database connection
	Find(out any, args ...any) Database             // Find records that match the conditions
	CreateInBatches(value any, batchSize int) error // Create records in batches
}

// Specific implementation of the Database interface using gorm
type GormDatabase struct {
	GormDb *gorm.DB
}

func (g *GormDatabase) Error() error {
	return g.GormDb.Error
}

func (g *GormDatabase) Model(value any) Database {
	return &GormDatabase{g.GormDb.Model(value)}
}

func (g *GormDatabase) Select(query string, args ...any) Database {
	return &GormDatabase{g.GormDb.Select(query, args...)}
}

func (g *GormDatabase) Where(query string, args ...any) Database {
	return &GormDatabase{g.GormDb.Where(query, args...)}
}

func (g *GormDatabase) Find(out any, conds ...any) Database {
	return &GormDatabase{g.GormDb.Find(out, conds...)}
}

func (g *GormDatabase) Scan(dest any) Database {
	return &GormDatabase{g.GormDb.Scan(dest)}
}

func (g *GormDatabase) Group(query string) Database {
	return &GormDatabase{g.GormDb.Group(query)}
}

func (g *GormDatabase) CreateInBatches(value any, batchSize int) error {
	return g.GormDb.CreateInBatches(value, batchSize).Error
}

// Global database connection.
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
	g.GormDb = db
	logger.L.Debug("Connected to database")
	return nil
}
