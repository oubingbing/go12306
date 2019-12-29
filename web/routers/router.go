package routers

import (
	"github.com/gin-gonic/gin"
	"newbug/controller"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	api := router.Group("/api")
	{
		api.GET("/test",controller.SavedeviceId)
	}

	return router
}
