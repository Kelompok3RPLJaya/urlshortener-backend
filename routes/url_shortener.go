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
		urlShortenerRoutes.GET("/me", middleware.Authenticate(jwtService, false), UrlShortenerController.GetMeUrlShortener)
		urlShortenerRoutes.GET("", UrlShortenerController.GetAllUrlShortener)
		urlShortenerRoutes.PUT("/:id", middleware.Authenticate(jwtService, false), UrlShortenerController.UpdateUrlShortener)
		urlShortenerRoutes.PUT("/private/:id", middleware.Authenticate(jwtService, false), UrlShortenerController.UpdatePrivate)
		urlShortenerRoutes.DELETE("/:id", middleware.Authenticate(jwtService, false), UrlShortenerController.DeleteUrlShortener)
	}
}
