package jwtservice

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	userentity "github.com/willjrcom/sales-backend-go/internal/domain/user"
)

var secretKey = "sua_chave_secreta"

func CreateToken(user *userentity.UserCommonAttributes) (string, error) {

	claims := jwt.MapClaims{
		"id":      user.Email,
		"schemas": user.Schemas,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	}

	// Criar um token JWT usando a chave secreta
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
