package model

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
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

type OrderTableCommonAttributes struct {
	Name    string           `bun:"name,notnull"`
	Contact string           `bun:"contact,notnull"`
	Status  StatusOrderTable `bun:"status,notnull"`
	OrderID uuid.UUID        `bun:"column:order_id,type:uuid,notnull"`
	TableID uuid.UUID        `bun:"column:table_id,type:uuid,notnull"`
}

type OrderTableTimeLogs struct {
	PendingAt *time.Time `bun:"pending_at"`
	ClosedAt  *time.Time `bun:"closed_at"`
}
