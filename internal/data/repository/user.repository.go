package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"tiktok-arena/internal/core/dtos"
	"tiktok-arena/internal/core/models"
	"tiktok-arena/internal/data/repository/scopes"
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

func (r *UserRepository) TotalUsers() (int64, error) {
	var totalUsers int64
	record := r.db.
		Model(&models.User{}).
		Count(&totalUsers)
	return totalUsers, record.Error
}

func (r *UserRepository) GetAllUsers(totalUsers int64, queries dtos.PaginationQueries) (dtos.UsersResponse, error) {
	var users []models.User
	record := r.db.
		Scopes(scopes.Search(queries.SearchText)).
		Scopes(scopes.Paginate(queries.Page, queries.Count)).
		Find(&users)
	return dtos.UsersResponse{UserCount: totalUsers, Users: users}, record.Error
}
