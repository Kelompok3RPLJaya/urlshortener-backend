package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"url-shortener-backend/common"
	"url-shortener-backend/dto"
	"url-shortener-backend/entity"
	"url-shortener-backend/service"

	"github.com/gin-gonic/gin"
)

type UrlShortenerController interface {
	CreateUrlShortener(ctx *gin.Context)
	GetMeUrlShortener(ctx *gin.Context)
	GetAllUrlShortener(ctx *gin.Context)
	UpdateUrlShortener(ctx *gin.Context)
	UpdatePrivate(ctx *gin.Context)
	DeleteUrlShortener(ctx *gin.Context)
	GetUrlShortenerByShortUrl(ctx *gin.Context)
}

type urlShortenerController struct {
	urlShortenerService service.UrlShortenerService
	jwtService          service.JWTService
}

func NewUrlShortenerController(us service.UrlShortenerService, js service.JWTService) UrlShortenerController {
	return &urlShortenerController{
		urlShortenerService: us,
		jwtService:          js,
	}
}

func (uc *urlShortenerController) CreateUrlShortener(ctx *gin.Context) {
	var urlShortener dto.UrlShortenerCreateDTO
	err := ctx.ShouldBind(&urlShortener)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Menambahkan Url Shortener", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	_, err = url.ParseRequestURI(urlShortener.LongUrl)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Menambahkan Url Shortener", "Url Tidak Valid", common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(urlShortener.ShortUrl) {
		res := common.BuildErrorResponse("Gagal Menambahkan Url Shortener", "Short Url Tidak Valid", common.EmptyObj{})
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

func (uc *urlShortenerController) GetMeUrlShortener(ctx *gin.Context) {
	search := ctx.Query("search")
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

	token := ctx.MustGet("token").(string)
	userID, err := uc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}
	result, err := uc.urlShortenerService.GetUrlShortenerByUserID(ctx.Request.Context(), userID.String(), search, pagination)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mendapatkan Url Shortener User", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := common.BuildResponse(true, "Berhasil Mendapatkan Url Shortener User", result)
	ctx.JSON(http.StatusOK, res)
}

func (uc *urlShortenerController) GetAllUrlShortener(ctx *gin.Context) {
	result, err := uc.urlShortenerService.GetAllUrlShortener(ctx.Request.Context())
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mendapatkan List Url Shortener", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	fmt.Println(result[0].ID)
	fmt.Println(result[0].CreatedAt)
	res := common.BuildResponse(true, "Berhasil Mendapatkan List Url Shortener", result)
	ctx.JSON(http.StatusOK, res)
}

func (uc *urlShortenerController) UpdateUrlShortener(ctx *gin.Context) {
	urlShortenerID := ctx.Param("id")
	var urlShortenerDTO dto.UrlShortenerUpdateDTO
	err := ctx.ShouldBind(&urlShortenerDTO)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mengupdate Url Shortener", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	if urlShortenerDTO.LongUrl != "" {
		_, err = url.ParseRequestURI(urlShortenerDTO.LongUrl)
		if err != nil {
			res := common.BuildErrorResponse("Gagal Menambahkan Url Shortener", "Url Tidak Valid", common.EmptyObj{})
			ctx.JSON(http.StatusBadRequest, res)
			return
		}
	}

	if urlShortenerDTO.ShortUrl != "" {
		if !regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(urlShortenerDTO.ShortUrl) {
			res := common.BuildErrorResponse("Gagal Menambahkan Url Shortener", "Short Url Tidak Valid", common.EmptyObj{})
			ctx.JSON(http.StatusBadRequest, res)
			return
		}
	}

	token := ctx.MustGet("token").(string)
	userID, err := uc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	checkUrlShortenerUser := uc.urlShortenerService.ValidateUrlShortenerUser(ctx, userID.String(), urlShortenerID)
	if !checkUrlShortenerUser {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Akun Anda Tidak Memiliki Akses Untuk Mengupdate Url Shortener Ini", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	checkDuplicateUrlShortener, _ := uc.urlShortenerService.ValidateShortUrl(ctx.Request.Context(), urlShortenerDTO.ShortUrl)
	if checkDuplicateUrlShortener.ShortUrl != "" {
		res := common.BuildErrorResponse("Gagal Menambahkan Url Shortener", "Short Url Sudah Terdaftar", common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	err = uc.urlShortenerService.UpdateUrlShortener(ctx, urlShortenerDTO, urlShortenerID)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mengupdate Url Shortener", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := common.BuildResponse(true, "Berhasil Mengupdate Url Shortener", common.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

func (uc *urlShortenerController) UpdatePrivate(ctx *gin.Context) {
	urlShortenerID := ctx.Param("id")

	var privateDTO dto.PrivateUpdateDTO
	err := ctx.ShouldBind(&privateDTO)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mengupdate Url Shortener", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	token := ctx.MustGet("token").(string)
	userID, err := uc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	checkUrlShortenerUser := uc.urlShortenerService.ValidateUrlShortenerUser(ctx, userID.String(), urlShortenerID)
	if !checkUrlShortenerUser {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Akun Anda Tidak Memiliki Akses Untuk Menghapus Url Shortener Ini", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	urlShortener, err := uc.urlShortenerService.GetUrlShortenerByID(ctx, urlShortenerID)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mengupdate Url Shortener", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if !*urlShortener.IsPrivate {
		if privateDTO.Password == "" {
			res := common.BuildErrorResponse("Gagal Mengupdate Url Shortener", "Url Shortener Private Harus Mengandung Password", common.EmptyObj{})
			ctx.JSON(http.StatusBadRequest, res)
			return
		} else {
			err = uc.urlShortenerService.UpdatePrivate(ctx, urlShortenerID, privateDTO)
		}
	} else {
		err = uc.urlShortenerService.UpdatePublic(ctx, urlShortenerID)
	}
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mengupdate Url Shortener", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := common.BuildResponse(true, "Berhasil Mengupdate Url Shortener", common.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

func (uc *urlShortenerController) DeleteUrlShortener(ctx *gin.Context) {
	urlShortenerID := ctx.Param("id")
	token := ctx.MustGet("token").(string)
	userID, err := uc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	checkURL := uc.urlShortenerService.ValidateUrlShortenerUser(ctx, userID.String(), urlShortenerID)
	if !checkURL {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Akun Anda Tidak Memiliki Akses Untuk Mengupdate Url Shortener Ini", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	err = uc.urlShortenerService.DeleteUrlShortener(ctx, urlShortenerID)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Menghapus Url Shortener", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := common.BuildResponse(true, "Berhasil Menghapus Url Shortener", common.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

func(uc *urlShortenerController) GetUrlShortenerByShortUrl(ctx *gin.Context) {
	shortUrl := ctx.Param("short_url")
	var private dto.PrivateUpdateDTO
	err := ctx.ShouldBind(&private)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mendapatkan Url", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	result, err := uc.urlShortenerService.GetUrlShortenerByShortUrl(ctx.Request.Context(), shortUrl, private)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mendapatkan Url", err.Error(), result)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := common.BuildResponse(true, "Berhasil Mendapatkan Url", result)
	ctx.JSON(http.StatusOK, res)
}