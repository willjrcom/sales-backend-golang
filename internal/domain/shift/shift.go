package shiftentity

import (
	"time"

	"github.com/google/uuid"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type Shift struct {
	entity.Entity
	ShiftTimeLogs
	ShiftCommonAttributes
}

type ShiftCommonAttributes struct {
	CurrentOrderNumber int
	Orders             []orderentity.Order
	Redeems            []string
	StartChange        float32
	EndChange          *float32
	AttendantID        *uuid.UUID
	Attendant          *employeeentity.Employee
}

type ShiftTimeLogs struct {
	OpenedAt *time.Time
	ClosedAt *time.Time
}

func NewShift(startChange float32) *Shift {
	shift := &Shift{
		Entity: entity.NewEntity(),
		ShiftCommonAttributes: ShiftCommonAttributes{
			CurrentOrderNumber: 0,
		},
		ShiftTimeLogs: ShiftTimeLogs{
			OpenedAt: &time.Time{},
		},
	}

	*shift.OpenedAt = time.Now().UTC()
	return shift
}

func (s *Shift) CloseShift(endChange float32) (err error) {
	s.EndChange = &endChange
	s.ClosedAt = &time.Time{}
	*s.ClosedAt = time.Now().UTC()
	return nil
}

func (s *Shift) IncrementCurrentOrder() {
	s.CurrentOrderNumber++
}

func (s *Shift) IsClosed() bool {
	return s.EndChange != nil
}
