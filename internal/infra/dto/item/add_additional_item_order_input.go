package itemdto

import (
	"errors"

	"github.com/google/uuid"
)

type AddAdditionalItemOrderInput struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}

func (a *AddAdditionalItemOrderInput) validate() error {
	if a.ProductID == uuid.Nil {
		return errors.New("product id is required")
	}

	if a.Quantity == 0 {
		return errors.New("quantity is required")
	}

	return nil
}

func (a *AddAdditionalItemOrderInput) ToModel() (productID uuid.UUID, quantity int, err error) {
	if err = a.validate(); err != nil {
		return
	}

	return a.ProductID, a.Quantity, nil
}
