package dtos

import (
	"github.com/google/uuid"
	"tiktok-arena/internal/core/models"
)

type UsersResponse struct {
	UserCount int64         `validate:"required" json:"userCount"`
	Users     []models.User `validate:"required" json:"users"`
}

type ChangePhotoURL struct {
	PhotoURL string `validate:"required" json:"photoURL"`
}

type AuthInput struct {
	Name     string `validate:"required" json:"name"`
	Password string `validate:"required" json:"password"`
}

type RegisterDetails struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Token string    `json:"token"`
}

type LoginDetails struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Token    string    `json:"token"`
	PhotoURL string    `json:"photoURL"`
}

type WhoAmI struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Token    string `json:"token"`
	PhotoURL string `json:"photoURL"`
}
