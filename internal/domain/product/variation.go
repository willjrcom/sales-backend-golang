package productentity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type ProductVariation struct {
	entity.Entity
	ProductID   uuid.UUID
	SizeID      uuid.UUID
	Size        *Size
	Price       decimal.Decimal
	Cost        decimal.Decimal
	IsAvailable bool
}

func NewProductVariation(productID uuid.UUID, sizeID uuid.UUID, price decimal.Decimal, cost decimal.Decimal) *ProductVariation {
	return &ProductVariation{
		Entity:      entity.NewEntity(),
		ProductID:   productID,
		SizeID:      sizeID,
		Price:       price,
		Cost:        cost,
		IsAvailable: true,
	}
}
