package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "tiktok-arena/docs"
)

func SetupRoutes(app *fiber.App,
	tournamentRouter func(router fiber.Router),
	authRouter func(router fiber.Router),
	userRouter func(router fiber.Router)) {
	api := app.Group("/api")

	//	Use 'swag init' to generate new /docs files, details: https://github.com/gofiber/swagger#usage
	api.Get("/docs/*", swagger.HandlerDefault)

	api.Route("/auth", authRouter)

	api.Route("/tournament", tournamentRouter)

	api.Route("/user", userRouter)
}
