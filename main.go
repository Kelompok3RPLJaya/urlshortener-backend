package main

import (
	"net/http"
	"os"
	"url-shortener-backend/common"
	"url-shortener-backend/config"
	"url-shortener-backend/controller"
	"url-shortener-backend/middleware"
	"url-shortener-backend/repository"
	"url-shortener-backend/routes"
	"url-shortener-backend/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		res := common.BuildErrorResponse("Gagal Terhubung ke Server", err.Error(), common.EmptyObj{})
		(*gin.Context).JSON((&gin.Context{}), http.StatusBadGateway, res)
		return
	}

	var (
		db *gorm.DB = config.SetupDatabaseConnection()

		jwtService service.JWTService = service.NewJWTService()

		privateRepository      repository.PrivateRepository      = repository.NewPrivateRepository(db)
		feedsRepository        repository.FeedsRepository        = repository.NewFeedsRepository(db)
		urlShortenerRepository repository.UrlShortenerRepository = repository.NewUrlShortenerRepository(db, feedsRepository)
		userRepository         repository.UserRepository         = repository.NewUserRepository(db)

		feedsService        service.FeedsService        = service.NewFeedsService(feedsRepository, urlShortenerRepository, userRepository)
		urlShortenerService service.UrlShortenerService = service.NewUrlShortenerService(urlShortenerRepository, privateRepository, userRepository)
		userService         service.UserService         = service.NewUserService(userRepository)

		feedsController        controller.FeedsController        = controller.NewFeedsController(feedsService)
		urlShortenerController controller.UrlShortenerController = controller.NewUrlShortenerController(urlShortenerService, jwtService)
		userController         controller.UserController         = controller.NewUserController(userService, jwtService)
	)

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	routes.UserRoutes(server, userController, jwtService)
	routes.UrlShortenerRoutes(server, urlShortenerController, jwtService)
	routes.FeedsRoutes(server, feedsController, jwtService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	server.Run("127.0.0.1:" + port)
}
