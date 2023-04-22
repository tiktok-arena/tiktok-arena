package services

import (
	"github.com/google/uuid"
	"tiktok-arena/internal/core/dtos"
	"tiktok-arena/internal/core/validator"
	"tiktok-arena/internal/data/repository"
)

type UserService struct {
	UserRepository       *repository.UserRepository
	TournamentRepository *repository.TournamentRepository
}

func NewUserService(userRepository *repository.UserRepository, tournamentRepository *repository.TournamentRepository) *UserService {
	return &UserService{UserRepository: userRepository, TournamentRepository: tournamentRepository}
}

func (s *UserService) TournamentsOfUser(id uuid.UUID, queries *dtos.PaginationQueries) (response dtos.TournamentsResponse, err error) {
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

func (s *UserService) ChangeUserPhoto(change *dtos.ChangePhotoURL, userId uuid.UUID) (err error) {
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
