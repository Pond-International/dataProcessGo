package controllers

import (
	"github.com/gin-gonic/gin"
	"pondDataProcessGo/services"
	"strconv"
)

type TwitterController struct {
	twitterService *services.TwitterService
}

func NewTwitterController() *TwitterController {
	return &TwitterController{
		twitterService: services.NewTwitterService(),
	}
}

func (mc *TwitterController) UserInfoById(ctx *gin.Context) {
	tidstr := ctx.Query("tid")
	tid, _ := strconv.ParseInt(tidstr, 10, 64)
	user := mc.twitterService.GetUserInfoByID([]int64{tid})
	if user != nil {
		ctx.JSON(200, user)
	} else {
		ctx.JSON(404, gin.H{"error": "User not found"})
	}
}

func (mc *TwitterController) Hello(ctx *gin.Context) {

	//mc.twitterService.GetFollowersByUserId("913499519593172992")
	ctx.JSON(200, "hello world!")
}
