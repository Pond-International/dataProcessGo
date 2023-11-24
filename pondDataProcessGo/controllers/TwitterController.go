package controllers

import (
	"github.com/gin-gonic/gin"
	"pondDataProcessGo/services"
	"pondDataProcessGo/utils"
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
	user := mc.twitterService.GetUserInfoByID([]int64{tid}, []string{})
	if user != nil {
		ctx.JSON(200, user)
	} else {
		ctx.JSON(404, gin.H{"error": "User not found"})
	}
}

func (mc *TwitterController) UserPicByIds(ctx *gin.Context) {
	ids := ctx.Query("ids")

	idsInt64, err := utils.StringSliceToIntSlice(ids)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "input format wrong"})
	}
	users := mc.twitterService.GetUserInfoByID(idsInt64, []string{"profile_image_url"})
	if users != nil {
		ctx.JSON(200, users)
	}
}

func (mc *TwitterController) Hello(ctx *gin.Context) {

	//mc.twitterService.GetFollowersByUserId("913499519593172992")
	ctx.JSON(200, "hello world!")
}
