package productcategorydto

import (
	"errors"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	ErrSKURequired          = errors.New("sku is required")
	ErrNameRequired         = errors.New("name is required")
	ErrCostGreaterThanPrice = errors.New("cost must be greater than Price")
	ErrCategoryRequired     = errors.New("category is required")
	ErrSizeRequired         = errors.New("size is required")
)

type ProductCreateDTO struct {
	SKU         string                      `json:"sku"`
	Name        string                      `json:"name"`
	Flavors     []string                    `json:"flavors"`
	Description string                      `json:"description"`
	IsActive    *bool                       `json:"is_active"`
	CategoryID  *uuid.UUID                  `json:"category_id"`
	ImagePath   string                      `json:"image_path"`
	Variations  []ProductVariationCreateDTO `json:"variations"`
}

func (p *ProductCreateDTO) validate() error {
	if p.SKU == "" {
		return ErrSKURequired
	}
	if p.Name == "" {
		return ErrNameRequired
	}
	if p.CategoryID == nil || len(p.CategoryID.String()) == 0 {
		return ErrCategoryRequired
	}

	if len(p.Variations) == 0 {
		return errors.New("at least one variation is required")
	}

	for _, v := range p.Variations {
		if v.Price.LessThan(v.Cost) {
			return ErrCostGreaterThanPrice
		}
		if v.SizeID == uuid.Nil {
			return ErrSizeRequired
		}
	}

	return nil
}

func (p *ProductCreateDTO) ToDomain() (*productentity.Product, error) {
	if err := p.validate(); err != nil {
		return nil, err
	}

	flavors := p.Flavors
	if len(flavors) == 0 {
		flavors = []string{}
	}

	isActive := true
	if p.IsActive != nil {
		isActive = *p.IsActive
	}

	productCommonAttributes := productentity.ProductCommonAttributes{
		SKU:         p.SKU,
		Name:        p.Name,
		Flavors:     flavors,
		Description: p.Description,
		IsActive:    isActive,
		CategoryID:  *p.CategoryID,
		ImagePath:   &p.ImagePath,
	}

	product := productentity.NewProduct(productCommonAttributes)

	for _, v := range p.Variations {
		product.AddVariation(v.ToDomain())
	}

	return product, nil
}
