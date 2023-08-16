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
	tournamentService := services.NewTournamentService(tournamentRepository, tiktokRepository, userRepository)

	// Create controller layer
	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)
	tournamentController := controllers.NewTournamentController(tournamentService)

	// Create routers for unprotected and protected routes
	authRouter := routers.NewAuthRouter(authController)
	tournamentRouter := routers.NewTournamentRouter(tournamentController)
	userRouter := routers.NewUserRouter(userController)

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
	userRouter(groupRoutes.UserGroup)
	tournamentRouter(groupRoutes.TournamentGroup)

	log.Fatal(app.Listen(":8000"))
}
