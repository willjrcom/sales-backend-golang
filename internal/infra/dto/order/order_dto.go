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
	GroupItems  []groupitemdto.GroupItemDTO `json:"group_items"`
	Payments    []PaymentOrderDTO           `json:"payments"`
	Fees        []AdditionalFee             `json:"fees"`
}

type OrderDetail struct {
	SubTotal      decimal.Decimal          `json:"sub_total"`
	Total         decimal.Decimal          `json:"total"`
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
	PendingAt   *time.Time `json:"pending_at"`
	FinishedAt  *time.Time `json:"finished_at"`
	ReadyAt     *time.Time `json:"ready_at"`
	CancelledAt *time.Time `json:"cancelled_at"`
	ArchivedAt  *time.Time `json:"archived_at"`
}

type AdditionalFee struct {
	Name  string          `json:"name"`
	Value decimal.Decimal `json:"value"`
}

func (o *OrderDTO) FromDomain(order *orderentity.Order) {
	if order == nil {
		return
	}
	*o = OrderDTO{
		OrderType: OrderType{},
		OrderTimeLogs: OrderTimeLogs{
			PendingAt:   order.PendingAt,
			ReadyAt:     order.ReadyAt,
			FinishedAt:  order.FinishedAt,
			CancelledAt: order.CancelledAt,
			ArchivedAt:  order.ArchivedAt,
		},
		OrderDetail: OrderDetail{
			SubTotal:      order.SubTotal,
			Total:         order.Total,
			TotalPaid:     order.TotalPaid,
			TotalChange:   order.TotalChange,
			QuantityItems: order.QuantityItems,
			Observation:   order.Observation,
			AttendantID:   order.AttendantID,
			ShiftID:       order.ShiftID,
		},
		ID:          order.ID,
		CreatedAt:   order.CreatedAt,
		OrderNumber: order.OrderNumber,
		Status:      order.Status,
	}

	if order.Delivery != nil {
		o.Delivery = &orderdeliverydto.OrderDeliveryDTO{}
		o.Delivery.FromDomain(order.Delivery)
	}
	if order.Table != nil {
		o.Table = &ordertabledto.OrderTableDTO{}
		o.Table.FromDomain(order.Table)
	}
	if order.Pickup != nil {
		o.Pickup = &orderpickupdto.OrderPickupDTO{}
		o.Pickup.FromDomain(order.Pickup)
	}
	if order.Attendant != nil {
		o.Attendant = &employeedto.EmployeeDTO{}
		o.Attendant.FromDomain(order.Attendant)
	}

	if len(order.GroupItems) > 0 {
		o.GroupItems = make([]groupitemdto.GroupItemDTO, len(order.GroupItems))
		for i := range order.GroupItems {
			o.GroupItems[i].FromDomain(&order.GroupItems[i])
		}
	}

	if len(order.Payments) > 0 {
		o.Payments = make([]PaymentOrderDTO, len(order.Payments))
		for i := range order.Payments {
			o.Payments[i].FromDomain(&order.Payments[i])
		}
	}

	if len(order.Fees) > 0 {
		o.Fees = make([]AdditionalFee, len(order.Fees))
		for i := range order.Fees {
			o.Fees[i] = AdditionalFee{
				Name:  string(order.Fees[i].Name),
				Value: order.Fees[i].Value,
			}
		}
	}

	if order.Delivery != nil {
		o.Delivery = &orderdeliverydto.OrderDeliveryDTO{}
		o.Delivery.FromDomain(order.Delivery)
	}
	if order.Table != nil {
		o.Table = &ordertabledto.OrderTableDTO{}
		o.Table.FromDomain(order.Table)
	}
	if order.Pickup != nil {
		o.Pickup = &orderpickupdto.OrderPickupDTO{}
		o.Pickup.FromDomain(order.Pickup)
	}
	if order.Attendant != nil {
		o.Attendant = &employeedto.EmployeeDTO{}
		o.Attendant.FromDomain(order.Attendant)
	}
}
