package db

import (
	"github.com/7yrionLannister/golang-technical-assesment/db/model"
	"github.com/stretchr/testify/mock"
)

// [mockDatabase] is a mock implementation of the [Database] interface
type mockDatabase struct {
	mock.Mock
}

// naive mock Find implementation that solely retrieves the expected output slice
func (m *mockDatabase) Find(out any, query string, args ...any) error {
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

// naive mock CreateInBatches implementation that solely returns the expected error
func (m *mockDatabase) CreateInBatches(value any, batchSize int) error {
	argsCall := m.Called(value, batchSize)
	return argsCall.Error(0)
}
