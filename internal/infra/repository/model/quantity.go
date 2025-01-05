package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrQuantityAlreadyExists = errors.New("quantity already exists")
)

type Quantity struct {
	entity.Entity
	bun.BaseModel `bun:"table:quantities"`
	QuantityCommonAttributes
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

type QuantityCommonAttributes struct {
	Quantity   float64   `bun:"quantity,notnull"`
	CategoryID uuid.UUID `bun:"column:category_id,type:uuid,notnull"`
}
