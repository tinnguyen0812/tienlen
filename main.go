package main

import (
	"log"
	"net/http"
	"tienlen-server/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ws", handlers.WebSocketHandler)

	log.Println("Server listening on :8080")
	http.ListenAndServe(":8080", r)
}
