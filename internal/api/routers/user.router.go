package routers

import (
	"github.com/gofiber/fiber/v2"
	"tiktok-arena/internal/api/controllers"
	"tiktok-arena/internal/api/middleware"
)

func NewUserRouter(c *controllers.UserController) func(router fiber.Router) {
	return func(router fiber.Router) {
		router.Put("/photo", middleware.Protected(), c.ChangeUserPhoto)
		router.Get("/profile/:userId", middleware.OptionalJWT(), c.UserInformation)
	}
}
