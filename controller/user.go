package controller

import (
	"net/http"
	"url-shortener-backend/common"
	"url-shortener-backend/dto"
	"url-shortener-backend/entity"
	"url-shortener-backend/service"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	RegisterUser(ctx *gin.Context)
	LoginUser(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
}

type userController struct {
	jwtService  service.JWTService
	userService service.UserService
}

func NewUserController(us service.UserService, jwts service.JWTService) UserController {
	return &userController{
		userService: us,
		jwtService:  jwts,
	}
}

func (uc *userController) RegisterUser(ctx *gin.Context) {
	var user dto.UserCreateDto
	err := ctx.ShouldBind(&user)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Menambahkan User", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	checkUser, _ := uc.userService.CheckUser(ctx.Request.Context(), user.Email)
	if checkUser {
		res := common.BuildErrorResponse("User Sudah Terdaftar", "false", common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	result, err := uc.userService.RegisterUser(ctx.Request.Context(), user)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Menambahkan User", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := common.BuildResponse(true, "Berhasil Menambahkan User", result)
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) LoginUser(ctx *gin.Context) {
	var userLoginDTO dto.UserLoginDTO
	err := ctx.ShouldBind(&userLoginDTO)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Login", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res, _ := uc.userService.Verify(ctx.Request.Context(), userLoginDTO.Email, userLoginDTO.Password)
	if !res {
		response := common.BuildErrorResponse("Gagal Login", "Email atau Password Salah", common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	user, err := uc.userService.FindUserByEmail(ctx.Request.Context(), userLoginDTO.Email)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Login", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	token := uc.jwtService.GenerateToken(user.ID, user.Role)
	userResponse := entity.Authorization{
		Token: token,
		Role:  user.Role,
	}

	response := common.BuildResponse(true, "Berhasil Login", userResponse)
	ctx.JSON(http.StatusOK, response)
}

func (uc *userController) UpdateUser(ctx *gin.Context) {
	var userUpdateDto dto.UserUpdateDto
	err := ctx.ShouldBind(&userUpdateDto)
	if err != nil {
		res := common.BuildErrorResponse("Gagal proses request update", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	token := ctx.MustGet("token").(string)
	userID, err := uc.jwtService.GetUserIDByToken(token)
	if err != nil {
		res := common.BuildErrorResponse("Gagal proses request token", "Token yang dimasukkan salah", common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	userUpdateDto.ID = userID
	err = uc.userService.UpdateUser(ctx, userUpdateDto)
	if err != nil {
		res := common.BuildErrorResponse("Gagal update user", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := common.BuildResponse(true, "Berhasil update user", common.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}
