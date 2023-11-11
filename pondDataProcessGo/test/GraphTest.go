package main

import (
	"pondDataProcessGo/repositories"
)

func main() {
	// 初始化 Gin 应用
	graphDb := repositories.NewGraphRepository()
	//user := models.User{}
	//user.Username = "ThinksDylan"
	//graphDb.UpdateComputeRelationStrength(user)
	graphDb.EasyTest()

}
