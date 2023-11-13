package sizedto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type SizeNameOutput struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (s *SizeNameOutput) FromModel(model *productentity.Size) {
	s.ID = model.ID
	s.Name = model.Name
}
