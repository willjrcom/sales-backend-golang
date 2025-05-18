package productcategorydto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type ProductUpdateDTO struct {
	Code        *string          `json:"code"`
	Name        *string          `json:"name"`
	Flavors     []string         `json:"flavors,omitempty"`
	ImagePath   *string          `json:"image_path"`
	Description *string          `json:"description"`
	Price       *decimal.Decimal `json:"price"`
	Cost        *decimal.Decimal `json:"cost"`
	IsAvailable *bool            `json:"is_available"`
	CategoryID  *uuid.UUID       `json:"category_id"`
	SizeID      *uuid.UUID       `json:"size_id"`
}

func (p *ProductUpdateDTO) Validate(product *productentity.Product) error {
	if product.Code == "" {
		return ErrCodeRequired
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
	if p.Code != nil {
		product.Code = *p.Code
	}
	if p.Name != nil {
		product.Name = *p.Name
	}
	if p.Flavors != nil {
		product.Flavors = p.Flavors
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
