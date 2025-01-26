package repository

import (
	"time"

	"github.com/7yrionLannister/golang-technical-assesment/db"
	"github.com/7yrionLannister/golang-technical-assesment/db/model"
	"github.com/7yrionLannister/golang-technical-assesment/util"
)

func GetEnergyConsumptionsByMeterIdBetweenDates(meterId uint, startDate time.Time, endDate time.Time) ([]model.EnergyConsumption, error) {
	energyConsumptions := make([]model.EnergyConsumption, 0)
	err := db.DB.Find(&energyConsumptions, "device_id = (?) AND created_at BETWEEN ? AND ?", meterId, startDate, endDate)
	if err != nil {
		return nil, util.HandleError(err, "Failed to query energy consumptions")
	}
	return energyConsumptions, nil
}
