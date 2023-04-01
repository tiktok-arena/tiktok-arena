package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"tiktok-arena/database"
	"tiktok-arena/models"
)

// TournamentsOfUser
//
//	@Summary		Get tournaments for user
//	@Description	Get tournaments for user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			page		query		string						false	"page number"
//	@Param			count		query		string						false	"page size"
//	@Param			sort_name	query		string						false	"sort page by name"
//	@Param			sort_size	query		string						false	"sort page by size"
//	@Param			search		query		string						false	"search"
//	@Success		200			{object}	models.TournamentsResponse	"Tournaments of user"
//	@Failure		400			{object}	MessageResponseType			"Couldn't get tournaments for specific user"
//	@Router			/user/tournaments [get]
func TournamentsOfUser(c *fiber.Ctx) error {
	userId, err := GetUserIdAndCheckJWT(c)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}
	p := new(models.PaginationQueries)
	if err := c.QueryParser(p); err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, "Failed to parse queries")
	}
	models.ValidatePaginationQueries(p)
	tournamentResponse, err := database.GetAllTournamentsForUserById(userId, *p)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, "Failed to get tournaments")
	}

	return c.Status(fiber.StatusOK).JSON(tournamentResponse)
}

// ChangeUserPhoto
//
//	@Summary		Change user photo
//	@Description	Change user photo for current user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			payload	body		models.ChangePhotoURL	true	"Data to change photo"
//	@Success		200		{object}	MessageResponseType		"Photo edited"
//	@Failure		400		{object}	MessageResponseType		"Error during photo change"
//	@Router			/user/photo [put]
func ChangeUserPhoto(c *fiber.Ctx) error {
	userId, err := GetUserIdAndCheckJWT(c)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}
	var payload *models.ChangePhotoURL
	err = c.BodyParser(&payload)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}

	err = models.ValidateStruct(payload)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}

	err = database.ChangeUserPhoto(payload.PhotoURL, userId)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}
	return MessageResponse(c, fiber.StatusOK,
		fmt.Sprintf("Successfully changed photo"))
}

func GetUserIdAndCheckJWT(c *fiber.Ctx) (uuid.UUID, error) {
	user := c.Locals("user")

	if user == nil {
		return uuid.UUID{}, MessageResponse(c, fiber.StatusBadRequest, "Empty jwt.token")
	}
	userJWT := user.(*jwt.Token)

	claims := userJWT.Claims.(jwt.MapClaims)

	userId, err := uuid.Parse(claims["sub"].(string))

	return userId, err
}
