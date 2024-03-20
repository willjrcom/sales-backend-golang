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
	ErrQuantityNotInteger = errors.New("quantity items in group is not integer")
	ErrGroupNotStaging    = errors.New("group not staging")
	ErrGroupNotPending    = errors.New("group not pending")
	ErrGroupNotStarted    = errors.New("group not started")
	ErrGroupNotReady      = errors.New("group not ready")
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
	Size             string                  `bun:"size,notnull" json:"size"`
	Status           StatusGroupItem         `bun:"status,notnull" json:"status"`
	Total            float64                 `bun:"total" json:"total"`
	Quantity         float64                 `bun:"quantity" json:"quantity"`
	NeedPrint        bool                    `bun:"need_print" json:"need_print"`
	CategoryID       uuid.UUID               `bun:"column:category_id,type:uuid,notnull" json:"category_id"`
	Category         *productentity.Category `bun:"rel:belongs-to" json:"category,omitempty"`
	ComplementItemID *uuid.UUID              `bun:"column:complement_item_id,type:uuid,notnull" json:"complement_item_id"`
	ComplementItem   *itementity.Item        `bun:"rel:belongs-to" json:"complement_item,omitempty"`
}

type GroupItemTimeLogs struct {
	PendingAt  *time.Time `bun:"pending_at" json:"pending_at,omitempty"`
	StartedAt  *time.Time `bun:"started_at" json:"started_at,omitempty"`
	ReadyAt    *time.Time `bun:"ready_at" json:"ready_at,omitempty"`
	CanceledAt *time.Time `bun:"canceled_at" json:"canceled_at,omitempty"`
}

func NewGroupItem(groupCommonAttributes GroupCommonAttributes) *GroupItem {
	groupCommonAttributes.Status = StatusGroupStaging

	return &GroupItem{
		Entity:                entity.NewEntity(),
		GroupCommonAttributes: groupCommonAttributes,
	}
}

func (i *GroupItem) PendingGroupItem() (err error) {
	if math.Mod(i.Quantity, 1) != 0 {
		return ErrQuantityNotInteger
	}

	if i.Status != StatusGroupStaging {
		return nil
	}

	for index := range i.Items {
		if err = i.Items[index].PendingItem(); err != nil {
			return err
		}
	}

	if i.ComplementItem != nil {
		if err = i.ComplementItem.PendingItem(); err != nil {
			return err
		}
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
	i.CanceledAt = &time.Time{}
	*i.CanceledAt = time.Now()

	i.ComplementItem.CancelItem()
}

func (i *GroupItem) CalculateTotalPrice() {
	qtdItems := 0.0
	totalPrice := 0.0

	for _, item := range i.Items {
		totalPrice = item.CalculateTotalPrice()
		qtdItems += item.Quantity
		totalPrice += item.Price
	}

	if i.ComplementItem != nil {
		totalPrice += i.ComplementItem.Price
	}

	i.GroupDetails.Quantity = qtdItems
	i.GroupDetails.Total = totalPrice
}

func (i *GroupItem) CanAddItems() bool {
	if i.Status != StatusGroupStaging && i.Status != StatusGroupPending {
		return false
	}

	return true
}
