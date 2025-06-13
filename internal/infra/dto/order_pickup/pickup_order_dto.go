package orderpickupdto

import (
	"time"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type OrderPickupDTO struct {
	ID uuid.UUID `json:"id"`
	PickupTimeLogs
	OrderPickupCommonAttributes
}

type OrderPickupCommonAttributes struct {
	Name    string                        `json:"name"`
	Status  orderentity.StatusOrderPickup `json:"status"`
	OrderID uuid.UUID                     `json:"order_id"`
}

type PickupTimeLogs struct {
	PendingAt   *time.Time `json:"pending_at"`
	ReadyAt     *time.Time `json:"ready_at"`
	DeliveredAt *time.Time `json:"delivered_at"`
}

func (o *OrderPickupDTO) FromDomain(pickup *orderentity.OrderPickup) {
	if pickup == nil {
		return
	}
	*o = OrderPickupDTO{
		ID: pickup.ID,
		OrderPickupCommonAttributes: OrderPickupCommonAttributes{
			Name:    pickup.Name,
			Status:  pickup.Status,
			OrderID: pickup.OrderID,
		},
		PickupTimeLogs: PickupTimeLogs{
			PendingAt:   pickup.PendingAt,
			ReadyAt:     pickup.ReadyAt,
			DeliveredAt: pickup.DeliveredAt,
		},
	}
}
