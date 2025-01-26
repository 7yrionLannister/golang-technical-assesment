package db

import (
	"errors"
	"testing"
	"time"

	"github.com/7yrionLannister/golang-technical-assesment/db/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	uuid1 = uuid.MustParse("62b7253d-4f6c-49d8-a0ec-c0452a832d10")
	uuid2 = uuid.MustParse("204757bc-fb23-49c4-9adc-41957ba19cb6")
)

func TestFind_Success(t *testing.T) {
	mockDB := new(mockDatabase)
	// For
	expectedData := []model.EnergyConsumption{
		{Id: uuid1, DeviceId: 123, Consumption: 1.0, CreatedAt: time.Now()},
		{Id: uuid2, DeviceId: 123, Consumption: 2.0, CreatedAt: time.Now()},
	}
	query := "device_id = (?)"
	args := []any{123}

	// If
	mockDB.On("Find", mock.Anything, query, args).
		// Then
		Return(expectedData, nil)

	// Test
	var result []model.EnergyConsumption
	err := mockDB.Find(&result, query, args...)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedData, result)
	mockDB.AssertExpectations(t)
}

func TestFind_Error(t *testing.T) {
	mockDB := new(mockDatabase)
	// For
	query := "device_id = (?)"
	args := []any{123}
	expectedErr := errors.New("database error")

	// If
	mockDB.On("Find", mock.Anything, query, args).
		// Then
		Return(nil, expectedErr)

	// Test
	var result []model.EnergyConsumption
	err := mockDB.Find(&result, query, args...)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockDB.AssertExpectations(t)
}

func TestCreateInBatches_Success(t *testing.T) {
	mockDB := new(mockDatabase)
	// For
	data := []model.EnergyConsumption{
		{Id: uuid1, DeviceId: 123, Consumption: 1.0, CreatedAt: time.Now()},
		{Id: uuid2, DeviceId: 123, Consumption: 2.0, CreatedAt: time.Now()},
	}
	batchSize := len(data)

	// If
	mockDB.On("CreateInBatches", data, batchSize).
		// Then
		Return(nil)

	// Test
	err := mockDB.CreateInBatches(data, batchSize)

	// Assert
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestCreateInBatches_Error(t *testing.T) {
	mockDB := new(mockDatabase)
	// For
	data := []model.EnergyConsumption{
		{Id: uuid1, DeviceId: 123, Consumption: 1.0, CreatedAt: time.Now()},
		{Id: uuid2, DeviceId: 123, Consumption: 2.0, CreatedAt: time.Now()},
	}
	batchSize := len(data)
	expectedErr := errors.New("failed to create in batches")

	// If
	mockDB.On("CreateInBatches", data, batchSize).
		// Then
		Return(expectedErr)

	// Test
	err := mockDB.CreateInBatches(data, batchSize)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockDB.AssertExpectations(t)
}
