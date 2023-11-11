package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"pondDataProcessGo/controllers"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
}

func main() {

	// 初始化 Gin 应用
	app := gin.Default()
	twitterControllers := controllers.NewTwitterController()
	processControllers := controllers.NewProcessController()
	graphControllers := controllers.NewGraphController()

	app.GET("/hello", twitterControllers.Hello)

	app.GET("/twitter/user", twitterControllers.UserInfoById)

	app.POST("/process/user", processControllers.NewUserProcess)

	app.GET("/graph/twitterInfo", graphControllers.GetTwitterAccountInfo)

	zap.L().Info("zap_log")
	app.Run(":8082")
}
