package db

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"

	"github.com/7yrionLannister/golang-technical-assesment/config/logger"
	"github.com/7yrionLannister/golang-technical-assesment/db/model"
	"github.com/7yrionLannister/golang-technical-assesment/util"
	"github.com/google/uuid"
)

const (
	dataFile  = "./data/test.csv"
	batchSize = 1000
)

// Read data from test.csv and import it into the database
func ImportTestData() error {
	// Read data from file
	logger.L.Debug("Importing data from file")
	file, err := os.Open(dataFile)
	if err != nil {
		return util.HandleError(err, "Failed to open data file")
	}
	defer file.Close()
	// Read all records from csv file
	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		return util.HandleError(err, "Failed to read data from file")
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

	err = DB.CreateInBatches(energyConsumptions, batchSize)
	if err != nil {
		return util.HandleError(err, "Failed to create in batches")
	}
	logger.L.Debug("Imported data from file")
	return nil
}
