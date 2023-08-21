package models

import (
	"github.com/google/uuid"
)

type Tournament struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name        string    `gorm:"not null;default:null" json:"name"`
	Size        int       `gorm:"not null" json:"size"`
	TimesPlayed int       `gorm:"not null" json:"timesPlayed"`
	UserID      uuid.UUID `gorm:"not null"  json:"userID"`
	User        User      `gorm:"foreignKey:UserID" json:"user"`
	IsPrivate   bool      `gorm:"not null;default:false" json:"isPrivate"`
	PhotoURL    string    `json:"photoURL"`
}
