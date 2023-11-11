package controllers

import (
	"github.com/gin-gonic/gin"
	"pondDataProcessGo/services"
)

type TwitterController struct {
	twitterService *services.TwitterService
}

func NewTwitterController() *TwitterController {
	return &TwitterController{
		twitterService: services.NewTwitterService(),
	}
}

func (mc *TwitterController) Hello(ctx *gin.Context) {
	//mc.twitterService.GetFollowersByUserId("913499519593172992")
	ctx.JSON(200, "hello world!")
}
