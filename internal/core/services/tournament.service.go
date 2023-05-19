package services

import (
	"github.com/google/uuid"
	"math"
	"tiktok-arena/internal/core/dtos"
	"tiktok-arena/internal/core/models"
	"tiktok-arena/internal/core/validator"
)

type TournamentServiceTournamentRepository interface {
	GetTournamentById(tournamentId uuid.UUID) (*models.Tournament, error)
	CheckIfTournamentExistsByName(name string) (bool, error)
	CheckIfNameIsTakenByOtherTournament(name string, id uuid.UUID) (bool, error)
	CheckIfTournamentExistsById(id uuid.UUID) (bool, error)
	CheckIfTournamentsExistsByIds(ids []string, userId uuid.UUID) (bool, error)
	CreateNewTournament(newTournament *models.Tournament) error
	EditTournament(t *models.Tournament) error
	DeleteTournamentById(id uuid.UUID, userId uuid.UUID) error
	DeleteTournamentsByIds(ids []string, userId uuid.UUID) error
	GetTournaments(totalTournaments int64, queries dtos.PaginationQueries) (dtos.TournamentsResponse, error)
	TotalTournaments() (int64, error)
	UpdateTournamentTimesPlayed(tournamentId uuid.UUID) error
}

type TournamentServiceTiktokRepository interface {
	CreateNewTiktok(newTiktok *models.Tiktok) error
	CreateNewTiktoks(t []models.Tiktok) error
	EditTiktok(t *models.Tiktok) error
	DeleteTiktoks(t *[]models.Tiktok) error
	DeleteTiktoksByIds(ids []string) error
	GetTournamentTiktoksById(tournamentId uuid.UUID) (*[]models.Tiktok, error)
	UpdateTiktokWins(tournamentId uuid.UUID, tiktokURL string) error
}

type TournamentService struct {
	TournamentRepository TournamentServiceTournamentRepository
	TiktokRepository     TournamentServiceTiktokRepository
}

func NewTournamentService(tournamentRepository TournamentServiceTournamentRepository,
	tiktokRepository TournamentServiceTiktokRepository) *TournamentService {
	return &TournamentService{TournamentRepository: tournamentRepository, TiktokRepository: tiktokRepository}
}

func (s *TournamentService) CreateTournament(create *dtos.CreateTournament, userId uuid.UUID) error {
	err := validator.ValidateStruct(create)
	if err != nil {
		return ValidateError{err}
	}

	if create.Size != len(create.Tiktoks) {
		return TournamentSizeAndTiktokCountMismatchError{create.Size, len(create.Tiktoks)}
	}

	tournamentExists, err := s.TournamentRepository.CheckIfTournamentExistsByName(create.Name)
	if err != nil {
		return RepositoryError{err}
	}
	if tournamentExists {
		return TournamentAlreadyExistsError{create.Name}
	}

	newTournamentId, err := uuid.NewRandom()
	if err != nil {
		return UUIDError{err}
	}

	newTournament := models.Tournament{
		ID:       &newTournamentId,
		Name:     create.Name,
		UserID:   &userId,
		Size:     create.Size,
		PhotoURL: create.PhotoURL,
	}
	err = s.TournamentRepository.CreateNewTournament(&newTournament)
	if err != nil {
		return RepositoryError{err}
	}

	for _, value := range create.Tiktoks {
		tiktok := models.Tiktok{
			TournamentID: &newTournamentId,
			Name:         value.Name,
			URL:          value.URL,
			Wins:         0,
		}
		err = s.TiktokRepository.CreateNewTiktok(&tiktok)
		if err != nil {
			return RepositoryError{err}
		}
	}

	return nil
}

