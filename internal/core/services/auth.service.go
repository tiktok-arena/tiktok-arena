package services

import (
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"tiktok-arena/configuration"
	"tiktok-arena/internal/core/dtos"
	"tiktok-arena/internal/core/models"
	"tiktok-arena/internal/core/validator"
	"time"
)

type AuthServiceUserRepository interface {
	GetUserByName(username string) (models.User, error)
	UserExists(username string) (bool, error)
	CreateUser(newUser *models.User) error
	GetUserPhoto(id string) (string, error)
}

type AuthService struct {
	UserRepository AuthServiceUserRepository
}

func NewAuthService(userRepository AuthServiceUserRepository) *AuthService {
	return &AuthService{UserRepository: userRepository}
}

func (s *AuthService) NewUser(auth *dtos.AuthInput) (details dtos.RegisterDetails, err error) {
	err = validator.ValidateStruct(auth)
	if err != nil {
		return details, ValidateError{err}
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(auth.Password), bcrypt.DefaultCost)
	if err != nil {
		return details, BcryptError{err}
	}
	exists, err := s.UserRepository.UserExists(auth.Name)
	if err != nil {
		return details, RepositoryError{err}
	}
	if exists {
		return details, UserAlreadyExistsError{auth.Name}
	}

	newUser := models.User{
		Name:     auth.Name,
		Password: string(hashedPassword),
	}
	err = s.UserRepository.CreateUser(&newUser)
	if err != nil {
		return details, RepositoryError{err}
	}

	token, err := UserJwtToken(&newUser)
	if err != nil {
		return details, JWTGenerateError{err}
	}

	return dtos.RegisterDetails{
		ID:       newUser.ID,
		Username: newUser.Name,
		Token:    token,
	}, err
}

func (s *AuthService) GetUserByNameAndPassword(input *dtos.AuthInput) (details dtos.LoginDetails, err error) {
	err = validator.ValidateStruct(input)
	if err != nil {
		return details, ValidateError{err}
	}

	user, err := s.UserRepository.GetUserByName(input.Name)
	if err != nil {
		return details, RepositoryError{err}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return details, BcryptError{err}
	}

	token, err := UserJwtToken(&user)

	if err != nil {
		return details, JWTGenerateError{err}
	}

	url, err := s.UserRepository.GetUserPhoto(user.ID.String())
	if err != nil {
		return details, RepositoryError{err}
	}
	return dtos.LoginDetails{
		ID:       user.ID.String(),
		Username: user.Name,
		Token:    token,
		PhotoURL: url,
	}, err

}

func (s *AuthService) WhoAmI(token *jwt.Token) (whoami dtos.WhoAmI, err error) {
	claims := token.Claims.(jwt.MapClaims)

	username := claims["name"].(string)
	id := claims["sub"].(string)

	exists, err := s.UserRepository.UserExists(username)
	if err != nil {
		return whoami, RepositoryError{err}
	}
	if !exists {
		return whoami, UserNotExistsError{Username: username}
	}
	url, err := s.UserRepository.GetUserPhoto(id)
	if err != nil {
		return whoami, RepositoryError{err}
	}
	return dtos.WhoAmI{
		ID:       id,
		Username: username,
		Token:    token.Raw,
		PhotoURL: url,
	}, err
}

func UserJwtToken(user *models.User) (string, error) {
	now := time.Now().UTC()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"name": user.Name,
		"exp":  now.Add(configuration.EnvConfig.JwtExpiresIn).Unix(),
		"iat":  now.Unix(),
		"nbf":  now.Unix(),
	})

	return token.SignedString([]byte(configuration.EnvConfig.JwtSecret))

}
