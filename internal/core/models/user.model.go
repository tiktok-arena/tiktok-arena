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
