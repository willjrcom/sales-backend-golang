package orderentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type OrderTable struct {
	entity.Entity
	bun.BaseModel `bun:"table:order_tables,alias:order_table"`
	OrderTableCommonAttributes
	OrderTableTimeLogs
}

type OrderTableCommonAttributes struct {
	Name    string           `bun:"name,notnull" json:"name,omitempty"`
	Contact string           `bun:"contact,notnull" json:"contact,omitempty"`
	Status  StatusOrderTable `bun:"status,notnull" json:"status"`
	OrderID uuid.UUID        `bun:"column:order_id,type:uuid,notnull" json:"order_id"`
	TableID uuid.UUID        `bun:"column:table_id,type:uuid,notnull" json:"table_id"`
}

type OrderTableTimeLogs struct {
	PendingAt *time.Time `bun:"pending_at" json:"pending_at,omitempty"`
	ClosedAt  *time.Time `bun:"closed_at" json:"closed_at,omitempty"`
}

func NewTable(orderTableCommonAttributes OrderTableCommonAttributes) *OrderTable {
	orderTableCommonAttributes.Status = OrderTableStatusStaging

	return &OrderTable{
		Entity:                     entity.NewEntity(),
		OrderTableCommonAttributes: orderTableCommonAttributes,
	}
}

func (t *OrderTable) Pend() error {
	if t.Status != OrderTableStatusStaging {
		return nil
	}

	t.Status = OrderTableStatusPending
	t.PendingAt = &time.Time{}
	*t.PendingAt = time.Now().UTC()
	return nil
}

func (t *OrderTable) Close() error {
	t.Status = OrderTableStatusClosed
	t.ClosedAt = &time.Time{}
	*t.ClosedAt = time.Now().UTC()
	return nil
}
