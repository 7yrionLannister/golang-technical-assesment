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
	"github.com/7yrionLannister/golang-technical-assesment/db/view"
	"github.com/7yrionLannister/golang-technical-assesment/test"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockDB *test.MockDatabase

var expectedResult = []dto.EnergyConsumptionDTO{
	{MeterId: 0, Address: mock.Anything, Active: []float64{1, 0}},
	{MeterId: 1, Address: mock.Anything, Active: []float64{0, 11}},
}

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

func mockForPeriods(metersIds []uint) {
	// Expect
	// Database state
	data := make([]view.EnergyConsumptionDTO, 2)
	for i := range 2 {
		data[i] = view.EnergyConsumptionDTO{
			MeterId:          metersIds[i],
			Address:          faker.GetRealAddress().Address,
			TotalConsumption: float64(i*10 + 1),
		}
	}

	// When
	mockDB.On("Scan", mock.Anything).
		// Then
		Return([]view.EnergyConsumptionDTO{data[0], {}}, nil).
		Once()
	mockDB.On("Error").
		// Then
		Return(nil).
		Once()
	mockDB.On("Scan", mock.Anything).
		// Then
		Return([]view.EnergyConsumptionDTO{{}, data[1]}, nil).
		Once()
	mockDB.On("Error").
		// Then
		Return(nil).
		Once()
}

func TestGetEnergyConsumptionsMonthly_Success(t *testing.T) {
	// For
	// Input
	metersIds := []uint{0, 1}
	startDate, _ := time.Parse("2006-01-02", "2025-01-01")
	endDate, _ := time.Parse("2006-01-02", "2025-02-28")

	mockForPeriods(metersIds)

	// Test
	result, err := GetEnergyConsumptions(metersIds, startDate, endDate, "monthly")

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

func TestGetEnergyConsumptionsWeekly_Success(t *testing.T) {
	// For
	// Input
	metersIds := []uint{0, 1}
	startDate, _ := time.Parse("2006-01-02", "2025-06-01")
	endDate, _ := time.Parse("2006-01-02", "2025-06-15")

	// Expect
	// Database state
	data := make([]view.EnergyConsumptionDTO, 2)
	for i := range 2 {
		data[i] = view.EnergyConsumptionDTO{
			MeterId:          metersIds[i],
			Address:          faker.GetRealAddress().Address,
			TotalConsumption: float64(i*10 + 1),
		}
	}

	mockForPeriods(metersIds)

	// Test
	result, err := GetEnergyConsumptions(metersIds, startDate, endDate, "weekly")

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
	assert.Equal(t, []string{"June 1 - June 7", "June 8 - June 14"}, result.Period)
	mockDB.AssertExpectations(t)
}

func TestGetEnergyConsumptionsDaily_Success(t *testing.T) {
	// For
	// Input
	metersIds := []uint{0, 1}
	startDate, _ := time.Parse("2006-01-02", "2025-07-01")
	endDate, _ := time.Parse("2006-01-02", "2025-07-03")

	// Expect
	// Database state
	data := make([]view.EnergyConsumptionDTO, 4)
	for i := range 4 {
		data[i] = view.EnergyConsumptionDTO{
			MeterId:          metersIds[i%2],
			Address:          faker.GetRealAddress().Address,
			TotalConsumption: float64(i*10 + 1),
		}
	}

	mockForPeriods(metersIds)

	// Test
	result, err := GetEnergyConsumptions(metersIds, startDate, endDate, "daily")

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
	assert.Equal(t, []string{"July 1", "July 2"}, result.Period)
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
	mockDB.On("Error").
		// Then
		Return(expectedErr).
		Once()
	// When
	mockDB.On("Scan", mock.Anything, mock.Anything, mock.Anything).
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

func TestGetEnergyConsumptionsWeekly_Error(t *testing.T) {
	// For
	// Input
	metersIds := []uint{0, 1}
	startDate, _ := time.Parse("2006-01-02", "2024-02-01")
	endDate, _ := time.Parse("2006-01-02", "2024-03-01")

	// Expect
	expectedErr := errors.New("database error")

	// When
	mockDB.On("Error").
		// Then
		Return(expectedErr).
		Once()
	// When
	mockDB.On("Scan", mock.Anything, mock.Anything, mock.Anything).
		// Then
		Return(nil, expectedErr)

	// Test
	result, err := GetEnergyConsumptions(metersIds, startDate, endDate, "weekly")

	// Assert
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	mockDB.AssertExpectations(t)
}

func TestGetEnergyConsumptionsDaily_Error(t *testing.T) {
	// For
	// Input
	metersIds := []uint{0, 1}
	startDate, _ := time.Parse("2006-01-02", "2024-02-01")
	endDate, _ := time.Parse("2006-01-02", "2024-03-01")

	// Expect
	expectedErr := errors.New("database error")

	// When
	mockDB.On("Error").
		// Then
		Return(expectedErr).
		Once()
	// When
	mockDB.On("Scan", mock.Anything, mock.Anything, mock.Anything).
		// Then
		Return(nil, expectedErr)

	// Test
	result, err := GetEnergyConsumptions(metersIds, startDate, endDate, "daily")

	// Assert
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	mockDB.AssertExpectations(t)
}
