package orderentity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	groupitementity "github.com/willjrcom/sales-backend-go/internal/domain/group_item"
)

var (
	ErrOrderMustBeFinishedOrCanceled = errors.New("order must be canceled or finished")
	ErrOrderWithoutItems             = errors.New("order must have at least one item")
	ErrOrderMustBePending            = errors.New("order must be pending")
	ErrOrderMustBeArchived           = errors.New("order must be archived")
	ErrOrderAlreadyFinished          = errors.New("order already finished")
	ErrOrderAlreadyCanceled          = errors.New("order already canceled")
	ErrOrderAlreadyArchived          = errors.New("order already archived")
)

type Order struct {
	entity.Entity
	bun.BaseModel `bun:"table:orders,alias:order"`
	OrderTimeLogs
	OrderCommonAttributes
}

type OrderCommonAttributes struct {
	OrderType
	OrderDetail
	OrderNumber int                         `bun:"order_number,notnull" json:"order_number"`
	Status      StatusOrder                 `bun:"status,notnull" json:"status"`
	Groups      []groupitementity.GroupItem `bun:"rel:has-many,join:id=order_id" json:"groups"`
	Payments    []PaymentOrder              `bun:"rel:has-many,join:id=order_id" json:"payments,omitempty"`
}

type OrderDetail struct {
	ScheduledOrder
	Observation string                   `bun:"observation" json:"observation"`
	AttendantID *uuid.UUID               `bun:"column:attendant_id,type:uuid,notnull" json:"attendant_id"`
	Attendant   *employeeentity.Employee `bun:"rel:belongs-to" json:"attendant,omitempty"`
	ShiftID     *uuid.UUID               `bun:"column:shift_id,type:uuid" json:"shift_id"`
}

type OrderType struct {
	Delivery *DeliveryOrder `bun:"rel:has-one,join:id=order_id" json:"delivery,omitempty"`
	Table    *TableOrder    `bun:"rel:has-one,join:id=order_id" json:"table,omitempty"`
}

type ScheduledOrder struct {
	StartAt *time.Time `bun:"start_at" json:"start_at,omitempty"`
}

type OrderTimeLogs struct {
	PendingAt   *time.Time `bun:"pending_at" json:"pending_at,omitempty"`
	FinishedAt  *time.Time `bun:"finished_at" json:"finished_at,omitempty"`
	CancelledAt *time.Time `bun:"cancelled_at" json:"cancelled_at,omitempty"`
	ArchivedAt  *time.Time `bun:"archived_at" json:"archived_at,omitempty"`
}

func NewDefaultOrder(shiftID *uuid.UUID, currentOrderNumber int, attendantID *uuid.UUID) *Order {
	order := &Order{
		Entity: entity.NewEntity(),
		OrderCommonAttributes: OrderCommonAttributes{
			OrderNumber: currentOrderNumber,
			OrderDetail: OrderDetail{
				ShiftID:     shiftID,
				AttendantID: attendantID,
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

	if len(o.Groups) == 0 {
		return ErrOrderWithoutItems
	}

	for i := range o.Groups {
		if err = o.Groups[i].PendingGroupItem(); err != nil {
			return err
		}
	}

	o.Status = OrderStatusPending
	o.PendingAt = &time.Time{}
	*o.PendingAt = time.Now()
	return nil
}

func (o *Order) FinishOrder() (err error) {
	if o.Status != OrderStatusPending {
		return ErrOrderMustBePending
	}

	if o.Status == OrderStatusFinished {
		return ErrOrderAlreadyFinished
	}

	o.Status = OrderStatusFinished
	o.FinishedAt = &time.Time{}
	*o.FinishedAt = time.Now()
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
	o.CancelledAt = &time.Time{}
	*o.CancelledAt = time.Now()
	return nil
}

func (o *Order) ArchiveOrder() (err error) {
	if o.Status != OrderStatusFinished && o.Status != OrderStatusCanceled {
		return ErrOrderMustBeFinishedOrCanceled
	}

	if o.Status == OrderStatusArchived {
		return ErrOrderAlreadyArchived
	}

	o.Status = OrderStatusArchived
	o.ArchivedAt = &time.Time{}
	*o.ArchivedAt = time.Now()
	return nil
}

func (o *Order) UnarchiveOrder() (err error) {
	if o.Status != OrderStatusArchived {
		return ErrOrderMustBeArchived
	}

	if o.CancelledAt != nil {
		o.Status = OrderStatusCanceled
		return
	}

	o.Status = OrderStatusFinished
	return
}

func (o *Order) ScheduleOrder(startAt *time.Time) {
	o.StartAt = startAt
}
