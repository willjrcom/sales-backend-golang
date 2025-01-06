package sizedto

import (
	"errors"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	ErrNameRequired     = errors.New("name is required")
	ErrCategoryRequired = errors.New("category is required")
)

type SizeCreateDTO struct {
	Name       string    `json:"name"`
	IsActive   *bool     `json:"is_active"`
	CategoryID uuid.UUID `json:"category_id"`
}

func (s *SizeCreateDTO) validate() error {
	if s.Name == "" {
		return ErrNameRequired
	}
	if s.CategoryID == uuid.Nil {
		return ErrCategoryRequired
	}
	if s.IsActive == nil {
		s.IsActive = new(bool)
		*s.IsActive = true
	}

	return nil
}

func (s *SizeCreateDTO) ToModel() (*productentity.Size, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}

	sizeCommonAttributes := productentity.SizeCommonAttributes{
		Name:       s.Name,
		IsActive:   s.IsActive,
		CategoryID: s.CategoryID,
	}

	return &productentity.Size{
		Entity:               entity.NewEntity(),
		SizeCommonAttributes: sizeCommonAttributes,
	}, nil
}
