package deliveryorderdto

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
	ClientID uuid.UUID `json:"client_id"`
}

func (o *CreateDeliveryOrderInput) validate() error {
	if o.ClientID == uuid.Nil {
		return ErrClientIDRequired
	}

	return nil
}

func (o *CreateDeliveryOrderInput) ToModel() (*orderentity.DeliveryOrder, error) {
	if err := o.validate(); err != nil {
		return nil, err
	}

	orderCommonAttributes := orderentity.DeliveryOrderCommonAttributes{
		ClientID: o.ClientID,
		Status:   orderentity.DeliveryOrderStatusPending,
	}

	return &orderentity.DeliveryOrder{
		Entity:                        entity.NewEntity(),
		DeliveryOrderCommonAttributes: orderCommonAttributes,
	}, nil
}
