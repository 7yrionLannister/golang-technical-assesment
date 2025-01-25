package service

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/7yrionLannister/golang-technical-assesment/config/logger"
	"github.com/7yrionLannister/golang-technical-assesment/controller/dto"
	"github.com/7yrionLannister/golang-technical-assesment/db"
	"github.com/go-faker/faker/v4"
)

const (
	daily   = "daily"
	weekly  = "weekly"
	monthly = "monthly"
)

func GetEnergyConsumptions(metersIds []uint, startDate time.Time, endDate time.Time, kindPeriod string) (*dto.PeriodicConsumptionDTO, error) {
	// monthly
	var periodDto = &dto.PeriodicConsumptionDTO{
		Period:    make([]string, 0),
		DataGraph: make([]*dto.EnergyConsumptionDTO, 0),
	}
	energyConsumptionDTOForMeter := make(map[uint]*dto.EnergyConsumptionDTO)
	for startDate.Before(endDate) {
		periodString := startDate.Format("January 2006") // TODO format as "JAN 2006"
		periodDto.Period = append(periodDto.Period, periodString)
		periodEndDate := startDate.AddDate(0, 1, 0)
		if periodEndDate.After(endDate) {
			periodEndDate = endDate
		}
		for _, meterId := range metersIds {
			energyConsumptions, err := db.GetEnergyConsumptionsByMeterIdBetweenDates(meterId, startDate, periodEndDate)
			if err != nil {
				msg := "Failed to fetch energy consumptions"
				e := fmt.Errorf("%s: %w", strings.ToLower(msg), err)
				logger.Error(msg, slog.Any("error", err))
				return nil, e
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
					Address: faker.GetRealAddress().Address, // TODO replace with address microservice
				}
				energyConsumptionDTOForMeter[meterId] = energyConsumptionDTO
				periodDto.DataGraph = append(periodDto.DataGraph, energyConsumptionDTO)
			}
			energyConsumptionDTO.Active = append(energyConsumptionDTO.Active, consumption)
		}
		startDate = periodEndDate
	}

	return periodDto, nil
}
