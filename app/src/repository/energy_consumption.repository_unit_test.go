package repository

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/7yrionLannister/golang-technical-assesment/config"
	"github.com/7yrionLannister/golang-technical-assesment/config/logger"
	"github.com/7yrionLannister/golang-technical-assesment/db"
	"github.com/7yrionLannister/golang-technical-assesment/db/view"
	"github.com/7yrionLannister/golang-technical-assesment/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockDB *test.MockDatabase

// Initialize the mock database and the logger.L.
func TestMain(m *testing.M) {
	logger.L.InitLogger(config.Env.LogLevel)
	db.DB = new(test.MockDatabase)
	db.DB.InitDatabaseConnection()
	mockDB = db.DB.(*test.MockDatabase)

	mockDB.On("Model", mock.Anything).Return(mockDB)
	mockDB.On("Select", mock.Anything, mock.Anything).Return(mockDB)
	mockDB.On("Where", mock.Anything, mock.Anything).Return(mockDB)
	mockDB.On("Group", mock.Anything).Return(mockDB)

	code := m.Run()
	os.Exit(code)
}

func TestGetEnergyConsumptionsByMeterIdBetweenDates_Success(t *testing.T) {
	// For
	meterId := uint(123)
	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now()

	// Expect
	expectedData := []view.EnergyConsumptionDTO{
		{MeterId: meterId, Address: "address", TotalConsumption: 1000.35},
	}

	// When
	mockDB.On("Scan", mock.Anything).
		// Then
		Return(expectedData, nil)
		// When
	mockDB.On("Error").
		// Then
		Return(nil).
		Once()

	// Test
	result, err := GetEnergyConsumptionsByMeterIdBetweenDates([]uint{meterId}, startDate, endDate)

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
	mockDB.On("Error").
		// Then
		Return(expectedErr).
		Once()
	mockDB.On("Scan", mock.Anything).
		// Then
		Return(nil, expectedErr)

	// Test
	result, err := GetEnergyConsumptionsByMeterIdBetweenDates([]uint{meterId}, startDate, endDate)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to query energy consumptions")
	mockDB.AssertExpectations(t)
}
