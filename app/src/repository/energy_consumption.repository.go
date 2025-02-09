package repository

import (
	"time"

	"github.com/7yrionLannister/golang-technical-assesment/config/logger"
	"github.com/7yrionLannister/golang-technical-assesment/db"
	"github.com/7yrionLannister/golang-technical-assesment/db/model"
	"github.com/7yrionLannister/golang-technical-assesment/db/view"
	"github.com/7yrionLannister/golang-technical-assesment/util"
)

// Queries the database for energy consumptions by meter ID between two dates.
func GetEnergyConsumptionsByMeterIdBetweenDates(metersIds []uint, startDate time.Time, endDate time.Time) ([]view.EnergyConsumptionDTO, error) {
	gormDb := db.DB.(*db.GormDatabase).GormDb
	var result []view.EnergyConsumptionDTO
	err := gormDb.
		Model(&model.EnergyConsumption{}).
		Select("device_id as meter_id, sum(consumption) as total_consumption").
		Where("device_id IN ? AND created_at BETWEEN ? AND ?", metersIds, startDate, endDate).
		Group("device_id").
		Scan(&result).Error
	if err != nil {
		return nil, util.HandleError(err, "Failed to query energy consumptions")
	}
	logger.L.Debug("query result", "dto", result)
	return result, nil
}
