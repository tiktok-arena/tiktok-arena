package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"tiktok-arena/configuration"
	"tiktok-arena/internal/api/controllers"
	"tiktok-arena/internal/api/middleware"
	"tiktok-arena/internal/api/routers"
	"tiktok-arena/internal/core/services"
	"tiktok-arena/internal/data/database"
	"tiktok-arena/internal/data/repository"
)

// @title			TikTok arena API
// @version		1.0
// @description	API for TikTok arena application
// @host			tiktok-arena.onrender.com
// @BasePath		/api
func Run(c *configuration.EnvConfigModel) {

	db := database.ConnectDB(c)

	userRepository := repository.NewUserRepository(db)
	tiktokRepository := repository.NewTiktokRepository(db)
	tournamentRepository := repository.NewTournamentRepository(db)

	userService := services.NewUserService(userRepository, tournamentRepository)
	authService := services.NewAuthService(userRepository)
	tournamentService := services.NewTournamentService(tournamentRepository, tiktokRepository)

	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)
	tournamentController := controllers.NewTournamentController(tournamentService)

	authRouter := routers.NewAuthRouter(authController)
	userRouter := routers.NewUserRouter(userController)
	tournamentRouter := routers.NewTournamentRouter(tournamentController)

	app := fiber.New(fiber.Config{ErrorHandler: middleware.ErrorHandler})

	//	Logger middleware for logging HTTP request/response details
	app.Use(logger.New())

	//	CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	routers.SetupRoutes(app, tournamentRouter, authRouter, userRouter)

	log.Fatal(app.Listen(":8000"))
}
