package jwtservice

import (
	"context"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
)

func CreateAccessToken(user *companyentity.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":                user.ID,
		"available_user_schemas": user.GetSchemas(),
		"sub":                    "access-token",
		"exp":                    time.Now().UTC().Add(time.Minute * 5).Unix(),
	}

	// Criar um token JWT usando a chave secreta
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := os.Getenv("JWT_SECRET_KEY")
	return token.SignedString([]byte(secretKey))
}

func CreateIDToken(accessToken *jwt.Token, schema string) (string, error) {
	oldClaims := accessToken.Claims.(jwt.MapClaims)

	claims := jwt.MapClaims{
		"user_id":        oldClaims["user_id"],
		"current_schema": schema,
		"sub":            "id-token",
		"exp":            time.Now().UTC().Add(time.Hour * 6).Unix(),
	}

	// Criar um token JWT usando a chave secreta
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := os.Getenv("JWT_SECRET_KEY")

	return token.SignedString([]byte(secretKey))
}

func ValidateToken(ctx context.Context, tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		secretKey := os.Getenv("JWT_SECRET_KEY")
		return []byte(secretKey), nil
	})
}

func GetSchemasFromToken(token *jwt.Token) []interface{} {
	return token.Claims.(jwt.MapClaims)["available_user_schemas"].([]interface{})
}

func GetSchemaFromToken(token *jwt.Token) string {
	return token.Claims.(jwt.MapClaims)["current_schema"].(string)
}

func GetUserIDFromToken(token *jwt.Token) *uuid.UUID {
	claims := token.Claims.(jwt.MapClaims)

	userID, ok := claims["user_id"]
	if !ok {
		return nil
	}

	userIDString := userID.(string)

	userIDUUID := uuid.MustParse(userIDString)
	return &userIDUUID
}
