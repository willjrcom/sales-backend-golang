package orderpickupdto

import (
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type UpdateOrderPickupInput struct {
	Name string `json:"name"`
}

func (o *UpdateOrderPickupInput) validate() error {
	if o.Name == "" {
		return ErrNameRequired
	}

	return nil
}

func (o *UpdateOrderPickupInput) ToModel() (*orderentity.OrderPickup, error) {
	if err := o.validate(); err != nil {
		return nil, err
	}

	return orderentity.NewOrderPickup(o.Name), nil
}
