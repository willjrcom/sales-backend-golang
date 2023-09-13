package itementity

import (
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type GroupItem struct {
	entity.Entity
	Items      []Item                  `bun:"rel:has-many,join:id=group_item_id"`
	CategoryID uuid.UUID               `bun:"column:category_id,type:uuid,notnull"`
	Category   *productentity.Category `bun:"rel:belongs-to"`
	Size       float64                 `bun:"size"`
	OrderID    uuid.UUID               `bun:"column:order_id,type:uuid,notnull"`
}

func (i *GroupItem) CalculateQuantity() float64 {
	total := 0.0
	for _, item := range i.Items {
		total += item.Quantity
	}
	return total
}

func (i *GroupItem) CalculatePrice() float64 {
	total := 0.0
	for _, item := range i.Items {
		total += item.Price * item.Quantity
	}
	return total
}
