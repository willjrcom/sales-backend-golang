package jwtservice

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dgrijalva/jwt-go"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
)

var secretKey = "sua_chave_secreta"

func CreateAccessToken(user *companyentity.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":                user.ID,
		"user":                   user,
		"available_user_schemas": user.GetSchemas(),
		"sub":                    "access-token",
		"exp":                    time.Now().UTC().Add(time.Minute * 5).Unix(),
	}

	// Criar um token JWT usando a chave secreta
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func CreateIDToken(accessToken *jwt.Token, schema string) (string, error) {
	oldClaims := accessToken.Claims.(jwt.MapClaims)

	claims := jwt.MapClaims{
		"user_id":        oldClaims["user_id"],
		"user":           oldClaims["user"],
		"current_schema": schema,
		"sub":            "id-token",
		"exp":            time.Now().UTC().Add(time.Hour * 6).Unix(),
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

func GetSchemaFromToken(token *jwt.Token) string {
	return token.Claims.(jwt.MapClaims)["current_schema"].(string)
}

func GetUserFromToken(token *jwt.Token) *companyentity.User {
	user := &companyentity.User{}
	claims := token.Claims.(jwt.MapClaims)

	userMap, ok := claims["user"]
	if !ok {
		return nil
	}

	userJson, err := json.Marshal(userMap) // .(interface{})
	if err != nil {
		return nil
	}

	if err := json.Unmarshal(userJson, user); err != nil {
		return nil
	}

	return user
}
