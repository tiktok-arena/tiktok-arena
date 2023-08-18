package dtos

import (
	"github.com/google/uuid"
	"tiktok-arena/internal/core/models"
)

type UsersResponse struct {
	UserCount int64         `validate:"required"`
	Users     []models.User `validate:"required"`
}

type ChangePhotoURL struct {
	PhotoURL string `validate:"required"`
}

type AuthInput struct {
	Name     string `validate:"required"`
	Password string `validate:"required"`
}

type RegisterDetails struct {
	ID    uuid.UUID
	Name  string
	Token string
}

type LoginDetails struct {
	ID       uuid.UUID
	Name     string
	Token    string
	PhotoURL string
}

type WhoAmI struct {
	ID       string
	Name     string
	Token    string
	PhotoURL string
}
