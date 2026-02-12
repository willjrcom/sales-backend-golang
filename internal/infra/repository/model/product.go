package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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
	SKU         string           `bun:"sku,notnull"`
	Name        string           `bun:"name,notnull"`
	Flavors     []string         `bun:"flavors,type:jsonb,notnull"`
	ImagePath   *string          `bun:"image_path"`
	Description string           `bun:"description"`
	Price       decimal.Decimal  `bun:"price,type:decimal(10,2),notnull"`
	Cost        decimal.Decimal  `bun:"cost,type:decimal(10,2)"`
	IsAvailable bool             `bun:"is_available"`
	IsActive    bool             `bun:"column:is_active,type:boolean"`
	CategoryID  uuid.UUID        `bun:"column:category_id,type:uuid,notnull"`
	Category    *ProductCategory `bun:"rel:belongs-to"`
	SizeID      uuid.UUID        `bun:"size_id,type:uuid,notnull"`
	Size        *Size            `bun:"rel:belongs-to"`
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
			Price:       product.Price,
			Cost:        product.Cost,
			IsAvailable: product.IsAvailable,
			IsActive:    product.IsActive,
			CategoryID:  product.CategoryID,
			Category:    &ProductCategory{},
			SizeID:      product.SizeID,
			Size:        &Size{},
		},
	}

	p.Category.FromDomain(product.Category)
	p.Size.FromDomain(product.Size)
}

func (p *Product) ToDomain() *productentity.Product {
	if p == nil {
		return nil
	}
	return &productentity.Product{
		Entity: p.Entity.ToDomain(),
		ProductCommonAttributes: productentity.ProductCommonAttributes{
			SKU:         p.SKU,
			Name:        p.Name,
			Flavors:     cloneFlavors(p.Flavors),
			ImagePath:   p.ImagePath,
			Description: p.Description,
			Price:       p.Price,
			Cost:        p.Cost,
			IsAvailable: p.IsAvailable,
			IsActive:    p.IsActive,
			CategoryID:  p.CategoryID,
			Category:    p.Category.ToDomain(),
			SizeID:      p.SizeID,
			Size:        p.Size.ToDomain(),
		},
	}
}

func cloneFlavors(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}

	cloned := make([]string, len(values))
	copy(cloned, values)
	return cloned
}
