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
	ErrDeliveryOrderMustBeStaging = errors.New("delivery order must be staging")
	ErrDeliveryOrderMustBePending = errors.New("delivery order must be pending")
	ErrDeliveryOrderMustBeShipped = errors.New("delivery order must be shipped")
)

type DeliveryOrder struct {
	entity.Entity
	bun.BaseModel `bun:"table:delivery_orders,alias:delivery"`
	DeliveryTimeLogs
	DeliveryOrderCommonAttributes
}

type DeliveryOrderCommonAttributes struct {
	Status      StatusDeliveryOrder      `bun:"status" json:"status"`
	DeliveryTax *float64                 `bun:"delivery_tax" json:"delivery_tax"`
	ClientID    uuid.UUID                `bun:"column:client_id,type:uuid,notnull" json:"client_id"`
	Client      *cliententity.Client     `bun:"rel:belongs-to" json:"client"`
	AddressID   uuid.UUID                `bun:"column:address_id,type:uuid,notnull" json:"address_id"`
	Address     *addressentity.Address   `bun:"rel:belongs-to" json:"address"`
	DriverID    *uuid.UUID               `bun:"column:driver_id,type:uuid" json:"driver_id"`
	Driver      *employeeentity.Employee `bun:"rel:belongs-to" json:"driver"`
	OrderID     uuid.UUID                `bun:"column:order_id,type:uuid,notnull" json:"order_id"`
}

type DeliveryTimeLogs struct {
	PendingAt   *time.Time `bun:"pending_at" json:"pending_at,omitempty"`
	ShippedAt   *time.Time `bun:"shipped_at" json:"shipped_at,omitempty"`
	DeliveredAt *time.Time `bun:"delivered_at" json:"delivered_at,omitempty"`
}

func NewDeliveryOrder(clientID uuid.UUID) *DeliveryOrder {
	orderCommonAttributes := DeliveryOrderCommonAttributes{
		ClientID: clientID,
		Status:   DeliveryOrderStatusStaging,
	}

	return &DeliveryOrder{
		Entity:                        entity.NewEntity(),
		DeliveryOrderCommonAttributes: orderCommonAttributes,
	}
}

func (d *DeliveryOrder) Pending() error {
	if d.Status != DeliveryOrderStatusStaging {
		return ErrDeliveryOrderMustBeStaging
	}
	d.PendingAt = &time.Time{}
	*d.PendingAt = time.Now()
	d.Status = DeliveryOrderStatusPending
	return nil
}

func (d *DeliveryOrder) Ship(driverID uuid.UUID) error {
	if d.Status != DeliveryOrderStatusPending {
		return ErrDeliveryOrderMustBePending
	}

	*d.DriverID = driverID
	d.ShippedAt = &time.Time{}
	*d.ShippedAt = time.Now()
	d.Status = DeliveryOrderStatusShipped
	return nil
}

func (d *DeliveryOrder) Delivery() error {
	if d.Status != DeliveryOrderStatusShipped {
		return ErrDeliveryOrderMustBeShipped
	}

	d.DeliveredAt = &time.Time{}
	*d.DeliveredAt = time.Now()
	d.Status = DeliveryOrderStatusDelivered
	return nil
}
