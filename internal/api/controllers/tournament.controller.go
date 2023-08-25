package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"tiktok-arena/internal/api/controllers/response"
	"tiktok-arena/internal/core/dtos"
	"tiktok-arena/internal/core/models"
	"tiktok-arena/internal/core/validator"
)

type TournamentService interface {
	CreateTournament(create dtos.CreateTournament, userId uuid.UUID) error
	EditTournament(edit dtos.EditTournament, userId uuid.UUID, tournamentIdString string) error
	DeleteTournament(userId uuid.UUID, tournamentIdString string) error
	DeleteTournaments(userId uuid.UUID, tournamentIds dtos.TournamentIds) error
	GetTournaments(queries dtos.PaginationQueries) (response dtos.TournamentsResponse, err error)
	GetTournament(tournamentIdString string) (tournament models.Tournament, err error)
	GetTournamentStats(tournamentIdString string) (tournamentStats dtos.TournamentStats, err error)
	TournamentWinner(tournamentIdString string, winner dtos.TournamentWinner) error
	GetTournamentContest(tournamentIdString string, contestType string) (bracket dtos.Contest, err error)
}

type TournamentController struct {
	TournamentService TournamentService
}

func NewTournamentController(tournamentService TournamentService) *TournamentController {
	return &TournamentController{TournamentService: tournamentService}
}

// CreateTournament
//
//	@Summary		Create new tournament
//	@Description	Create new tournament for current user
//	@Tags			tournament
//	@Accept			json
//	@Produce		json
//	@Security		JWT
//	@Param			payload	body		dtos.CreateTournament		true	"Data to create tournament"
//	@Success		200		{object}	dtos.MessageResponseType	"Tournament created"
//	@Failure		400		{object}	dtos.MessageResponseType	"Error during tournament creation"
//	@Router			/api/tournament/create [post]
func (cr *TournamentController) CreateTournament(c *fiber.Ctx) error {
	user := c.Locals("user")
	userId, err := validator.GetUserIdAndCheckJWT(user)
	if err != nil {
		return err
	}

	var payload dtos.CreateTournament
	err = c.BodyParser(&payload)
	if err != nil {
		return err
	}

	err = cr.TournamentService.CreateTournament(payload, userId)
	if err != nil {
		return err
	}

	return response.MessageResponse(c, fiber.StatusCreated, fmt.Sprintf("Tournament created %v", payload.Name))
}

// EditTournament
//
//	@Summary		Edit tournament
//	@Description	Edit tournament for current user
//	@Tags			tournament
//	@Accept			json
//	@Produce		json
//	@Security		JWT
//	@Param			tournamentId	path		string						true	"Tournament id"
//	@Param			payload			body		dtos.EditTournament			true	"Data to edit tournament"
//	@Success		200				{object}	dtos.MessageResponseType	"Tournament edited"
//	@Failure		400				{object}	dtos.MessageResponseType	"Error during tournament edition"
//	@Router			/api/tournament/edit/{tournamentId} [put]
func (cr *TournamentController) EditTournament(c *fiber.Ctx) error {
	user := c.Locals("user")
	userId, err := validator.GetUserIdAndCheckJWT(user)
	if err != nil {
		return err
	}

	var payload dtos.EditTournament
	err = c.BodyParser(&payload)
	if err != nil {
		return err
	}

	tournamentIdString := c.Params("tournamentId")
	err = cr.TournamentService.EditTournament(payload, userId, tournamentIdString)
	if err != nil {
		return err
	}

	return response.MessageResponse(c, fiber.StatusOK, fmt.Sprintf("Successfully edited tournament %s", payload.Name))
}

// DeleteTournament
//
//	@Summary		Delete tournament
//	@Description	Delete tournament for current user
//	@Tags			tournament
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	dtos.MessageResponseType	"Tournament deleted"
//	@Failure		400	{object}	dtos.MessageResponseType	"Error during tournament deletion"
//	@Router			/api/tournament/delete/{tournamentId} [delete]
func (cr *TournamentController) DeleteTournament(c *fiber.Ctx) error {
	user := c.Locals("user")
	userId, err := validator.GetUserIdAndCheckJWT(user)
	if err != nil {
		return err
	}

	tournamentIdString := c.Params("tournamentId")
	err = cr.TournamentService.DeleteTournament(userId, tournamentIdString)
	if err != nil {
		return err
	}

	return response.MessageResponse(c, fiber.StatusOK,
		fmt.Sprintf("Successfully deleted tournament %s", tournamentIdString))
}

