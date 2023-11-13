package productdto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type UpdateProductInput struct {
	Code        *string    `json:"code"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	SizeID      *uuid.UUID `json:"size_id"`
	Price       *float64   `json:"price"`
	Cost        *float64   `json:"cost"`
	CategoryID  *uuid.UUID `json:"category_id"`
	IsAvailable *bool      `json:"is_available"`
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

	return nil
}

func (p *UpdateProductInput) UpdateModel(product *productentity.Product) (err error) {
	if p.Code != nil {
		product.Code = *p.Code
	}
	if p.Name != nil {
		product.Name = *p.Name
	}
	if p.Description != nil {
		product.Description = *p.Description
	}
	if p.SizeID != nil {
		product.SizeID = *p.SizeID
	}
	if p.Price != nil {
		product.Price = *p.Price
	}
	if p.Cost != nil {
		product.Cost = *p.Cost
	}
	if p.CategoryID != nil {
		product.Category.ID = *p.CategoryID
	}
	if p.IsAvailable != nil {
		product.IsAvailable = *p.IsAvailable
	}

	if err = p.Validate(product); err != nil {
		return err
	}

	return nil
}
