package dtos

type CreateTiktok struct {
	Name string `validate:"required"`
	URL  string `validate:"required"`
}

type TiktokStats struct {
	Name string `gorm:"not null;default:null"`
	URL  string `gorm:"not null;primaryKey;default:null"`
	Wins int
}
