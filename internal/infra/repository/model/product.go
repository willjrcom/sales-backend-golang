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
	Code        string           `bun:"code,notnull"`
	Name        string           `bun:"name,notnull"`
	Flavors     []string         `bun:"flavors,type:jsonb"`
	ImagePath   *string          `bun:"image_path"`
	Description string           `bun:"description"`
	Price       decimal.Decimal  `bun:"price,type:decimal(10,2),notnull"`
	Cost        decimal.Decimal  `bun:"cost,type:decimal(10,2)"`
	IsAvailable bool             `bun:"is_available"`
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
			Code:        product.Code,
			Name:        product.Name,
			Flavors:     product.Flavors,
			ImagePath:   product.ImagePath,
			Description: product.Description,
			Price:       product.Price,
			Cost:        product.Cost,
			IsAvailable: product.IsAvailable,
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
			Code:        p.Code,
			Name:        p.Name,
			Flavors:     p.Flavors,
			ImagePath:   p.ImagePath,
			Description: p.Description,
			Price:       p.Price,
			Cost:        p.Cost,
			IsAvailable: p.IsAvailable,
			CategoryID:  p.CategoryID,
			Category:    p.Category.ToDomain(),
			SizeID:      p.SizeID,
			Size:        p.Size.ToDomain(),
		},
	}
}
