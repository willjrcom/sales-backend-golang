package sizedto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type SizeDTO struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	IsActive   *bool     `json:"is_active"`
	CategoryID uuid.UUID `json:"category_id"`
}

func (s *SizeDTO) FromDomain(model *productentity.Size) {
	*s = SizeDTO{
		ID:         model.ID,
		Name:       model.Name,
		IsActive:   model.IsActive,
		CategoryID: model.CategoryID,
	}
}
