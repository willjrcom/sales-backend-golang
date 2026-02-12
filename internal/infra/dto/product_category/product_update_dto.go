package productcategorydto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type ProductUpdateDTO struct {
	SKU         *string          `json:"sku"`
	Name        *string          `json:"name"`
	Flavors     []string         `json:"flavors,omitempty"`
	ImagePath   *string          `json:"image_path"`
	Description *string          `json:"description"`
	Price       *decimal.Decimal `json:"price"`
	Cost        *decimal.Decimal `json:"cost"`
	IsAvailable *bool            `json:"is_available"`
	IsActive    *bool            `json:"is_active"`
	CategoryID  *uuid.UUID       `json:"category_id"`
	SizeID      *uuid.UUID       `json:"size_id"`
}

func (p *ProductUpdateDTO) Validate(product *productentity.Product) error {
	if product.SKU == "" {
		return ErrSKURequired
	}
	if product.Name == "" {
		return ErrNameRequired
	}
	if product.Price.LessThan(product.Cost) {
		return ErrCostGreaterThanPrice
	}

	return nil
}

func (p *ProductUpdateDTO) UpdateDomain(product *productentity.Product) (err error) {
	if p.SKU != nil {
		product.SKU = *p.SKU
	}
	if p.Name != nil {
		product.Name = *p.Name
	}
	if p.Flavors != nil {
		if len(p.Flavors) == 0 {
			product.Flavors = []string{}
		} else {
			product.Flavors = append([]string{}, p.Flavors...)
		}
	}
	if p.ImagePath != nil {
		product.ImagePath = p.ImagePath
	}
	if p.Description != nil {
		product.Description = *p.Description
	}
	if p.Price != nil {
		product.Price = *p.Price
	}
	if p.Cost != nil {
		product.Cost = *p.Cost
	}
	if p.IsAvailable != nil {
		product.IsAvailable = *p.IsAvailable
	}
	if p.IsActive != nil {
		product.IsActive = *p.IsActive
	}
	if p.CategoryID != nil {
		product.Category.ID = *p.CategoryID
	}
	if p.SizeID != nil {
		product.SizeID = *p.SizeID
	}

	if err = p.Validate(product); err != nil {
		return err
	}

	return nil
}
