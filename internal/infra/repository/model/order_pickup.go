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
	PendingAt *time.Time `bun:"pending_at"`
	ReadyAt   *time.Time `bun:"ready_at"`
}

func (p *OrderPickup) FromDomain(pickup *orderentity.OrderPickup) {
	*p = OrderPickup{
		Entity: entitymodel.FromDomain(pickup.Entity),
		PickupTimeLogs: PickupTimeLogs{
			PendingAt: pickup.PendingAt,
			ReadyAt:   pickup.ReadyAt,
		},
		OrderPickupCommonAttributes: OrderPickupCommonAttributes{
			Name:    pickup.Name,
			Status:  string(pickup.Status),
			OrderID: pickup.OrderID,
		},
	}
}

func (p *OrderPickup) ToDomain() *orderentity.OrderPickup {
	return &orderentity.OrderPickup{
		Entity: p.Entity.ToDomain(),
		PickupTimeLogs: orderentity.PickupTimeLogs{
			PendingAt: p.PendingAt,
			ReadyAt:   p.ReadyAt,
		},
		OrderPickupCommonAttributes: orderentity.OrderPickupCommonAttributes{
			Name:    p.Name,
			Status:  orderentity.StatusOrderPickup(p.Status),
			OrderID: p.OrderID,
		},
	}
}
