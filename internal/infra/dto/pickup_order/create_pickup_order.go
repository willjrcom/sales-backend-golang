package pickuporderdto

import (
	"errors"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

var (
	ErrNameRequired = errors.New("name is required")
)

type CreatePickupOrderInput struct {
	Name string `json:"name"`
}

func (o *CreatePickupOrderInput) validate() error {
	if o.Name == "" {
		return ErrNameRequired
	}

	return nil
}

func (o *CreatePickupOrderInput) ToModel() (*orderentity.PickupOrder, error) {
	if err := o.validate(); err != nil {
		return nil, err
	}

	return orderentity.NewPickupOrder(o.Name), nil
}
