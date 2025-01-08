package orderentity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrOrderWithoutItems            = errors.New("order must have at least one item")
	ErrOrderMustBePending           = errors.New("order must be pending")
	ErrOrderMustBePendingOrReady    = errors.New("order must be pending or ready")
	ErrOrderMustBeCanceled          = errors.New("order must be canceled")
	ErrOrderMustBeArchived          = errors.New("order must be archived")
	ErrOrderAlreadyFinished         = errors.New("order already finished")
	ErrOrderAlreadyCanceled         = errors.New("order already canceled")
	ErrOrderAlreadyArchived         = errors.New("order already archived")
	ErrOrderPaidMoreThanTotal       = errors.New("order paid more than total")
	ErrOrderPaidLessThanTotal       = errors.New("order paid less than total")
	ErrDeliveryOrderMustBeDelivered = errors.New("order delivery must be delivered")
)

type Order struct {
	entity.Entity
	OrderTimeLogs
	OrderCommonAttributes
}

type OrderCommonAttributes struct {
	OrderType
	OrderDetail
	OrderNumber int
	Status      StatusOrder
	Groups      []GroupItem
	Payments    []PaymentOrder
}

type OrderDetail struct {
	TotalPayable  float64
	TotalPaid     float64
	TotalChange   float64
	QuantityItems float64
	Observation   string
	AttendantID   *uuid.UUID
	Attendant     *employeeentity.Employee
	ShiftID       uuid.UUID
}

type OrderType struct {
	Delivery *OrderDelivery
	Table    *OrderTable
	Pickup   *OrderPickup
}

type OrderTimeLogs struct {
	PendingAt  *time.Time
	FinishedAt *time.Time
	CanceledAt *time.Time
	ArchivedAt *time.Time
}

func NewDefaultOrder(shiftID uuid.UUID, currentOrderNumber int, attendantID *uuid.UUID) *Order {
	order := &Order{
		Entity: entity.NewEntity(),
		OrderCommonAttributes: OrderCommonAttributes{
			OrderNumber: currentOrderNumber,
			OrderDetail: OrderDetail{
				ShiftID:     shiftID,
				AttendantID: attendantID,
				TotalPaid:   0,
				TotalChange: 0,
			},
		},
	}

	order.StagingOrder()
	return order
}

func (o *Order) StagingOrder() {
	o.Status = OrderStatusStaging
}

func (o *Order) PendingOrder() (err error) {
	if o.Status == OrderStatusFinished {
		return ErrOrderAlreadyFinished
	}

	if o.Status == OrderStatusCanceled {
		return ErrOrderAlreadyCanceled
	}

	if o.Status == OrderStatusArchived {
		return ErrOrderAlreadyArchived
	}

	if len(o.Groups) == 0 {
		return ErrOrderWithoutItems
	}

	for i := range o.Groups {
		if err = o.Groups[i].PendingGroupItem(); err != nil {
			return err
		}
	}

	o.Status = OrderStatusPending

	if o.PendingAt == nil {
		o.PendingAt = &time.Time{}
		*o.PendingAt = time.Now().UTC()
	}

	if o.Delivery != nil {
		if err := o.Delivery.Pend(); err != nil {
			return err
		}
	} else if o.Pickup != nil {
		if err := o.Pickup.Pend(); err != nil {
			return err
		}
	} else if o.Table != nil {
		if err := o.Table.Pend(); err != nil {
			return err
		}
	}

	return nil
}

func (o *Order) ReadyOrder() (err error) {
	if o.Status != OrderStatusPending {
		return ErrOrderMustBePending
	}

	o.Status = OrderStatusReady
	return nil
}

func (o *Order) FinishOrder() (err error) {
	if o.Status != OrderStatusReady && o.Status != OrderStatusPending {
		return ErrOrderMustBePendingOrReady
	}

	if o.Delivery != nil && o.Delivery.Status != OrderDeliveryStatusDelivered {
		return ErrDeliveryOrderMustBeDelivered
	}

	totalPaid := 0.00
	for _, payment := range o.Payments {
		totalPaid += payment.TotalPaid
	}

	if totalPaid < o.TotalPayable {
		return ErrOrderPaidLessThanTotal
	}

	o.Status = OrderStatusFinished
	o.FinishedAt = &time.Time{}
	*o.FinishedAt = time.Now().UTC()
	return nil
}

func (o *Order) CancelOrder() (err error) {
	if o.Status == OrderStatusCanceled {
		return ErrOrderAlreadyCanceled
	}

	if o.Status == OrderStatusArchived {
		return ErrOrderAlreadyArchived
	}

	for i := range o.Groups {
		o.Groups[i].CancelGroupItem()
	}

	o.Status = OrderStatusCanceled
	o.CanceledAt = &time.Time{}
	*o.CanceledAt = time.Now().UTC()
	return nil
}

func (o *Order) ArchiveOrder() (err error) {
	if o.Status != OrderStatusCanceled {
		return ErrOrderMustBeCanceled
	}

	if o.Status == OrderStatusArchived {
		return ErrOrderAlreadyArchived
	}

	o.Status = OrderStatusArchived
	o.ArchivedAt = &time.Time{}
	*o.ArchivedAt = time.Now().UTC()
	return nil
}

func (o *Order) UnarchiveOrder() (err error) {
	if o.Status != OrderStatusArchived {
		return ErrOrderMustBeArchived
	}

	if o.CanceledAt != nil {
		o.Status = OrderStatusCanceled
		return
	}

	o.Status = OrderStatusCanceled
	return
}

func (o *Order) ValidatePayments() error {
	if o.TotalPayable <= o.TotalPaid {
		return ErrOrderPaidMoreThanTotal
	}

	return nil
}

func (o *Order) AddPayment(payment *PaymentOrder) {
	o.TotalPaid += payment.TotalPaid
	o.Payments = append(o.Payments, *payment)
}

func (o *Order) CalculateTotalPrice() {
	o.TotalPayable = 0.00
	o.QuantityItems = 0.00

	for i := range o.Groups {
		o.Groups[i].CalculateTotalPrice()
		o.TotalPayable += o.Groups[i].TotalPrice
		o.QuantityItems += o.Groups[i].Quantity
	}

	o.TotalPaid = 0.00
	for _, payment := range o.Payments {
		o.TotalPaid += payment.TotalPaid
	}

	if o.Delivery != nil && o.Delivery.DeliveryTax != nil {
		o.TotalPayable += *o.Delivery.DeliveryTax
	}

	if o.TotalPayable < o.TotalPaid {
		o.TotalChange = o.TotalPaid - o.TotalPayable
	} else {
		o.TotalChange = 0
	}
}
