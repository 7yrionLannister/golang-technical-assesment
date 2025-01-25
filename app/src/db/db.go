package db

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/7yrionLannister/golang-technical-assesment/config"
	"github.com/7yrionLannister/golang-technical-assesment/config/logger"
	"github.com/7yrionLannister/golang-technical-assesment/db/model"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
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

// Global gorm database connection.
// Call [InitDatabaseConnection] once at the beggining of the application to initialize the connection.
var GormDb *gorm.DB

func InitDatabaseConnection() error {
	// Connect gorm to database
	logger.Debug("Connecting to database at " + config.Env.DataBaseUrl)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: "pgx",
		DSN:        config.Env.DataBaseUrl,
	}), gormConfig)
	if err != nil {
		msg := "Failed to connect to database"
		e := fmt.Errorf("%s: %w", strings.ToLower(msg), err)
		logger.Error(msg, slog.Any("error", err))
		return e
	}
	GormDb = db
	logger.Debug("Connected to database")
	return nil
}

// Read data from test.csv and import it into the database
func ImportTestData() error {
	// Read data from file
	logger.Debug("Importing data from file")
	file, err := os.Open(dataFile)
	if err != nil {
		msg := "Failed to open data file"
		e := fmt.Errorf("%s: %w", strings.ToLower(msg), err)
		logger.Error(msg, slog.Any("error", err))
		return e
	}
	defer file.Close()
	// Read all records from csv file
	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		msg := "Failed to read data from file"
		e := fmt.Errorf("%s: %w", strings.ToLower(msg), err)
		logger.Error(msg, slog.Any("error", err))
		return e
	}
	// Store records as model.EnergyConsumption slice
	energyConsumptions := make([]model.EnergyConsumption, 0)
	for _, record := range records {
		deviceId, _ := strconv.Atoi(record[1])
		consumption, _ := strconv.ParseFloat(record[2], 64)
		createdAt, _ := time.Parse("2006-01-02 15:04:05+00", record[3])
		energyConsumptions = append(energyConsumptions, model.EnergyConsumption{
			Id:          uuid.MustParse(record[0]),
			DeviceId:    uint(deviceId),
			Consumption: consumption,
			CreatedAt:   createdAt,
		})
	}
	// Batch insert for efficiency
	GormDb.CreateInBatches(energyConsumptions, len(energyConsumptions))
	logger.Debug("Imported data from file")
	return nil
}

func GetEnergyConsumptionsByMeterIdBetweenDates(meterId uint, startDate time.Time, endDate time.Time) ([]model.EnergyConsumption, error) {
	energyConsumptions := make([]model.EnergyConsumption, 0)
	err := GormDb.Where("device_id = (?) AND created_at BETWEEN ? AND ?", meterId, startDate, endDate).Find(&energyConsumptions).Error
	if err != nil {
		msg := "Failed to query energy consumptions"
		e := fmt.Errorf("%s: %w", strings.ToLower(msg), err)
		logger.Error(msg, slog.Any("error", err))
		return nil, e
	}
	return energyConsumptions, nil
}
