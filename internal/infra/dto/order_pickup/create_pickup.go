package orderpickupdto

import (
	"errors"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

var (
	ErrNameRequired = errors.New("name is required")
)

type CreateOrderPickupInput struct {
	Name string `json:"name"`
}

func (o *CreateOrderPickupInput) validate() error {
	if o.Name == "" {
		return ErrNameRequired
	}

	return nil
}

func (o *CreateOrderPickupInput) ToModel() (*orderentity.OrderPickup, error) {
	if err := o.validate(); err != nil {
		return nil, err
	}

	return orderentity.NewOrderPickup(o.Name), nil
}
