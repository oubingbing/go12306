package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"newbug/controller"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	api := router.Group("/api")
	{
		api.GET("/answer",controller.Answer)
	}

	view := router.Group("/view")
	{
		view.GET("/device_info",func(c *gin.Context) {
			c.HTML(http.StatusOK, "device.html", gin.H{
				"title": "Main website",
			})
			fmt.Println("ceshi")
		})
	}


	return router
}
