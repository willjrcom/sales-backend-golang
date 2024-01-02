package quantitydto

import (
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type UpdateQuantityInput struct {
	productentity.PatchQuantity
}

func (s *UpdateQuantityInput) validate() error {
	if s.Quantity != nil && *s.Quantity <= 0 {
		return ErrQuantityRequired
	}

	return nil
}
func (s *UpdateQuantityInput) UpdateModel(model *productentity.Quantity) (err error) {
	if err = s.validate(); err != nil {
		return err
	}

	if *s.Quantity > 0.0 {
		model.Quantity = *s.Quantity
	}

	return nil
}
