package validator

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func GetUserIdAndCheckJWT(user interface{}) (uuid.UUID, error) {
	if user == nil {
		return uuid.UUID{}, EmptyJWTError{}
	}
	userJWT := user.(*jwt.Token)

	claims := userJWT.Claims.(jwt.MapClaims)

	return uuid.Parse(claims["sub"].(string))
}
