package orderentity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrTableOrderMustBeStaging = errors.New("table order must be staging")
	ErrTableOrderMustBePending = errors.New("table order must be pending")
)

type TableOrder struct {
	entity.Entity
	bun.BaseModel `bun:"table:table_orders,alias:table_order"`
	TableOrderCommonAttributes
	TableOrderTimeLogs
}

type TableOrderCommonAttributes struct {
	Name    string           `bun:"name,notnull" json:"name,omitempty"`
	Contact string           `bun:"contact,notnull" json:"contact,omitempty"`
	Status  StatusTableOrder `bun:"status,notnull" json:"status"`
	OrderID uuid.UUID        `bun:"column:order_id,type:uuid,notnull" json:"order_id"`
	TableID uuid.UUID        `bun:"column:table_id,type:uuid,notnull" json:"table_id"`
}

type TableOrderTimeLogs struct {
	PendingAt *time.Time `bun:"pending_at" json:"pending_at,omitempty"`
	ClosedAt  *time.Time `bun:"closed_at" json:"closed_at,omitempty"`
}

func NewTable(tableOrderCommonAttributes TableOrderCommonAttributes) *TableOrder {
	tableOrderCommonAttributes.Status = TableOrderStatusStaging

	return &TableOrder{
		Entity:                     entity.NewEntity(),
		TableOrderCommonAttributes: tableOrderCommonAttributes,
	}
}

func (t *TableOrder) Pend() error {
	if t.Status != TableOrderStatusStaging {
		return ErrTableOrderMustBeStaging
	}

	t.Status = TableOrderStatusPending
	t.PendingAt = &time.Time{}
	*t.PendingAt = time.Now()
	return nil
}

func (t *TableOrder) Close() error {
	if t.Status != TableOrderStatusPending {
		return ErrTableOrderMustBePending
	}

	t.Status = TableOrderStatusClosed
	t.ClosedAt = &time.Time{}
	*t.ClosedAt = time.Now()
	return nil
}
