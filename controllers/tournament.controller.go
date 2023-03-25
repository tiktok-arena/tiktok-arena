package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"math"
	"tiktok-arena/database"
	"tiktok-arena/models"
)

// CreateTournament
//
//	@Summary		Create new tournament
//	@Description	Create new tournament for current user
//	@Tags			tournament
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			payload	body		models.CreateTournament	true	"Data to create tournament"
//	@Success		200		{object}	MessageResponseType		"Tournament created"
//	@Failure		400		{object}	MessageResponseType		"Error during tournament creation"
//	@Router			/tournament/create [post]
func CreateTournament(c *fiber.Ctx) error {
	userId, err := GetUserIdAndCheckJWT(c)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}
	var payload *models.CreateTournament
	err = c.BodyParser(&payload)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}

	err = models.ValidateStruct(payload)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if payload.Size != len(payload.Tiktoks) {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("Tournament size and count of tiktoks mismatch (%d != %d)",
				payload.Size,
				len(payload.Tiktoks)),
		)
	}

	if database.CheckIfTournamentExistsByName(payload.Name) {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("Tournament %s already exists", payload.Name))
	}

	newTournamentId, err := uuid.NewRandom()
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}

	newTournament := models.Tournament{
		ID:     &newTournamentId,
		Name:   payload.Name,
		UserID: &userId,
		Size:   payload.Size,
	}
	err = database.CreateNewTournament(&newTournament)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}

	for _, value := range payload.Tiktoks {
		tiktok := models.Tiktok{
			TournamentID: &newTournamentId,
			Name:         value.Name,
			URL:          value.URL,
			Wins:         0,
		}
		err = database.CreateNewTiktok(&tiktok)
		if err != nil {
			return MessageResponse(c, fiber.StatusBadRequest, err.Error())
		}
	}

	return MessageResponse(c, fiber.StatusOK,
		fmt.Sprintf("Successfully created tournament %s", payload.Name))
}

// EditTournament
//
//	@Summary		Edit tournament
//	@Description	Edit tournament for current user
//	@Tags			tournament
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			payload	body		models.EditTournament	true	"Data to edit tournament"
//	@Success		200		{object}	MessageResponseType		"Tournament edited"
//	@Failure		400		{object}	MessageResponseType		"Error during tournament edition"
//	@Router			/tournament/edit/{tournamentId} [post]
func EditTournament(c *fiber.Ctx) error {
	userId, err := GetUserIdAndCheckJWT(c)

	var payload *models.EditTournament
	err = c.BodyParser(&payload)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}

	err = models.ValidateStruct(payload)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if payload.Size != len(payload.Tiktoks) {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("Tournament size and count of tiktoks mismatch (%d != %d)",
				payload.Size,
				len(payload.Tiktoks)),
		)
	}

	tournamentIdString := c.Params("tournamentId")

	if tournamentIdString == "" {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("%s is not a valid tournament id", tournamentIdString))
	}

	tournamentId, err := uuid.Parse(tournamentIdString)

	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("Could not parse id %s", tournamentIdString))
	}

	if !database.CheckIfTournamentExistsById(tournamentId) {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("Tournament with id:%s doesn't exist", tournamentIdString))
	}

	if database.CheckIfNameIsTakenByOtherTournament(payload.Name, tournamentId) {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("Tournament name:%s is taken by other tournament", payload.Name))
	}
	// Get tiktoks to edit
	oldS, err := database.GetTournamentTiktoksById(tournamentId)

	editedTournament := models.Tournament{
		ID:     &tournamentId,
		Name:   payload.Name,
		UserID: &userId,
		Size:   payload.Size,
	}

	err = database.EditTournament(&editedTournament)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}

	var newS []models.Tiktok
	for _, value := range payload.Tiktoks {
		tiktok := models.Tiktok{
			TournamentID: &tournamentId,
			Name:         value.Name,
			URL:          value.URL,
			Wins:         0,
		}
		if containsTiktok(oldS, tiktok) {
			err = database.EditTiktok(&tiktok)
			if err != nil {
				return MessageResponse(c, fiber.StatusBadRequest, err.Error())
			}
		}
		newS = append(newS, tiktok)
	}
	tiktoksToDelete := findDifferenceOfTwoTiktokSlices(oldS, newS)

	if len(tiktoksToDelete) != 0 {
		err = database.DeleteTiktoks(tiktoksToDelete)
		if err != nil {
			return MessageResponse(c, fiber.StatusBadRequest, err.Error())
		}
	}

	tiktoksToCreate := findDifferenceOfTwoTiktokSlices(newS, oldS)

	if len(tiktoksToCreate) != 0 {
		err = database.CreateNewTiktoks(tiktoksToCreate)
		if err != nil {
			return MessageResponse(c, fiber.StatusBadRequest, err.Error())
		}
	}

	return MessageResponse(c, fiber.StatusOK,
		fmt.Sprintf("Successfully edited tournament %s", payload.Name))
}

