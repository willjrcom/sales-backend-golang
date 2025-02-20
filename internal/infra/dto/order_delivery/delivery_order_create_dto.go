package orderdeliverydto

import (
	"errors"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

var (
	ErrOrderIDRequired   = errors.New("order id is required")
	ErrClientIDRequired  = errors.New("client id is required")
	ErrAddressIDRequired = errors.New("address id is required")
)

type DeliveryOrderCreateDTO struct {
	ClientID uuid.UUID `json:"client_id"`
}

func (o *DeliveryOrderCreateDTO) validate() error {
	if o.ClientID == uuid.Nil {
		return ErrClientIDRequired
	}

	return nil
}

func (o *DeliveryOrderCreateDTO) ToDomain() (*orderentity.OrderDelivery, error) {
	if err := o.validate(); err != nil {
		return nil, err
	}

	return orderentity.NewOrderDelivery(o.ClientID), nil
}
