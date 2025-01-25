package controller

import (
	"log/slog"
	"time"

	"github.com/7yrionLannister/golang-technical-assesment/config/logger"
	"github.com/7yrionLannister/golang-technical-assesment/service"
	"github.com/gin-gonic/gin"
)

type consumptionQueryParams struct {
	MeterIds   []uint    `form:"meters_ids" binding:"required"`
	StartDate  time.Time `form:"start_date" binding:"required" time_format:"2006-01-02"`
	EndDate    time.Time `form:"end_date" binding:"required" time_format:"2006-01-02"`
	KindPeriod string    `form:"kind_period" binding:"required"`
}

func GetConsumption(c *gin.Context) {
	logger.Debug("Querying consumption data")
	var params consumptionQueryParams
	err := c.BindQuery(&params)
	if err != nil {
		logger.Error("Error binding query params", slog.Any("error", err))
		c.JSON(400, gin.H{"error": "Invalid query params"})
		return
	}
	logger.Debug("Query params", slog.Any("params", params))
	// Query data
	periodDto, err := service.GetEnergyConsumptions(params.MeterIds, params.StartDate, params.EndDate, params.KindPeriod)
	if err != nil {
		logger.Error("Error querying consumption data", slog.Any("error", err))
		c.JSON(500, gin.H{"error": "Internal server error"})
	}
	c.JSON(200, periodDto)
}
