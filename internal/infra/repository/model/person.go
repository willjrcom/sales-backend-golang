package model

import (
	"time"
)

type Person struct {
	Name     string     `bun:"name,notnull"`
	Email    string     `bun:"email"`
	Cpf      string     `bun:"cpf"`
	Birthday *time.Time `bun:"birthday"`
	Contact  *Contact   `bun:"rel:has-one,join:id=object_id,notnull"`
	Address  *Address   `bun:"rel:has-one,join:id=object_id,notnull"`
}
