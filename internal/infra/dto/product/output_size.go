package productdto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type SizeOutput struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Active bool      `json:"active"`
}

func (s *SizeOutput) FromModel(model *productentity.Size) {
	s.ID = model.ID
	s.Name = model.Name
	s.Active = model.Active
}
