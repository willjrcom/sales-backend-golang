package orderentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type DeliveryOrder struct {
	entity.Entity
	bun.BaseModel `bun:"table:delivery_orders"`
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
	LaunchedAt  *time.Time `bun:"launched_at" json:"launched_at,omitempty"`
	DeliveredAt *time.Time `bun:"delivered_at" json:"delivered_at,omitempty"`
}

func NewDeliveryOrder(clientID uuid.UUID) *DeliveryOrder {
	orderCommonAttributes := DeliveryOrderCommonAttributes{
		ClientID: clientID,
		Status:   DeliveryOrderStatusPending,
	}

	return &DeliveryOrder{
		Entity:                        entity.NewEntity(),
		DeliveryOrderCommonAttributes: orderCommonAttributes,
	}
}

func (d *DeliveryOrder) LaunchDelivery(driverID uuid.UUID) {
	*d.DriverID = driverID
	d.LaunchedAt = &time.Time{}
	*d.LaunchedAt = time.Now()
	d.Status = DeliveryOrderStatusShipped
}

func (d *DeliveryOrder) FinishDelivery() {
	d.DeliveredAt = &time.Time{}
	*d.DeliveredAt = time.Now()
	d.Status = DeliveryOrderStatusDelivered
}
