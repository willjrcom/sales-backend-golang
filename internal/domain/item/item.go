package itementity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrItemNotPending    = errors.New("item not pending")
	ErrItemNotStarted    = errors.New("item not started")
	ErrItemNotReady      = errors.New("item not ready")
	ErrOnlyCancelAllowed = errors.New("only cancel allowed")
)

type Item struct {
	entity.Entity
	bun.BaseModel `bun:"table:items"`
	ItemTimeLogs
	ItemCommonAttributes
}

type ItemCommonAttributes struct {
	Name            string     `bun:"name,notnull" json:"name"`
	Status          StatusItem `bun:"status,notnull" json:"status"`
	Description     string     `bun:"description" json:"description"`
	Observation     string     `bun:"observation" json:"observation"`
	Price           float64    `bun:"price,notnull" json:"price"`
	TotalPrice      float64    `bun:"total_price,notnull" json:"total_price"`
	Size            string     `bun:"size,notnull" json:"size"`
	Quantity        float64    `bun:"quantity,notnull" json:"quantity"`
	GroupItemID     uuid.UUID  `bun:"group_item_id,type:uuid" json:"group_item_id"`
	AdditionalItems []Item     `bun:"m2m:item_to_additional,join:Item=AdditionalItem" json:"item_to_additional,omitempty"`
	ProductID       uuid.UUID  `bun:"product_id,type:uuid" json:"product_id"`
}

type ItemTimeLogs struct {
	PendingAt  *time.Time `bun:"pending_at" json:"pending_at,omitempty"`
	StartedAt  *time.Time `bun:"started_at" json:"started_at,omitempty"`
	ReadyAt    *time.Time `bun:"ready_at" json:"ready_at,omitempty"`
	CanceledAt *time.Time `bun:"canceled_at" json:"canceled_at,omitempty"`
}

func NewItem(name string, price float64, quantity float64, size string, status StatusItem, productID uuid.UUID) *Item {
	itemAdditionalCommonAttributes := ItemCommonAttributes{
		Name:       name + " (" + size + ")",
		Status:     status,
		Price:      price,
		TotalPrice: price * quantity,
		Size:       size,
		Quantity:   quantity,
		ProductID:  productID,
	}

	return &Item{
		Entity:               entity.NewEntity(),
		ItemCommonAttributes: itemAdditionalCommonAttributes,
	}
}

func (i *Item) PendingItem() (err error) {
	if i.Status != StatusItemStaging {
		return nil
	}

	for index := range i.AdditionalItems {
		if err = i.AdditionalItems[index].PendingItem(); err != nil {
			return err
		}
	}

	i.Status = StatusItemPending
	i.PendingAt = &time.Time{}
	*i.PendingAt = time.Now()
	return nil
}

func (i *Item) StartItem() (err error) {
	if i.Status != StatusItemPending {
		return ErrItemNotPending
	}

	i.Status = StatusItemStarted
	i.StartedAt = &time.Time{}
	*i.StartedAt = time.Now()
	return nil
}

func (i *Item) ReadyItem() (err error) {
	if i.Status != StatusItemStarted {
		return ErrItemNotStarted
	}

	i.Status = StatusItemReady
	i.ReadyAt = &time.Time{}
	*i.ReadyAt = time.Now()
	return nil
}

func (i *Item) CancelItem() {
	i.Status = StatusItemCanceled
	i.CanceledAt = &time.Time{}
	*i.CanceledAt = time.Now()
}

func (i *Item) CanAddAdditionalItems() bool {
	if i.Status != StatusItemStaging && i.Status != StatusItemPending {
		return false
	}

	return true
}

func (i *Item) CalculateTotalPrice() float64 {
	totalPriceItemAndAdditionals := i.TotalPrice

	for _, additionalItem := range i.AdditionalItems {
		totalPriceItemAndAdditionals += additionalItem.TotalPrice
	}

	return totalPriceItemAndAdditionals
}
