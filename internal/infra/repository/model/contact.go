package model

import (
	"errors"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrContactInvalid = errors.New("contact format invalid")
)

type Contact struct {
	entity.Entity
	bun.BaseModel `bun:"table:contacts"`
	ContactCommonAttributes
	ObjectID uuid.UUID `bun:"object_id,type:uuid,notnull"`
}

type ContactCommonAttributes struct {
	Ddd    string      `bun:"ddd,notnull"`
	Number string      `bun:"number,notnull"`
	Type   ContactType `bun:"type,notnull"`
}

type ContactType string

const (
	ContactTypeClient   ContactType = "Client"
	ContactTypeEmployee ContactType = "Employee"
)
