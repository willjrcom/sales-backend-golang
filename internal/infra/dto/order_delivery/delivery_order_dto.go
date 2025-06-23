package orderdeliverydto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	addressdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/address"
	clientdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/client"
	deliverydriverdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/delivery_driver"
)

type OrderDeliveryDTO struct {
	ID uuid.UUID `json:"id"`
	DeliveryTimeLogs
	OrderDeliveryCommonAttributes
}

type OrderDeliveryCommonAttributes struct {
	Status        orderentity.StatusOrderDelivery      `json:"status"`
	DeliveryTax   *decimal.Decimal                     `json:"delivery_tax"`
	Change        decimal.Decimal                      `json:"change"`
	PaymentMethod string                               `json:"payment_method"`
	ClientID      uuid.UUID                            `json:"client_id"`
	Client        *clientdto.ClientDTO                 `json:"client"`
	AddressID     uuid.UUID                            `json:"address_id"`
	Address       *addressdto.AddressDTO               `json:"address"`
	DriverID      *uuid.UUID                           `json:"driver_id"`
	Driver        *deliverydriverdto.DeliveryDriverDTO `json:"driver"`
	OrderID       uuid.UUID                            `json:"order_id"`
	OrderNumber   int                                  `json:"order_number"`
}

type DeliveryTimeLogs struct {
	PendingAt   *time.Time `json:"pending_at"`
	ReadyAt     *time.Time `json:"ready_at"`
	ShippedAt   *time.Time `json:"shipped_at"`
	DeliveredAt *time.Time `json:"delivered_at"`
}

func (o *OrderDeliveryDTO) FromDomain(delivery *orderentity.OrderDelivery) {
	if delivery == nil {
		return
	}
	*o = OrderDeliveryDTO{
		ID: delivery.ID,
		OrderDeliveryCommonAttributes: OrderDeliveryCommonAttributes{
			Status:        delivery.Status,
			DeliveryTax:   delivery.DeliveryTax,
			Change:        delivery.Change,
			PaymentMethod: string(delivery.PaymentMethod),
			ClientID:      delivery.ClientID,
			Client:        &clientdto.ClientDTO{},
			AddressID:     delivery.AddressID,
			Address:       &addressdto.AddressDTO{},
			DriverID:      delivery.DriverID,
			Driver:        &deliverydriverdto.DeliveryDriverDTO{},
			OrderID:       delivery.OrderID,
			OrderNumber:   delivery.OrderNumber,
		},
		DeliveryTimeLogs: DeliveryTimeLogs{
			PendingAt:   delivery.PendingAt,
			ReadyAt:     delivery.ReadyAt,
			ShippedAt:   delivery.ShippedAt,
			DeliveredAt: delivery.DeliveredAt,
		},
	}

	o.Client.FromDomain(delivery.Client)
	o.Address.FromDomain(delivery.Address)
	o.Driver.FromDomain(delivery.Driver)

	if delivery.Client == nil {
		o.Client = nil
	}
	if delivery.Address == nil {
		o.Address = nil
	}
	if delivery.Driver == nil {
		o.Driver = nil
	}
}
