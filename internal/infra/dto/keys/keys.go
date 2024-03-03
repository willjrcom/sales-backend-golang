package keysdto

import "errors"

var (
	ErrInvalidQuery = errors.New("invalid query")
)

type KeysInput struct {
	Query  string `json:"query"`
	Ddd    string `json:"ddd"`
	Number string `json:"number"`
}
