package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"pondDataProcessGo/services"
	"pondDataProcessGo/utils"
	"strconv"
	"strings"
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
	tidList := strings.Split(tidstr, ",")
	var tidIntList []int64

	for _, tid := range tidList {
		tidInt, err := strconv.ParseInt(tid, 10, 64)
		if err != nil {
			fmt.Println("Error converting", tid, "to int64:", err)
			continue
		}
		tidIntList = append(tidIntList, tidInt)
	}
	user := mc.twitterService.GetUserInfoByID(tidIntList, []string{})
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
