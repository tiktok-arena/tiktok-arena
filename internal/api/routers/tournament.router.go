package routers

import (
	"github.com/gofiber/fiber/v2"
	"tiktok-arena/internal/api/controllers"
	"tiktok-arena/internal/api/middleware"
)

func NewTournamentRouter(c *controllers.TournamentController) func(router fiber.Router) {
	return func(router fiber.Router) {
		router.Get("/tournaments", c.GetAllTournaments)
		router.Get("/contest/:tournamentId", c.GetTournamentContest)
		router.Get("/tiktoks/:tournamentId", c.GetTournamentStats)
		router.Get("/details/:tournamentId", c.GetTournamentDetails)
		router.Put("/winner/:tournamentId", c.TournamentWinner)

		router.Post("/create", middleware.Protected(), c.CreateTournament)
		router.Put("/edit/:tournamentId", middleware.Protected(), c.EditTournament)
		router.Delete("/delete/:tournamentId", middleware.Protected(), c.DeleteTournament)
		router.Delete("/delete", middleware.Protected(), c.DeleteTournaments)
	}
}
