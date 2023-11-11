package controllers

import (
	"github.com/gin-gonic/gin"
	"pondDataProcessGo/models"
	"pondDataProcessGo/services"
	"strconv"
)

type ProcessController struct {
	twitterService *services.TwitterService
	graphService   *services.GraphService
}

func NewProcessController() *ProcessController {
	return &ProcessController{
		twitterService: services.NewTwitterService(),
		graphService:   services.NewGraphService(),
	}
}

func (mc *ProcessController) NewUserProcess(ctx *gin.Context) {
	//新用户接入
	//需要完成:获取该用户的infor，获取该用户follower，followering的infor
	//TwitterAccount，TwitterFollowing，Twitter2Person，Person2Person，Compute relationship strength，Update Algorithm on Read Instances
	twitterId := ctx.PostForm("twitterId")
	twitterIdInt64, _ := strconv.ParseInt(twitterId, 10, 64)
	userInfo := mc.twitterService.GetUserInfoByID([]int64{twitterIdInt64})[0]
	followerIds := mc.twitterService.GetFollowIdsByUserId(twitterId, true)
	followingIds := mc.twitterService.GetFollowIdsByUserId(twitterId, false)
	followersInfo := mc.twitterService.GetUserInfoByID(followerIds)
	followingsInfo := mc.twitterService.GetUserInfoByID(followingIds)

	//add twitterAccount
	mc.graphService.MergeTwitterAccount([]models.User{userInfo})
	mc.graphService.MergeTwitterAccount(followersInfo)
	mc.graphService.MergeTwitterAccount(followingsInfo)

	//merge twitterFollowing
	mc.graphService.MergeTwitterFollowing(userInfo, followersInfo, true)
	mc.graphService.MergeTwitterFollowing(userInfo, followingsInfo, false)

	//add person
	uIds := mc.graphService.MergeTwitter2Person([]models.User{userInfo})
	ferIds := mc.graphService.MergeTwitter2Person(followersInfo)
	fingIds := mc.graphService.MergeTwitter2Person(followingsInfo)
	uIds = append(uIds, ferIds...)
	uIds = append(uIds, fingIds...)

	// person 2 person merge
	mc.graphService.MergePerson2Person(userInfo, true)
	mc.graphService.MergePerson2Person(userInfo, false)

	sources, targets, weights := mc.graphService.UpdateComputeRelationStrength(userInfo)

	//update Algorithm on Read Instances
	mc.graphService.AddNodesFromWithNodes(uIds)
	mc.graphService.AddNodesFromSourcesTargetsWeights(sources, targets, weights)
}
