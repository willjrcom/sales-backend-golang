package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type Shift struct {
	entity.Entity
	bun.BaseModel `bun:"table:shifts"`
	ShiftTimeLogs
	ShiftCommonAttributes
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

type ShiftCommonAttributes struct {
	CurrentOrderNumber int                      `bun:"current_order_number,notnull"`
	Orders             []orderentity.Order      `bun:"rel:has-many,join:id=shift_id"`
	Redeems            []string                 `bun:"redeems,type:json"`
	StartChange        float32                  `bun:"start_change"`
	EndChange          *float32                 `bun:"end_change"`
	AttendantID        *uuid.UUID               `bun:"column:attendant_id,type:uuid"`
	Attendant          *employeeentity.Employee `bun:"rel:belongs-to"`
}

type ShiftTimeLogs struct {
	OpenedAt *time.Time `bun:"opened_at"`
	ClosedAt *time.Time `bun:"closed_at"`
}
