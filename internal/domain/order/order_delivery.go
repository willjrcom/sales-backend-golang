package orderentity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrOrderDeliveryMustBeStaging = errors.New("order delivery must be staging")
	ErrOrderDeliveryMustBePending = errors.New("order delivery must be pending")
	ErrOrderDeliveryMustBeShipped = errors.New("order delivery must be shipped")
)

type OrderDelivery struct {
	entity.Entity
	bun.BaseModel `bun:"table:order_deliveries,alias:delivery"`
	DeliveryTimeLogs
	OrderDeliveryCommonAttributes
}

type OrderDeliveryCommonAttributes struct {
	Status      StatusOrderDelivery      `bun:"status" json:"status"`
	DeliveryTax *float64                 `bun:"delivery_tax" json:"delivery_tax"`
	ClientID    uuid.UUID                `bun:"column:client_id,type:uuid,notnull" json:"client_id"`
	Client      *cliententity.Client     `bun:"rel:belongs-to" json:"client"`
	AddressID   uuid.UUID                `bun:"column:address_id,type:uuid,notnull" json:"address_id"`
	Address     *addressentity.Address   `bun:"rel:belongs-to" json:"address"`
	DriverID    uuid.UUID                `bun:"column:driver_id,type:uuid" json:"driver_id,omitempty"`
	Driver      *employeeentity.Employee `bun:"rel:belongs-to" json:"driver"`
	OrderID     uuid.UUID                `bun:"column:order_id,type:uuid,notnull" json:"order_id"`
}

type DeliveryTimeLogs struct {
	PendingAt   *time.Time `bun:"pending_at" json:"pending_at,omitempty"`
	ShippedAt   *time.Time `bun:"shipped_at" json:"shipped_at,omitempty"`
	DeliveredAt *time.Time `bun:"delivered_at" json:"delivered_at,omitempty"`
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
		return ErrOrderDeliveryMustBeStaging
	}
	d.PendingAt = &time.Time{}
	*d.PendingAt = time.Now()
	d.Status = OrderDeliveryStatusPending
	return nil
}

func (d *OrderDelivery) Ship(driverID uuid.UUID) error {
	if d.Status != OrderDeliveryStatusPending {
		return ErrOrderDeliveryMustBePending
	}

	d.DriverID = driverID
	d.ShippedAt = &time.Time{}
	*d.ShippedAt = time.Now()
	d.Status = OrderDeliveryStatusShipped
	return nil
}

func (d *OrderDelivery) Delivery() error {
	if d.Status != OrderDeliveryStatusShipped {
		return ErrOrderDeliveryMustBeShipped
	}

	d.DeliveredAt = &time.Time{}
	*d.DeliveredAt = time.Now()
	d.Status = OrderDeliveryStatusDelivered
	return nil
}
