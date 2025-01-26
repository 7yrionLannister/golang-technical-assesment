package middleware

import (
	"github.com/7yrionLannister/golang-technical-assesment/docs"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(app *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/"
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