// DeleteTournaments
//
//	@Summary		Delete tournaments
//	@Description	Delete tournaments for current user
//	@Tags			tournament
//	@Accept			json
//	@Produce		json
//	@Security		JWT
//	@Param			payload	body		dtos.TournamentIds			true	"Data to delete tournaments"
//	@Success		200		{object}	dtos.MessageResponseType	"Tournaments deleted"
//	@Failure		400		{object}	dtos.MessageResponseType	"Error during tournaments deletion"
//	@Router			/api/tournament/delete [delete]
func (cr *TournamentController) DeleteTournaments(c *fiber.Ctx) error {
	user := c.Locals("user")
	userId, err := validator.GetUserIdAndCheckJWT(user)
	if err != nil {
		return err
	}

	var payload dtos.TournamentIds
	err = c.BodyParser(&payload)
	if err != nil {
		return err
	}

	err = cr.TournamentService.DeleteTournaments(userId, payload)
	if err != nil {
		return err
	}

	return response.MessageResponse(c, fiber.StatusOK, "Successfully deleted tournaments")
}

// GetAllTournaments
//
//	@Summary		All tournaments
//	@Description	Get all tournaments
//	@Tags			tournament
//	@Accept			json
//	@Produce		json
//	@Param			page								query		string						false	"page number"
//	@Param			count								query		string						false	"page size"
//	@Param			search								query		string						false	"search"
//	@Success		200									{array}		dtos.TournamentsResponse	"All tournaments"
//	@Failure		400									{object}	dtos.MessageResponseType	"Failed to get all tournaments"
//	@Router			/api/tournament/tournaments [get]																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																												[get]
func (cr *TournamentController) GetAllTournaments(c *fiber.Ctx) error {
	q := new(dtos.PaginationQueries)
	if err := c.QueryParser(q); err != nil {
		return err
	}
	dtos.ValidatePaginationQueries(q)
	tournamentResponse, err := cr.TournamentService.GetTournaments(*q)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(tournamentResponse)
}

// GetTournamentContest
//
//	@Summary		Tournament contests
//	@Description	Get tournament contests
//	@Tags			tournament
//	@Accept			json
//	@Produce		json
//	@Param			tournamentId	path		string						true	"Tournament id"
//	@Param			payload			query		dtos.ContestPayload			true	"Contest type"
//	@Success		200				{object}	dtos.Contest				"Contest bracket"
//	@Failure		400				{object}	dtos.MessageResponseType	"Failed to return tournament contests"
//	@Router			/api/tournament/contest/{tournamentId} [get]
func (cr *TournamentController) GetTournamentContest(c *fiber.Ctx) error {
	tournamentIdString := c.Params("tournamentId")
	contestType := c.Query("type")
	bracket, err := cr.TournamentService.GetTournamentContest(tournamentIdString, contestType)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(bracket)
}

// GetTournamentStats
//
//	@Summary		Tournament tiktoks
//	@Description	Get tournament tiktoks
//	@Tags			tournament
//	@Accept			json
//	@Produce		json
//	@Param			tournamentId	path		string						true	"Tournament id"
//	@Success		200				{array}		models.Tiktok				"Tournament tiktoks"
//	@Failure		400				{object}	dtos.MessageResponseType	"Tournament not found"
//	@Router			/api/tournament/tiktoks/{tournamentId} [get]
func (cr *TournamentController) GetTournamentStats(c *fiber.Ctx) error {
	tournamentIdString := c.Params("tournamentId")
	tiktoks, err := cr.TournamentService.GetTournamentStats(tournamentIdString)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(tiktoks)
}

// GetTournamentDetails
//
//	@Summary		Tournament details
//	@Description	Get tournament details by its id
//	@Tags			tournament
//	@Accept			json
//	@Produce		json
//	@Param			tournamentId	path		string						true	"Tournament id"
//	@Success		200				{object}	models.Tournament			"Tournament"
//	@Failure		400				{object}	dtos.MessageResponseType	"Tournament not found"
//	@Router			/api/tournament/details/{tournamentId} [get]
func (cr *TournamentController) GetTournamentDetails(c *fiber.Ctx) error {
	tournamentIdString := c.Params("tournamentId")
	tournament, err := cr.TournamentService.GetTournament(tournamentIdString)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(tournament)
}

// TournamentWinner
//
//	@Summary		Update tournament winner statistics
//	@Description	Increment wins and increment times_played
//	@Tags			tournament
//	@Accept			json
//	@Produce		json
//	@Param			tournamentId	path		string						true	"Tournament id"
//	@Param			payload			body		dtos.TournamentWinner		true	"Data to update tournament winner"
//	@Success		200				{object}	dtos.MessageResponseType	"Winner updated"
//	@Failure		400				{object}	dtos.MessageResponseType	"Error during winner updating"
//	@Router			/api/tournament/winner/{tournamentId} [put]
func (cr *TournamentController) TournamentWinner(c *fiber.Ctx) error {
	tournamentIdString := c.Params("tournamentId")

	var payload dtos.TournamentWinner
	err := c.BodyParser(&payload)
	if err != nil {
		return err
	}

	err = cr.TournamentService.TournamentWinner(tournamentIdString, payload)
	if err != nil {
		return err
	}

	return response.MessageResponse(c, fiber.StatusOK,
		fmt.Sprintf("Successfully registered winner for tournament %s", tournamentIdString))
}
