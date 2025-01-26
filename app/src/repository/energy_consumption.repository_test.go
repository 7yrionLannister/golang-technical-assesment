package repository

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/7yrionLannister/golang-technical-assesment/config"
	"github.com/7yrionLannister/golang-technical-assesment/config/logger"
	"github.com/7yrionLannister/golang-technical-assesment/db"
	"github.com/7yrionLannister/golang-technical-assesment/db/model"
	"github.com/7yrionLannister/golang-technical-assesment/test"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	uuid1 = uuid.MustParse("62b7253d-4f6c-49d8-a0ec-c0452a832d10")
	uuid2 = uuid.MustParse("204757bc-fb23-49c4-9adc-41957ba19cb6")
)

var mockDB *test.MockDatabase

// Initialize the mock database and the logger.
func TestMain(m *testing.M) {
	logger.InitLogger(config.Env.LogLevel)
	db.DB = new(test.MockDatabase)
	db.DB.InitDatabaseConnection()
	mockDB = db.DB.(*test.MockDatabase)

	code := m.Run()
	os.Exit(code)
}

func TestGetEnergyConsumptionsByMeterIdBetweenDates_Success(t *testing.T) {
	// For
	meterId := uint(123)
	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now()

	// Expect
	expectedData := []model.EnergyConsumption{
		{Id: uuid1, DeviceId: 123, Consumption: 1.0, CreatedAt: endDate},
		{Id: uuid2, DeviceId: 123, Consumption: 2.0, CreatedAt: endDate},
	}

	// When
	mockDB.On("Find", mock.Anything, "device_id = (?) AND created_at BETWEEN ? AND ?", []any{meterId, startDate, endDate}).
		// Then
		Return(expectedData, nil)

	// Test
	result, err := GetEnergyConsumptionsByMeterIdBetweenDates(meterId, startDate, endDate)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedData, result)
	mockDB.AssertExpectations(t)
}

func TestGetEnergyConsumptionsByMeterIdBetweenDates_Error(t *testing.T) {
	// For
	meterId := uint(123)
	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now()

	// Expect
	expectedErr := errors.New("database error")

	// When
	mockDB.On("Find", mock.Anything, "device_id = (?) AND created_at BETWEEN ? AND ?", []any{meterId, startDate, endDate}).
		// Then
		Return(nil, expectedErr)

	// Test
	result, err := GetEnergyConsumptionsByMeterIdBetweenDates(meterId, startDate, endDate)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to query energy consumptions")
	mockDB.AssertExpectations(t)
}
