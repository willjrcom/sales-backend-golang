package orderentity

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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
	Price           decimal.Decimal
	TotalPrice      decimal.Decimal
	Size            string
	Quantity        float64
	GroupItemID     uuid.UUID
	CategoryID      uuid.UUID
	IsAdditional    bool
	AdditionalItems []Item
	RemovedItems    []string
	ProductID       uuid.UUID
}

// NewItem creates a new order item with initial price and total
func NewItem(name string, price decimal.Decimal, quantity float64, size string, productID uuid.UUID, categoryID uuid.UUID) *Item {
	itemAdditionalCommonAttributes := ItemCommonAttributes{
		Name:       name,
		Price:      price,
		TotalPrice: price.Mul(decimal.NewFromFloat(quantity)),
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

// CalculateTotalPrice computes the total price including additional items, updates TotalPrice, and returns it
func (i *Item) CalculateTotalPrice() decimal.Decimal {
	total := i.Price.Mul(decimal.NewFromFloat(i.Quantity))

	for _, additionalItem := range i.AdditionalItems {
		total = total.Add(additionalItem.TotalPrice)
	}
	i.TotalPrice = total
	return total
}
