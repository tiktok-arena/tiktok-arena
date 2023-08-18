package dtos

import (
	"github.com/google/uuid"
	"tiktok-arena/internal/core/models"
)

type TournamentsResponse struct {
	TournamentCount int64               `validate:"required"`
	Tournaments     []models.Tournament `validate:"required"`
}

type TournamentsResponseWithUser struct {
	TournamentCount int64                   `validate:"required"`
	Tournaments     []TournamentWithoutUser `validate:"required"`
	User            models.User             `validate:"required"`
}

type TournamentWinner struct {
	TiktokURL string `validate:"required"`
}

type TournamentIds struct {
	TournamentIds []string `validate:"required"`
}

type CreateTournament struct {
	Name      string         `validate:"required"`
	PhotoURL  string         `validate:"required"`
	Size      int            `validate:"gte=4,lte=64"`
	Tiktoks   []CreateTiktok `validate:"required"`
	IsPrivate bool           `validate:"required"`
}

type EditTournament struct {
	Name      string         `validate:"required"`
	PhotoURL  string         `validate:"required"`
	Size      int            `validate:"gte=4,lte=64"`
	Tiktoks   []CreateTiktok `validate:"required"`
	IsPrivate bool           `validate:"required"`
}

type TournamentWithoutUser struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name        string    `gorm:"not null;default:null"`
	Size        int       `gorm:"not null"`
	TimesPlayed int       `gorm:"not null"`
	IsPrivate   bool      `gorm:"not null;default:false"`
	PhotoURL    string
}

type TournamentStats struct {
	TournamentId uuid.UUID
	TiktoksStats []TiktokStats
}
