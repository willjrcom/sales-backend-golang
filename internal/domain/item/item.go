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
	bun.BaseModel `bun:"table:order_items"`
	ItemCommonAttributes
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

type ItemCommonAttributes struct {
	Name            string    `bun:"name,notnull" json:"name"`
	Observation     string    `bun:"observation" json:"observation"`
	Price           float64   `bun:"price,notnull" json:"price"`
	TotalPrice      float64   `bun:"total_price,notnull" json:"total_price"`
	Size            string    `bun:"size,notnull" json:"size"`
	Quantity        float64   `bun:"quantity,notnull" json:"quantity"`
	GroupItemID     uuid.UUID `bun:"group_item_id,type:uuid" json:"group_item_id"`
	CategoryID      uuid.UUID `bun:"column:category_id,type:uuid,notnull" json:"category_id"`
	AdditionalItems []Item    `bun:"m2m:item_to_additional,join:Item=AdditionalItem" json:"item_to_additional,omitempty"`
	RemovedItems    []string  `bun:"removed_items,type:jsonb" json:"removed_items,omitempty"`
	ProductID       uuid.UUID `bun:"product_id,type:uuid" json:"product_id"`
}

func NewItem(name string, price float64, quantity float64, size string, productID uuid.UUID, categoryID uuid.UUID) *Item {
	itemAdditionalCommonAttributes := ItemCommonAttributes{
		Name:       name + " (" + size + ")",
		Price:      price,
		TotalPrice: price * quantity,
		Size:       size,
		Quantity:   quantity,
		ProductID:  productID,
		CategoryID: categoryID,
	}

	return &Item{
		Entity:               entity.NewEntity(),
		ItemCommonAttributes: itemAdditionalCommonAttributes,
	}
}

func (i *Item) AddSizeToName() {
	i.Name += " (" + i.Size + ")"
}

func (i *Item) AddRemovedItem(name string) {
	i.RemovedItems = append(i.RemovedItems, name)
}

func (i *Item) RemoveRemovedItem(name string) {
	for index, item := range i.RemovedItems {
		if item == name {
			i.RemovedItems = append(i.RemovedItems[:index], i.RemovedItems[index+1:]...)
		}
	}
}

func (i *Item) CalculateTotalPrice() float64 {
	totalPriceItemAndAdditionals := i.TotalPrice

	for _, additionalItem := range i.AdditionalItems {
		totalPriceItemAndAdditionals += additionalItem.TotalPrice
	}

	return totalPriceItemAndAdditionals
}
