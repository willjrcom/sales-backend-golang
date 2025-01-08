package orderdeliverydto

import (
	"time"

	"github.com/google/uuid"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	clientdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/client"
	deliverydriverdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/delivery_driver"
)

type OrderDeliveryDTO struct {
	ID uuid.UUID `json:"id"`
	DeliveryTimeLogs
	OrderDeliveryCommonAttributes
}

type OrderDeliveryCommonAttributes struct {
	Status      orderentity.StatusOrderDelivery      `json:"status"`
	DeliveryTax *float64                             `json:"delivery_tax"`
	ClientID    uuid.UUID                            `json:"client_id"`
	Client      *clientdto.ClientDTO                 `json:"client"`
	AddressID   uuid.UUID                            `json:"address_id"`
	Address     *addressentity.Address               `json:"address"`
	DriverID    *uuid.UUID                           `json:"driver_id"`
	Driver      *deliverydriverdto.DeliveryDriverDTO `json:"driver"`
	OrderID     uuid.UUID                            `json:"order_id"`
}

type DeliveryTimeLogs struct {
	PendingAt   *time.Time `json:"pending_at"`
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
			Status:      delivery.Status,
			DeliveryTax: delivery.DeliveryTax,
			ClientID:    delivery.ClientID,
			Client:      &clientdto.ClientDTO{},
			AddressID:   delivery.AddressID,
			Address:     delivery.Address,
			DriverID:    delivery.DriverID,
			Driver:      &deliverydriverdto.DeliveryDriverDTO{},
			OrderID:     delivery.OrderID,
		},
		DeliveryTimeLogs: DeliveryTimeLogs{
			PendingAt:   delivery.PendingAt,
			ShippedAt:   delivery.ShippedAt,
			DeliveredAt: delivery.DeliveredAt,
		},
	}

	o.Client.FromDomain(delivery.Client)
	o.Driver.FromDomain(delivery.Driver)

	if delivery.Client == nil {
		o.Client = nil
	}
	if delivery.Driver == nil {
		o.Driver = nil
	}
}
