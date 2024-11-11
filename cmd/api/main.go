package main

import (
	"cast/internal/api/handler"
	"cast/internal/api/sse"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	registerRoutes(r)
	if err := r.Run(":8080"); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}

func registerRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	v1.POST("/init", handler.InitHandler)
	v1.GET("/events", sse.SSEHandler)
	v1.POST("/stop", handler.StopHandler)
	v1.POST("/send", handler.SendMessage)
}
