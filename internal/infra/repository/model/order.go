package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Order struct {
	entity.Entity
	bun.BaseModel `bun:"table:orders,alias:order"`
	OrderTimeLogs
	OrderCommonAttributes
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

type OrderCommonAttributes struct {
	OrderType
	OrderDetail
	OrderNumber int            `bun:"order_number,notnull"`
	Status      StatusOrder    `bun:"status,notnull"`
	Groups      []GroupItem    `bun:"rel:has-many,join:id=order_id"`
	Payments    []PaymentOrder `bun:"rel:has-many,join:id=order_id"`
}

type OrderDetail struct {
	TotalPayable  float64                  `bun:"total_payable"`
	TotalPaid     float64                  `bun:"total_paid"`
	TotalChange   float64                  `bun:"total_change"`
	QuantityItems float64                  `bun:"quantity_items"`
	Observation   string                   `bun:"observation"`
	AttendantID   *uuid.UUID               `bun:"column:attendant_id,type:uuid,notnull"`
	Attendant     *employeeentity.Employee `bun:"rel:belongs-to"`
	ShiftID       *uuid.UUID               `bun:"column:shift_id,type:uuid"`
}

type OrderType struct {
	Delivery *OrderDelivery `bun:"rel:has-one,join:id=order_id"`
	Table    *OrderTable    `bun:"rel:has-one,join:id=order_id"`
	Pickup   *OrderPickup   `bun:"rel:has-one,join:id=order_id"`
}

type OrderTimeLogs struct {
	PendingAt  *time.Time `bun:"pending_at"`
	FinishedAt *time.Time `bun:"finished_at"`
	CanceledAt *time.Time `bun:"canceled_at"`
	ArchivedAt *time.Time `bun:"archived_at"`
}
