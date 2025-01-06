package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type Shift struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:shifts"`
	ShiftTimeLogs
	ShiftCommonAttributes
}

type ShiftCommonAttributes struct {
	CurrentOrderNumber int        `bun:"current_order_number,notnull"`
	Orders             []Order    `bun:"rel:has-many,join:id=shift_id"`
	Redeems            []string   `bun:"redeems,type:json"`
	StartChange        float32    `bun:"start_change"`
	EndChange          *float32   `bun:"end_change"`
	AttendantID        *uuid.UUID `bun:"column:attendant_id,type:uuid"`
	Attendant          *Employee  `bun:"rel:belongs-to"`
}

type ShiftTimeLogs struct {
	OpenedAt *time.Time `bun:"opened_at"`
	ClosedAt *time.Time `bun:"closed_at"`
}
