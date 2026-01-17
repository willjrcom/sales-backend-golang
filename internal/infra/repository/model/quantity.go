package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type Quantity struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:quantities"`
	QuantityCommonAttributes
}

type QuantityCommonAttributes struct {
	Quantity   float64   `bun:"quantity,notnull"`
	IsActive   bool      `bun:"column:is_active,type:boolean"`
	CategoryID uuid.UUID `bun:"column:category_id,type:uuid,notnull"`
}

func (q *Quantity) FromDomain(model *productentity.Quantity) {
	if model == nil {
		return
	}
	*q = Quantity{
		Entity: entitymodel.FromDomain(model.Entity),
		QuantityCommonAttributes: QuantityCommonAttributes{
			Quantity:   model.Quantity,
			IsActive:   model.IsActive,
			CategoryID: model.CategoryID,
		},
	}
}

func (q *Quantity) ToDomain() *productentity.Quantity {
	if q == nil {
		return nil
	}
	return &productentity.Quantity{
		Entity: q.Entity.ToDomain(),
		QuantityCommonAttributes: productentity.QuantityCommonAttributes{
			Quantity:   q.Quantity,
			IsActive:   q.IsActive,
			CategoryID: q.CategoryID,
		},
	}
}
