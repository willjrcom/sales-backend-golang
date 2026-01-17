package quantitydto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type QuantityDTO struct {
	ID         uuid.UUID `json:"id"`
	Quantity   float64   `json:"quantity"`
	IsActive   bool      `json:"is_active"`
	CategoryID uuid.UUID `json:"category_id"`
}

func (s *QuantityDTO) FromDomain(quantity *productentity.Quantity) {
	if quantity == nil {
		return
	}
	*s = QuantityDTO{
		ID:         quantity.ID,
		Quantity:   quantity.Quantity,
		IsActive:   quantity.IsActive,
		CategoryID: quantity.CategoryID,
	}
}
