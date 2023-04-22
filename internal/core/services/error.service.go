package services

import (
	"fmt"
	"github.com/google/uuid"
)

type ValidateError struct {
	error
}

func (e ValidateError) Error() string {
	return fmt.Sprintf("Validate error: %v", e.error)
}

type UserNotExistsError struct {
	Username string
}

func (e UserNotExistsError) Error() string {
	return fmt.Sprintf("User %s already exists", e.Username)
}

type UserAlreadyExistsError struct {
	Username string
}

func (e UserAlreadyExistsError) Error() string {
	return fmt.Sprintf("User %s already exists", e.Username)
}

type JWTGenerateError struct {
	error
}

func (e JWTGenerateError) Error() string {
	return fmt.Sprintf("Generating JWT Token failed: %v", e.error)
}

type RepositoryError struct {
	error
}

func (e RepositoryError) Error() string {
	return fmt.Sprintf("Repository Error: %v", e.error)
}

type BcryptError struct {
	error
}

func (e BcryptError) Error() string {
	return fmt.Sprintf("Bcrypt error: %v", e.error)
}

type UUIDError struct {
	error
}

func (e UUIDError) Error() string {
	return fmt.Sprintf("UUID error: %v", e.error)
}

type TournamentSizeAndTiktokCountMismatchError struct {
	TournamentSize int
	TiktokCount    int
}

func (e TournamentSizeAndTiktokCountMismatchError) Error() string {
	return fmt.Sprintf("Tournament size and count of tiktoks mismatch (%d != %d)",
		e.TournamentSize,
		e.TiktokCount)
}

type TournamentAlreadyExistsError struct {
	TournamentName string
}

func (e TournamentAlreadyExistsError) Error() string {
	return fmt.Sprintf("Tournament %s already exists", e.TournamentName)
}

type TournamentNotExistsError struct {
	TournamentId uuid.UUID
}

func (e TournamentNotExistsError) Error() string {
	return fmt.Sprintf("Tournament with id: %s already exists", e.TournamentId)
}

type TournamentNameIsTakenError struct {
	TournamentName string
}

func (e TournamentNameIsTakenError) Error() string {
	return fmt.Sprintf("Tournament name: %s is taken by other tournament", e.TournamentName)
}

type EmptyTournamentIdError struct{}

func (e EmptyTournamentIdError) Error() string {
	return "Provided empty tournament ID"
}

type EmptyTiktokURLError struct{}

func (e EmptyTiktokURLError) Error() string {
	return "Provided empty Tiktok URL"
}

type NotAllowedTournamentTypeError struct {
	ContestType string
}

func (e NotAllowedTournamentTypeError) Error() string {
	return fmt.Sprintf("Provided not allowed tournament type: %s", e.ContestType)
}
