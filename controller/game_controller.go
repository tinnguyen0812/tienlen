package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tienlen-server/websocket"
)

func ServeWs(c *gin.Context) {
	websocket.HandleConnections(c.Writer, c.Request)
}
