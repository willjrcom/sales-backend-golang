package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
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
	bun.BaseModel `bun:"table:order_deliveries,alias:delivery"`
	DeliveryTimeLogs
	OrderDeliveryCommonAttributes
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

type OrderDeliveryCommonAttributes struct {
	Status      StatusOrderDelivery    `bun:"status"`
	DeliveryTax *float64               `bun:"delivery_tax"`
	ClientID    uuid.UUID              `bun:"column:client_id,type:uuid,notnull"`
	Client      *cliententity.Client   `bun:"rel:belongs-to"`
	AddressID   uuid.UUID              `bun:"column:address_id,type:uuid,notnull"`
	Address     *addressentity.Address `bun:"rel:belongs-to"`
	DriverID    *uuid.UUID             `bun:"column:driver_id,type:uuid"`
	Driver      *DeliveryDriver        `bun:"rel:belongs-to"`
	OrderID     uuid.UUID              `bun:"column:order_id,type:uuid,notnull"`
}

type DeliveryTimeLogs struct {
	PendingAt   *time.Time `bun:"pending_at"`
	ShippedAt   *time.Time `bun:"shipped_at"`
	DeliveredAt *time.Time `bun:"delivered_at"`
}
