package headerservice

import (
	"errors"
	"net/http"
)

func GetAnyToken(r *http.Request) (string, error) {
	if accessToken, err := GetAccessTokenFromHeader(r); err == nil {
		return accessToken, nil
	}

	return GetIDTokenFromHeader(r)
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
