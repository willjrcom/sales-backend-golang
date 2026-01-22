package orderdto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	employeedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/employee"
	groupitemdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/group_item"
	orderdeliverydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_delivery"
	orderpickupdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_pickup"
	ordertabledto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_table"
)

type OrderDTO struct {
	OrderType
	OrderTimeLogs
	OrderDetail
	ID          uuid.UUID                   `json:"id"`
	CreatedAt   time.Time                   `json:"created_at"`
	OrderNumber int                         `json:"order_number"`
	Status      orderentity.StatusOrder     `json:"status"`
	GroupsItems []groupitemdto.GroupItemDTO `json:"group_items"`
	Payments    []PaymentOrderDTO           `json:"payments"`
}

type OrderDetail struct {
	TotalPayable  decimal.Decimal          `json:"total_payable"`
	TotalPaid     decimal.Decimal          `json:"total_paid"`
	TotalChange   decimal.Decimal          `json:"total_change"`
	QuantityItems float64                  `json:"quantity_items"`
	Observation   string                   `json:"observation"`
	AttendantID   *uuid.UUID               `json:"attendant_id"`
	Attendant     *employeedto.EmployeeDTO `json:"attendant"`
	ShiftID       uuid.UUID                `json:"shift_id"`
}

type OrderType struct {
	Delivery *orderdeliverydto.OrderDeliveryDTO `json:"delivery"`
	Table    *ordertabledto.OrderTableDTO       `json:"table"`
	Pickup   *orderpickupdto.OrderPickupDTO     `json:"pickup"`
}

type OrderTimeLogs struct {
	PendingAt  *time.Time `json:"pending_at"`
	FinishedAt *time.Time `json:"finished_at"`
	ReadyAt    *time.Time `json:"ready_at"`
	CanceledAt *time.Time `json:"canceled_at"`
	ArchivedAt *time.Time `json:"archived_at"`
}

func (o *OrderDTO) FromDomain(order *orderentity.Order) {
	if order == nil {
		return
	}
	*o = OrderDTO{
		OrderType: OrderType{
			Delivery: &orderdeliverydto.OrderDeliveryDTO{},
			Table:    &ordertabledto.OrderTableDTO{},
			Pickup:   &orderpickupdto.OrderPickupDTO{},
		},
		OrderTimeLogs: OrderTimeLogs{
			PendingAt:  order.PendingAt,
			ReadyAt:    order.ReadyAt,
			FinishedAt: order.FinishedAt,
			CanceledAt: order.CanceledAt,
			ArchivedAt: order.ArchivedAt,
		},
		OrderDetail: OrderDetail{
			TotalPayable:  order.TotalPayable,
			TotalPaid:     order.TotalPaid,
			TotalChange:   order.TotalChange,
			QuantityItems: order.QuantityItems,
			Observation:   order.Observation,
			AttendantID:   order.AttendantID,
			Attendant:     &employeedto.EmployeeDTO{},
			ShiftID:       order.ShiftID,
		},
		ID:          order.ID,
		CreatedAt:   order.CreatedAt,
		OrderNumber: order.OrderNumber,
		Status:      order.Status,
		GroupsItems: []groupitemdto.GroupItemDTO{},
		Payments:    []PaymentOrderDTO{},
	}

	o.Delivery.FromDomain(order.Delivery)
	o.Table.FromDomain(order.Table)
	o.Pickup.FromDomain(order.Pickup)
	o.Attendant.FromDomain(order.Attendant)

	for _, group := range order.GroupItems {
		groupItemDTO := groupitemdto.GroupItemDTO{}
		groupItemDTO.FromDomain(&group)
		o.GroupsItems = append(o.GroupsItems, groupItemDTO)
	}

	for _, payment := range order.Payments {
		paymentOrderDTO := PaymentOrderDTO{}
		paymentOrderDTO.FromDomain(&payment)
		o.Payments = append(o.Payments, paymentOrderDTO)
	}

	if order.Delivery == nil {
		o.Delivery = nil
	}
	if order.Table == nil {
		o.Table = nil
	}
	if order.Pickup == nil {
		o.Pickup = nil
	}
	if order.Attendant == nil {
		o.Attendant = nil
	}

	if len(order.GroupItems) == 0 {
		o.GroupsItems = nil
	}
	if len(order.Payments) == 0 {
		o.Payments = nil
	}

}
