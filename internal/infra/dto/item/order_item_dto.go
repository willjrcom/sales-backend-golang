package itemdto

import (
	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type ItemDTO struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Observation     string    `json:"observation"`
	Price           float64   `json:"price"`
	TotalPrice      float64   `json:"total_price"`
	Size            string    `json:"size"`
	Quantity        float64   `json:"quantity"`
	GroupItemID     uuid.UUID `json:"group_item_id"`
	CategoryID      uuid.UUID `json:"category_id"`
	AdditionalItems []ItemDTO `json:"additional_items"`
	RemovedItems    []string  `json:"removed_items"`
	ProductID       uuid.UUID `json:"product_id"`
}

func (i *ItemDTO) FromDomain(item *orderentity.Item) {
	if item == nil {
		return
	}
	*i = ItemDTO{
		ID:              item.ID,
		Name:            item.Name,
		Observation:     item.Observation,
		Price:           item.Price,
		TotalPrice:      item.TotalPrice,
		Size:            item.Size,
		Quantity:        item.Quantity,
		GroupItemID:     item.GroupItemID,
		CategoryID:      item.CategoryID,
		RemovedItems:    item.RemovedItems,
		ProductID:       item.ProductID,
		AdditionalItems: []ItemDTO{},
	}

	for _, additionalItem := range item.AdditionalItems {
		itemDTO := ItemDTO{}
		itemDTO.FromDomain(&additionalItem)
		i.AdditionalItems = append(i.AdditionalItems, itemDTO)
	}

	if len(item.AdditionalItems) == 0 {
		i.AdditionalItems = nil
	}
}
