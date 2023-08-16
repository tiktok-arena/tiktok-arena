package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"tiktok-arena/internal/api/controllers/response"
	"tiktok-arena/internal/core/dtos"
	"tiktok-arena/internal/core/validator"
)

type UserService interface {
	TournamentsOfUser(id uuid.UUID, queries dtos.PaginationQueries) (response dtos.TournamentsResponse, err error)
	ChangeUserPhoto(change dtos.ChangePhotoURL, userId uuid.UUID) (err error)
}

type UserController struct {
	UserService UserService
}

func NewUserController(userService UserService) *UserController {
	return &UserController{UserService: userService}
}

// TournamentsOfUser
//
//	@Summary		Get tournaments for user
//	@Description	Get tournaments for user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		JWT
//	@Param			page		query		string						false	"page number"
//	@Param			count		query		string						false	"page size"
//	@Param			sort_name	query		string						false	"sort page by name"
//	@Param			sort_size	query		string						false	"sort page by size"
//	@Param			search		query		string						false	"search"
//	@Success		200			{object}	dtos.TournamentsResponse	"Tournaments of user"
//	@Failure		400			{object}	dtos.MessageResponseType	"Couldn't get tournaments for specific user"
//	@Router			/api/user/tournaments [get]
func (cr *UserController) TournamentsOfUser(c *fiber.Ctx) error {
	user := c.Locals("user")
	userId, err := validator.GetUserIdAndCheckJWT(user)
	if err != nil {
		return err
	}
	payload := new(dtos.PaginationQueries)
	if err = c.QueryParser(payload); err != nil {
		return err
	}
	dtos.ValidatePaginationQueries(payload)

	tournamentsResponse, err := cr.UserService.TournamentsOfUser(userId, *payload)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(tournamentsResponse)
}

// ChangeUserPhoto
//
//	@Summary		Change user photo
//	@Description	Change user photo for current user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		JWT
//	@Param			payload	body		dtos.ChangePhotoURL			true	"Data to change photo"
//	@Success		200		{object}	dtos.MessageResponseType	"Photo edited"
//	@Failure		400		{object}	dtos.MessageResponseType	"Error during photo change"
//	@Router			/api/user/photo [put]
func (cr *UserController) ChangeUserPhoto(c *fiber.Ctx) error {
	user := c.Locals("user")
	userId, err := validator.GetUserIdAndCheckJWT(user)
	if err != nil {
		return err
	}
	var payload dtos.ChangePhotoURL
	err = c.BodyParser(&payload)
	if err != nil {
		return err
	}
	err = cr.UserService.ChangeUserPhoto(payload, userId)
	if err != nil {
		return err
	}
	return response.MessageResponse(c, fiber.StatusCreated, "User Photo successfully added")
}
