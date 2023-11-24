package controllers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	zap.L().Info("NewUserProcess/startNewUserProcess", zap.String("twitterId", twitterId))
	twitterIdInt64, _ := strconv.ParseInt(twitterId, 10, 64)
	userInfo := mc.twitterService.GetUserInfoByID([]int64{twitterIdInt64}, []string{})[0]
	zap.L().Info("NewUserProcess/getUserInfo", zap.String("username", userInfo.Username))
	followerIds := mc.twitterService.GetFollowIdsByUserId(twitterId, true)
	followingIds := mc.twitterService.GetFollowIdsByUserId(twitterId, false)
	followersInfo := mc.twitterService.GetUserInfoByID(followerIds, []string{})
	followingsInfo := mc.twitterService.GetUserInfoByID(followingIds, []string{})
	zap.L().Info("NewUserProcess/getFollowInfo", zap.String("username", userInfo.Username), zap.Int("followerCount", len(followerIds)), zap.Int("followingCount", len(followingIds)))

	//add twitterAccount
	mc.graphService.MergeTwitterAccount([]models.User{userInfo})
	mc.graphService.MergeTwitterAccount(followersInfo)
	mc.graphService.MergeTwitterAccount(followingsInfo)
	zap.L().Info("NewUserProcess/MergeTwitterAccount", zap.String("username", userInfo.Username))
	//merge twitterFollowing
	mc.graphService.MergeTwitterFollowing(userInfo, followersInfo, true)
	mc.graphService.MergeTwitterFollowing(userInfo, followingsInfo, false)
	zap.L().Info("NewUserProcess/MergeTwitterFollowing", zap.String("username", userInfo.Username))
	//add person
	uIds := mc.graphService.MergeTwitter2Person([]models.User{userInfo})
	ferIds := mc.graphService.MergeTwitter2Person(followersInfo)
	fingIds := mc.graphService.MergeTwitter2Person(followingsInfo)

	zap.L().Info("NewUserProcess/MergeTwitter2Person", zap.String("username", userInfo.Username))

	uIds = append(uIds, ferIds...)
	uIds = append(uIds, fingIds...)

	// person 2 person merge
	mc.graphService.MergePerson2Person(userInfo, true)
	mc.graphService.MergePerson2Person(userInfo, false)
	zap.L().Info("NewUserProcess/MergePerson2Person", zap.String("username", userInfo.Username))

	sources, targets, weights := mc.graphService.UpdateComputeRelationStrength(userInfo)
	zap.L().Info("NewUserProcess/UpdateComputeRelationStrength", zap.String("username", userInfo.Username))

	//update Algorithm on Read Instances
	mc.graphService.AddNodesFromWithNodes(uIds)
	mc.graphService.AddNodesFromSourcesTargetsWeights(sources, targets, weights)
	zap.L().Info("NewUserProcess/AddNodes", zap.String("username", userInfo.Username))

}
