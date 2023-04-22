package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"tiktok-arena/internal/core/models"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db.Model(&models.User{})}
}

func (r *UserRepository) GetUserByName(username string) (models.User, error) {
	var user models.User
	record := r.db.
		First(&user, "name = ?", username)
	return user, record.Error
}

func (r *UserRepository) UserExists(username string) (bool, error) {
	var user models.User
	record := r.db.
		Select("id").
		First(&user, "name = ?", username)
	return user.ID != nil, record.Error
}

func (r *UserRepository) CreateUser(newUser *models.User) error {
	record := r.db.
		Create(&newUser)
	return record.Error
}

func (r *UserRepository) ChangeUserPhoto(url string, id uuid.UUID) error {
	record := r.db.
		Where("id = ?", id).
		Update("photo_url", url)
	return record.Error
}

func (r *UserRepository) GetUserPhoto(id string) (string, error) {
	var url string
	record := r.db.
		Where("id = ?", id).
		Select("photo_url").
		Find(&url)
	return url, record.Error
}
