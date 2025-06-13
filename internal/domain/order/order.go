package orderentity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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
	ErrOrderTableMustBeClosed       = errors.New("order table must be closed")
	ErrOrderPickupMustBeReady       = errors.New("order pickup must be ready")
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
	GroupItems  []GroupItem
	Payments    []PaymentOrder
}

type OrderDetail struct {
	TotalPayable  decimal.Decimal
	TotalPaid     decimal.Decimal
	TotalChange   decimal.Decimal
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
				TotalPaid:   decimal.Zero,
				TotalChange: decimal.Zero,
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

	if len(o.GroupItems) == 0 {
		return ErrOrderWithoutItems
	}

	for i := range o.GroupItems {
		if err = o.GroupItems[i].PendingGroupItem(); err != nil {
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

	if o.Delivery != nil {
		if err := o.Delivery.Ready(); err != nil {
			return err
		}
	} else if o.Pickup != nil {
		if err := o.Pickup.Ready(); err != nil {
			return err
		}
	}

	return nil
}

func (o *Order) FinishOrder() (err error) {
	if o.Status != OrderStatusReady && o.Status != OrderStatusPending {
		return ErrOrderMustBePendingOrReady
	}

	if o.Delivery != nil && o.Delivery.Status != OrderDeliveryStatusDelivered {
		return ErrDeliveryOrderMustBeDelivered
	} else if o.Table != nil && o.Table.Status != OrderTableStatusClosed {
		return ErrOrderTableMustBeClosed
	} else if o.Pickup != nil && o.Pickup.Status != OrderPickupStatusReady {
		return ErrOrderPickupMustBeReady
	}

	totalPaid := decimal.Zero
	for _, payment := range o.Payments {
		totalPaid = totalPaid.Add(payment.TotalPaid)
	}

	if totalPaid.LessThan(o.TotalPayable) {
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

	for i := range o.GroupItems {
		o.GroupItems[i].CancelGroupItem()
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
	// Error if the total paid exceeds the total payable
	if o.TotalPaid.GreaterThan(o.TotalPayable) {
		return ErrOrderPaidMoreThanTotal
	}

	return nil
}

func (o *Order) AddPayment(payment *PaymentOrder) {
	o.TotalPaid = o.TotalPaid.Add(payment.TotalPaid)
	o.Payments = append(o.Payments, *payment)
}

func (o *Order) CalculateTotalPrice() {
	o.TotalPayable = decimal.Zero
	o.QuantityItems = 0.0

	for i := range o.GroupItems {
		o.GroupItems[i].CalculateTotalPrice()
		o.TotalPayable = o.TotalPayable.Add(o.GroupItems[i].TotalPrice)
		o.QuantityItems += o.GroupItems[i].Quantity
	}

	o.TotalPaid = decimal.Zero
	for _, payment := range o.Payments {
		o.TotalPaid = o.TotalPaid.Add(payment.TotalPaid)
	}

	if o.Table != nil && !o.Table.TaxRate.IsZero() {
		taxRate := o.TotalPayable.Mul(o.Table.TaxRate.Div(decimal.NewFromInt(100)))
		o.TotalPayable = o.TotalPayable.Add(taxRate)
	}

	if o.Delivery != nil && o.Delivery.DeliveryTax != nil {
		o.TotalPayable = o.TotalPayable.Add(*o.Delivery.DeliveryTax)
	}

	if o.TotalPaid.GreaterThan(o.TotalPayable) {
		o.TotalChange = o.TotalPaid.Sub(o.TotalPayable)
	} else {
		o.TotalChange = decimal.Zero
	}
}
