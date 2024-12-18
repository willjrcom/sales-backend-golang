package sizedto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type SizeOutput struct {
	ID uuid.UUID `json:"id"`
	productentity.SizeCommonAttributes
}

func (s *SizeOutput) FromModel(model *productentity.Size) {
	s.ID = model.ID
	s.SizeCommonAttributes = model.SizeCommonAttributes
}
