package routes

import (
	"url-shortener-backend/controller"
	"url-shortener-backend/middleware"
	"url-shortener-backend/service"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, UserController controller.UserController, jwtService service.JWTService) {
	userRoutes := router.Group("/api/user")
	{
		userRoutes.POST("", UserController.RegisterUser)
		userRoutes.POST("/login", UserController.LoginUser)
		userRoutes.POST("/edit", middleware.Authenticate(jwtService, false), UserController.UpdateUser)
		userRoutes.GET("/me", middleware.Authenticate(jwtService, false), UserController.DetailUserAccount)
	}
}
