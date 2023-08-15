package routers

import (
	"github.com/gofiber/fiber/v2"
	"tiktok-arena/internal/api/controllers"
)

func NewTournamentRouter(c *controllers.TournamentController) func(router fiber.Router) {
	return func(router fiber.Router) {
		router.Get("", c.GetAllTournaments)
		router.Get("/contest/:tournamentId", c.GetTournamentContest)
		router.Get("/tiktoks/:tournamentId", c.GetTournamentStats)
		router.Get("/:tournamentId", c.GetTournamentDetails)
		router.Put("/winner/:tournamentId", c.TournamentWinner)
	}
}

func NewTournamentProtectedRouter(c *controllers.TournamentController) func(router fiber.Router) {
	return func(router fiber.Router) {
		router.Post("/create", c.CreateTournament)
		router.Put("/edit/:tournamentId", c.EditTournament)
		router.Delete("/delete/:tournamentId", c.DeleteTournament)
		router.Delete("/delete", c.DeleteTournaments)
	}
}
