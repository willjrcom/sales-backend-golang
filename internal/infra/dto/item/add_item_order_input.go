package itemdto

import (
	"errors"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type AddItemOrderInput struct {
	OrderID     uuid.UUID  `json:"order_id"`
	ProductID   uuid.UUID  `json:"product_id"`
	GroupItemID *uuid.UUID `json:"group_item_id"`
	Quantity    *float64   `json:"quantity"`
	Observation string     `json:"observation"`
}

func (a *AddItemOrderInput) validate() error {
	if a.OrderID == uuid.Nil {
		return errors.New("order id is required")
	}

	if a.ProductID == uuid.Nil {
		return errors.New("item id is required")
	}

	if a.Quantity == nil {
		return errors.New("quantity is required")
	}

	return nil
}

func (a *AddItemOrderInput) ToModel(product *productentity.Product) (item *itementity.Item, err error) {
	if err = a.validate(); err != nil {
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
