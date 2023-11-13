package quantitydto

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

type RegisterQuantityInput struct {
	Name       string    `json:"name"`
	CategoryID uuid.UUID `json:"category_id"`
}

func (s *RegisterQuantityInput) validate() error {
	if s.Name == "" {
		return ErrNameRequired
	}
	if s.CategoryID == uuid.Nil {
		return ErrCategoryRequired
	}

	return nil
}

func (s *RegisterQuantityInput) ToModel() (*productentity.Quantity, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}

	return &productentity.Quantity{
		Entity:     entity.NewEntity(),
		Name:       s.Name,
		CategoryID: s.CategoryID,
	}, nil
}
