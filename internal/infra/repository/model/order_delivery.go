package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type OrderDelivery struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:order_deliveries,alias:delivery"`
	DeliveryTimeLogs
	OrderDeliveryCommonAttributes
}

type OrderDeliveryCommonAttributes struct {
	Status      string          `bun:"status"`
	DeliveryTax *float64        `bun:"delivery_tax"`
	ClientID    uuid.UUID       `bun:"column:client_id,type:uuid,notnull"`
	Client      *Client         `bun:"rel:belongs-to"`
	AddressID   uuid.UUID       `bun:"column:address_id,type:uuid,notnull"`
	Address     *Address        `bun:"rel:belongs-to"`
	DriverID    *uuid.UUID      `bun:"column:driver_id,type:uuid"`
	Driver      *DeliveryDriver `bun:"rel:belongs-to"`
	OrderID     uuid.UUID       `bun:"column:order_id,type:uuid,notnull"`
}

type DeliveryTimeLogs struct {
	PendingAt   *time.Time `bun:"pending_at"`
	ShippedAt   *time.Time `bun:"shipped_at"`
	DeliveredAt *time.Time `bun:"delivered_at"`
}
