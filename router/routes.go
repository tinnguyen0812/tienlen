package router

import (
	"github.com/gin-gonic/gin"
	"tienlen-server/controller"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/ws", controller.ServeWs)
}
