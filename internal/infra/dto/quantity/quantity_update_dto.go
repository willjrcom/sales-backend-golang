package quantitydto

import (
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type QuantityUpdateDTO struct {
	Quantity *float64 `json:"quantity"`
	IsActive *bool    `json:"is_active"`
}

func (s *QuantityUpdateDTO) validate() error {
	if s.Quantity != nil && *s.Quantity <= 0 {
		return ErrQuantityRequired
	}

	return nil
}

func (s *QuantityUpdateDTO) UpdateDomain(quantity *productentity.Quantity) (err error) {
	if err = s.validate(); err != nil {
		return err
	}

	if s.Quantity != nil {
		quantity.Quantity = *s.Quantity
	}

	if s.IsActive != nil {
		quantity.IsActive = *s.IsActive
	}

	return nil
}
