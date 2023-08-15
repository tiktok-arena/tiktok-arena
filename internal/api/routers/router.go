package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "tiktok-arena/docs"
)

type GroupRoutes struct {
	AuthGroup       fiber.Router
	UserGroup       fiber.Router
	TournamentGroup fiber.Router
}

func GetGroupRoutes(app *fiber.App) GroupRoutes {

	api := app.Group("/api")

	//	Use 'swag init' to generate new /docs files, details: https://github.com/gofiber/swagger#usage
	api.Get("/docs/*", swagger.HandlerDefault)
	// Redirect to docs
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Redirect("/api/docs/")
	})

	authGroup := api.Group("/auth")
	userGroup := api.Group("/user")
	tournamentGroup := api.Group("/tournament")

	return GroupRoutes{
		AuthGroup:       authGroup,
		UserGroup:       userGroup,
		TournamentGroup: tournamentGroup,
	}
}
