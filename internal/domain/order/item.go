package orderentity

import (
	"errors"
	"strings"

	"github.com/google/uuid"
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
	ItemCommonAttributes
}

type ItemCommonAttributes struct {
	Name            string
	Observation     string
	Price           float64
	TotalPrice      float64
	Size            string
	Quantity        float64
	GroupItemID     uuid.UUID
	CategoryID      uuid.UUID
	AdditionalItems []Item
	RemovedItems    []string
	ProductID       uuid.UUID
}

func NewItem(name string, price float64, quantity float64, size string, productID uuid.UUID, categoryID uuid.UUID) *Item {
	itemAdditionalCommonAttributes := ItemCommonAttributes{
		Name:       name,
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
	if strings.Contains(i.Name, "("+i.Size+")") {
		return
	}

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
