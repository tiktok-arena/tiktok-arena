package routers

import (
	"github.com/gofiber/fiber/v2"
	"tiktok-arena/internal/api/controllers"
	"tiktok-arena/internal/api/middleware"
)

func NewAuthRouter(c *controllers.AuthController) func(router fiber.Router) {
	return func(router fiber.Router) {
		router.Post("/register", c.RegisterUser)
		router.Post("/login", c.LoginUser)
		router.Get("/whoami", middleware.Protected(), c.WhoAmI)
	}
}
