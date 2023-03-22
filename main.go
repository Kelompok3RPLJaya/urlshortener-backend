package main

import (
	"net/http"
	"os"
	"url-shortener-backend/common"
	"url-shortener-backend/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		res := common.BuildErrorResponse("Gagal Terhubung ke Server", err.Error(), common.EmptyObj{})
		(*gin.Context).JSON((&gin.Context{}), http.StatusBadGateway, res)
		return
	}

	var (
		// db *gorm.DB = config.SetupDatabaseConnection()
		
		// jwtService service.JWTService = service.NewJWTService()
	)

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	// routes.UserRoutes(server, userController, jwtService)
	// routes.UrlShortenerRoutes(server, urlShortenerController, jwtService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	server.Run(":" + port)
}