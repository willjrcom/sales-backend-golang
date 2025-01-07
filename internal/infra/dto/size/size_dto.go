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

func (s *SizeDTO) FromDomain(size *productentity.Size) {
	if size == nil {
		return
	}
	*s = SizeDTO{
		ID:         size.ID,
		Name:       size.Name,
		IsActive:   size.IsActive,
		CategoryID: size.CategoryID,
	}
}
