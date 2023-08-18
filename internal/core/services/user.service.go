package services

import (
	"github.com/google/uuid"
	"tiktok-arena/internal/core/dtos"
	"tiktok-arena/internal/core/models"
	"tiktok-arena/internal/core/validator"
)

type UserServiceTournamentRepository interface {
	TotalTournamentsByUserId(id uuid.UUID, isPrivate bool) (int64, error)
	GetTournamentsByUserID(id uuid.UUID, totalTournaments int64, queries dtos.PaginationQueries, isPrivate bool) (dtos.TournamentsResponseWithUser, error)
}

type UserServiceUserRepository interface {
	ChangeUserPhoto(url string, id uuid.UUID) error
	GetUserByID(id uuid.UUID) (user models.User, err error)
	TotalUsers() (int64, error)
	GetAllUsers(totalUsers int64, queries dtos.PaginationQueries) (dtos.UsersResponse, error)
}

type UserService struct {
	UserRepository       UserServiceUserRepository
	TournamentRepository UserServiceTournamentRepository
}

func NewUserService(userRepository UserServiceUserRepository, tournamentRepository UserServiceTournamentRepository) *UserService {
	return &UserService{UserRepository: userRepository, TournamentRepository: tournamentRepository}
}

func (s *UserService) GetUsers(queries dtos.PaginationQueries) (response dtos.UsersResponse, err error) {
	countUsers, err := s.UserRepository.TotalUsers()
	if err != nil {
		return response, RepositoryError{err}
	}
	response, err = s.UserRepository.GetAllUsers(countUsers, queries)
	if err != nil {
		return response, RepositoryError{err}
	}
	return
}

func (s *UserService) TournamentsOfUser(id uuid.UUID, queries dtos.PaginationQueries, hasAccessToPrivate bool) (response dtos.TournamentsResponseWithUser, err error) {
	countTournamentsForUser, err := s.TournamentRepository.TotalTournamentsByUserId(id, hasAccessToPrivate)
	if err != nil {
		return response, RepositoryError{err}
	}
	response, err = s.TournamentRepository.GetTournamentsByUserID(id, countTournamentsForUser, queries, hasAccessToPrivate)
	if err != nil {
		return response, RepositoryError{err}
	}
	response.User, err = s.UserRepository.GetUserByID(id)
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
