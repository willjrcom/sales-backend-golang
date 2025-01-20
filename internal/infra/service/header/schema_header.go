package headerservice

import (
	"context"
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	jwtservice "github.com/willjrcom/sales-backend-go/internal/infra/service/jwt"
)

func GetAnyValidToken(ctx context.Context, r *http.Request) (*jwt.Token, error) {
	if idToken, err := GetIDTokenFromHeader(r); err == nil {
		validToken, err := jwtservice.ValidateToken(ctx, idToken)

		if err == nil {
			return validToken, nil
		}
	}

	accessToken, err := GetAccessTokenFromHeader(r)
	if err != nil {
		return nil, err
	}

	validToken, err := jwtservice.ValidateToken(ctx, accessToken)

	if err != nil {
		return nil, err
	}

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
