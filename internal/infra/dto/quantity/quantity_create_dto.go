package quantitydto

import (
	"errors"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	ErrQuantityRequired = errors.New("quantity is required")
	ErrCategoryRequired = errors.New("category is required")
)

type QuantityCreateDTO struct {
	Quantity   float64   `json:"quantity"`
	CategoryID uuid.UUID `json:"category_id"`
}

func (s *QuantityCreateDTO) validate() error {
	if s.Quantity <= 0 {
		return ErrQuantityRequired
	}
	if s.CategoryID == uuid.Nil {
		return ErrCategoryRequired
	}

	return nil
}

func (s *QuantityCreateDTO) ToDomain() (*productentity.Quantity, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}

	quantityCommonAttributes := productentity.QuantityCommonAttributes{
		Quantity:   s.Quantity,
		CategoryID: s.CategoryID,
	}

	return productentity.NewQuantity(quantityCommonAttributes), nil
}
