package database

import (
	"github.com/google/uuid"
	"tiktok-arena/models"
)

func GetUserByName(username string) (models.User, error) {
	var user models.User
	record := usersTable.First(&user, "name = ?", username)
	return user, record.Error
}

func UserExists(username string) (bool, error) {
	var user models.User
	record := usersTable.Select("id").First(&user, "name = ?", username)
	return user.ID != nil, record.Error
}

func CreateNewUser(newUser *models.User) error {
	record := usersTable.Create(&newUser)
	return record.Error
}

func ChangeUserPhoto(url string, id uuid.UUID) error {
	record := usersTable.Where("id = ?", id).Update("photo_url", url)
	return record.Error
}

func GetUserPhoto(id string) (string, error) {
	var url string
	record := usersTable.Where("id = ?", id).Select("photo_url").Find(&url)
	return url, record.Error
}
