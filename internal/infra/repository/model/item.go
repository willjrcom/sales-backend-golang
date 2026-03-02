package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type Item struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:order_items,alias:item"`
	ItemCommonAttributes
}

type ItemCommonAttributes struct {
	Name               string           `bun:"name,notnull"`
	Observation        string           `bun:"observation"`
	SubTotal           *decimal.Decimal `bun:"sub_total,type:decimal(10,2),notnull"`
	Total              *decimal.Decimal `bun:"total,type:decimal(10,2),notnull"`
	Size               string           `bun:"size,notnull"`
	Quantity           float64          `bun:"quantity,notnull"`
	GroupItemID        uuid.UUID        `bun:"group_item_id,type:uuid"`
	CategoryID         uuid.UUID        `bun:"column:category_id,type:uuid,notnull"`
	IsAdditional       bool             `bun:"is_additional"`
	AdditionalItems    []Item           `bun:"m2m:item_to_additional,join:Item=AdditionalItem"`
	RemovedItems       []string         `bun:"removed_items,type:jsonb"`
	ProductID          uuid.UUID        `bun:"product_id,type:uuid,notnull"`
	Product            *Product         `bun:"rel:has-one,join:product_id=id"`
	ProductVariationID uuid.UUID        `bun:"product_variation_id,type:uuid,notnull"`
	Flavor             *string          `bun:"flavor"`
}

func (i *Item) FromDomain(item *orderentity.Item) {
	if item == nil {
		return
	}
	*i = Item{
		Entity: entitymodel.FromDomain(item.Entity),
		ItemCommonAttributes: ItemCommonAttributes{
			Name:               item.Name,
			Observation:        item.Observation,
			SubTotal:           &item.SubTotal,
			Total:              &item.Total,
			Size:               item.Size,
			Quantity:           item.Quantity,
			GroupItemID:        item.GroupItemID,
			CategoryID:         item.CategoryID,
			IsAdditional:       item.IsAdditional,
			AdditionalItems:    []Item{},
			RemovedItems:       item.RemovedItems,
			ProductID:          item.ProductID,
			ProductVariationID: item.ProductVariationID,
			Product:            &Product{},
			Flavor:             item.Flavor,
		},
	}

	for _, additionalItem := range item.AdditionalItems {
		ai := Item{}
		ai.FromDomain(&additionalItem)
		i.AdditionalItems = append(i.AdditionalItems, ai)
	}

	i.Product.FromDomain(item.Product)
}

func (i *Item) ToDomain() *orderentity.Item {
	if i == nil {
		return nil
	}
	item := &orderentity.Item{
		Entity: i.Entity.ToDomain(),
		ItemCommonAttributes: orderentity.ItemCommonAttributes{
			Name:               i.Name,
			Observation:        i.Observation,
			SubTotal:           i.GetSubTotal(),
			Total:              i.GetTotal(),
			Size:               i.Size,
			Quantity:           i.Quantity,
			GroupItemID:        i.GroupItemID,
			CategoryID:         i.CategoryID,
			IsAdditional:       i.IsAdditional,
			AdditionalItems:    []orderentity.Item{},
			RemovedItems:       i.RemovedItems,
			ProductID:          i.ProductID,
			ProductVariationID: i.ProductVariationID,
			Product:            i.Product.ToDomain(),
			Flavor:             i.Flavor,
		},
	}

	for _, additionalItem := range i.AdditionalItems {
		item.AdditionalItems = append(item.AdditionalItems, *additionalItem.ToDomain())
	}

	return item
}

func (i *Item) GetSubTotal() decimal.Decimal {
	if i.SubTotal == nil {
		return decimal.Zero
	}
	return *i.SubTotal
}

func (i *Item) GetTotal() decimal.Decimal {
	if i.Total == nil {
		return decimal.Zero
	}
	return *i.Total
}
