package dtos

import "github.com/google/uuid"

type ChangePhotoURL struct {
	PhotoURL string `validate:"required"`
}

type AuthInput struct {
	Name     string `validate:"required"`
	Password string `validate:"required"`
}

type RegisterDetails struct {
	ID       *uuid.UUID
	Username string
	Token    string
}

type LoginDetails struct {
	ID       *uuid.UUID
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
