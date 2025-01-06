package keysdto

import "errors"

var (
	ErrInvalidQuery = errors.New("invalid query")
)

type KeysDTO struct {
	Query string `json:"query"`
}
