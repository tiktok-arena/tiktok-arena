package routers

import (
	"github.com/gofiber/fiber/v2"
	"tiktok-arena/internal/api/controllers"
)

func NewAuthRouter(c *controllers.AuthController) func(router fiber.Router) {
	return func(router fiber.Router) {
		router.Post("/register", c.RegisterUser)
		router.Post("/login", c.LoginUser)
	}
}

func NewAuthProtectedRouter(c *controllers.AuthController) func(router fiber.Router) {
	return func(router fiber.Router) {
		router.Get("/whoami", c.WhoAmI)
	}
}
