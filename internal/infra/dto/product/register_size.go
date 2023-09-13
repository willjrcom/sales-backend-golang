package productdto

import (
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type RegisterSizeInput struct {
	Name       string    `json:"name"`
	Active     *bool     `json:"active"`
	CategoryID uuid.UUID `json:"category_id"`
}

func (s *RegisterSizeInput) validate() error {
	if s.Name == "" {
		return ErrNameIsEmpty
	}
	if s.CategoryID == uuid.Nil {
		return ErrCategoryRequired
	}
	if s.Active == nil {
		s.Active = new(bool)
		*s.Active = true
	}

	return nil
}

func (s *RegisterSizeInput) ToModel() (*productentity.Size, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}

	return &productentity.Size{
		Entity:     entity.NewEntity(),
		Name:       s.Name,
		Active:     *s.Active,
		CategoryID: s.CategoryID,
	}, nil
}
