package controller

import (
	"net/http"
	"strconv"
	"url-shortener-backend/common"
	"url-shortener-backend/entity"
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
	var pagination entity.Pagination
	page, _ := strconv.Atoi(ctx.Query("page"))
	if page <= 0 {
		page = 1
	}
	pagination.Page = page
	perPage, _ := strconv.Atoi(ctx.Query("per_page"))
	if perPage <= 0 {
		perPage = 10
	}
	pagination.PerPage = perPage
	
	result, err := fc.feedsService.GetAllFeeds(ctx.Request.Context(), pagination)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mendapatkan List Feeds", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := common.BuildResponse(true, "Berhasil Mendapatkan List Feeds", result)
	ctx.JSON(http.StatusOK, res)
}