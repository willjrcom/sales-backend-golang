package itemdto

import (
	"errors"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	ErrGroupItemNotStaging      = errors.New("group item not staging")
	ErrGroupItemCategoryInvalid = errors.New("group item category invalid")
	ErrGroupItemSizeInvalid     = errors.New("group item size invalid")
)

type OrderItemCreateDTO struct {
	OrderID     uuid.UUID  `json:"order_id"`
	ProductID   uuid.UUID  `json:"product_id"`
	QuantityID  uuid.UUID  `json:"quantity_id"`
	GroupItemID *uuid.UUID `json:"group_item_id"`
	Observation string     `json:"observation"`
}

func (a *OrderItemCreateDTO) Validate() error {
	if a.OrderID == uuid.Nil {
		return errors.New("order id is required")
	}

	if a.ProductID == uuid.Nil {
		return errors.New("item id is required")
	}

	if a.QuantityID == uuid.Nil {
		return errors.New("quantity id is required")
	}

	return nil
}

func (a *OrderItemCreateDTO) validateInternal(product *productentity.Product, groupItem *orderentity.GroupItem, quantity *productentity.Quantity) error {
	if a.QuantityID != quantity.ID {
		return errors.New("quantity id is invalid")
	}

	if groupItem.Status != orderentity.StatusGroupStaging {
		return ErrGroupItemNotStaging
	}

	if groupItem.CategoryID != product.CategoryID {
		return ErrGroupItemCategoryInvalid
	}

	if groupItem.Size != product.Size.Name {
		return ErrGroupItemSizeInvalid
	}
	return nil
}

func (a *OrderItemCreateDTO) ToDomain(product *productentity.Product, groupItem *orderentity.GroupItem, quantity *productentity.Quantity) (item *orderentity.Item, err error) {
	if err = a.validateInternal(product, groupItem, quantity); err != nil {
		return
	}

	item = orderentity.NewItem(product.Name, product.Price, quantity.Quantity, product.Size.Name, product.ID, product.CategoryID)
	item.AddSizeToName()
	item.GroupItemID = *a.GroupItemID
	item.Observation = a.Observation
	return
}
