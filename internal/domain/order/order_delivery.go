package orderentity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrOrderDeliveryMustBePending = errors.New("order delivery must be pending")
	ErrOrderDeliveryMustBeShipped = errors.New("order delivery must be shipped")
)

type OrderDelivery struct {
	entity.Entity
	DeliveryTimeLogs
	OrderDeliveryCommonAttributes
}

type OrderDeliveryCommonAttributes struct {
	Status      StatusOrderDelivery
	DeliveryTax *float64
	ClientID    uuid.UUID
	Client      *cliententity.Client
	AddressID   uuid.UUID
	Address     *addressentity.Address
	DriverID    *uuid.UUID
	Driver      *DeliveryDriver
	OrderID     uuid.UUID
}

type DeliveryTimeLogs struct {
	PendingAt   *time.Time
	ShippedAt   *time.Time
	DeliveredAt *time.Time
}

func NewOrderDelivery(clientID uuid.UUID) *OrderDelivery {
	orderCommonAttributes := OrderDeliveryCommonAttributes{
		ClientID: clientID,
		Status:   OrderDeliveryStatusStaging,
	}

	return &OrderDelivery{
		Entity:                        entity.NewEntity(),
		OrderDeliveryCommonAttributes: orderCommonAttributes,
	}
}

func (d *OrderDelivery) Pend() error {
	if d.Status != OrderDeliveryStatusStaging {
		return nil
	}

	d.PendingAt = &time.Time{}
	*d.PendingAt = time.Now().UTC()
	d.Status = OrderDeliveryStatusPending
	return nil
}

func (d *OrderDelivery) Ship(driverID *uuid.UUID) error {
	if d.Status != OrderDeliveryStatusPending {
		return ErrOrderDeliveryMustBePending
	}

	d.DriverID = driverID
	d.ShippedAt = &time.Time{}
	*d.ShippedAt = time.Now().UTC()
	d.Status = OrderDeliveryStatusShipped
	return nil
}

func (d *OrderDelivery) Delivery() error {
	if d.Status != OrderDeliveryStatusShipped {
		return ErrOrderDeliveryMustBeShipped
	}

	d.DeliveredAt = &time.Time{}
	*d.DeliveredAt = time.Now().UTC()
	d.Status = OrderDeliveryStatusDelivered
	return nil
}
