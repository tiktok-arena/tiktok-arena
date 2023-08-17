package middleware

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"tiktok-arena/configuration"
)

func Protected() func(*fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(configuration.EnvConfig.JwtSecret),
		ErrorHandler: jwtError,
	})
}

func OptionalJWT() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")

		if token == "" {
			// No JWT provided, continue without authentication
			return c.Next()
		}

		// Validate and process the JWT if provided
		return jwtware.New(jwtware.Config{
			SigningKey:   []byte(configuration.EnvConfig.JwtSecret),
			ErrorHandler: jwtError,
		})(c)

	}
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Missing or malformed JWT",
			"data":    nil,
		})

	} else {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Invalid or expired JWT",
			"data":    nil,
		})
	}
}
