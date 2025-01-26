package router

import "github.com/gin-gonic/gin"

func Setup(app *gin.Engine) {
	ConsumptionRouter(app)
}
