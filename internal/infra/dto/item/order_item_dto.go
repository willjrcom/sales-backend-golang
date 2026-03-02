package itemdto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	productcategorydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product_category"
)

type ItemDTO struct {
	ID              uuid.UUID                      `json:"id"`
	Name            string                         `json:"name"`
	Observation     string                         `json:"observation"`
	SubTotal        decimal.Decimal                `json:"sub_total"`
	Total           decimal.Decimal                `json:"total"`
	Size            string                         `json:"size"`
	Quantity        float64                        `json:"quantity"`
	GroupItemID     uuid.UUID                      `json:"group_item_id"`
	CategoryID      uuid.UUID                      `json:"category_id"`
	AdditionalItems []ItemDTO                      `json:"additional_items"`
	RemovedItems    []string                       `json:"removed_items"`
	ProductID       uuid.UUID                      `json:"product_id"`
	Product         *productcategorydto.ProductDTO `json:"product"`
	Flavor          *string                        `json:"flavor,omitempty"`
}

func (i *ItemDTO) FromDomain(item *orderentity.Item) {
	if item == nil {
		return
	}
	*i = ItemDTO{
		ID:              item.ID,
		Name:            item.Name,
		Observation:     item.Observation,
		SubTotal:        item.SubTotal,
		Total:           item.Total,
		Size:            item.Size,
		Quantity:        item.Quantity,
		GroupItemID:     item.GroupItemID,
		CategoryID:      item.CategoryID,
		RemovedItems:    item.RemovedItems,
		ProductID:       item.ProductID,
		Flavor:          item.Flavor,
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

	if item.Product != nil {
		i.Product = &productcategorydto.ProductDTO{}
		i.Product.FromDomain(item.Product)
	}
}
