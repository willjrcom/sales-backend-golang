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
	CurrentOrderNumber int                      `bun:"current_order_number,notnull" json:"current_order_number"`
	Orders             []orderentity.Order      `bun:"rel:has-many,join:id=shift_id" json:"orders,omitempty"`
	Redeems            []string                 `bun:"redeems,type:json" json:"redeems,omitempty"`
	StartChange        float32                  `bun:"start_change" json:"start_change"`
	EndChange          *float32                 `bun:"end_change" json:"end_change,omitempty"`
	AttendantID        *uuid.UUID               `bun:"column:attendant_id,type:uuid" json:"attendant_id"`
	Attendant          *employeeentity.Employee `bun:"rel:belongs-to" json:"attendant"`
}

type OrderTimeLogs struct {
	OpenedAt *time.Time `bun:"opened_at" json:"opened_at,omitempty"`
	ClosedAt *time.Time `bun:"finished_at" json:"finished_at,omitempty"`
}

func (s *Shift) OpenShift() {
	s.CurrentOrderNumber = 0
	s.OpenedAt = &time.Time{}
	*s.OpenedAt = time.Now()
}

func (s *Shift) CloseShift(endChange float32) (err error) {
	s.EndChange = &endChange
	s.ClosedAt = &time.Time{}
	*s.ClosedAt = time.Now()
	return nil
}

func (s *Shift) IncrementCurrentOrder() {
	s.CurrentOrderNumber++
}

func (s *Shift) IsClosed() bool {
	return s.EndChange != nil
}
