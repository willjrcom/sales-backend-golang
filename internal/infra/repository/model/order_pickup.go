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
	Name    string    `bun:"name,notnull"`
	Status  string    `bun:"status"`
	OrderID uuid.UUID `bun:"column:order_id,type:uuid,notnull"`
}

type PickupTimeLogs struct {
	PendingAt  *time.Time `bun:"pending_at"`
	ReadyAt    *time.Time `bun:"ready_at"`
	CanceledAt *time.Time `bun:"canceled_at"`
}

func (p *OrderPickup) FromDomain(pickup *orderentity.OrderPickup) {
	if pickup == nil {
		return
	}
	*p = OrderPickup{
		Entity: entitymodel.FromDomain(pickup.Entity),
		PickupTimeLogs: PickupTimeLogs{
			PendingAt:  pickup.PendingAt,
			ReadyAt:    pickup.ReadyAt,
			CanceledAt: pickup.CanceledAt,
		},
		OrderPickupCommonAttributes: OrderPickupCommonAttributes{
			Name:    pickup.Name,
			Status:  string(pickup.Status),
			OrderID: pickup.OrderID,
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
			PendingAt:  p.PendingAt,
			ReadyAt:    p.ReadyAt,
			CanceledAt: p.CanceledAt,
		},
		OrderPickupCommonAttributes: orderentity.OrderPickupCommonAttributes{
			Name:    p.Name,
			Status:  orderentity.StatusOrderPickup(p.Status),
			OrderID: p.OrderID,
		},
	}
}
