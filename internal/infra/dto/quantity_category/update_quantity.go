package quantitydto

import (
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type UpdateQuantityInput struct {
	Quantity float64 `json:"quantity"`
}

func (s *UpdateQuantityInput) validate() error {
	if s.Quantity <= 0 {
		return ErrQuantityRequired
	}

	return nil
}
func (s *UpdateQuantityInput) UpdateModel(model *productentity.Quantity) (err error) {
	if err = s.validate(); err != nil {
		return err
	}

	if s.Quantity != 0 {
		model.Quantity = s.Quantity
	}

	return nil
}
