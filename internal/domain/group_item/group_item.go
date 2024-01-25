package groupitementity

import (
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	ErrQuantityNotInteger   = errors.New("quantity items in group isnot integer")
	ErrGroupNotPending      = errors.New("group not pending")
	ErrGroupNotStarted      = errors.New("group not started")
	ErrGroupNotReady        = errors.New("group not ready")
	ErrGroupAlreadyFinished = errors.New("group already finished")
	ErrOnlyCancelAllowed    = errors.New("only cancel allowed")
)

type GroupItem struct {
	entity.Entity
	bun.BaseModel `bun:"table:group_items"`
	GroupCommonAttributes
}

type GroupCommonAttributes struct {
	GroupDetails
	GroupItemTimeLogs
	Items   []itementity.Item `bun:"rel:has-many,join:id=group_item_id" json:"items"`
	OrderID uuid.UUID         `bun:"column:order_id,type:uuid,notnull" json:"order_id"`
}

type GroupDetails struct {
	Size       string                  `bun:"size,notnull" json:"size"`
	Status     StatusGroupItem         `bun:"status,notnull" json:"status"`
	Price      float64                 `bun:"price" json:"price"`
	Quantity   float64                 `bun:"quantity" json:"quantity"`
	CategoryID uuid.UUID               `bun:"column:category_id,type:uuid,notnull" json:"category_id"`
	Category   *productentity.Category `bun:"rel:belongs-to" json:"category,omitempty"`
}

type GroupItemTimeLogs struct {
	PendingAt   *time.Time `bun:"pending_at" json:"pending_at,omitempty"`
	StartedAt   *time.Time `bun:"started_at" json:"started_at,omitempty"`
	ReadyAt     *time.Time `bun:"ready_at" json:"ready_at,omitempty"`
	CancelledAt *time.Time `bun:"cancelled_at" json:"cancelled_at,omitempty"`
}

func (i *GroupItem) PendingGroupItem() (err error) {
	if math.Mod(i.Quantity, 1) != 0 {
		return ErrQuantityNotInteger
	}

	if i.Status != StatusGroupStaging {
		return nil
	}

	i.Status = StatusGroupPending
	i.PendingAt = &time.Time{}
	*i.PendingAt = time.Now()
	return nil
}

func (i *GroupItem) StartGroupItem() (err error) {
	if i.Status != StatusGroupPending {
		return ErrGroupNotPending
	}

	i.Status = StatusGroupStarted
	i.StartedAt = &time.Time{}
	*i.StartedAt = time.Now()
	return nil
}

func (i *GroupItem) ReadyGroupItem() (err error) {
	if i.Status != StatusGroupStarted {
		return ErrGroupNotStarted
	}

	i.Status = StatusGroupReady
	i.ReadyAt = &time.Time{}
	*i.ReadyAt = time.Now()
	return nil
}

func (i *GroupItem) CancelGroupItem() {
	i.Status = StatusGroupCanceled
	i.CancelledAt = &time.Time{}
	*i.CancelledAt = time.Now()
}

func (i *GroupItem) DeleteGroupItem() (err error) {
	if i.Status != StatusGroupPending && i.Status != StatusGroupStaging {
		return ErrOnlyCancelAllowed
	}

	return nil
}
func (i *GroupItem) CalculateTotalValues() {
	qtd := 0.0
	price := 0.0

	for _, item := range i.Items {
		qtd += item.Quantity
		price += item.Price * item.Quantity
	}

	i.GroupDetails.Quantity = qtd
	i.GroupDetails.Price = price
}