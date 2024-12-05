package groupitementity

import (
	"context"
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
	bun.BaseModel `bun:"table:order_group_items"`
	GroupCommonAttributes
}

type GroupCommonAttributes struct {
	GroupDetails
	GroupItemTimeLogs
	Items   []itementity.Item `bun:"rel:has-many,join:id=group_item_id" json:"items"`
	OrderID uuid.UUID         `bun:"column:order_id,type:uuid,notnull" json:"order_id"`
}

type GroupDetails struct {
	Size             string                         `bun:"size,notnull" json:"size"`
	Status           StatusGroupItem                `bun:"status,notnull" json:"status"`
	TotalPrice       float64                        `bun:"total_price" json:"total_price"`
	Quantity         float64                        `bun:"quantity" json:"quantity"`
	NeedPrint        bool                           `bun:"need_print" json:"need_print"`
	Observation      string                         `bun:"observation" json:"observation"`
	CategoryID       uuid.UUID                      `bun:"column:category_id,type:uuid,notnull" json:"category_id"`
	Category         *productentity.ProductCategory `bun:"rel:belongs-to" json:"category,omitempty"`
	ComplementItemID *uuid.UUID                     `bun:"column:complement_item_id,type:uuid" json:",omitempty"`
	ComplementItem   *itementity.Item               `bun:"rel:belongs-to" json:"complement_item,omitempty"`
}

type GroupItemTimeLogs struct {
	StartAt    *time.Time `bun:"start_at" json:"start_at,omitempty"`
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

func (i *GroupItem) BeforeUpdate(ctx context.Context, query *bun.UpdateQuery) error {
	if model, ok := query.GetModel().Value().(*GroupItem); ok {
		model.CalculateTotalPrice()
	}
	return nil
}

func (i *GroupItem) Schedule(startAt *time.Time) (err error) {
	i.StartAt = startAt
	return nil
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
	i.CanceledAt = &time.Time{}
	*i.CanceledAt = time.Now()
}

func (i *GroupItem) CalculateTotalPrice() {
	qtdItems := 0.0
	totalPrice := 0.0

	for _, item := range i.Items {
		totalPrice += item.CalculateTotalPrice() // item + additionals
		qtdItems += item.Quantity
	}

	if i.ComplementItem != nil {
		i.ComplementItem.Quantity = qtdItems
		i.ComplementItem.TotalPrice = i.ComplementItem.Price * qtdItems
		totalPrice += i.ComplementItem.TotalPrice
	}

	i.GroupDetails.Quantity = qtdItems
	i.GroupDetails.TotalPrice = totalPrice
}

func (i *GroupItem) CanAddItems() (bool, error) {
	if i.Status != StatusGroupStaging && i.Status != StatusGroupPending {
		return false, errors.New("group not staging or pending")
	}

	return true, nil
}

func (i *GroupItem) GetDistinctProductIDs() ([]uuid.UUID, error) {
	if len(i.Items) == 0 {
		return []uuid.UUID{}, errors.New("no items in group")
	}

	productIDs := []uuid.UUID{}
	mapProduct := map[uuid.UUID]uuid.UUID{}
	for _, item := range i.Items {
		if _, ok := mapProduct[item.ProductID]; ok {
			continue
		}

		productIDs = append(productIDs, item.ProductID)
		mapProduct[item.ProductID] = item.ProductID
	}

	return productIDs, nil
}
