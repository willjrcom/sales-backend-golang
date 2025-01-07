package groupitemdto

import (
	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	itemdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/item"
	productcategorydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product_category"
)

type GroupItemDTO struct {
	ID               uuid.UUID                       `json:"id"`
	Size             string                          `json:"size"`
	Status           orderentity.StatusGroupItem     `json:"status"`
	TotalPrice       float64                         `json:"total_price"`
	Quantity         float64                         `json:"quantity"`
	NeedPrint        bool                            `json:"need_print"`
	UseProcessRule   bool                            `json:"use_process_rule"`
	Observation      string                          `json:"observation"`
	CategoryID       uuid.UUID                       `json:"category_id"`
	Category         *productcategorydto.CategoryDTO `json:"category"`
	ComplementItemID *uuid.UUID                      `json:"complement_item_id"`
	ComplementItem   *itemdto.ItemDTO                `json:"complement_item"`
	Items            []itemdto.ItemDTO               `json:"items"`
	OrderID          uuid.UUID                       `json:"order_id"`
}

func (i *GroupItemDTO) FromDomain(groupItem *orderentity.GroupItem) {
	if groupItem == nil {
		return
	}
	*i = GroupItemDTO{
		ID:               groupItem.ID,
		Size:             groupItem.Size,
		Status:           groupItem.Status,
		TotalPrice:       groupItem.TotalPrice,
		Quantity:         groupItem.Quantity,
		NeedPrint:        groupItem.NeedPrint,
		UseProcessRule:   groupItem.UseProcessRule,
		Observation:      groupItem.Observation,
		CategoryID:       groupItem.CategoryID,
		ComplementItemID: groupItem.ComplementItemID,
		OrderID:          groupItem.OrderID,
	}

	i.Category.FromDomain(groupItem.Category)
	i.ComplementItem.FromDomain(groupItem.ComplementItem)

	for i, additionalItem := range i.Items {
		additionalItem.FromDomain(&groupItem.Items[i])
	}
}
