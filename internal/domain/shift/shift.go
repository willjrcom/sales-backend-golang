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
	CurrentOrderNumber     int
	Orders                 []orderentity.Order
	Redeems                []Redeem
	StartChange            decimal.Decimal
	EndChange              *decimal.Decimal
	AttendantID            *uuid.UUID
	Attendant              *employeeentity.Employee
	TotalOrdersFinished    int
	TotalOrdersCanceled    int
	TotalSales             decimal.Decimal
	SalesByCategory        map[string]decimal.Decimal
	ProductsSoldByCategory map[string]float64
	TotalItemsSold         float64         // soma de todas as quantidades de itens, para medir o “pulo de prato”
	AverageOrderValue      decimal.Decimal // TotalSales ÷ TotalOrders, para análise de ticket médio
	Payments               []orderentity.PaymentOrder
	DeliveryDrivers        []DeliveryDriverTax
}

type ShiftTimeLogs struct {
	OpenedAt *time.Time
	ClosedAt *time.Time
}

// NewShift creates a new shift with initial start change
func NewShift(startChange decimal.Decimal) *Shift {
	newEntity := entity.NewEntity()
	now := time.Now().UTC()
	newEntity.CreatedAt = now

	shift := &Shift{
		Entity: newEntity,
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
func (s *Shift) CloseShift(endChange decimal.Decimal) {
	now := time.Now().UTC()
	s.EndChange = &endChange
	s.ClosedAt = &now

	s.Load()
}

func (s *Shift) Load() {
	// compute analytics for reporting
	s.TotalOrdersFinished = 0
	s.TotalOrdersCanceled = 0
	s.TotalSales = decimal.Zero

	// initialize maps
	s.SalesByCategory = make(map[string]decimal.Decimal)
	s.ProductsSoldByCategory = make(map[string]float64)
	s.TotalItemsSold = 0
	s.Payments = make([]orderentity.PaymentOrder, 0)
	s.DeliveryDrivers = make([]DeliveryDriverTax, 0)

	// aggregate orders data
	for _, o := range s.Orders {
		if o.Status == orderentity.OrderStatusCanceled {
			s.TotalOrdersCanceled++
			continue
		}

		if o.Status != orderentity.OrderStatusFinished {
			continue
		}

		s.TotalOrdersFinished++
		s.Payments = append(s.Payments, o.Payments...)

		if o.Delivery != nil && o.Delivery.Driver != nil && o.Delivery.Driver.Employee != nil && o.Delivery.Driver.Employee.User != nil {
			deliveryDriverTax := DeliveryDriverTax{
				DeliveryDriverID:   o.Delivery.Driver.ID,
				DeliveryDriverName: o.Delivery.Driver.Employee.User.Name,
				OrderNumber:        o.OrderNumber,
				DeliveryID:         o.Delivery.ID,
				DeliveryTax:        *o.Delivery.DeliveryTax,
			}

			s.DeliveryDrivers = append(s.DeliveryDrivers, deliveryDriverTax)
		}

		// ensure totals are up to date
		s.TotalSales = s.TotalSales.Add(o.TotalPayable)
		for _, g := range o.GroupItems {
			cat := ""
			if g.Category != nil {
				cat = g.Category.Name
			}

			// sum revenue by category
			rev := g.TotalPrice
			if prev, ok := s.SalesByCategory[cat]; ok {
				s.SalesByCategory[cat] = prev.Add(rev)
			} else {
				s.SalesByCategory[cat] = rev
			}

			// sum quantities by category and total items
			qty := g.Quantity
			s.ProductsSoldByCategory[cat] += qty
			s.TotalItemsSold += qty
		}
	}

	s.AverageOrderValue = decimal.Zero
	// average order value
	if s.TotalOrdersFinished > 0 {
		s.AverageOrderValue = s.TotalSales.Div(decimal.NewFromInt(int64(s.TotalOrdersFinished)))
	}
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
