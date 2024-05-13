package productcategoryquantitydto

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

type RegisterQuantityInput struct {
	productentity.QuantityCommonAttributes
}

func (s *RegisterQuantityInput) validate() error {
	if s.Quantity <= 0 {
		return ErrQuantityRequired
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

	quantityCommonAttributes := productentity.QuantityCommonAttributes{
		Quantity:   s.Quantity,
		CategoryID: s.CategoryID,
	}

	return &productentity.Quantity{
		Entity:                   entity.NewEntity(),
		QuantityCommonAttributes: quantityCommonAttributes,
	}, nil
}
