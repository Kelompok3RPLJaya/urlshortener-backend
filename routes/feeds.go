package routes

import (
	"url-shortener-backend/controller"
	"url-shortener-backend/service"

	"github.com/gin-gonic/gin"
)

func FeedsRoutes(router *gin.Engine, FeedsController controller.FeedsController, jwtService service.JWTService) {
	FeedsRoutes := router.Group("/api/feeds")
	{
		FeedsRoutes.GET("", FeedsController.GetAllFeeds)
	}
}