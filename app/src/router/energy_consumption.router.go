package router

import (
	"github.com/7yrionLannister/golang-technical-assesment/controller"
	"github.com/gin-gonic/gin"
)

func ConsumptionRouter(app *gin.Engine) {
	consumptionRouter := app.Group("/consumption")
	consumptionRouter.GET("", controller.GetConsumption)
}
