package service

import (
	"time"

	"github.com/7yrionLannister/golang-technical-assesment/controller/dto"
	"github.com/7yrionLannister/golang-technical-assesment/repository"
	"github.com/7yrionLannister/golang-technical-assesment/util"
	"github.com/go-faker/faker/v4"
)

const (
	daily   = "daily"
	weekly  = "weekly"
	monthly = "monthly"
)

// Returns the energy consumption for each meter in the metersIds slice for the period between startDate and endDate.
func GetEnergyConsumptions(metersIds []uint, startDate time.Time, endDate time.Time, kindPeriod string) (*dto.PeriodicConsumptionDTO, error) {
	var periodDto = &dto.PeriodicConsumptionDTO{
		Period:    make([]string, 0),
		DataGraph: make([]*dto.EnergyConsumptionDTO, 0),
	}
	err := stepThroughPeriod(periodDto, metersIds, startDate, endDate, kindPeriod)
	if err != nil {
		return nil, err
	}
	return periodDto, nil
}

// Iterates through the period between startDate and endDate, incrementing the date by the kindPeriod to find the sub-periods.
// For each sub-period, it computes the energy consumption for each meter in the metersIds slice.
func stepThroughPeriod(periodDto *dto.PeriodicConsumptionDTO, metersIds []uint, startDate time.Time, endDate time.Time, kindPeriod string) error {
	// map to keep track of the *dto.EnergyConsumptionDTO that is part of the DataGraph
	energyConsumptionDTOForMeter := make(map[uint]*dto.EnergyConsumptionDTO)
	for startDate.Before(endDate) {
		periodEndDate, periodString := stepDateAndGetPeriodString(kindPeriod, startDate)
		periodDto.Period = append(periodDto.Period, periodString)
		if periodEndDate.After(endDate) {
			periodEndDate = endDate
		}
		err := populateDataGraphForPeriod(periodDto, metersIds, energyConsumptionDTOForMeter, startDate, periodEndDate)
		if err != nil {
			return err
		}
		startDate = periodEndDate
	}
	return nil
}

// Increments the date by the kindPeriod.
// Gets the period string for the kindPeriod.
func stepDateAndGetPeriodString(kindPeriod string, initialDate time.Time) (newDate time.Time, periodString string) {
	switch kindPeriod {
	case daily:
		return initialDate.AddDate(0, 0, 1), initialDate.Format("January 2") // TODO format as "JAN 2"
	case weekly:
		periodString = initialDate.Format("January 2") + " - " + initialDate.AddDate(0, 0, 6).Format("January 2") // TODO format as "JAN 2 - JAN 2"
		return initialDate.AddDate(0, 0, 7), periodString
	case monthly:
		return initialDate.AddDate(0, 1, 0), initialDate.Format("January 2006") // TODO format as "JAN 2006"
	default:
		return initialDate, "TODO"
	}
}

// Computes the energy consumption for each meter in the metersIds slice for the period between periodStartDate and periodEndDate
func populateDataGraphForPeriod(periodDto *dto.PeriodicConsumptionDTO, metersIds []uint, energyConsumptionDTOForMeter map[uint]*dto.EnergyConsumptionDTO, periodStartDate time.Time, periodEndDate time.Time) error {
	for _, meterId := range metersIds {
		energyConsumptions, err := repository.GetEnergyConsumptionsByMeterIdBetweenDates(meterId, periodStartDate, periodEndDate)
		if err != nil {
			return util.HandleError(err, "Failed to fetch energy consumptions")
		}
		consumption := 0.0
		for _, energyConsumption := range energyConsumptions {
			consumption += energyConsumption.Consumption
		}
		energyConsumptionDTO, present := energyConsumptionDTOForMeter[meterId]
		if !present {
			energyConsumptionDTO = &dto.EnergyConsumptionDTO{
				MeterId: meterId,
				Active:  make([]float64, 0),
				Address: faker.GetRealAddress().Address, // Asume faker to be an http client that gets the address for the meter
			}
			energyConsumptionDTOForMeter[meterId] = energyConsumptionDTO
			periodDto.DataGraph = append(periodDto.DataGraph, energyConsumptionDTO)
		}
		energyConsumptionDTO.Active = append(energyConsumptionDTO.Active, consumption)
	}
	return nil
}
