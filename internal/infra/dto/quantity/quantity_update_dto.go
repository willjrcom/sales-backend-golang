package quantitydto

import (
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type QuantityUpdateDTO struct {
	Quantity *float64 `json:"quantity"`
}

func (s *QuantityUpdateDTO) validate() error {
	if s.Quantity != nil && *s.Quantity <= 0 {
		return ErrQuantityRequired
	}

	return nil
}

func (s *QuantityUpdateDTO) UpdateDomain(model *productentity.Quantity) (err error) {
	if err = s.validate(); err != nil {
		return err
	}

	if s.Quantity != nil {
		model.Quantity = *s.Quantity
	}

	return nil
}
