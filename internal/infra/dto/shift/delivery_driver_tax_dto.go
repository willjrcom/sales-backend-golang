package shiftdto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	shiftentity "github.com/willjrcom/sales-backend-go/internal/domain/shift"
)

type DeliveryDriverTaxDTO struct {
	DeliveryDriverID   uuid.UUID       `json:"delivery_driver_id"`
	DeliveryDriverName string          `json:"delivery_driver_name"`
	DeliveryID         uuid.UUID       `json:"delivery_id"`
	DeliveryTax        decimal.Decimal `json:"delivery_tax"`
}

func (d *DeliveryDriverTaxDTO) FromDomain(deliveryDriverTax *shiftentity.DeliveryDriverTax) {
	if deliveryDriverTax == nil {
		return
	}
	*d = DeliveryDriverTaxDTO{
		DeliveryDriverID:   deliveryDriverTax.DeliveryDriverID,
		DeliveryDriverName: deliveryDriverTax.DeliveryDriverName,
		DeliveryID:         deliveryDriverTax.DeliveryID,
		DeliveryTax:        deliveryDriverTax.DeliveryTax,
	}
}
