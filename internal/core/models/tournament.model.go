package models

import (
	"github.com/google/uuid"
)

type Tournament struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name        string    `gorm:"not null;default:null"`
	Size        int       `gorm:"not null"`
	TimesPlayed int       `gorm:"not null"`
	UserID      uuid.UUID `gorm:"not null"`
	User        *User     `gorm:"foreignKey:UserID"`
	PhotoURL    string
}
