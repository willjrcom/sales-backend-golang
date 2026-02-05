package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type OrderPickup struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:order_pickups,alias:pickup"`
	PickupTimeLogs
	OrderPickupCommonAttributes
}

type OrderPickupCommonAttributes struct {
	Name        string    `bun:"name,notnull"`
	Status      string    `bun:"status"`
	OrderID     uuid.UUID `bun:"column:order_id,type:uuid,notnull"`
	OrderNumber int       `bun:"order_number,notnull"`
}

type PickupTimeLogs struct {
	PendingAt   *time.Time `bun:"pending_at"`
	ReadyAt     *time.Time `bun:"ready_at"`
	DeliveredAt *time.Time `bun:"delivered_at"`
	CancelledAt *time.Time `bun:"cancelled_at"`
}

func (p *OrderPickup) FromDomain(pickup *orderentity.OrderPickup) {
	if pickup == nil {
		return
	}
	*p = OrderPickup{
		Entity: entitymodel.FromDomain(pickup.Entity),
		PickupTimeLogs: PickupTimeLogs{
			PendingAt:   pickup.PendingAt,
			ReadyAt:     pickup.ReadyAt,
			DeliveredAt: pickup.DeliveredAt,
			CancelledAt: pickup.CancelledAt,
		},
		OrderPickupCommonAttributes: OrderPickupCommonAttributes{
			Name:        pickup.Name,
			Status:      string(pickup.Status),
			OrderID:     pickup.OrderID,
			OrderNumber: pickup.OrderNumber,
		},
	}
}

func (p *OrderPickup) ToDomain() *orderentity.OrderPickup {
	if p == nil {
		return nil
	}
	return &orderentity.OrderPickup{
		Entity: p.Entity.ToDomain(),
		PickupTimeLogs: orderentity.PickupTimeLogs{
			PendingAt:   p.PendingAt,
			ReadyAt:     p.ReadyAt,
			DeliveredAt: p.DeliveredAt,
			CancelledAt: p.CancelledAt,
		},
		OrderPickupCommonAttributes: orderentity.OrderPickupCommonAttributes{
			Name:        p.Name,
			Status:      orderentity.StatusOrderPickup(p.Status),
			OrderID:     p.OrderID,
			OrderNumber: p.OrderNumber,
		},
	}
}
