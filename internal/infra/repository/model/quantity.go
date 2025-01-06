package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type Quantity struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:quantities"`
	QuantityCommonAttributes
}

type QuantityCommonAttributes struct {
	Quantity   float64   `bun:"quantity,notnull"`
	CategoryID uuid.UUID `bun:"column:category_id,type:uuid,notnull"`
}
