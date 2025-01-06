package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	shiftentity "github.com/willjrcom/sales-backend-go/internal/domain/shift"
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

func (s *Shift) FromDomain(shift *shiftentity.Shift) {
	*s = Shift{
		Entity: entitymodel.FromDomain(shift.Entity),
		ShiftTimeLogs: ShiftTimeLogs{
			OpenedAt: shift.OpenedAt,
			ClosedAt: shift.ClosedAt,
		},
		ShiftCommonAttributes: ShiftCommonAttributes{
			CurrentOrderNumber: shift.CurrentOrderNumber,
			StartChange:        shift.StartChange,
			EndChange:          shift.EndChange,
			AttendantID:        shift.AttendantID,
			Redeems:            shift.Redeems,
			Orders:             []Order{},
			Attendant:          &Employee{},
		},
	}

	for _, order := range shift.Orders {
		o := Order{}
		o.FromDomain(&order)
		s.Orders = append(s.Orders, o)
	}

	s.Attendant.FromDomain(shift.Attendant)
}

func (s *Shift) ToDomain() *shiftentity.Shift {
	if s == nil {
		return nil
	}
	shift := &shiftentity.Shift{
		Entity: s.Entity.ToDomain(),
		ShiftTimeLogs: shiftentity.ShiftTimeLogs{
			OpenedAt: s.OpenedAt,
			ClosedAt: s.ClosedAt,
		},
		ShiftCommonAttributes: shiftentity.ShiftCommonAttributes{
			CurrentOrderNumber: s.CurrentOrderNumber,
			StartChange:        s.StartChange,
			EndChange:          s.EndChange,
			AttendantID:        s.AttendantID,
			Redeems:            s.Redeems,
			Orders:             []orderentity.Order{},
			Attendant:          s.Attendant.ToDomain(),
		},
	}

	for _, order := range s.Orders {
		shift.Orders = append(shift.Orders, *order.ToDomain())
	}

	return shift
}