// DeleteTournament
//
//	@Summary		Delete tournament
//	@Description	Delete tournament for current user
//	@Tags			tournament
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	MessageResponseType	"Tournament deleted"
//	@Failure		400	{object}	MessageResponseType	"Error during tournament deletion"
//	@Router			/tournament/delete/{tournamentId} [delete]
func DeleteTournament(c *fiber.Ctx) error {
	userId, err := GetUserIdAndCheckJWT(c)

	tournamentIdString := c.Params("tournamentId")

	if tournamentIdString == "" {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("%s is not a valid tournament id", tournamentIdString))
	}

	tournamentId, err := uuid.Parse(tournamentIdString)

	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("Could not parse id %s", tournamentIdString))
	}

	if !database.CheckIfTournamentExistsById(tournamentId) {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("Tournament with id:%s doesn't exist", tournamentIdString))
	}

	// Get tiktoks to delete
	tiktoksToDelete, err := database.GetTournamentTiktoksById(tournamentId)

	err = database.DeleteTiktoks(tiktoksToDelete)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}

	err = database.DeleteTournamentById(tournamentId, userId)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return MessageResponse(c, fiber.StatusOK,
		fmt.Sprintf("Successfully deleted tournament %s", tournamentIdString))
}

// DeleteTournaments
//
//	@Summary		Delete tournaments
//	@Description	Delete tournaments for current user
//	@Tags			tournament
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			payload	body		models.TournamentIds	true	"Data to delete tournaments"
//	@Success		200		{object}	MessageResponseType		"Tournaments deleted"
//	@Failure		400		{object}	MessageResponseType		"Error during tournaments deletion"
//	@Router			/tournament/delete [delete]
func DeleteTournaments(c *fiber.Ctx) error {
	userId, err := GetUserIdAndCheckJWT(c)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}
	var payload *models.TournamentIds
	err = c.BodyParser(&payload)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}

	err = models.ValidateStruct(payload)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}
	ids := payload.TournamentIds
	b, err := database.CheckIfTournamentsExistsByIds(ids, userId)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}
	if !b {
		return MessageResponse(c, fiber.StatusBadRequest, fmt.Sprintf("One or more tournaments doesn't exist"))
	}
	err = database.DeleteTiktoksByIds(ids)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}
	err = database.DeleteTournamentsByIds(ids, userId)

	return MessageResponse(c, fiber.StatusOK,
		fmt.Sprintf("Successfully deleted tournaments"))
}

// GetAllTournaments
//
//	@Summary		All tournaments
//	@Description	Get all tournaments
//	@Tags			tournament
//	@Accept			json
//	@Produce		json
//	@Param			page				query		string						false	"page number"
//	@Param			count				query		string						false	"page size"
//	@Param			sort_name			query		string						false	"sort page by name"
//	@Param			sort_size			query		string						false	"sort page by size"
//	@Param			search				query		string						false	"search"
//	@Success		200					{array}		models.TournamentsResponse	"Contest bracket"
//	@Failure		400					{object}	MessageResponseType			"Failed to return tournament contest"
//	@Router			/tournament [get]																																																																												[get]
func GetAllTournaments(c *fiber.Ctx) error {
	p := new(models.PaginationQueries)
	if err := c.QueryParser(p); err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, "Failed to parse queries")
	}
	models.ValidatePaginationQueries(p)
	tournamentResponse, err := database.GetTournaments(*p)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, "Failed to get tournaments")
	}

	if tournamentResponse.TournamentCount == 0 {
		return MessageResponse(c, fiber.StatusMovedPermanently, "There is no page")
	}

	return c.Status(fiber.StatusOK).JSON(tournamentResponse)
}

