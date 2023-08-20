package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name     string    `gorm:"not null;default:null" json:"name"`
	Password string    `gorm:"not null;default:null" json:"password"`
	PhotoURL string    `json:"photoURL"`
}
