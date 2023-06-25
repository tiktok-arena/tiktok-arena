package validator

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func GetUserIdAndCheckJWT(user interface{}) (id uuid.UUID, err error) {
	if user == nil {
		return id, EmptyJWTError{}
	}
	userJWT := user.(*jwt.Token)

	claims := userJWT.Claims.(jwt.MapClaims)

	id, err = uuid.Parse(claims["sub"].(string))
	return
}
