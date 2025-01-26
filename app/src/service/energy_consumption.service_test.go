package service

import (
	"errors"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/7yrionLannister/golang-technical-assesment/config"
	"github.com/7yrionLannister/golang-technical-assesment/config/logger"
	"github.com/7yrionLannister/golang-technical-assesment/controller/dto"
	"github.com/7yrionLannister/golang-technical-assesment/db"
	"github.com/7yrionLannister/golang-technical-assesment/db/model"
	"github.com/7yrionLannister/golang-technical-assesment/test"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestGetEnergyConsumptionsMonthly_Success(t *testing.T) {
	// For
	// Input
	metersIds := []uint{0, 1}
	startDate, _ := time.Parse("2006-01-02", "2025-01-01")
	endDate, _ := time.Parse("2006-01-02", "2025-02-28")

	// Expect
	// Database state
	data := make([]model.EnergyConsumption, 4)
	for i := range 4 {
		createdAt := startDate
		if i%2 == 0 {
			// even index records the next month
			createdAt = createdAt.AddDate(0, 1, 0)
		}
		randomId, _ := uuid.NewRandom()
		data[i] = model.EnergyConsumption{
			Id:          randomId,
			DeviceId:    uint(i % 2), // 0 or 1
			Consumption: float64(i*10 + 1),
			CreatedAt:   createdAt,
		}
	}

	// When
	mockDB.On("Find", mock.Anything, "device_id = (?) AND created_at BETWEEN ? AND ?", []any{metersIds[1], startDate, startDate.AddDate(0, 1, 0)}).
		// Then
		Return([]model.EnergyConsumption{data[1], data[3]}, nil)
	mockDB.On("Find", mock.Anything, "device_id = (?) AND created_at BETWEEN ? AND ?", []any{metersIds[0], startDate.AddDate(0, 1, 0), endDate}).
		// Then
		Return([]model.EnergyConsumption{data[0], data[2]}, nil)
	mockDB.On("Find", mock.Anything, "device_id = (?) AND created_at BETWEEN ? AND ?", []any{metersIds[1], startDate.AddDate(0, 1, 0), endDate}).
		// Then
		Return([]model.EnergyConsumption{}, nil)
	mockDB.On("Find", mock.Anything, "device_id = (?) AND created_at BETWEEN ? AND ?", []any{metersIds[0], startDate, startDate.AddDate(0, 1, 0)}).
		// Then
		Return([]model.EnergyConsumption{}, nil)

	// Test
	result, err := GetEnergyConsumptions(metersIds, startDate, endDate, "monthly")

	expectedResult := []dto.EnergyConsumptionDTO{
		{MeterId: 0, Address: mock.Anything, Active: []float64{0, 22}},
		{MeterId: 1, Address: mock.Anything, Active: []float64{42, 0}},
	}
	// Assert
	assert.NoError(t, err)
	assert.Len(t, result.DataGraph, 2)
	assert.Condition(t, func() bool {
		for index, expected := range expectedResult {
			actual := result.DataGraph[index]
			// Ignore address in comparison
			equal := expected.MeterId == actual.MeterId && reflect.DeepEqual(expected.Active, actual.Active)
			if !equal {
				return false
			}
		}
		return true
	}, "result.DataGraph is not equal to expectedResult")
	assert.Equal(t, []string{"January 2025", "February 2025"}, result.Period)
	mockDB.AssertExpectations(t)
}

func TestGetEnergyConsumptionsMonthly_Error(t *testing.T) {
	// For
	// Input
	metersIds := []uint{0, 1}
	startDate, _ := time.Parse("2006-01-02", "2024-01-01")
	endDate, _ := time.Parse("2006-01-02", "2024-02-27")

	// Expect
	expectedErr := errors.New("database error")

	// When
	mockDB.On("Find", mock.Anything, mock.Anything, mock.Anything).
		// Then
		Return(nil, expectedErr)

	// Test
	result, err := GetEnergyConsumptions(metersIds, startDate, endDate, "monthly")

	// Assert
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	mockDB.AssertExpectations(t)
}