// GetTournamentDetails
//
//	@Summary		Tournament details
//	@Description	Get tournament details by its id
//	@Tags			tournament
//	@Accept			json
//	@Produce		json
//	@Param			tournamentId	path		string				true	"Tournament id"
//	@Success		200				{object}	models.Tournament	"Tournament"
//	@Failure		400				{object}	MessageResponseType	"Tournament not found"
//	@Router			/tournament/{tournamentId} [get]
func GetTournamentDetails(c *fiber.Ctx) error {
	tournamentId := c.Params("tournamentId")
	if tournamentId == "" {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("%s is not a valid tournament id", tournamentId))
	}
	tournament, err := database.GetTournamentById(tournamentId)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(tournament)
}

// GetTournamentTiktoks
//
//	@Summary		Tournament tiktoks
//	@Description	Get tournament tiktoks
//	@Tags			tournament
//	@Accept			json
//	@Produce		json
//	@Param			tournamentId	path		string				true	"Tournament id"
//	@Success		200				{array}		models.Tiktok		"Tournament tiktoks"
//	@Failure		400				{object}	MessageResponseType	"Tournament not found"
//	@Router			/tournament/tiktoks/{tournamentId} [get]
func GetTournamentTiktoks(c *fiber.Ctx) error {
	tournamentIdString := c.Params("tournamentId")
	if tournamentIdString == "" {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("%s is not a valid tournament id", tournamentIdString))
	}
	tournamentId, err := uuid.Parse(tournamentIdString)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}
	tiktoks, err := database.GetTournamentTiktoksById(tournamentId)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("Could not get tiktoks for tournament with id %s", tournamentId))
	}
	return c.Status(fiber.StatusOK).JSON(tiktoks)
}

// TournamentWinner
//
//	@Summary		Update tournament winner statistics
//	@Description	Increment wins and increment times_played
//	@Tags			tournament
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			payload	body		models.CreateTournament	true	"Data to update tournament winner"
//	@Success		200		{object}	MessageResponseType		"Winner updated"
//	@Failure		400		{object}	MessageResponseType		"Error during winner updating"
//	@Router			/tournament/{tournamentId} [post]
func TournamentWinner(c *fiber.Ctx) error {
	_, err := GetUserIdAndCheckJWT(c)

	tournamentIdString := c.Params("tournamentId")

	if tournamentIdString == "" {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("%s is not a valid tournament id", tournamentIdString))
	}

	tournamentId, err := uuid.Parse(tournamentIdString)

	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("Could not parse id %s", tournamentIdString))
	}

	var payload *models.TournamentWinner
	err = c.BodyParser(&payload)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}

	err = models.ValidateStruct(payload)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}

	err = database.RegisterTiktokWinner(tournamentId, payload.TiktokURL)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return MessageResponse(c, fiber.StatusOK,
		fmt.Sprintf("Successfully registered winner for tournament %s", tournamentId))
}

// GetTournamentContest
//
//	@Summary		Tournament contest
//	@Description	Get tournament contest
//	@Tags			tournament
//	@Accept			json
//	@Produce		json
//	@Param			tournamentId	path		string					true	"Tournament id"
//	@Param			payload			query		models.ContestPayload	true	"Contest type"
//	@Success		200				{object}	models.Bracket			"Contest bracket"
//	@Failure		400				{object}	MessageResponseType		"Failed to return tournament contest"
//	@Router			/tournament/contest/{tournamentId} [get]
func GetTournamentContest(c *fiber.Ctx) error {
	tournamentIdString := c.Params("tournamentId")
	if tournamentIdString == "" {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("%s is not a valid tournament id", tournamentIdString))
	}

	contestType := c.Query("type")
	if !models.CheckIfAllowedTournamentType(contestType) {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("%s is not allowed tournament format", contestType),
		)
	}
	tournamentId, err := uuid.Parse(tournamentIdString)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}

	tiktoks, err := database.GetTournamentTiktoksById(tournamentId)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("Could not get tiktoks for tournament with id %s", tournamentId))
	}
	shuffleTiktok(tiktoks)
	if contestType == models.SingleElimination {
		return c.Status(fiber.StatusOK).JSON(SingleElimination(tiktoks))
	}
	if contestType == models.KingOfTheHill {
		return c.Status(fiber.StatusOK).JSON(KingOfTheHill(tiktoks))
	}
	return MessageResponse(c, fiber.StatusBadRequest, "Unknown error")
}

