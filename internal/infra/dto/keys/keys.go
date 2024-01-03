package keysdto

import "errors"

var (
	ErrInvalidQuery = errors.New("invalid query")
)

type KeysInput struct {
	Query string `json:"query"`
}