func (s *TournamentService) EditTournament(edit *dtos.EditTournament, userId uuid.UUID, tournamentIdString string) error {
	err := validator.ValidateStruct(edit)
	if err != nil {
		return ValidateError{err}
	}

	if edit.Size != len(edit.Tiktoks) {
		return TournamentSizeAndTiktokCountMismatchError{edit.Size, len(edit.Tiktoks)}
	}

	if tournamentIdString == "" {
		return EmptyTournamentIdError{}
	}

	tournamentIdUUID, err := uuid.Parse(tournamentIdString)

	if err != nil {
		return UUIDError{err}
	}

	exists, err := s.TournamentRepository.CheckIfTournamentExistsById(tournamentIdUUID)
	if err != nil {
		return RepositoryError{err}
	}
	if !exists {
		return TournamentNotExistsError{tournamentIdUUID}
	}

	nameIsTakenByOtherTournament, err := s.TournamentRepository.CheckIfNameIsTakenByOtherTournament(edit.Name, tournamentIdUUID)
	if err != nil {
		return RepositoryError{err}
	}
	if nameIsTakenByOtherTournament {
		return TournamentNameIsTakenError{TournamentName: edit.Name}
	}

	// Get tiktoks to edit
	oldS, err := s.TiktokRepository.GetTournamentTiktoksById(tournamentIdUUID)
	if err != nil {
		return RepositoryError{err}
	}

	editedTournament := models.Tournament{
		ID:       &tournamentIdUUID,
		Name:     edit.Name,
		UserID:   &userId,
		Size:     edit.Size,
		PhotoURL: edit.PhotoURL,
	}

	err = s.TournamentRepository.EditTournament(&editedTournament)
	if err != nil {
		return RepositoryError{err}
	}

	var newS []models.Tiktok
	for _, value := range edit.Tiktoks {
		tiktok := models.Tiktok{
			TournamentID: &tournamentIdUUID,
			Name:         value.Name,
			URL:          value.URL,
			Wins:         0,
		}
		if models.ContainsTiktok(*oldS, tiktok) {
			err = s.TiktokRepository.EditTiktok(&tiktok)
			if err != nil {
				return RepositoryError{err}
			}
		}
		newS = append(newS, tiktok)
	}

	tiktoksToDelete := models.FindDifferenceOfTwoTiktokSlices(*oldS, newS)
	if len(tiktoksToDelete) != 0 {
		err = s.TiktokRepository.DeleteTiktoks(&tiktoksToDelete)
		if err != nil {
			return RepositoryError{err}
		}
	}

	tiktoksToCreate := models.FindDifferenceOfTwoTiktokSlices(newS, *oldS)
	if len(tiktoksToCreate) != 0 {
		err = s.TiktokRepository.CreateNewTiktoks(tiktoksToCreate)
		if err != nil {
			return RepositoryError{err}
		}
	}

	return nil
}

func (s *TournamentService) DeleteTournament(userId uuid.UUID, tournamentIdString string) error {
	if tournamentIdString == "" {
		return EmptyTournamentIdError{}
	}

	tournamentIdUUID, err := uuid.Parse(tournamentIdString)

	if err != nil {
		return UUIDError{err}
	}

	exists, err := s.TournamentRepository.CheckIfTournamentExistsById(tournamentIdUUID)
	if err != nil {
		return RepositoryError{err}
	}
	if !exists {
		return TournamentNotExistsError{}
	}

	// Get tiktoks to delete
	tiktoksToDelete, err := s.TiktokRepository.GetTournamentTiktoksById(tournamentIdUUID)
	if err != nil {
		return RepositoryError{err}
	}

	err = s.TiktokRepository.DeleteTiktoks(tiktoksToDelete)
	if err != nil {
		return RepositoryError{err}
	}

	err = s.TournamentRepository.DeleteTournamentById(tournamentIdUUID, userId)
	if err != nil {
		return RepositoryError{err}
	}

	return nil
}

func (s *TournamentService) DeleteTournaments(userId uuid.UUID, tournamentIds *dtos.TournamentIds) error {
	err := validator.ValidateStruct(tournamentIds)
	if err != nil {
		return ValidateError{err}
	}
	ids := tournamentIds.TournamentIds
	exists, err := s.TournamentRepository.CheckIfTournamentsExistsByIds(ids, userId)
	if err != nil {
		return RepositoryError{err}
	}
	if !exists {
		return TournamentNotExistsError{}
	}
	err = s.TiktokRepository.DeleteTiktoksByIds(ids)
	if err != nil {
		return RepositoryError{err}
	}
	err = s.TournamentRepository.DeleteTournamentsByIds(ids, userId)
	if err != nil {
		return RepositoryError{err}
	}
	return nil
}

