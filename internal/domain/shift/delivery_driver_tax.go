package shiftentity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type DeliveryDriverTax struct {
	DeliveryDriverID   uuid.UUID
	DeliveryDriverName string
	OrderNumber        int
	DeliveryID         uuid.UUID
	DeliveryTax        decimal.Decimal
}
