package persondto

import (
	"errors"
)

var (
	ErrContactInvalid = errors.New("contact format invalid")
)

type Contact struct {
	Ddd    string `bun:"ddd,notnull"`
	Number string `bun:"number,notnull"`
}
