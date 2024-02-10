package headerservice

import (
	"errors"
	"net/http"
)

func GetAccessTokenHeader(r *http.Request) (string, error) {
	accessToken := r.Header.Get("access-token")

	if accessToken == "" {
		return "", errors.New("access-token is required")
	}

	return accessToken, nil
}

func GetIDTokenHeader(r *http.Request) (string, error) {
	idToken := r.Header.Get("id-token")

	if idToken == "" {
		return "", errors.New("id-token is required")
	}

	return idToken, nil
}
