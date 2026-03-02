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
	ErrOrderMustBeCancelled         = errors.New("order must be cancelled")
	ErrOrderMustBeArchived          = errors.New("order must be archived")
	ErrOrderAlreadyFinished         = errors.New("order already finished")
	ErrOrderAlreadyCancelled        = errors.New("order already cancelled")
	ErrOrderAlreadyArchived         = errors.New("order already archived")
	ErrOrderPaidMoreThanTotal       = errors.New("order paid more than total")
	ErrOrderPaidLessThanTotal       = errors.New("order paid less than total")
	ErrDeliveryOrderMustBeDelivered = errors.New("order delivery must be delivered")
	ErrOrderTableMustBeClosed       = errors.New("order table must be closed")
	ErrOrderPickupMustBeReady       = errors.New("order pickup must be ready")
	ErrOrderPickupMustBeDelivered   = errors.New("order pickup must be delivered")
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
	Fees        []AdditionalFee
}

type OrderDetail struct {
	SubTotal      decimal.Decimal
	Total         decimal.Decimal
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
	PendingAt   *time.Time
	ReadyAt     *time.Time
	FinishedAt  *time.Time
	CancelledAt *time.Time
	ArchivedAt  *time.Time
}

type AdditionalFeeName string

const (
	AdditionalFeeTypeTableTax    AdditionalFeeName = "table_tax"
	AdditionalFeeTypeDeliveryFee AdditionalFeeName = "delivery_fee"
)

type AdditionalFee struct {
	Name  AdditionalFeeName
	Value decimal.Decimal
}

func NewDefaultOrder(shiftID uuid.UUID, currentOrderNumber int, attendantID *uuid.UUID) *Order {
	order := &Order{
		Entity: entity.NewEntity(),
		OrderCommonAttributes: OrderCommonAttributes{
			OrderNumber: currentOrderNumber,
			Fees:        []AdditionalFee{},
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

	if o.Status == OrderStatusCancelled {
		return ErrOrderAlreadyCancelled
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

	o.ReadyAt = &time.Time{}
	*o.ReadyAt = time.Now().UTC()
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
	} else if o.Pickup != nil && o.Pickup.Status != OrderPickupStatusDelivered {
		return ErrOrderPickupMustBeDelivered
	}

	totalPaid := decimal.Zero
	for _, payment := range o.Payments {
		totalPaid = totalPaid.Add(payment.TotalPaid)
	}

	if totalPaid.LessThan(o.Total) {
		return ErrOrderPaidLessThanTotal
	}

	o.Status = OrderStatusFinished
	o.FinishedAt = &time.Time{}
	*o.FinishedAt = time.Now().UTC()
	return nil
}

func (o *Order) CancelOrder() (err error) {
	if o.Status == OrderStatusCancelled {
		return ErrOrderAlreadyCancelled
	}

	if o.Status == OrderStatusArchived {
		return ErrOrderAlreadyArchived
	}

	for i := range o.GroupItems {
		o.GroupItems[i].CancelGroupItem()
	}

	o.Status = OrderStatusCancelled
	o.CancelledAt = &time.Time{}
	*o.CancelledAt = time.Now().UTC()
	return nil
}

func (o *Order) ArchiveOrder() (err error) {
	if o.Status != OrderStatusCancelled {
		return ErrOrderMustBeCancelled
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

	if o.CancelledAt != nil {
		o.Status = OrderStatusCancelled
		return
	}

	o.Status = OrderStatusCancelled
	return
}

func (o *Order) ValidatePayments() error {
	// Error if the total paid exceeds the total payable
	if o.TotalPaid.GreaterThan(o.Total) {
		return ErrOrderPaidMoreThanTotal
	}

	return nil
}

func (o *Order) AddPayment(payment *PaymentOrder) {
	o.TotalPaid = o.TotalPaid.Add(payment.TotalPaid)
	o.Payments = append(o.Payments, *payment)
}

func (o *Order) CalculateTotalOrder() {
	o.CalculateSubTotal()
	o.CalculateFees()
	o.CalculateTotal()
	o.CalculateTotalPaid()
}

func (o *Order) CalculateSubTotal() {
	o.SubTotal = decimal.Zero
	o.QuantityItems = 0.0

	for i := range o.GroupItems {
		o.GroupItems[i].CalculateTotal()
		o.SubTotal = o.SubTotal.Add(o.GroupItems[i].Total)
		o.QuantityItems += o.GroupItems[i].Quantity
	}

	o.SubTotal = o.SubTotal.Round(2)
}

func (o *Order) CalculateFees() {
	o.Fees = []AdditionalFee{}

	if o.Table != nil && !o.Table.TaxRate.IsZero() {
		taxRate := o.SubTotal.Mul(o.Table.TaxRate.Div(decimal.NewFromInt(100)))
		o.Fees = append(o.Fees, AdditionalFee{
			Name:  AdditionalFeeTypeTableTax,
			Value: taxRate.Round(2),
		})
	}

	if o.Delivery != nil && o.Delivery.DeliveryTax != nil && !o.Delivery.IsDeliveryFree {
		o.Fees = append(o.Fees, AdditionalFee{
			Name:  AdditionalFeeTypeDeliveryFee,
			Value: o.Delivery.DeliveryTax.Round(2),
		})
	}
}

func (o *Order) CalculateTotal() {
	totalFees := decimal.Zero
	for _, fee := range o.Fees {
		totalFees = totalFees.Add(fee.Value)
	}

	o.Total = o.SubTotal.Add(totalFees).Round(2)
}

func (o *Order) CalculateTotalPaid() {
	o.TotalPaid = decimal.Zero
	for _, payment := range o.Payments {
		o.TotalPaid = o.TotalPaid.Add(payment.TotalPaid)
	}

	o.TotalPaid = o.TotalPaid.Round(2)

	if o.TotalPaid.GreaterThan(o.Total) {
		o.TotalChange = o.TotalPaid.Sub(o.Total)
	} else {
		o.TotalChange = decimal.Zero
	}

	o.TotalChange = o.TotalChange.Round(2)
}