// SingleElimination
// https://en.wikipedia.org/wiki/Single-elimination_tournament
func SingleElimination(t []models.Tiktok) models.Bracket {
	countTiktok := len(t)
	countRound := int(math.Ceil(math.Log2(float64(countTiktok))))
	countSecondRoundParticipators := 1 << (countRound - 1) // Equivalent to int(math.Pow(2, float64(countRound)) / 2)
	countFirstRoundMatches := countTiktok - int(math.Pow(2, float64(countRound)-1))
	countFirstRoundTiktoks := countFirstRoundMatches * 2

	rounds := make([]models.Round, 0, countRound)

	firstRoundMatches := make([]models.Match, 0, countFirstRoundMatches)
	secondRoundMatches := make([]models.Match, 0, countSecondRoundParticipators/2)

	secondRoundParticipators := make([]models.Option, 0, countSecondRoundParticipators) // This slice should store MatchOption or TiktokOption

	// Filling first round with firstRoundMatches and appending MatchOptions to second round participators
	for j := 0; j < countFirstRoundTiktoks; j += 2 {
		matchID := uuid.NewString()
		firstRoundMatches = append(firstRoundMatches, models.Match{
			MatchID: matchID,
			FirstOption: models.TiktokOption{
				TiktokURL: t[j].URL,
			},
			SecondOption: models.TiktokOption{
				TiktokURL: t[j+1].URL,
			},
		})
		secondRoundParticipators = append(secondRoundParticipators,
			models.MatchOption{MatchID: matchID})
	}
	// Appending first round firstRoundMatches to rounds
	rounds = append(rounds, models.Round{
		Round:   1,
		Matches: firstRoundMatches,
	})
	// Appending TiktokOptions to second round participators
	for _, tiktok := range t[countFirstRoundTiktoks:] {
		secondRoundParticipators = append(secondRoundParticipators,
			models.TiktokOption{TiktokURL: tiktok.URL})
	}
	// Generating second round firstRoundMatches
	for i := 0; i < int(countSecondRoundParticipators); i += 2 {
		match := models.Match{
			MatchID:      uuid.NewString(),
			FirstOption:  secondRoundParticipators[i],
			SecondOption: secondRoundParticipators[i+1],
		}
		secondRoundMatches = append(secondRoundMatches, match)
	}
	// Generating second round
	secondRound := models.Round{
		Round:   2,
		Matches: secondRoundMatches}
	rounds = append(rounds, secondRound)

	previousRoundMatches := secondRoundMatches
	for roundID := 3; roundID <= countRound; roundID++ {
		// Generating Nth round matches (where N > 2)
		var currentRoundMatches []models.Match
		for matchID := 0; matchID < len(previousRoundMatches); matchID += 2 {
			match := models.Match{
				MatchID: uuid.NewString(),
				FirstOption: models.MatchOption{
					MatchID: previousRoundMatches[matchID].MatchID,
				},
				SecondOption: models.MatchOption{
					MatchID: previousRoundMatches[matchID+1].MatchID,
				},
			}
			currentRoundMatches = append(currentRoundMatches, match)
		}
		// Generating Nth round (where N > 2)
		round := models.Round{
			Round:   roundID,
			Matches: currentRoundMatches,
		}
		rounds = append(rounds, round)

		previousRoundMatches = currentRoundMatches
	}
	return models.Bracket{
		CountMatches: countTiktok - 1,
		Rounds:       rounds,
	}

}

// KingOfTheHill
// First match decided randomly between two participators.
// Loser of match leaves the game, winner will go to next match, next opponent decided randomly from standings.
// Procedure continues until last standing.
func KingOfTheHill(t []models.Tiktok) models.Bracket {
	countTiktok := len(t)
	rounds := make([]models.Round, 0, countTiktok-1)
	match := models.Match{
		MatchID:      uuid.NewString(),
		FirstOption:  models.TiktokOption{TiktokURL: t[0].URL},
		SecondOption: models.TiktokOption{TiktokURL: t[1].URL},
	}
	rounds = append(rounds, models.Round{
		Round:   1,
		Matches: []models.Match{match},
	})
	previousMatch := match
	for i := 2; i < countTiktok-1; i++ {
		match = models.Match{
			MatchID:      uuid.NewString(),
			FirstOption:  models.MatchOption{MatchID: previousMatch.MatchID},
			SecondOption: models.TiktokOption{TiktokURL: t[i].URL},
		}
		rounds = append(rounds, models.Round{
			Round:   i,
			Matches: []models.Match{match},
		})
		previousMatch = match
	}
	return models.Bracket{
		CountMatches: countTiktok - 1,
		Rounds:       rounds,
	}
}
