package main

import (
	"github.com/gin-gonic/gin"
	"pondDataProcessGo/controllers"
)

func main() {
	// 初始化 Gin 应用
	app := gin.Default()
	twitterControllers := controllers.NewTwitterController()
	processControllers := controllers.NewProcessController()
	app.GET("/hello", twitterControllers.Hello)
	app.POST("/user/new", processControllers.NewUserProcess)
	app.Use(gin.Logger())
	app.Run(":8082")
}
