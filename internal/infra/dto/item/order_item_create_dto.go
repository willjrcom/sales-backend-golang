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
	VariationID uuid.UUID  `json:"variation_id"`
	Quantity    float64    `json:"quantity"`
	GroupItemID *uuid.UUID `json:"group_item_id"`
	Observation string     `json:"observation"`
	Flavor      *string    `json:"flavor,omitempty"`
}

func (a *OrderItemCreateDTO) Validate() error {
	if a.OrderID == uuid.Nil {
		return errors.New("order id is required")
	}

	if a.ProductID == uuid.Nil {
		return errors.New("product id is required")
	}

	if a.VariationID == uuid.Nil {
		return errors.New("variation id is required")
	}

	if a.Quantity == 0 {
		return errors.New("quantity is required")
	}

	return nil
}

func (a *OrderItemCreateDTO) validateInternal(product *productentity.Product, variation *productentity.ProductVariation, groupItem *orderentity.GroupItem) error {
	if groupItem.Status != orderentity.StatusGroupStaging {
		return ErrGroupItemNotStaging
	}

	if groupItem.CategoryID != product.CategoryID {
		return ErrGroupItemCategoryInvalid
	}

	if groupItem.Size != variation.Size.Name {
		return ErrGroupItemSizeInvalid
	}

	flavor, err := NormalizeFlavor(a.Flavor, product.Flavors)
	if err != nil {
		return err
	}
	a.Flavor = flavor

	return nil
}

func (a *OrderItemCreateDTO) ToDomain(product *productentity.Product, variation *productentity.ProductVariation, groupItem *orderentity.GroupItem, quantity float64) (item *orderentity.Item, err error) {
	if err = a.validateInternal(product, variation, groupItem); err != nil {
		return
	}

	item = orderentity.NewItem(product.Name, variation.Price, quantity, variation.Size.Name, product.ID, variation.ID, product.CategoryID, a.Flavor)
	item.AddSizeToName()
	item.GroupItemID = *a.GroupItemID
	item.Observation = a.Observation
	return
}
