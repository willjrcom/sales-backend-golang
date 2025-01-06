package model

import (
	"errors"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

var (
	ErrContactInvalid = errors.New("contact format invalid")
)

type Contact struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:contacts"`
	ContactCommonAttributes
	ObjectID uuid.UUID `bun:"object_id,type:uuid,notnull"`
}

type ContactCommonAttributes struct {
	Ddd    string `bun:"ddd,notnull"`
	Number string `bun:"number,notnull"`
	Type   string `bun:"type,notnull"`
}
