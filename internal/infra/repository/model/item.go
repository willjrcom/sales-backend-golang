package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type Item struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:order_items"`
	ItemCommonAttributes
}

type ItemCommonAttributes struct {
	Name            string    `bun:"name,notnull"`
	Observation     string    `bun:"observation"`
	Price           float64   `bun:"price,notnull"`
	TotalPrice      float64   `bun:"total_price,notnull"`
	Size            string    `bun:"size,notnull"`
	Quantity        float64   `bun:"quantity,notnull"`
	GroupItemID     uuid.UUID `bun:"group_item_id,type:uuid"`
	CategoryID      uuid.UUID `bun:"column:category_id,type:uuid,notnull"`
	AdditionalItems []Item    `bun:"m2m:item_to_additional,join:Item=AdditionalItem"`
	RemovedItems    []string  `bun:"removed_items,type:jsonb"`
	ProductID       uuid.UUID `bun:"product_id,type:uuid"`
}

func (i *Item) FromDomain(item *orderentity.Item) {
	if item == nil {
		return
	}
	*i = Item{
		Entity: entitymodel.Entity{ID: item.ID},
		ItemCommonAttributes: ItemCommonAttributes{
			Name:            item.Name,
			Observation:     item.Observation,
			Price:           item.Price,
			TotalPrice:      item.TotalPrice,
			Size:            item.Size,
			Quantity:        item.Quantity,
			GroupItemID:     item.GroupItemID,
			CategoryID:      item.CategoryID,
			AdditionalItems: []Item{},
			RemovedItems:    item.RemovedItems,
			ProductID:       item.ProductID,
		},
	}

	for _, additionalItem := range item.AdditionalItems {
		ai := Item{}
		ai.FromDomain(&additionalItem)
		i.AdditionalItems = append(i.AdditionalItems, ai)
	}
}

func (i *Item) ToDomain() *orderentity.Item {
	if i == nil {
		return nil
	}
	return &orderentity.Item{
		Entity: i.Entity.ToDomain(),
		ItemCommonAttributes: orderentity.ItemCommonAttributes{
			Name:            i.Name,
			Observation:     i.Observation,
			Price:           i.Price,
			TotalPrice:      i.TotalPrice,
			Size:            i.Size,
			Quantity:        i.Quantity,
			GroupItemID:     i.GroupItemID,
			CategoryID:      i.CategoryID,
			AdditionalItems: []orderentity.Item{},
			RemovedItems:    i.RemovedItems,
			ProductID:       i.ProductID,
		},
	}
}
