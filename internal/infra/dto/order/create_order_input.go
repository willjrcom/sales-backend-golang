package orderdto

import (
	"errors"

	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

var (
	ErrNameRequired = errors.New("name is required")
)

type CreateOrderInput struct {
	Name string `json:"name"`
}

func (o *CreateOrderInput) validate() error {
	if o.Name == "" {
		return ErrNameRequired
	}

	return nil
}

func (o *CreateOrderInput) ToModel() (*orderentity.Order, error) {
	if err := o.validate(); err != nil {
		return nil, err
	}

	orderCommonAttributes := orderentity.OrderCommonAttributes{
		OrderDetail: orderentity.OrderDetail{Name: o.Name},
		Status:      orderentity.OrderStatusStaging,
	}

	return &orderentity.Order{
		Entity:                entity.NewEntity(),
		OrderCommonAttributes: orderCommonAttributes,
	}, nil
}