func (s *TournamentService) GetAllTournaments(queries *dtos.PaginationQueries) (response dtos.TournamentsResponse, err error) {
	countTournaments, err := s.TournamentRepository.TotalTournaments()
	if err != nil {
		return response, RepositoryError{err}
	}
	response, err = s.TournamentRepository.GetTournaments(countTournaments, *queries)
	if err != nil {
		return response, RepositoryError{err}
	}
	return
}

func (s *TournamentService) GetTournamentDetails(tournamentIdString string) (tournament *models.Tournament, err error) {
	if tournamentIdString == "" {
		return tournament, EmptyTournamentIdError{}
	}
	tournamentIdUUID, err := uuid.Parse(tournamentIdString)
	if err != nil {
		return tournament, UUIDError{err}
	}
	tournament, err = s.TournamentRepository.GetTournamentById(tournamentIdUUID)
	if err != nil {
		return tournament, RepositoryError{err}
	}
	return
}

func (s *TournamentService) GetTournamentTiktoks(tournamentIdString string) (tiktoks *[]models.Tiktok, err error) {
	if tournamentIdString == "" {
		return tiktoks, EmptyTournamentIdError{}
	}
	tournamentIdUUID, err := uuid.Parse(tournamentIdString)
	if err != nil {
		return tiktoks, UUIDError{err}
	}
	tiktoks, err = s.TiktokRepository.GetTournamentTiktoksById(tournamentIdUUID)
	if err != nil {
		return tiktoks, RepositoryError{err}
	}
	return
}

func (s *TournamentService) TournamentWinner(tournamentIdString string, winner *dtos.TournamentWinner) error {
	if tournamentIdString == "" {
		return EmptyTournamentIdError{}
	}

	tournamentId, err := uuid.Parse(tournamentIdString)

	err = validator.ValidateStruct(winner)
	if err != nil {
		return ValidateError{err}
	}

	if winner.TiktokURL == "" {
		return EmptyTiktokURLError{}
	}

	err = s.TournamentRepository.UpdateTournamentTimesPlayed(tournamentId)
	if err != nil {
		return RepositoryError{err}
	}

	err = s.TiktokRepository.UpdateTiktokWins(tournamentId, winner.TiktokURL)
	if err != nil {
		return RepositoryError{err}
	}
	return nil
}

func (s *TournamentService) GetTournamentContest(tournamentIdString string, bracketType string) (bracket *dtos.Bracket, err error) {
	if tournamentIdString == "" {
		return bracket, EmptyTournamentIdError{}
	}

	if !dtos.CheckIfAllowedBracketType(bracketType) {
		return bracket, NotAllowedBracketTypeError{bracketType}
	}
	tournamentId, err := uuid.Parse(tournamentIdString)
	if err != nil {
		return bracket, UUIDError{err}
	}

	tiktoks, err := s.TiktokRepository.GetTournamentTiktoksById(tournamentId)
	if err != nil {
		return bracket, RepositoryError{err}
	}
	models.ShuffleTiktok(*tiktoks)
	if bracketType == dtos.SingleElimination {
		return SingleElimination(*tiktoks), err
	}
	if bracketType == dtos.KingOfTheHill {
		return KingOfTheHill(*tiktoks), err
	}
	return
}

