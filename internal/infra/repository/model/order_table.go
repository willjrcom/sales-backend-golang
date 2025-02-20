package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type OrderTable struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:order_tables,alias:order_table"`
	OrderTableCommonAttributes
	OrderTableTimeLogs
}

type OrderTableCommonAttributes struct {
	Name    string    `bun:"name,notnull"`
	Contact string    `bun:"contact,notnull"`
	Status  string    `bun:"status,notnull"`
	OrderID uuid.UUID `bun:"column:order_id,type:uuid,notnull"`
	TableID uuid.UUID `bun:"column:table_id,type:uuid,notnull"`
}

type OrderTableTimeLogs struct {
	PendingAt *time.Time `bun:"pending_at"`
	ClosedAt  *time.Time `bun:"closed_at"`
}

func (t *OrderTable) FromDomain(table *orderentity.OrderTable) {
	if table == nil {
		return
	}
	*t = OrderTable{
		Entity: entitymodel.FromDomain(table.Entity),
		OrderTableCommonAttributes: OrderTableCommonAttributes{
			Name:    table.Name,
			Contact: table.Contact,
			Status:  string(table.Status),
			OrderID: table.OrderID,
			TableID: table.TableID,
		},
		OrderTableTimeLogs: OrderTableTimeLogs{
			PendingAt: table.PendingAt,
			ClosedAt:  table.ClosedAt,
		},
	}
}

func (t *OrderTable) ToDomain() *orderentity.OrderTable {
	if t == nil {
		return nil
	}
	return &orderentity.OrderTable{
		Entity: t.Entity.ToDomain(),
		OrderTableCommonAttributes: orderentity.OrderTableCommonAttributes{
			Name:    t.Name,
			Contact: t.Contact,
			Status:  orderentity.StatusOrderTable(t.Status),
			OrderID: t.OrderID,
			TableID: t.TableID,
		},
		OrderTableTimeLogs: orderentity.OrderTableTimeLogs{
			PendingAt: t.PendingAt,
			ClosedAt:  t.ClosedAt,
		},
	}
}
