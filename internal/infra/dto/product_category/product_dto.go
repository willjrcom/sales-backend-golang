package productcategorydto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type ProductDTO struct {
	ID          uuid.UUID             `json:"id"`
	SKU         string                `json:"sku"`
	Name        string                `json:"name"`
	Flavors     []string              `json:"flavors"`
	ImagePath   *string               `json:"image_path"`
	Description string                `json:"description"`
	IsActive    bool                  `json:"is_active"`
	CategoryID  uuid.UUID             `json:"category_id"`
	Category    *CategoryDTO          `json:"category"`
	Variations  []ProductVariationDTO `json:"variations"`
}

func (p *ProductDTO) FromDomain(product *productentity.Product) {
	if product == nil {
		return
	}

	*p = ProductDTO{
		ID:          product.ID,
		SKU:         product.SKU,
		Name:        product.Name,
		Flavors:     append([]string{}, product.Flavors...),
		ImagePath:   product.ImagePath,
		Description: product.Description,
		IsActive:    product.IsActive,
		CategoryID:  product.CategoryID,
		Category:    &CategoryDTO{},
		Variations:  []ProductVariationDTO{},
	}

	p.Category.FromDomain(product.Category)

	for _, v := range product.Variations {
		dto := ProductVariationDTO{}
		dto.FromDomain(v)
		p.Variations = append(p.Variations, dto)
	}

	if len(p.Flavors) == 0 {
		p.Flavors = []string{}
	}

	if product.Category == nil {
		p.Category = nil
	}
}
