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

	TotalOrdersFinished    int                        `bun:"total_orders_finished"`
	TotalOrdersCancelled   int                        `bun:"total_orders_cancelled"`
	TotalSales             *decimal.Decimal           `bun:"total_sales,type:decimal(10,2)"`
	SalesByCategory        map[string]decimal.Decimal `bun:"sales_by_category,type:jsonb"`
	ProductsSoldByCategory map[string]float64         `bun:"products_sold_by_category,type:jsonb"`
	TotalItemsSold         float64                    `bun:"total_items_sold"`
	AverageOrderValue      *decimal.Decimal           `bun:"average_order_value,type:decimal(10,2)"`
	Payments               []PaymentOrder             `bun:"payments,type:jsonb"`
	DeliveryDrivers        []DeliveryDriverTax        `bun:"delivery_drivers,type:jsonb"`
}

type Redeem struct {
	Name  string           `bun:"name,notnull"`
	Value *decimal.Decimal `bun:"value,type:decimal(10,2),notnull"`
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
			CurrentOrderNumber:     shift.CurrentOrderNumber,
			StartChange:            shift.StartChange,
			EndChange:              shift.EndChange,
			AttendantID:            shift.AttendantID,
			Redeems:                []Redeem{},
			Orders:                 []Order{},
			Attendant:              &Employee{},
			TotalOrdersFinished:    shift.TotalOrdersFinished,
			TotalOrdersCancelled:   shift.TotalOrdersCancelled,
			TotalSales:             &shift.TotalSales,
			SalesByCategory:        shift.SalesByCategory,
			ProductsSoldByCategory: shift.ProductsSoldByCategory,
			TotalItemsSold:         shift.TotalItemsSold,
			AverageOrderValue:      &shift.AverageOrderValue,
			Payments:               []PaymentOrder{},
			DeliveryDrivers:        []DeliveryDriverTax{},
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
			Value: &redeem.Value,
		}
		s.Redeems = append(s.Redeems, r)
	}

	for _, payment := range shift.Payments {
		p := PaymentOrder{}
		p.FromDomain(&payment)
		s.Payments = append(s.Payments, p)
	}

	for _, driver := range shift.DeliveryDrivers {
		d := DeliveryDriverTax{}
		d.FromDomain(&driver)
		s.DeliveryDrivers = append(s.DeliveryDrivers, d)
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
			CurrentOrderNumber:     s.CurrentOrderNumber,
			StartChange:            s.StartChange,
			EndChange:              s.EndChange,
			AttendantID:            s.AttendantID,
			Redeems:                []shiftentity.Redeem{},
			Orders:                 []orderentity.Order{},
			Attendant:              s.Attendant.ToDomain(),
			TotalOrdersFinished:    s.TotalOrdersFinished,
			TotalOrdersCancelled:   s.TotalOrdersCancelled,
			TotalSales:             s.getTotalSales(),
			SalesByCategory:        s.SalesByCategory,
			ProductsSoldByCategory: s.ProductsSoldByCategory,
			TotalItemsSold:         s.TotalItemsSold,
			AverageOrderValue:      s.getAverageOrderValue(),
			Payments:               []orderentity.PaymentOrder{},
			DeliveryDrivers:        []shiftentity.DeliveryDriverTax{},
		},
	}

	for _, order := range s.Orders {
		shift.Orders = append(shift.Orders, *order.ToDomain())
	}

	for _, redeem := range s.Redeems {
		shift.Redeems = append(shift.Redeems, shiftentity.Redeem{
			Name:  redeem.Name,
			Value: redeem.getValue(),
		})
	}

	for _, payment := range s.Payments {
		shift.Payments = append(shift.Payments, *payment.ToDomain())
	}

	for _, driver := range s.DeliveryDrivers {
		shift.DeliveryDrivers = append(shift.DeliveryDrivers, *driver.ToDomain())
	}

	return shift
}

func (s *Shift) getTotalSales() decimal.Decimal {
	if s.TotalSales == nil {
		return decimal.Zero
	}
	return *s.TotalSales
}

func (s *Shift) getAverageOrderValue() decimal.Decimal {
	if s.AverageOrderValue == nil {
		return decimal.Zero
	}
	return *s.AverageOrderValue
}

func (r *Redeem) getValue() decimal.Decimal {
	if r.Value == nil {
		return decimal.Zero
	}
	return *r.Value
}
