// Package test contains mocks for testing purposes.
// Do not use these mocks in production code.
package test

import (
	"github.com/7yrionLannister/golang-technical-assesment/db/model"
	"github.com/stretchr/testify/mock"
)

// [MockDatabase] is a mock implementation of the [Database] interface
type MockDatabase struct {
	mock.Mock
}

// naive mock [db.Database.InitDatabaseConnection] implementation.
// It does not perform any database connection, but simply returns nil error.
func (m *MockDatabase) InitDatabaseConnection() error {
	return nil
}

// naive mock [db.Database.Find] implementation that solely retrieves the expected output slice
// or returns the expected error.
func (m *MockDatabase) Find(out any, query string, args ...any) error {
	argsCall := m.Called(out, query, args)
	if argsCall.Get(1) == nil {
		switch result := out.(type) {
		case *[]model.EnergyConsumption:
			*result = argsCall.Get(0).([]model.EnergyConsumption)
		default:
			panic("unsupported type for Find output")
		}
	}
	return argsCall.Error(1)
}

// naive mock [db.Database.CreateInBatches] implementation that solely returns the expected error
func (m *MockDatabase) CreateInBatches(value any, batchSize int) error {
	argsCall := m.Called(value, batchSize)
	return argsCall.Error(0)
}

var MockDB *MockDatabase
