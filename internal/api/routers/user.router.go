package routers

import (
	"github.com/gofiber/fiber/v2"
	"tiktok-arena/internal/api/controllers"
)

func NewUserProtectedRouter(c *controllers.UserController) func(router fiber.Router) {
	return func(router fiber.Router) {
		router.Put("/photo", c.ChangeUserPhoto)
		router.Get("/tournaments", c.TournamentsOfUser)
	}
}
