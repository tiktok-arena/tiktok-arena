package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"tiktok-arena/internal/api/controllers/response"
	"tiktok-arena/internal/core/dtos"
	"tiktok-arena/internal/core/validator"
)

type UserService interface {
	GetUsers(queries dtos.PaginationQueries) (response dtos.UsersResponse, err error)
	TournamentsOfUser(id uuid.UUID, queries dtos.PaginationQueries, hasAccessToPrivate bool) (response dtos.TournamentsResponseWithUser, err error)
	ChangeUserPhoto(change dtos.ChangePhotoURL, userId uuid.UUID) (err error)
}

type UserController struct {
	UserService UserService
}

func NewUserController(userService UserService) *UserController {
	return &UserController{UserService: userService}
}

// GetAllUsers
//
//	@Summary		All users
//	@Description	Get all users
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			page					query		string						false	"page number"
//	@Param			count					query		string						false	"page size"
//	@Param			search					query		string						false	"search"
//	@Success		200						{array}		dtos.UsersResponse			"All users"
//	@Failure		400						{object}	dtos.MessageResponseType	"Failed to get all users"
//	@Router			/api/user/users [get]																																																																																																																																																																																																																																																																																																																																																																																																																																																																																				[get]
func (cr *UserController) GetAllUsers(c *fiber.Ctx) error {
	q := new(dtos.PaginationQueries)
	if err := c.QueryParser(q); err != nil {
		return err
	}
	dtos.ValidatePaginationQueries(q)
	userResponse, err := cr.UserService.GetUsers(*q)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(userResponse)
}

// UserInformation
//
//	@Summary		Get user information
//	@Description	Get user information (tournaments, photo and etc.)
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		JWT
//	@Param			userId		path		string								true	"User id"
//	@Param			page		query		string								false	"page number"
//	@Param			count		query		string								false	"page size"
//	@Param			sort_name	query		string								false	"sort page by name"
//	@Param			sort_size	query		string								false	"sort page by size"
//	@Param			search		query		string								false	"search"
//	@Success		200			{object}	dtos.TournamentsResponseWithUser	"User information"
//	@Failure		400			{object}	dtos.MessageResponseType			"Couldn't user information for specific user"
//	@Router			/api/user/profile/{userId} [get]
func (cr *UserController) UserInformation(c *fiber.Ctx) error {
	var hasAccessToPrivate bool
	userIdStringFromPath := c.Params("userId")
	userIdFromPath, err := uuid.Parse(userIdStringFromPath)
	if err != nil {
		return err
	}

	user := c.Locals("user")
	userId, _ := validator.GetUserIdAndCheckJWT(user) // All errors are emitted because JWT is OPTIONAL
	if userId == userIdFromPath {
		hasAccessToPrivate = true
	}

	payload := new(dtos.PaginationQueries)
	if err = c.QueryParser(payload); err != nil {
		return err
	}
	dtos.ValidatePaginationQueries(payload)

	tournamentsResponse, err := cr.UserService.TournamentsOfUser(userIdFromPath, *payload, hasAccessToPrivate)
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
