package itemdto

import (
	"errors"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
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
	GroupItemID *uuid.UUID `json:"group_item_id"`
	Quantity    *float64   `json:"quantity"`
	Observation string     `json:"observation"`
}

func (a *AddItemOrderInput) validate(product *productentity.Product, groupItem *groupitementity.GroupItem) error {
	if a.OrderID == uuid.Nil {
		return errors.New("order id is required")
	}

	if a.ProductID == uuid.Nil {
		return errors.New("item id is required")
	}

	if a.Quantity == nil {
		return errors.New("quantity is required")
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

func (a *AddItemOrderInput) ToModel(product *productentity.Product, groupItem *groupitementity.GroupItem) (item *itementity.Item, err error) {
	if err = a.validate(product, groupItem); err != nil {
		return
	}

	item = &itementity.Item{
		Entity:      entity.NewEntity(),
		Name:        product.Name,
		Price:       product.Price * (*a.Quantity),
		Description: product.Description,
		Quantity:    *a.Quantity,
		Observation: a.Observation,
		GroupItemID: *a.GroupItemID,
	}

	return
}
