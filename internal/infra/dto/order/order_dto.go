package orderdto

import (
	"time"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	employeedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/employee"
	groupitemdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/group_item"
	orderdeliverydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_delivery"
	orderpickupdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_pickup"
	ordertabledto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_table"
)

type OrderDTO struct {
	OrderType
	OrderDetail
	ID          uuid.UUID                   `json:"id"`
	OrderNumber int                         `json:"order_number"`
	Status      orderentity.StatusOrder     `json:"status"`
	Groups      []groupitemdto.GroupItemDTO `json:"groups"`
	Payments    []PaymentOrderDTO           `json:"payments"`
}

type OrderDetail struct {
	TotalPayable  float64                  `json:"total_payable"`
	TotalPaid     float64                  `json:"total_paid"`
	TotalChange   float64                  `json:"total_change"`
	QuantityItems float64                  `json:"quantity_items"`
	Observation   string                   `json:"observation"`
	AttendantID   *uuid.UUID               `json:"attendant_id"`
	Attendant     *employeedto.EmployeeDTO `json:"attendant"`
	ShiftID       *uuid.UUID               `json:"shift_id"`
}

type OrderType struct {
	Delivery *orderdeliverydto.OrderDeliveryDTO `json:"delivery"`
	Table    *ordertabledto.OrderTableDTO       `json:"table"`
	Pickup   *orderpickupdto.OrderPickupDTO     `json:"pickup"`
}

type OrderTimeLogs struct {
	PendingAt  *time.Time `json:"pending_at"`
	FinishedAt *time.Time `json:"finished_at"`
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
		OrderNumber: order.OrderNumber,
		Status:      order.Status,
		Groups:      []groupitemdto.GroupItemDTO{},
		Payments:    []PaymentOrderDTO{},
	}

	o.Delivery.FromDomain(order.Delivery)
	o.Table.FromDomain(order.Table)
	o.Pickup.FromDomain(order.Pickup)
	o.Attendant.FromDomain(order.Attendant)

	for _, group := range order.Groups {
		groupItemDTO := groupitemdto.GroupItemDTO{}
		groupItemDTO.FromDomain(&group)
		o.Groups = append(o.Groups, groupItemDTO)
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

	if len(order.Groups) == 0 {
		o.Groups = nil
	}
	if len(order.Payments) == 0 {
		o.Payments = nil
	}

}
