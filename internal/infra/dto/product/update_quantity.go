package productdto

import (
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type UpdateQuantityInput struct {
	Name *string `json:"name"`
}

func (s *UpdateQuantityInput) validate() error {
	if s.Name == nil {
		return ErrNameIsEmpty
	}

	return nil
}
func (s *UpdateQuantityInput) UpdateModel(model *productentity.Quantity) (err error) {
	if err = s.validate(); err != nil {
		return err
	}

	if s.Name != nil {
		model.Name = *s.Name
	}

	return nil
}
