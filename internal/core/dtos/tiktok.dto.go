package dtos

type CreateTiktok struct {
	Name string `validate:"required" json:"name"`
	URL  string `validate:"required" json:"url"`
}

type TiktokStats struct {
	Name string `gorm:"not null;default:null" json:"name"`
	URL  string `gorm:"not null;primaryKey;default:null" json:"url"`
	Wins int    `json:"wins"`
}
