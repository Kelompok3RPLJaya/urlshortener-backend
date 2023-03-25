package controller

import (
	"net/http"
	"url-shortener-backend/common"
	"url-shortener-backend/dto"
	"url-shortener-backend/service"

	"github.com/gin-gonic/gin"
)

type UrlShortenerController interface {
	CreateUrlShortener(ctx *gin.Context)
	GetMeUrlShortener(ctx *gin.Context)
	GetAllUrlShortener(ctx *gin.Context)
}

type urlShortenerController struct {
	urlShortenerService service.UrlShortenerService
	jwtService service.JWTService
}

func NewUrlShortenerController(us service.UrlShortenerService, js service.JWTService) UrlShortenerController {
	return &urlShortenerController{
		urlShortenerService: us,
		jwtService: js,
	}
}

func(uc *urlShortenerController) CreateUrlShortener(ctx *gin.Context) {
	var urlShortener dto.UrlShortenerCreateDTO
	err := ctx.ShouldBind(&urlShortener)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Menambahkan Url Shortener", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if ctx.Request.Header["Authorization"] != nil {
		token := ctx.MustGet("token").(string)
		userID, err := uc.jwtService.GetUserIDByToken(token)
		if err != nil {
			response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		urlShortener.UserID = userID
	}
	
	checkUrlShortener, _ := uc.urlShortenerService.ValidateShortUrl(ctx.Request.Context(), urlShortener.ShortUrl)
	if checkUrlShortener.ShortUrl != "" {
		res := common.BuildErrorResponse("Gagal Menambahkan Url Shortener", "Short Url Sudah Terdaftar", common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if *urlShortener.IsPrivate && urlShortener.Password == "" {
		res := common.BuildErrorResponse("Gagal Menambahkan Url Shortener", "Url Shortener Private Harus Mengandung Password", common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := uc.urlShortenerService.CreateUrlShortener(ctx.Request.Context(), urlShortener)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Menambahkan Url Shortener", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := common.BuildResponse(true, "Berhasil Menambahkan Url Shortener", result)
	ctx.JSON(http.StatusOK, res)
}

func(uc *urlShortenerController) GetMeUrlShortener(ctx *gin.Context) {
	token := ctx.MustGet("token").(string)
	userID, err := uc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}
	result, err := uc.urlShortenerService.GetUrlShortenerByUserID(ctx.Request.Context(), userID.String())
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mendapatkan Url Shortener User", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := common.BuildResponse(true, "Berhasil Mendapatkan Url Shortener User", result)
	ctx.JSON(http.StatusOK, res)
}

func(uc *urlShortenerController) GetAllUrlShortener(ctx *gin.Context) {
	result, err := uc.urlShortenerService.GetAllUrlShortener(ctx.Request.Context())
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mendapatkan List Url Shortener", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := common.BuildResponse(true, "Berhasil Mendapatkan List Url Shortener", result)
	ctx.JSON(http.StatusOK, res)
}