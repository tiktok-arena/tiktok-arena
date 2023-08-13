package services

import (
	"github.com/google/uuid"
	"tiktok-arena/internal/core/dtos"
	"tiktok-arena/internal/core/validator"
)

type UserServiceTournamentRepository interface {
	TotalTournamentsByUserId(id uuid.UUID) (int64, error)
	GetAllTournamentsForUserById(id uuid.UUID, totalTournaments int64, queries dtos.PaginationQueries) (dtos.TournamentsResponse, error)
}

type UserServiceUserRepository interface {
	ChangeUserPhoto(url string, id uuid.UUID) error
}

type UserService struct {
	UserRepository       UserServiceUserRepository
	TournamentRepository UserServiceTournamentRepository
}

func NewUserService(userRepository UserServiceUserRepository, tournamentRepository UserServiceTournamentRepository) *UserService {
	return &UserService{UserRepository: userRepository, TournamentRepository: tournamentRepository}
}

func (s *UserService) TournamentsOfUser(id uuid.UUID, queries dtos.PaginationQueries) (response dtos.TournamentsResponse, err error) {
	countTournamentsForUser, err := s.TournamentRepository.TotalTournamentsByUserId(id)
	if err != nil {
		return response, RepositoryError{err}
	}
	response, err = s.TournamentRepository.GetAllTournamentsForUserById(id, countTournamentsForUser, queries)
	if err != nil {
		return response, RepositoryError{err}
	}
	return
}

func (s *UserService) ChangeUserPhoto(change dtos.ChangePhotoURL, userId uuid.UUID) (err error) {
	err = validator.ValidateStruct(change)
	if err != nil {
		return ValidateError{err}
	}
	err = s.UserRepository.ChangeUserPhoto(change.PhotoURL, userId)
	if err != nil {
		return RepositoryError{err}
	}
	return
}
