package itemdto

import (
	"errors"

	"github.com/google/uuid"
)

type AddAdditionalItemOrderInput struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  uuid.UUID `json:"quantity_id"`
}

func (a *AddAdditionalItemOrderInput) validate() error {
	if a.ProductID == uuid.Nil {
		return errors.New("product id is required")
	}

	if a.Quantity == uuid.Nil {
		return errors.New("quantity id is required")
	}

	return nil
}

func (a *AddAdditionalItemOrderInput) ToModel() (productID uuid.UUID, quantity uuid.UUID, err error) {
	if err = a.validate(); err != nil {
		return
	}

	return a.ProductID, a.Quantity, nil
}
