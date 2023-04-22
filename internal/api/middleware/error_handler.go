package middleware

import (
	"github.com/gofiber/fiber/v2"
	"tiktok-arena/internal/core/services"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	switch e := err.(type) {
	case services.ValidateError:
		code = fiber.StatusBadRequest
		message = e.Error()
	case services.UserAlreadyExistsError:
		code = fiber.StatusConflict
		message = e.Error()
	case services.JWTGenerateError:
		code = fiber.StatusBadRequest
		message = e.Error()
	case services.RepositoryError:
		code = fiber.StatusBadRequest
		message = e.Error()
	case services.BcryptError:
		code = fiber.StatusBadRequest
		message = e.Error()
	case services.UUIDError:
		code = fiber.StatusBadRequest
		message = e.Error()
	case services.TournamentSizeAndTiktokCountMismatchError:
		code = fiber.StatusBadRequest
		message = e.Error()
	case services.TournamentAlreadyExistsError:
		code = fiber.StatusConflict
		message = e.Error()
	case services.TournamentNotExistsError:
		code = fiber.StatusBadRequest
		message = e.Error()
	case services.TournamentNameIsTakenError:
		code = fiber.StatusBadRequest
		message = e.Error()
	case services.EmptyTournamentIdError:
		code = fiber.StatusBadRequest
		message = e.Error()
	case services.EmptyTiktokURLError:
		code = fiber.StatusBadRequest
		message = e.Error()
	case services.NotAllowedTournamentTypeError:
		code = fiber.StatusBadRequest
		message = e.Error()
	default:
		message = err.Error()
	}

	err = c.Status(code).JSON(fiber.Map{"message": message})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Cannot send error JSON message")
	}

	return nil
}
