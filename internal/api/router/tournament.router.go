package router

import (
	"github.com/gofiber/fiber/v2"
	"tiktok-arena/internal/api/controllers"
	"tiktok-arena/internal/api/middleware"
)

func NewTournamentRouter(c *controllers.TournamentController) func(router fiber.Router) {
	return func(router fiber.Router) {
		router.Get("", c.GetAllTournaments)
		router.Get("/contest/:tournamentId", c.GetTournamentContest)
		router.Post("/create", middleware.Protected(), c.CreateTournament)
		router.Put("/edit/:tournamentId", middleware.Protected(), c.EditTournament)
		router.Delete("/delete/:tournamentId", middleware.Protected(), c.DeleteTournament)
		router.Delete("/delete", middleware.Protected(), c.DeleteTournaments)
		router.Get("/tiktoks/:tournamentId", c.GetTournamentTiktoks)
		router.Get("/:tournamentId", c.GetTournamentDetails)
		router.Put("/:tournamentId", c.TournamentWinner)
	}
}
