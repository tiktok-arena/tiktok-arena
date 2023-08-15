package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"tiktok-arena/internal/core/dtos"
)

type AuthService interface {
	NewUser(auth dtos.AuthInput) (details dtos.RegisterDetails, err error)
	GetUserByNameAndPassword(input dtos.AuthInput) (details dtos.LoginDetails, err error)
	WhoAmI(token jwt.Token) (whoami dtos.WhoAmI, err error)
}

type AuthController struct {
	AuthService AuthService
}

func NewAuthController(authService AuthService) *AuthController {
	return &AuthController{AuthService: authService}
}

// RegisterUser
//
//	@Summary		Register user
//	@Description	Register new user with given credentials
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload			body		dtos.AuthInput				true	"Data to register user"
//	@Success		200				{object}	dtos.RegisterDetails		"Register success"
//	@Failure		400				{object}	dtos.MessageResponseType	"Failed to register user"
//	@Router			/auth/register	[post]
func (cr *AuthController) RegisterUser(c *fiber.Ctx) error {
	var payload dtos.AuthInput

	err := c.BodyParser(&payload)
	if err != nil {
		return err
	}

	details, err := cr.AuthService.NewUser(payload)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(details)
}

// LoginUser
//
//	@Summary		Login user
//	@Description	Login user with given credentials
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload			body		dtos.AuthInput				true	"Data to login user"
//	@Success		200				{object}	dtos.RegisterDetails		"Login success"
//	@Failure		400				{object}	dtos.MessageResponseType	"Error logging in"
//	@Router			/auth/login    	[post]
func (cr *AuthController) LoginUser(c *fiber.Ctx) error {
	var payload dtos.AuthInput

	err := c.BodyParser(&payload)
	if err != nil {
		return err
	}

	details, err := cr.AuthService.GetUserByNameAndPassword(payload)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(details)
}

// WhoAmI
//
//	@Summary		Authenticated user details
//	@Description	Get current user id and name
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	dtos.WhoAmI					"User details"
//	@Failure		400	{object}	dtos.MessageResponseType	"Error getting user data"
//	@Router			/auth/whoami [get]
func (cr *AuthController) WhoAmI(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)

	whoami, err := cr.AuthService.WhoAmI(*token)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(whoami)
}
