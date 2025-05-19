package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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
	CurrentOrderNumber int              `bun:"current_order_number,notnull"`
	Orders             []Order          `bun:"rel:has-many,join:id=shift_id"`
	Redeems            []Redeem         `bun:"redeems,type:jsonb"`
	StartChange        decimal.Decimal  `bun:"start_change,type:decimal(10,2)"`
	EndChange          *decimal.Decimal `bun:"end_change,type:decimal(10,2)"`
	AttendantID        *uuid.UUID       `bun:"column:attendant_id,type:uuid"`
	Attendant          *Employee        `bun:"rel:belongs-to"`
}

type Redeem struct {
	Name  string          `bun:"name,notnull"`
	Value decimal.Decimal `bun:"value,type:decimal(10,2),notnull"`
}

type ShiftTimeLogs struct {
	OpenedAt *time.Time `bun:"opened_at"`
	ClosedAt *time.Time `bun:"closed_at"`
}

func (s *Shift) FromDomain(shift *shiftentity.Shift) {
	if shift == nil {
		return
	}
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
			Redeems:            []Redeem{},
			Orders:             []Order{},
			Attendant:          &Employee{},
		},
	}

	for _, order := range shift.Orders {
		o := Order{}
		o.FromDomain(&order)
		s.Orders = append(s.Orders, o)
	}

	for _, redeem := range shift.Redeems {
		r := Redeem{
			Name:  redeem.Name,
			Value: redeem.Value,
		}
		s.Redeems = append(s.Redeems, r)
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
			Redeems:            []shiftentity.Redeem{},
			Orders:             []orderentity.Order{},
			Attendant:          s.Attendant.ToDomain(),
		},
	}

	for _, order := range s.Orders {
		shift.Orders = append(shift.Orders, *order.ToDomain())
	}

	for _, redeem := range s.Redeems {
		shift.Redeems = append(shift.Redeems, shiftentity.Redeem{
			Name:  redeem.Name,
			Value: redeem.Value,
		})
	}

	return shift
}
