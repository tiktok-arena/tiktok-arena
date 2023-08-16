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
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByName(username string) (models.User, error) {
	var user models.User
	record := r.db.
		Where("name = ?", username).
		First(&user)
	return user, record.Error
}

func (r *UserRepository) UserExists(username string) (bool, error) {
	var user models.User
	record := r.db.
		Select("id").
		First(&user, "name = ?", username)
	if record.Error == gorm.ErrRecordNotFound {
		return false, nil
	}
	return user.ID != uuid.Nil, record.Error
}

func (r *UserRepository) CreateUser(newUser *models.User) error {
	record := r.db.
		Create(&newUser)
	return record.Error
}

func (r *UserRepository) ChangeUserPhoto(url string, id uuid.UUID) error {
	record := r.db.
		Model(&models.User{}).
		Where("id = ?", id).
		Update("photo_url", url)
	return record.Error
}

func (r *UserRepository) GetUserPhoto(id string) (string, error) {
	var url string
	record := r.db.
		Model(&models.User{}).
		Where("id = ?", id).
		Select("photo_url").
		Find(&url)
	return url, record.Error
}

func (r *UserRepository) GetUserByID(id uuid.UUID) (user models.User, err error) {
	err = r.db.
		Model(&models.User{}).
		Where("id = ?", id).
		Find(&user).Error
	return
}
