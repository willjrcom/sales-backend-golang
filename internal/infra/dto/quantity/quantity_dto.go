package quantitydto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type QuantityDTO struct {
	ID         uuid.UUID `json:"id"`
	Quantity   float64   `json:"quantity"`
	CategoryID uuid.UUID `json:"category_id"`
}

func (s *QuantityDTO) FromDomain(model *productentity.Quantity) {
	*s = QuantityDTO{
		ID:         model.ID,
		Quantity:   model.Quantity,
		CategoryID: model.CategoryID,
	}
}
