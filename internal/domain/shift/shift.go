package shiftentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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
	Redeems            []Redeem
	StartChange        decimal.Decimal
	EndChange          *decimal.Decimal
	AttendantID        *uuid.UUID
	Attendant          *employeeentity.Employee
}

type Redeem struct {
	Name  string
	Value decimal.Decimal
}

type ShiftTimeLogs struct {
	OpenedAt *time.Time
	ClosedAt *time.Time
}

// NewShift creates a new shift with initial start change
func NewShift(startChange decimal.Decimal) *Shift {
	now := time.Now().UTC()
	shift := &Shift{
		Entity: entity.NewEntity(),
		ShiftCommonAttributes: ShiftCommonAttributes{
			CurrentOrderNumber: 0,
			StartChange:        startChange,
		},
		ShiftTimeLogs: ShiftTimeLogs{
			OpenedAt: &now,
		},
	}
	return shift
}

// CloseShift closes the shift with final change
func (s *Shift) CloseShift(endChange decimal.Decimal) (err error) {
	now := time.Now().UTC()
	s.EndChange = &endChange
	s.ClosedAt = &now
	return nil
}

func (s *Shift) IncrementCurrentOrder() {
	s.CurrentOrderNumber++
}

func (s *Shift) IsClosed() bool {
	return s.EndChange != nil
}

func (s *Shift) AddRedeem(redeem *Redeem) {
	s.Redeems = append(s.Redeems, *redeem)
}
