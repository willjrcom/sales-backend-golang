package headerservice

import (
	"context"
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	jwtservice "github.com/willjrcom/sales-backend-go/internal/infra/service/jwt"
)

func GetAnyValidToken(ctx context.Context, r *http.Request) (*jwt.Token, error) {
	// Tenta primeiro com id-token
	if idToken, err := GetIDTokenFromHeader(r); err == nil {
		// fmt.Println("Tentando validar ID token...")
		validToken, err := jwtservice.ValidateToken(ctx, idToken)

		if err == nil {
			// fmt.Println("ID token validado com sucesso")
			return validToken, nil
		}
		// fmt.Println("ID token inválido, tentando access token:", err.Error())
	}

	// Se não tem id-token ou falhou, tenta access-token
	accessToken, err := GetAccessTokenFromHeader(r)
	if err != nil {
		// fmt.Println("Erro ao obter access token:", err.Error())
		return nil, err
	}

	// fmt.Println("Validando access token...")
	validToken, err := jwtservice.ValidateToken(ctx, accessToken)

	if err != nil {
		// fmt.Println("Erro ao validar access token:", err.Error())
		return nil, err
	}

	// fmt.Println("Access token validado com sucesso")
	return validToken, nil
}

func GetAccessTokenFromHeader(r *http.Request) (string, error) {
	accessToken := r.Header.Get("access-token")

	if accessToken == "" {
		return "", errors.New("access-token is required")
	}

	return accessToken, nil
}

func GetIDTokenFromHeader(r *http.Request) (string, error) {
	idToken := r.Header.Get("id-token")

	if idToken == "" {
		return "", errors.New("id-token is required")
	}

	return idToken, nil
}
