package orderentity

import (
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
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
	GroupItemTimeLogs
	GroupCommonAttributes
}

type GroupCommonAttributes struct {
	GroupDetails
	Items   []Item
	OrderID uuid.UUID
}

type GroupDetails struct {
	Size             string
	Status           StatusGroupItem
	TotalPrice       decimal.Decimal
	Quantity         float64
	NeedPrint        bool
	PrinterName      string
	UseProcessRule   bool
	Observation      string
	CategoryID       uuid.UUID
	Category         *productentity.ProductCategory
	ComplementItemID *uuid.UUID
	ComplementItem   *Item
}

type GroupItemTimeLogs struct {
	StartAt    *time.Time
	PendingAt  *time.Time
	StartedAt  *time.Time
	ReadyAt    *time.Time
	CanceledAt *time.Time
}

func NewGroupItem(groupCommonAttributes GroupCommonAttributes) *GroupItem {
	groupCommonAttributes.Status = StatusGroupStaging

	return &GroupItem{
		Entity:                entity.NewEntity(),
		GroupCommonAttributes: groupCommonAttributes,
	}
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
	*i.PendingAt = time.Now().UTC()
	return nil
}

func (i *GroupItem) StartGroupItem() (err error) {
	if i.Status != StatusGroupPending {
		return ErrGroupNotPending
	}

	i.Status = StatusGroupStarted
	i.StartedAt = &time.Time{}
	*i.StartedAt = time.Now().UTC()
	return nil
}

func (i *GroupItem) ReadyGroupItem() (err error) {
	if i.Status != StatusGroupStarted {
		return ErrGroupNotStarted
	}

	i.Status = StatusGroupReady
	i.ReadyAt = &time.Time{}
	*i.ReadyAt = time.Now().UTC()
	return nil
}

func (i *GroupItem) CancelGroupItem() {
	i.Status = StatusGroupCanceled
	i.CanceledAt = &time.Time{}
	*i.CanceledAt = time.Now().UTC()
}

func (i *GroupItem) CalculateTotalPrice() {
	qtdItems := 0.0
	totalPrice := decimal.Zero

	if i.Status == StatusGroupCanceled {
		i.GroupDetails.TotalPrice = decimal.Zero
		i.GroupDetails.Quantity = 0.0
		return
	}

	for _, item := range i.Items {
		totalPrice = totalPrice.Add(item.CalculateTotalPrice())
		qtdItems += item.Quantity
	}

	if i.ComplementItem != nil {
		i.ComplementItem.Quantity = qtdItems
		// price * quantity
		compTotal := i.ComplementItem.TotalPrice.Mul(decimal.NewFromFloat(qtdItems))
		i.ComplementItem.TotalPrice = compTotal
		totalPrice = totalPrice.Add(compTotal)
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
