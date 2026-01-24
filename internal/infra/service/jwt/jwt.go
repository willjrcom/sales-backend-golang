package jwtservice

import (
	"context"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
)

func CreateBasicAccessToken(user *companyentity.User) (string, error) {
	userID := user.ID.String()
	claims := jwt.MapClaims{
		"user_id": userID,
		"sub":     "access-token",
		"exp":     time.Now().UTC().Add(30 * time.Minute).Unix(),
	}

	// Criar um token JWT usando a chave secreta
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := os.Getenv("JWT_SECRET_KEY")
	return token.SignedString([]byte(secretKey))
}

func CreateFullAccessToken(accessToken *jwt.Token, schema string) (string, error) {
	oldClaims := accessToken.Claims.(jwt.MapClaims)

	claims := jwt.MapClaims{
		"user_id":        oldClaims["user_id"],
		"current_schema": schema,
		"sub":            "access-token",
		"exp":            time.Now().UTC().Add(2 * time.Hour).Unix(),
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

func ValidateTokenWithoutExpiry(ctx context.Context, tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		secretKey := os.Getenv("JWT_SECRET_KEY")
		return []byte(secretKey), nil
	})
}

func GetSchemaFromAccessToken(token *jwt.Token) string {
	claims := token.Claims.(jwt.MapClaims)
	currentSchema, ok := claims["current_schema"]

	if !ok {
		return ""
	}

	return currentSchema.(string)
}

func GetUserIDFromToken(token *jwt.Token) string {
	claims := token.Claims.(jwt.MapClaims)

	userID, ok := claims["user_id"]
	if !ok {
		return ""
	}

	return userID.(string)
}

// Cria um JWT para redefinição de senha, contendo apenas o e-mail do usuário
func CreatePasswordResetToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"sub":   "password-reset",
		"exp":   time.Now().UTC().Add(30 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := os.Getenv("JWT_SECRET_KEY")
	return token.SignedString([]byte(secretKey))
}

// Valida o JWT de redefinição de senha e retorna o e-mail se válido
func ValidatePasswordResetToken(ctx context.Context, tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		secretKey := os.Getenv("JWT_SECRET_KEY")
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email, ok := claims["email"].(string)
		if !ok {
			return "", jwt.ErrInvalidKey
		}
		return email, nil
	}
	return "", jwt.ErrInvalidKey
}
