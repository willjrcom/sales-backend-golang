package headerservice

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	jwtservice "github.com/willjrcom/sales-backend-go/internal/infra/service/jwt"
)

func GetAccessTokenFromHeader(ctx context.Context, r *http.Request) (*jwt.Token, error) {
	accessToken := r.Header.Get("access-token")

	if accessToken == "" {
		return nil, errors.New("access-token is required")
	}

	validToken, err := jwtservice.ValidateToken(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("access-token invalid: %v", err)
	}

	return validToken, nil
}
