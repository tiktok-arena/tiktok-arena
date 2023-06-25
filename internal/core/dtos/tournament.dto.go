package dtos

import (
	"tiktok-arena/internal/core/models"
)

type TournamentsResponse struct {
	TournamentCount int64               `validate:"required"`
	Tournaments     []models.Tournament `validate:"required"`
}

type TournamentWinner struct {
	TiktokURL string `validate:"required"`
}

type TournamentIds struct {
	TournamentIds []string `validate:"required"`
}

type CreateTournament struct {
	Name     string         `validate:"required"`
	PhotoURL string         `validate:"required"`
	Size     int            `validate:"gte=4,lte=64"`
	Tiktoks  []CreateTiktok `validate:"required"`
}

type EditTournament struct {
	Name     string         `validate:"required"`
	PhotoURL string         `validate:"required"`
	Size     int            `validate:"gte=4,lte=64"`
	Tiktoks  []CreateTiktok `validate:"required"`
}
