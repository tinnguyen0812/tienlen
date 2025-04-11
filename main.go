package main

import (
	"github.com/gin-gonic/gin"
	"tienlen-server/router"
)

func main() {
	r := gin.Default()
	router.SetupRoutes(r)
	r.Run(":8080")
}
