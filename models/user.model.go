package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID       *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name     string     `gorm:"not null;default:null"`
	Password string     `gorm:"not null;default:null"`
	PhotoURL string
}

type ChangePhotoURL struct {
	PhotoURL string `validate:"required"`
}

type AuthInput struct {
	Name     string `validate:"required"`
	Password string `validate:"required"`
}

type RegisterDetails struct {
	ID       string
	Username string
	Token    string
}

type LoginDetails struct {
	ID       string
	Username string
	Token    string
	PhotoURL string
}

type WhoAmI struct {
	ID       string
	Username string
	Token    string
	PhotoURL string
}
