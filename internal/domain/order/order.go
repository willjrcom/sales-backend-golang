package orderentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	groupitementity "github.com/willjrcom/sales-backend-go/internal/domain/group_item"
)

type Order struct {
	entity.Entity
	bun.BaseModel `bun:"table:orders"`
	OrderTimeLogs
	OrderCommonAttributes
}

type OrderCommonAttributes struct {
	OrderType
	OrderDetail
	Payment     *PaymentOrder
	OrderNumber int                         `bun:"order_number,notnull" json:"order_number"`
	Status      StatusOrder                 `bun:"status,notnull" json:"status"`
	Groups      []groupitementity.GroupItem `bun:"rel:has-many,join:id=order_id" json:"groups"`
}

type OrderDetail struct {
	Name        string                   `bun:"name" json:"name"`
	Observation string                   `bun:"observation" json:"observation"`
	AttendantID *uuid.UUID               `bun:"column:attendant_id,type:uuid" json:"attendant_id"`
	Attendant   *employeeentity.Employee `bun:"rel:belongs-to" json:"attendant"`
	LaunchedAt  *time.Time               `bun:"launched_at" json:"launched_at"`
	ShiftID     *uuid.UUID               `bun:"column:shift_id,type:uuid" json:"shift_id"`
}

type OrderType struct {
	Delivery *DeliveryOrder `bun:"rel:has-one,join:id=order_id" json:"delivery"`
	Table    *TableOrder    `bun:"rel:has-one,join:id=order_id" json:"table"`
}

type OrderTimeLogs struct {
	FinishedAt  *time.Time `bun:"finished_at" json:"finished_at"`
	CancelledAt *time.Time `bun:"cancelled_at" json:"cancelled_at"`
	ArchivedAt  *time.Time `bun:"archived_at" json:"archived_at"`
}

func NewDefaultOrder() *Order {
	orderCommonAttributes := OrderCommonAttributes{
		Status: OrderStatusStaging,
	}

	return &Order{
		Entity:                entity.NewEntity(),
		OrderCommonAttributes: orderCommonAttributes,
	}
}

func (o *Order) StagingOrder() {
	o.Status = OrderStatusStaging
}

func (o *Order) LaunchOrder() {
	o.Status = OrderStatusPending
	*o.LaunchedAt = time.Now()
}

func (o *Order) ReadyOrder() {
	o.Status = OrderStatusFinished
}

func (o *Order) FinishOrder() {
	o.Status = OrderStatusFinished

	if o.Delivery != nil {
		*(*o.Delivery).DeliveredAt = time.Now()
	}
}

func (o *Order) CancelOrder() {
	o.Status = OrderStatusCanceled
}

func (o *Order) ArchiveOrder() {
	o.Status = OrderStatusArchived
}
