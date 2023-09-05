package productdto

import (
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type UpdateProductInput struct {
	Code        *string  `json:"code"`
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Size        *string  `json:"size"`
	Price       *float64 `json:"price"`
	Cost        *float64 `json:"cost"`
	Category    *string  `json:"category"`
	IsAvailable *bool    `json:"is_available"`
}

func (p *UpdateProductInput) Validate(product *productentity.Product) error {
	if product.Code == "" {
		return ErrCodeRequired
	}
	if product.Name == "" {
		return ErrNameRequired
	}
	if product.Price < product.Cost {
		return ErrCostGreaterThanPrice
	}

	if product.Category.Name == "" {
		return ErrCategoryRequired
	}

	return nil
}

func (p *UpdateProductInput) UpdateModel(product *productentity.Product) error {
	if p.Code != nil {
		product.Code = *p.Code
	}
	if p.Name != nil {
		product.Name = *p.Name
	}
	if p.Description != nil {
		product.Description = *p.Description
	}
	if p.Size != nil {
		product.Size = *p.Size
	}
	if p.Price != nil {
		product.Price = *p.Price
	}
	if p.Cost != nil {
		product.Cost = *p.Cost
	}
	if p.Category != nil {
		product.Category.Name = *p.Category
	}
	if p.IsAvailable != nil {
		product.IsAvailable = *p.IsAvailable
	}

	if err := p.Validate(product); err != nil {
		return err
	}

	return nil
}
