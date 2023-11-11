package controllers

import (
	"github.com/gin-gonic/gin"
	"pondDataProcessGo/services"
)

type GraphController struct {
	graphService *services.GraphService
}

func NewGraphController() *GraphController {
	return &GraphController{
		graphService: services.NewGraphService(),
	}
}

func (mc *GraphController) GetTwitterAccountInfo(ctx *gin.Context) {
	tname := ctx.Query("tname")
	ctx.JSON(200, mc.graphService.GetTwitterAccountInfo(tname))
}
