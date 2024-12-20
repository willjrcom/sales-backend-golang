package quantitydto

import (
	"errors"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	ErrQuantityRequired = errors.New("quantity is required")
	ErrCategoryRequired = errors.New("category is required")
)

type CreateQuantityInput struct {
	Quantity   float64   `json:"quantity"`
	CategoryID uuid.UUID `json:"category_id"`
}

func (s *CreateQuantityInput) validate() error {
	if s.Quantity <= 0 {
		return ErrQuantityRequired
	}
	if s.CategoryID == uuid.Nil {
		return ErrCategoryRequired
	}

	return nil
}

func (s *CreateQuantityInput) ToModel() (*productentity.Quantity, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}

	quantityCommonAttributes := productentity.QuantityCommonAttributes{
		Quantity:   s.Quantity,
		CategoryID: s.CategoryID,
	}

	return &productentity.Quantity{
		Entity:                   entity.NewEntity(),
		QuantityCommonAttributes: quantityCommonAttributes,
	}, nil
}
