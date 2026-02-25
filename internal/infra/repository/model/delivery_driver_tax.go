package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	shiftentity "github.com/willjrcom/sales-backend-go/internal/domain/shift"
)

type DeliveryDriverTax struct {
	DeliveryDriverID   uuid.UUID        `bun:"column:delivery_driver_id,type:uuid,notnull"`
	DeliveryDriverName string           `bun:"delivery_driver_name"`
	OrderNumber        int              `bun:"order_number,notnull"`
	DeliveryID         uuid.UUID        `bun:"column:delivery_id,type:uuid,notnull"`
	DeliveryTax        *decimal.Decimal `bun:"delivery_tax,type:decimal(10,2),notnull"`
}

func (d *DeliveryDriverTax) FromDomain(deliveryDriverTax *shiftentity.DeliveryDriverTax) {
	if deliveryDriverTax == nil {
		return
	}
	*d = DeliveryDriverTax{
		DeliveryDriverID:   deliveryDriverTax.DeliveryDriverID,
		DeliveryDriverName: deliveryDriverTax.DeliveryDriverName,
		OrderNumber:        deliveryDriverTax.OrderNumber,
		DeliveryID:         deliveryDriverTax.DeliveryID,
		DeliveryTax:        &deliveryDriverTax.DeliveryTax,
	}
}

func (d *DeliveryDriverTax) ToDomain() *shiftentity.DeliveryDriverTax {
	if d == nil {
		return nil
	}
	return &shiftentity.DeliveryDriverTax{
		DeliveryDriverID:   d.DeliveryDriverID,
		DeliveryDriverName: d.DeliveryDriverName,
		OrderNumber:        d.OrderNumber,
		DeliveryID:         d.DeliveryID,
		DeliveryTax:        d.GetDeliveryTax(),
	}
}

func (d *DeliveryDriverTax) GetDeliveryTax() decimal.Decimal {
	if d.DeliveryTax == nil {
		return decimal.Zero
	}
	return *d.DeliveryTax
}
