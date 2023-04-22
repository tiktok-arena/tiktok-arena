package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "tiktok-arena/docs"
	"tiktok-arena/internal/api/controllers"
	"tiktok-arena/internal/api/middleware"
)

func SetupRoutes(app *fiber.App,
	tournamentController controllers.TournamentController,
	authController controllers.AuthController,
	userController controllers.UserController) {
	api := app.Group("/api")

	//	Use 'swag init' to generate new /docs files, details: https://github.com/gofiber/swagger#usage
	api.Get("/docs/*", swagger.HandlerDefault)

	api.Route("/auth", func(router fiber.Router) {
		router.Post("/register", authController.RegisterUser)
		router.Post("/login", authController.LoginUser)
		router.Get("/whoami", middleware.Protected(), authController.WhoAmI)
	})

	api.Route("/tournament", func(router fiber.Router) {
		router.Get("", tournamentController.GetAllTournaments)
		router.Get("/contest/:tournamentId", tournamentController.GetTournamentContest)
		router.Post("/create", middleware.Protected(), tournamentController.CreateTournament)
		router.Put("/edit/:tournamentId", middleware.Protected(), tournamentController.EditTournament)
		router.Delete("/delete/:tournamentId", middleware.Protected(), tournamentController.DeleteTournament)
		router.Delete("/delete", middleware.Protected(), tournamentController.DeleteTournaments)
		router.Get("/tiktoks/:tournamentId", tournamentController.GetTournamentTiktoks)
		router.Get("/:tournamentId", tournamentController.GetTournamentDetails)
		router.Put("/:tournamentId", tournamentController.TournamentWinner)
	})

	api.Route("/user", func(router fiber.Router) {
		router.Put("/photo", middleware.Protected(), userController.ChangeUserPhoto)
		router.Get("/tournaments", middleware.Protected(), userController.TournamentsOfUser)
	})
}