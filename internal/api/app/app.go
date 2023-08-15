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

//	@title			TikTok arena API
//	@version		1.0
//	@description	API for TikTok arena application
//	@host			tiktok-arena.onrender.com
//	@BasePath		/api
func Run(c *configuration.EnvConfigModel) {
	// Create connection to DB
	db := database.ConnectDB(c)

	// Create repositories to access DB
	userRepository := repository.NewUserRepository(db)
	tiktokRepository := repository.NewTiktokRepository(db)
	tournamentRepository := repository.NewTournamentRepository(db)

	// Create service layer
	userService := services.NewUserService(userRepository, tournamentRepository)
	authService := services.NewAuthService(userRepository)
	tournamentService := services.NewTournamentService(tournamentRepository, tiktokRepository)

	// Create controller layer
	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)
	tournamentController := controllers.NewTournamentController(tournamentService)

	// Create routers for unprotected and protected routes
	authRouter := routers.NewAuthRouter(authController)
	tournamentRouter := routers.NewTournamentRouter(tournamentController)

	authProtectedRouter := routers.NewAuthProtectedRouter(authController)
	userProtectedRouter := routers.NewUserProtectedRouter(userController)
	tournamentProtectedRouter := routers.NewTournamentProtectedRouter(tournamentController)

	// ErrorHandler middleware
	app := fiber.New(fiber.Config{ErrorHandler: middleware.ErrorHandler})

	//	Logger middleware for logging HTTP request/response details
	app.Use(logger.New())

	//	CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	// Get group routes
	groupRoutes := routers.GetGroupRoutes(app)

	// Setup unprotected routes
	authRouter(groupRoutes.AuthGroup)
	tournamentRouter(groupRoutes.TournamentGroup)

	// Setup JWT middleware
	app.Use(middleware.Protected())

	// Setup protected routes
	authProtectedRouter(groupRoutes.AuthGroup)
	userProtectedRouter(groupRoutes.UserGroup)
	tournamentProtectedRouter(groupRoutes.TournamentGroup)

	log.Fatal(app.Listen(":8000"))
}