// SingleElimination
// https://en.wikipedia.org/wiki/Single-elimination_tournament
func SingleElimination(t []models.Tiktok) *dtos.Bracket {
	countTiktok := len(t)
	countRound := int(math.Ceil(math.Log2(float64(countTiktok))))
	countSecondRoundParticipators := 1 << (countRound - 1) // Equivalent to int(math.Pow(2, float64(countRound)) / 2)
	countFirstRoundMatches := countTiktok - int(math.Pow(2, float64(countRound)-1))
	countFirstRoundTiktoks := countFirstRoundMatches * 2

	rounds := make([]dtos.Round, 0, countRound)

	firstRoundMatches := make([]dtos.Match, 0, countFirstRoundMatches)
	secondRoundMatches := make([]dtos.Match, 0, countSecondRoundParticipators/2)

	secondRoundParticipators := make([]dtos.Option, 0, countSecondRoundParticipators) // This slice should store MatchOption or TiktokOption

	// Filling first round with firstRoundMatches and appending MatchOptions to second round participators
	for j := 0; j < countFirstRoundTiktoks; j += 2 {
		matchID := uuid.NewString()
		firstRoundMatches = append(firstRoundMatches, dtos.Match{
			MatchID: matchID,
			FirstOption: dtos.TiktokOption{
				TiktokURL: t[j].URL,
			},
			SecondOption: dtos.TiktokOption{
				TiktokURL: t[j+1].URL,
			},
		})
		secondRoundParticipators = append(secondRoundParticipators,
			dtos.MatchOption{MatchID: matchID})
	}
	// Appending first round firstRoundMatches to rounds
	rounds = append(rounds, dtos.Round{
		Round:   1,
		Matches: firstRoundMatches,
	})
	// Appending TiktokOptions to second round participators
	for _, tiktok := range t[countFirstRoundTiktoks:] {
		secondRoundParticipators = append(secondRoundParticipators,
			dtos.TiktokOption{TiktokURL: tiktok.URL})
	}
	// Generating second round firstRoundMatches
	for i := 0; i < int(countSecondRoundParticipators); i += 2 {
		match := dtos.Match{
			MatchID:      uuid.NewString(),
			FirstOption:  secondRoundParticipators[i],
			SecondOption: secondRoundParticipators[i+1],
		}
		secondRoundMatches = append(secondRoundMatches, match)
	}
	// Generating second round
	secondRound := dtos.Round{
		Round:   2,
		Matches: secondRoundMatches}
	rounds = append(rounds, secondRound)

	previousRoundMatches := secondRoundMatches
	for roundID := 3; roundID <= countRound; roundID++ {
		// Generating Nth round matches (where N > 2)
		var currentRoundMatches []dtos.Match
		for matchID := 0; matchID < len(previousRoundMatches); matchID += 2 {
			match := dtos.Match{
				MatchID: uuid.NewString(),
				FirstOption: dtos.MatchOption{
					MatchID: previousRoundMatches[matchID].MatchID,
				},
				SecondOption: dtos.MatchOption{
					MatchID: previousRoundMatches[matchID+1].MatchID,
				},
			}
			currentRoundMatches = append(currentRoundMatches, match)
		}
		// Generating Nth round (where N > 2)
		round := dtos.Round{
			Round:   roundID,
			Matches: currentRoundMatches,
		}
		rounds = append(rounds, round)

		previousRoundMatches = currentRoundMatches
	}
	return &dtos.Bracket{
		CountMatches: countTiktok - 1,
		Rounds:       rounds,
	}

}

// KingOfTheHill
// First match decided randomly between two participators.
// Loser of match leaves the game, winner will go to next match, next opponent decided randomly from standings.
// Procedure continues until last standing.
func KingOfTheHill(t []models.Tiktok) *dtos.Bracket {
	countTiktok := len(t)
	rounds := make([]dtos.Round, 0, countTiktok-1)
	match := dtos.Match{
		MatchID:      uuid.NewString(),
		FirstOption:  dtos.TiktokOption{TiktokURL: t[0].URL},
		SecondOption: dtos.TiktokOption{TiktokURL: t[1].URL},
	}
	rounds = append(rounds, dtos.Round{
		Round:   1,
		Matches: []dtos.Match{match},
	})
	previousMatch := match
	for i := 2; i < countTiktok-1; i++ {
		match = dtos.Match{
			MatchID:      uuid.NewString(),
			FirstOption:  dtos.MatchOption{MatchID: previousMatch.MatchID},
			SecondOption: dtos.TiktokOption{TiktokURL: t[i].URL},
		}
		rounds = append(rounds, dtos.Round{
			Round:   i,
			Matches: []dtos.Match{match},
		})
		previousMatch = match
	}
	return &dtos.Bracket{
		CountMatches: countTiktok - 1,
		Rounds:       rounds,
	}
}
