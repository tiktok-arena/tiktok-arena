package dtos

import (
	"github.com/google/uuid"
	"tiktok-arena/internal/core/models"
)

type TournamentsResponse struct {
	TournamentCount int64               `validate:"required" json:"tournamentCount"`
	Tournaments     []models.Tournament `validate:"required" json:"tournaments"`
}

type TournamentsResponseWithUser struct {
	TournamentCount int64                   `validate:"required" json:"tournamentCount"`
	Tournaments     []TournamentWithoutUser `validate:"required" json:"tournaments"`
	User            models.User             `validate:"required" json:"user"`
}

type TournamentWinner struct {
	TiktokURL string `validate:"required" json:"tiktokURL"`
}

type TournamentIds struct {
	TournamentIds []string `validate:"required" json:"tournamentIds"`
}

type CreateTournament struct {
	Name      string         `validate:"required" json:"name"`
	PhotoURL  string         `validate:"required" json:"photoURL"`
	Size      int            `validate:"gte=4,lte=64" json:"size"`
	Tiktoks   []CreateTiktok `validate:"required" json:"tiktoks"`
	IsPrivate bool           `json:"isPrivate"` // by default public, so we don't need this field to be required
}

type EditTournament struct {
	Name      string         `validate:"required" json:"name"`
	PhotoURL  string         `validate:"required" json:"photoURL"`
	Size      int            `validate:"gte=4,lte=64" json:"size"`
	Tiktoks   []CreateTiktok `validate:"required" json:"tiktoks"`
	IsPrivate bool           `json:"isPrivate"` // by default public, so we don't need this field to be required
}

type TournamentWithoutUser struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name        string    `gorm:"not null;default:null" json:"name"`
	Size        int       `gorm:"not null" json:"size"`
	TimesPlayed int       `gorm:"not null" json:"timesPlayed"`
	IsPrivate   bool      `gorm:"not null;default:false" json:"isPrivate"`
	PhotoURL    string    `json:"photoURL"`
}

type TournamentStats struct {
	TournamentId uuid.UUID     `json:"tournamentId"`
	TiktoksStats []TiktokStats `json:"tiktoksStats"`
}
