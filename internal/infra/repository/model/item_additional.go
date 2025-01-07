package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type ItemToAdditional struct {
	bun.BaseModel    `bun:"table:item_to_additional"`
	ItemID           uuid.UUID `bun:"item_id,type:uuid,pk"`
	Item             *Item     `bun:"rel:belongs-to,join:item_id=id"`
	AdditionalItemID uuid.UUID `bun:"type:uuid,pk"`
	AdditionalItem   *Item     `bun:"rel:belongs-to,join:additional_item_id=id"`
	ProductID        uuid.UUID `bun:"type:uuid,pk"`
}

func (i *ItemToAdditional) FromDomain(itemToAdditional *orderentity.ItemToAdditional) {
	if itemToAdditional == nil {
		return
	}
	*i = ItemToAdditional{
		ItemID:           itemToAdditional.ItemID,
		AdditionalItemID: itemToAdditional.AdditionalItemID,
		ProductID:        itemToAdditional.ProductID,
		Item:             &Item{},
		AdditionalItem:   &Item{},
	}

	i.Item.FromDomain(itemToAdditional.Item)
	i.AdditionalItem.FromDomain(itemToAdditional.AdditionalItem)
}

func (i *ItemToAdditional) ToDomain() *orderentity.ItemToAdditional {
	if i == nil {
		return nil
	}
	return &orderentity.ItemToAdditional{
		ItemID:           i.ItemID,
		Item:             i.Item.ToDomain(),
		AdditionalItemID: i.AdditionalItemID,
		AdditionalItem:   i.AdditionalItem.ToDomain(),
		ProductID:        i.ProductID,
	}
}
