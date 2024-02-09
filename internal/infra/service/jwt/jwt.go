package jwtservice

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	userentity "github.com/willjrcom/sales-backend-go/internal/domain/user"
)

var secretKey = "sua_chave_secreta"

func CreateAccessToken(user *userentity.User) (string, error) {

	claims := jwt.MapClaims{
		"user_id":                user.ID,
		"user_email":             user.Email,
		"available_user_schemas": user.Schemas,
		"sub":                    "access-token",
		"exp":                    time.Now().Add(time.Minute * 5).Unix(),
	}

	// Criar um token JWT usando a chave secreta
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func CreateIDToken(accessToken *jwt.Token, schema string) (string, error) {

	oldClaims := accessToken.Claims.(jwt.MapClaims)

	claims := jwt.MapClaims{
		"user_id":        oldClaims["user_id"],
		"user_email":     oldClaims["user_email"],
		"current_schema": schema,
		"sub":            "id-token",
		"exp":            time.Now().Add(time.Hour * 6).Unix(),
	}

	// Criar um token JWT usando a chave secreta
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func ValidateToken(ctx context.Context, tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secretKey), nil
	})
}

func GetSchemasFromToken(token *jwt.Token) []interface{} {
	return token.Claims.(jwt.MapClaims)["available_user_schemas"].([]interface{})
}
