package controller

import (
	"net/http"
	"url-shortener-backend/common"
	"url-shortener-backend/service"

	"github.com/gin-gonic/gin"
)

type FeedsController interface {
	GetAllFeeds(ctx *gin.Context)
}

type feedsController struct {
	feedsService service.FeedsService
}

func NewFeedsController(fs service.FeedsService) FeedsController {
	return &feedsController{
		feedsService: fs,
	}
}

func(fc *feedsController) GetAllFeeds(ctx *gin.Context) {
	result, err := fc.feedsService.GetAllFeeds(ctx.Request.Context())
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mendapatkan List Feeds", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := common.BuildResponse(true, "Berhasil Mendapatkan List Feeds", result)
	ctx.JSON(http.StatusOK, res)
}