package shiftdto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	shiftentity "github.com/willjrcom/sales-backend-go/internal/domain/shift"
	employeedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/employee"
	orderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order"
)

type ShiftDTO struct {
	ID uuid.UUID `json:"id"`
	ShiftTimeLogs
	ShiftCommonAttributes
}

type ShiftCommonAttributes struct {
	CurrentOrderNumber int                      `json:"current_order_number"`
	Orders             []orderdto.OrderDTO      `json:"orders"`
	Redeems            []RedeemDTO              `json:"redeems"`
	StartChange        decimal.Decimal          `json:"start_change"`
	EndChange          *decimal.Decimal         `json:"end_change"`
	AttendantID        *uuid.UUID               `json:"attendant_id"`
	Attendant          *employeedto.EmployeeDTO `json:"attendant"`
}

type ShiftTimeLogs struct {
	OpenedAt *time.Time `json:"opened_at"`
	ClosedAt *time.Time `json:"closed_at"`
}

func (s *ShiftDTO) FromDomain(shift *shiftentity.Shift) {
	if shift == nil {
		return
	}
	*s = ShiftDTO{
		ID: shift.ID,
		ShiftTimeLogs: ShiftTimeLogs{
			OpenedAt: shift.OpenedAt,
			ClosedAt: shift.ClosedAt,
		},
		ShiftCommonAttributes: ShiftCommonAttributes{
			CurrentOrderNumber: shift.CurrentOrderNumber,
			Orders:             []orderdto.OrderDTO{},
			Redeems:            []RedeemDTO{},
			StartChange:        shift.StartChange,
			EndChange:          shift.EndChange,
			AttendantID:        shift.AttendantID,
			Attendant:          &employeedto.EmployeeDTO{},
		},
	}

	for _, order := range shift.Orders {
		o := orderdto.OrderDTO{}
		o.FromDomain(&order)
		s.Orders = append(s.Orders, o)
	}

	for _, redeem := range shift.Redeems {
		r := RedeemDTO{}
		r.FromDomain(&redeem)
		s.Redeems = append(s.Redeems, r)
	}

	s.Attendant.FromDomain(shift.Attendant)

	if len(shift.Orders) == 0 {
		s.Orders = nil
	}
	if len(shift.Redeems) == 0 {
		s.Redeems = nil
	}
	if shift.Attendant == nil {
		s.Attendant = nil
	}
}
