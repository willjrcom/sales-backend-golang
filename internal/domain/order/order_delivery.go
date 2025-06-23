package orderentity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrOrderDeliveryMustBePending = errors.New("order delivery must be pending")
	ErrOrderDeliveryMustBeReady   = errors.New("order delivery must be ready")
	ErrOrderDeliveryMustBeShipped = errors.New("order delivery must be shipped")
)

type OrderDelivery struct {
	entity.Entity
	DeliveryTimeLogs
	OrderDeliveryCommonAttributes
}

type OrderDeliveryCommonAttributes struct {
	Status         StatusOrderDelivery
	DeliveryTax    *decimal.Decimal
	IsDeliveryFree bool
	Change         decimal.Decimal
	PaymentMethod  PayMethod
	ClientID       uuid.UUID
	Client         *cliententity.Client
	AddressID      uuid.UUID
	Address        *addressentity.Address
	DriverID       *uuid.UUID
	Driver         *DeliveryDriver
	OrderID        uuid.UUID
	OrderNumber    int
}

type DeliveryTimeLogs struct {
	PendingAt   *time.Time
	ReadyAt     *time.Time
	ShippedAt   *time.Time
	DeliveredAt *time.Time
	CanceledAt  *time.Time
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

// AddChange updates the uncollected change and payment method
func (d *OrderDelivery) AddChange(change decimal.Decimal, paymentMethod PayMethod) {
	d.Change = change
	d.PaymentMethod = paymentMethod
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

func (d *OrderDelivery) Ready() error {
	if d.Status != OrderDeliveryStatusPending {
		return nil
	}

	d.ReadyAt = &time.Time{}
	*d.ReadyAt = time.Now().UTC()
	d.Status = OrderDeliveryStatusReady
	return nil
}

func (d *OrderDelivery) Ship(driverID *uuid.UUID) error {
	if d.Status != OrderDeliveryStatusReady {
		return ErrOrderDeliveryMustBeReady
	}

	d.DriverID = driverID
	d.ShippedAt = &time.Time{}
	*d.ShippedAt = time.Now().UTC()
	d.Status = OrderDeliveryStatusShipped
	return nil
}

func (d *OrderDelivery) Cancel() error {
	d.CanceledAt = &time.Time{}
	*d.CanceledAt = time.Now().UTC()
	d.Status = OrderDeliveryStatusCanceled
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
