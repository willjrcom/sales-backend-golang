package itementity

import (
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Items struct {
	entity.Entity
	Products []Item
	Category string
	Size     float64
	OrderID  uuid.UUID
}

func (i *Items) calculateQuantity() float64 {
	total := 0.0
	for _, item := range i.Products {
		total += item.Quantity
	}
	return total
}
