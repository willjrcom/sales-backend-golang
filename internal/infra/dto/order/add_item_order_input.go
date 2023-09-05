package orderdto

import (
	itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type AddItemOrderInput struct {
	Quantity    float64               `json:"quantity"`
	Status      itementity.StatusItem `json:"status"`
	Observation string                `json:"observation"`
}

func (a *AddItemOrderInput) ToModel(product *productentity.Product) *itementity.Item {
	return &itementity.Item{
		Name:        product.Name,
		Price:       product.Price,
		Description: product.Description,
		Quantity:    a.Quantity,
		Status:      a.Status,
		Observation: a.Observation,
	}
}
