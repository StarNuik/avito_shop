package main

import (
	"github.com/avito_shop/internal/handler"
	"github.com/gin-gonic/gin"
)

func Run() {
	engine := gin.Default()

	engine.POST("/api/auth", handler.Auth)
	engine.GET("/api/info")
	engine.GET("/api/buy/{item}")
	engine.POST("/api/sendCoin")

	_ = engine.Run()
}

func main() {
	Run()
}
