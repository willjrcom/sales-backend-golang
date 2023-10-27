package orderdto

import (
	"errors"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

var (
	ErrOrderIDRequired   = errors.New("order id is required")
	ErrClientIDRequired  = errors.New("client id is required")
	ErrAddressIDRequired = errors.New("address id is required")
)

type CreateDeliveryOrderInput struct {
	OrderID   uuid.UUID `json:"order_id"`
	ClientID  uuid.UUID `json:"client_id"`
	AddressID uuid.UUID `json:"address_id"`
}

func (o *CreateDeliveryOrderInput) Validate() error {
	if o.OrderID == uuid.Nil {
		return ErrOrderIDRequired
	}

	if o.ClientID == uuid.Nil {
		return ErrClientIDRequired
	}

	if o.AddressID == uuid.Nil {
		return ErrAddressIDRequired
	}

	return nil
}

func (o *CreateDeliveryOrderInput) ToModel() (*orderentity.DeliveryOrder, error) {
	if err := o.Validate(); err != nil {
		return nil, err
	}

	return &orderentity.DeliveryOrder{
		Entity:      entity.NewEntity(),
		OrderID:     o.OrderID,
		ClientID:    o.ClientID,
		AddressID:   o.AddressID,
		IsCompleted: false,
	}, nil
}
