package model

import (
	"time"

	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Person struct {
	entity.Entity
	PersonCommonAttributes
}

type PersonCommonAttributes struct {
	Name     string                 `bun:"name,notnull"`
	Email    string                 `bun:"email"`
	Cpf      string                 `bun:"cpf"`
	Birthday *time.Time             `bun:"birthday"`
	Contact  *Contact               `bun:"rel:has-one,join:id=object_id,notnull"`
	Address  *addressentity.Address `bun:"rel:has-one,join:id=object_id,notnull"`
}
