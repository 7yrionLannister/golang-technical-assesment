package controller

import (
	"log/slog"
	"time"

	"github.com/7yrionLannister/golang-technical-assesment/config/logger"
	"github.com/7yrionLannister/golang-technical-assesment/service"
	"github.com/7yrionLannister/golang-technical-assesment/util"
	"github.com/gin-gonic/gin"
)

// struct used to parse the query parameters.
type consumptionQueryParams struct {
	MeterIdsString string `form:"meters_ids" binding:"required"` // Do not use this field, use [MeterIds] instead
	MeterIds       []uint
	StartDate      time.Time `form:"start_date" binding:"required" time_format:"2006-01-02"`
	EndDate        time.Time `form:"end_date" binding:"required" time_format:"2006-01-02"`
	KindPeriod     string    `form:"kind_period" binding:"required"`
}

// GetConsumption godoc
// @Summary      Get consumption report
// @Description  Returns the energy consumption for each meter for the period between startDate and endDate.
// @Tags         consumption
// @Produce      json
// @Param        meters_ids  query  []uint  true  "Meter ids"
// @Param        start_date  query  string  true  "Start date (format: YYYY-MM-DD)"
// @Param        end_date    query  string  true  "End date (format: YYYY-MM-DD)"
// @Param        kind_period query  string  true  "Kind period (any of [monthly, weekly, daily])"
// @Success      200  {object}  dto.PeriodicConsumptionDTO  "Consumption report"
// @Failure      400  {object}  map[string]string  "Invalid query params"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Router       /consumption [get]
func GetConsumption(c *gin.Context) {
	logger.L.Debug("Querying consumption data")
	var params consumptionQueryParams
	err := c.BindQuery(&params)
	if err != nil {
		logger.L.Error("Error binding query params", slog.Any("error", err))
		c.JSON(400, gin.H{"error": "Invalid query params", "cause": err})
		return
	}
	metersIds, err := util.String2UintSlice(params.MeterIdsString)
	if err != nil {
		logger.L.Error("Error converting meter ids to slice", slog.Any("error", err))
		c.JSON(400, gin.H{"error": "Failed to convert meter ids to slice", "cause": err})
		return
	}
	params.MeterIds = metersIds
	logger.L.Debug("Query params", slog.Any("params", params))
	// Query data from the service
	periodDto, err := service.GetEnergyConsumptions(params.MeterIds, params.StartDate, params.EndDate, params.KindPeriod)
	if err != nil {
		logger.L.Error("Error querying consumption data", slog.Any("error", err))
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	logger.L.Info("Consumption data retrieved successfully", slog.Any("response", periodDto))
	c.JSON(200, periodDto)
}
