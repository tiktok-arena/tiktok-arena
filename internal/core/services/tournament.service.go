package services

import (
	"github.com/google/uuid"
	"tiktok-arena/internal/core/contests"
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

func (s *TournamentService) GetTournamentContest(tournamentIdString string, bracketType string) (bracket *dtos.Contest, err error) {
	if tournamentIdString == "" {
		return bracket, EmptyTournamentIdError{}
	}

	if !dtos.CheckIfAllowedContestType(bracketType) {
		return bracket, NotAllowedContestTypeError{bracketType}
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
		return contests.SingleElimination(*tiktoks), err
	}
	if bracketType == dtos.KingOfTheHill {
		return contests.KingOfTheHill(*tiktoks), err
	}
	return
}
