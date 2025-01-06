package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type Order struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:orders,alias:order"`
	OrderTimeLogs
	OrderCommonAttributes
}

type OrderCommonAttributes struct {
	OrderType
	OrderDetail
	OrderNumber int            `bun:"order_number,notnull"`
	Status      string         `bun:"status,notnull"`
	Groups      []GroupItem    `bun:"rel:has-many,join:id=order_id"`
	Payments    []PaymentOrder `bun:"rel:has-many,join:id=order_id"`
}

type OrderDetail struct {
	TotalPayable  float64    `bun:"total_payable"`
	TotalPaid     float64    `bun:"total_paid"`
	TotalChange   float64    `bun:"total_change"`
	QuantityItems float64    `bun:"quantity_items"`
	Observation   string     `bun:"observation"`
	AttendantID   *uuid.UUID `bun:"column:attendant_id,type:uuid,notnull"`
	Attendant     *Employee  `bun:"rel:belongs-to"`
	ShiftID       *uuid.UUID `bun:"column:shift_id,type:uuid"`
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
