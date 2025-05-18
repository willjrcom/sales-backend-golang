package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type OrderDelivery struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:order_deliveries,alias:delivery"`
	DeliveryTimeLogs
	OrderDeliveryCommonAttributes
}

type OrderDeliveryCommonAttributes struct {
	Status        string           `bun:"status"`
	DeliveryTax   *decimal.Decimal `bun:"delivery_tax"`
	Change        decimal.Decimal  `bun:"change"`
	PaymentMethod string           `bun:"payment_method"`
	ClientID      uuid.UUID        `bun:"column:client_id,type:uuid,notnull"`
	Client        *Client          `bun:"rel:belongs-to"`
	AddressID     uuid.UUID        `bun:"column:address_id,type:uuid,notnull"`
	Address       *Address         `bun:"rel:belongs-to"`
	DriverID      *uuid.UUID       `bun:"column:driver_id,type:uuid"`
	Driver        *DeliveryDriver  `bun:"rel:belongs-to"`
	OrderID       uuid.UUID        `bun:"column:order_id,type:uuid,notnull"`
}

type DeliveryTimeLogs struct {
	PendingAt   *time.Time `bun:"pending_at"`
	ShippedAt   *time.Time `bun:"shipped_at"`
	DeliveredAt *time.Time `bun:"delivered_at"`
}

func (d *OrderDelivery) FromDomain(delivery *orderentity.OrderDelivery) {
	if delivery == nil {
		return
	}
	*d = OrderDelivery{
		Entity: entitymodel.FromDomain(delivery.Entity),
		OrderDeliveryCommonAttributes: OrderDeliveryCommonAttributes{
			Status:        string(delivery.Status),
			DeliveryTax:   delivery.DeliveryTax,
			Change:        delivery.Change,
			PaymentMethod: string(delivery.PaymentMethod),
			ClientID:      delivery.ClientID,
			Client:        &Client{},
			AddressID:     delivery.AddressID,
			Address:       &Address{},
			DriverID:      delivery.DriverID,
			Driver:        &DeliveryDriver{},
			OrderID:       delivery.OrderID,
		},
		DeliveryTimeLogs: DeliveryTimeLogs{
			PendingAt:   delivery.PendingAt,
			ShippedAt:   delivery.ShippedAt,
			DeliveredAt: delivery.DeliveredAt,
		},
	}

	d.Address.FromDomain(delivery.Address)
	d.Client.FromDomain(delivery.Client)
	d.Driver.FromDomain(delivery.Driver)
}

func (d *OrderDelivery) ToDomain() *orderentity.OrderDelivery {
	if d == nil {
		return nil
	}
	return &orderentity.OrderDelivery{
		Entity: d.Entity.ToDomain(),
		OrderDeliveryCommonAttributes: orderentity.OrderDeliveryCommonAttributes{
			Status:        orderentity.StatusOrderDelivery(d.Status),
			DeliveryTax:   d.DeliveryTax,
			Change:        d.Change,
			PaymentMethod: orderentity.PayMethod(d.PaymentMethod),
			ClientID:      d.ClientID,
			Client:        d.Client.ToDomain(),
			AddressID:     d.AddressID,
			Address:       d.Address.ToDomain(),
			DriverID:      d.DriverID,
			Driver:        d.Driver.ToDomain(),
			OrderID:       d.OrderID,
		},
		DeliveryTimeLogs: orderentity.DeliveryTimeLogs{
			PendingAt:   d.PendingAt,
			ShippedAt:   d.ShippedAt,
			DeliveredAt: d.DeliveredAt,
		},
	}
}
