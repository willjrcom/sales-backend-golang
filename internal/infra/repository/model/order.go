package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type Order struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:orders,alias:order"`
	OrderTimeLogs
	OrderCommonAttributes
}

type OrderCommonAttributes struct {
	OrderType
	OrderDetail
	OrderNumber int            `bun:"order_number,notnull"`
	Status      string         `bun:"status,notnull"`
	GroupItems  []GroupItem    `bun:"rel:has-many,join:id=order_id"`
	Payments    []PaymentOrder `bun:"rel:has-many,join:id=order_id"`
}

type OrderDetail struct {
	TotalPayable  *decimal.Decimal `bun:"total_payable,type:decimal(10,2)"`
	TotalPaid     *decimal.Decimal `bun:"total_paid,type:decimal(10,2)"`
	TotalChange   *decimal.Decimal `bun:"total_change,type:decimal(10,2)"`
	QuantityItems float64          `bun:"quantity_items"`
	Observation   string           `bun:"observation"`
	AttendantID   *uuid.UUID       `bun:"column:attendant_id,type:uuid"`
	Attendant     *Employee        `bun:"rel:belongs-to"`
	ShiftID       uuid.UUID        `bun:"column:shift_id,type:uuid,notnull"`
}

type OrderType struct {
	Delivery *OrderDelivery `bun:"rel:has-one,join:id=order_id"`
	Table    *OrderTable    `bun:"rel:has-one,join:id=order_id"`
	Pickup   *OrderPickup   `bun:"rel:has-one,join:id=order_id"`
}

type OrderTimeLogs struct {
	PendingAt   *time.Time `bun:"pending_at"`
	FinishedAt  *time.Time `bun:"finished_at"`
	ReadyAt     *time.Time `bun:"ready_at"`
	CancelledAt *time.Time `bun:"cancelled_at"`
	ArchivedAt  *time.Time `bun:"archived_at"`
}

func (o *Order) FromDomain(order *orderentity.Order) {
	if order == nil {
		return
	}
	*o = Order{
		Entity: entitymodel.FromDomain(order.Entity),
		OrderCommonAttributes: OrderCommonAttributes{
			OrderNumber: order.OrderNumber,
			Status:      string(order.Status),
			GroupItems:  []GroupItem{},
			Payments:    []PaymentOrder{},
			OrderType: OrderType{
				Delivery: &OrderDelivery{},
				Table:    &OrderTable{},
				Pickup:   &OrderPickup{},
			},
			OrderDetail: OrderDetail{
				Attendant:     &Employee{},
				TotalPayable:  &order.TotalPayable,
				TotalPaid:     &order.TotalPaid,
				TotalChange:   &order.TotalChange,
				QuantityItems: order.QuantityItems,
				Observation:   order.Observation,
				AttendantID:   order.AttendantID,
				ShiftID:       order.ShiftID,
			},
		},
		OrderTimeLogs: OrderTimeLogs{
			PendingAt:   order.PendingAt,
			FinishedAt:  order.FinishedAt,
			ReadyAt:     order.ReadyAt,
			CancelledAt: order.CancelledAt,
			ArchivedAt:  order.ArchivedAt,
		},
	}

	for i := range order.GroupItems {
		o.GroupItems = append(o.GroupItems, GroupItem{})
		o.GroupItems[i].FromDomain(&order.GroupItems[i])
	}

	for i := range order.Payments {
		o.Payments = append(o.Payments, PaymentOrder{})
		o.Payments[i].FromDomain(&order.Payments[i])
	}

	o.OrderType.Delivery.FromDomain(order.Delivery)
	o.OrderType.Table.FromDomain(order.Table)
	o.OrderType.Pickup.FromDomain(order.Pickup)
	o.OrderDetail.Attendant.FromDomain(order.Attendant)

	if order.Delivery == nil {
		o.OrderType.Delivery = nil
	}
	if order.Table == nil {
		o.OrderType.Table = nil
	}
	if order.Pickup == nil {
		o.OrderType.Pickup = nil
	}
	if order.Attendant == nil {
		o.OrderDetail.Attendant = nil
	}
}

func (o *Order) ToDomain() *orderentity.Order {
	if o == nil {
		return nil
	}
	order := &orderentity.Order{
		Entity: o.Entity.ToDomain(),
		OrderCommonAttributes: orderentity.OrderCommonAttributes{
			OrderNumber: o.OrderNumber,
			Status:      orderentity.StatusOrder(o.Status),
			GroupItems:  []orderentity.GroupItem{},
			Payments:    []orderentity.PaymentOrder{},
			OrderType: orderentity.OrderType{
				Delivery: &orderentity.OrderDelivery{},
				Table:    &orderentity.OrderTable{},
				Pickup:   &orderentity.OrderPickup{},
			},
			OrderDetail: orderentity.OrderDetail{
				TotalPayable:  o.GetTotalPayable(),
				TotalPaid:     o.GetTotalPaid(),
				TotalChange:   o.GetTotalChange(),
				QuantityItems: o.QuantityItems,
				Observation:   o.Observation,
				AttendantID:   o.AttendantID,
				ShiftID:       o.ShiftID,
			},
		},
		OrderTimeLogs: orderentity.OrderTimeLogs{
			PendingAt:   o.PendingAt,
			ReadyAt:     o.ReadyAt,
			FinishedAt:  o.FinishedAt,
			CancelledAt: o.CancelledAt,
			ArchivedAt:  o.ArchivedAt,
		},
	}

	for i := range o.GroupItems {
		order.GroupItems = append(order.GroupItems, *o.GroupItems[i].ToDomain())
	}

	for i := range o.Payments {
		order.Payments = append(order.Payments, *o.Payments[i].ToDomain())
	}

	order.Delivery = o.Delivery.ToDomain()
	order.Table = o.Table.ToDomain()
	order.Pickup = o.Pickup.ToDomain()
	order.Attendant = o.Attendant.ToDomain()
	return order
}

func (o *Order) GetTotalPayable() decimal.Decimal {
	if o.TotalPayable == nil {
		return decimal.Zero
	}
	return *o.TotalPayable
}

func (o *Order) GetTotalPaid() decimal.Decimal {
	if o.TotalPaid == nil {
		return decimal.Zero
	}
	return *o.TotalPaid
}

func (o *Order) GetTotalChange() decimal.Decimal {
	if o.TotalChange == nil {
		return decimal.Zero
	}
	return *o.TotalChange
}
