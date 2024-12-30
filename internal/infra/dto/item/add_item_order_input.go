package itemdto

import (
	"errors"

	"github.com/google/uuid"
	groupitementity "github.com/willjrcom/sales-backend-go/internal/domain/group_item"
	itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	ErrGroupItemNotStaging      = errors.New("group item not staging")
	ErrGroupItemCategoryInvalid = errors.New("group item category invalid")
	ErrGroupItemSizeInvalid     = errors.New("group item size invalid")
)

type AddItemOrderInput struct {
	OrderID     uuid.UUID  `json:"order_id"`
	ProductID   uuid.UUID  `json:"product_id"`
	QuantityID  uuid.UUID  `json:"quantity_id"`
	GroupItemID *uuid.UUID `json:"group_item_id"`
	Observation string     `json:"observation"`
}

func (a *AddItemOrderInput) Validate() error {
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

func (a *AddItemOrderInput) validateInternal(product *productentity.Product, groupItem *groupitementity.GroupItem, quantity *productentity.Quantity) error {
	if a.QuantityID != quantity.ID {
		return errors.New("quantity id is invalid")
	}

	if groupItem.Status != groupitementity.StatusGroupStaging {
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

func (a *AddItemOrderInput) ToModel(product *productentity.Product, groupItem *groupitementity.GroupItem, quantity *productentity.Quantity) (item *itementity.Item, err error) {
	if err = a.validateInternal(product, groupItem, quantity); err != nil {
		return
	}

	item = itementity.NewItem(product.Name, product.Price, quantity.Quantity, product.Size.Name, product.ID, product.CategoryID)
	item.AddSizeToName()
	item.GroupItemID = *a.GroupItemID
	item.Observation = a.Observation
	item.Description = product.Description
	return
}
