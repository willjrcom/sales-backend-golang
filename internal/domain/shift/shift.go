package shiftentity

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
	OrderTimeLogs
	ShiftCommonAttributes
}

type ShiftCommonAttributes struct {
	Day         *time.Time               `bun:"day,notnull" json:"day"`
	Orders      []orderentity.Order      `bun:"rel:has-many,join:id=shift_id" json:"orders"`
	Redeem      []string                 `bun:"redeem,type:json" json:"redeem"`
	StartChange float32                  `bun:"start_change" json:"start_change"`
	EndChange   *float32                 `bun:"end_change" json:"end_change"`
	AttendantID *uuid.UUID               `bun:"column:attendant_id,type:uuid" json:"attendant_id"`
	Attendant   *employeeentity.Employee `bun:"rel:belongs-to" json:"attendant"`
}

type OrderTimeLogs struct {
	OpenedAt *time.Time `bun:"opened_at" json:"opened_at,omitempty"`
	ClosedAt *time.Time `bun:"finished_at" json:"finished_at,omitempty"`
}

func (s *Shift) OpenShift() (err error) {
	s.OpenedAt = &time.Time{}
	*s.OpenedAt = time.Now()
	return nil
}

func (s *Shift) CloseShift(endChange float32) (err error) {
	s.EndChange = &endChange
	s.ClosedAt = &time.Time{}
	*s.ClosedAt = time.Now()
	return nil
}
