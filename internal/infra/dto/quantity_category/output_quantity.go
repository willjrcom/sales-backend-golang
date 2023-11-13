package quantitydto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type QuantityOutput struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Active bool      `json:"active"`
}

func (s *QuantityOutput) FromModel(model *productentity.Quantity) {
	s.ID = model.ID
	s.Name = model.Name
}
