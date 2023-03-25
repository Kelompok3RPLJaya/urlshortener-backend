package routes

import (
	"url-shortener-backend/controller"
	"url-shortener-backend/middleware"
	"url-shortener-backend/service"

	"github.com/gin-gonic/gin"
)

func UrlShortenerRoutes(router *gin.Engine, UrlShortenerController controller.UrlShortenerController, jwtService service.JWTService) {
	urlShortenerRoutes := router.Group("/api/url_shortener")
	{
		urlShortenerRoutes.POST("", middleware.CreateShortUrlAuthenticate(jwtService, false), UrlShortenerController.CreateUrlShortener)
	}
}