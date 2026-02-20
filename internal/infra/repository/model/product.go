package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type Product struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:products,alias:product"`
	ProductCommonAttributes
}

type ProductCommonAttributes struct {
	SKU         string              `bun:"sku,notnull"`
	Name        string              `bun:"name,notnull"`
	Flavors     []string            `bun:"flavors,type:jsonb,notnull"`
	ImagePath   *string             `bun:"image_path"`
	Description string              `bun:"description"`
	IsActive    bool                `bun:"column:is_active,type:boolean"`
	CategoryID  uuid.UUID           `bun:"column:category_id,type:uuid,notnull"`
	Category    *ProductCategory    `bun:"rel:belongs-to"`
	Variations  []*ProductVariation `bun:"rel:has-many,join:id=product_id"`
}

func (p *Product) FromDomain(product *productentity.Product) {
	if product == nil {
		return
	}
	*p = Product{
		Entity: entitymodel.FromDomain(product.Entity),
		ProductCommonAttributes: ProductCommonAttributes{
			SKU:         product.SKU,
			Name:        product.Name,
			Flavors:     cloneFlavors(product.Flavors),
			ImagePath:   product.ImagePath,
			Description: product.Description,
			IsActive:    product.IsActive,
			CategoryID:  product.CategoryID,
			Category:    &ProductCategory{},
		},
	}

	for _, v := range product.Variations {
		variation := &ProductVariation{}
		variation.FromDomain(v)
		p.Variations = append(p.Variations, variation)
	}

	if product.Category != nil {
		p.Category.FromDomain(product.Category)
	}
}

func (p *Product) ToDomain() *productentity.Product {
	if p == nil {
		return nil
	}

	domain := &productentity.Product{
		Entity: p.Entity.ToDomain(),
		ProductCommonAttributes: productentity.ProductCommonAttributes{
			SKU:         p.SKU,
			Name:        p.Name,
			Flavors:     cloneFlavors(p.Flavors),
			ImagePath:   p.ImagePath,
			Description: p.Description,
			IsActive:    p.IsActive,
			CategoryID:  p.CategoryID,
			Category:    p.Category.ToDomain(),
		},
		Variations: []productentity.ProductVariation{},
	}

	for _, v := range p.Variations {
		domain.Variations = append(domain.Variations, v.ToDomain())
	}

	return domain
}

func cloneFlavors(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}

	cloned := make([]string, len(values))
	copy(cloned, values)
	return cloned
}
